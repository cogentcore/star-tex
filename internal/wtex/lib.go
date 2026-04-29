// Copyright ©2026 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wtex

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"git.sr.ht/~sbinet/overlayfs"
	wrap "star-tex.org/x/tex/internal/wtex/internal/wrap"
	"star-tex.org/x/tex/kpath"
)

type memory struct {
	mem *[]byte
}

func (m *memory) Slice() *[]byte { return m.mem }
func (m *memory) Grow(delta, max int64) int64 {
	panic("not implemented")
}

type xlib struct {
	mem   []byte
	now   func() time.Time
	files []fdescr
	ktx   kpath.Context
	input []byte

	stdout *fdescr

	initex bool

	tmp string
	fs  fs.FS
}

var (
	_ wrap.Xenv     = (*xlib)(nil)
	_ wrap.Xlibrary = (*xlib)(nil)
)

func newLib(core, input []byte, texmf string) *xlib {
	if texmf == "" {
		panic(fmt.Errorf("invalid TEXMF directory: %q", texmf))
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	ufs := overlayfs.From(
		os.DirFS(cwd),
		wrap.FS,
		os.DirFS(texmf),
	)

	ktx, err := kpath.NewFromFS(overlayfs.From(ufs))
	if err != nil {
		panic(err)
	}

	//if len(wrap.Fmt) != 2500*65536 {
	//	panic("boo")
	//}

	var mem = make([]byte, npages*65536)
	switch len(core) {
	case 0:
		// ok.
	default:
		copy(mem, core)
		//mem = bytes.Clone(core)
	}
	if got, want := len(mem), npages*65536; got != want {
		panic(fmt.Errorf("invalid memory size: got=%d, want=%d", got, want))
	}

	tmp, err := os.MkdirTemp("", "wtex-tmp-")
	if err != nil {
		panic(err)
	}

	return &xlib{
		mem: mem,
		now: func() time.Time {
			return time.Now()
		},
		files: []fdescr{},
		ktx:   ktx,
		input: input,

		stdout: &fdescr{
			name:   "stdout",
			stdout: true,
			w:      os.Stdout,
		},

		tmp: tmp,
		fs:  os.DirFS(cwd),
	}
}

func (lib *xlib) Xmemory() wrap.Memory {
	return &memory{&lib.mem}
}

func (lib *xlib) fd(id int32) *fdescr {
	if id < 0 {
		return lib.stdout
	}
	return &lib.files[id]
}

func (lib *xlib) fname(fname string) string {
	fname = strings.TrimRight(fname, " \x00")
	if strings.HasPrefix(fname, "{") {
		fname = strings.TrimPrefix(fname, "{")
		fname = strings.TrimSuffix(fname, "}")
	}

	if strings.HasPrefix(fname, "\"") {
		fname = strings.TrimPrefix(fname, "\"")
		fname = strings.TrimSuffix(fname, "\"")
	}
	fname = strings.TrimRight(fname, " ")
	fname = strings.TrimPrefix(fname, "*")

	fname = strings.TrimPrefix(fname, "TeXfonts:")
	fname = strings.TrimPrefix(fname, "TeXformats:")

	switch fname {
	case "TEX.POOL":
		return "tex.pool"
	}

	return fname
}

func (lib *xlib) find(name string) (string, error) {
	paths, err := lib.ktx.FindAll(name)
	if err == nil {
		return paths[0], nil
	}

	// FIXME(sbinet): ugly hack to get initex+latex to work.
	if lib.initex {
		switch name {
		case "babel-latex.cfg",
			"il2enc.dfu",
			"omlenc.dfu", "omxenc.dfu", "uenc.dfu":
			name = filepath.Join(lib.tmp, name)
			err := os.WriteFile(name, nil, 0644)
			if err != nil {
				return "", err
			}
			return name, nil

		case "./texsys.aux":
			// FIXME(sbinet): figure out how this is actually created.
			err := os.WriteFile(name, []byte(lib.now().Format("2006/01/02:15:04")+"\x0a\x0a"), 0644)
			if err != nil {
				return "", err
			}
			return name, nil
		}
	}

	switch name {
	// FIXME(sbinet): streamline handling of all these
	// "non-existent during first run" files .
	case "texput.aux":
		_, err := fs.ReadFile(lib.fs, name)
		if err != nil {
			_ = os.WriteFile(name, nil, 0644)
		}
	}

	f, err2 := lib.fs.Open(name)
	if err2 != nil {
		return "", err
	}
	defer f.Close()

	return name, nil
}

func (lib *xlib) u8slice(addr, size int32) []byte {
	return lib.mem[addr : addr+size]
}

func (lib *xlib) u32slice(addr, size int32) []uint32 {
	return unsafe.Slice((*uint32)(unsafe.Pointer(&lib.mem[addr])), size)
}

func (lib *xlib) XprintInteger(fd, v int32) {
	f := lib.fd(fd)
	//log.Printf("xxx-int: [%s](%d) → [%s]", f.name, v, []byte(strconv.FormatInt(int64(v), 10)))
	f.Write([]byte(strconv.FormatInt(int64(v), 10)))
}

func (lib *xlib) XprintChar(fd, v int32) {
	f := lib.fd(fd)
	//if !f.stdout {
	//	log.Printf("xxx-char: [%s](%d) → [%s|%04x] | %d", f.name, v, string(rune(v)), v, f.pos)
	//}
	f.Write([]byte{byte(v)})
}

func (lib *xlib) XprintString(fd, addr int32) {
	var (
		vlen = lib.u8slice(addr, 1)[0]
		vbuf = lib.u8slice(addr+1, int32(vlen))
		str  = string(vbuf)
	)
	f := lib.fd(fd)
	//log.Printf("xxx-str: [%s](%d)", f.name, addr)
	f.Write([]byte(str))
}

func (lib *xlib) XprintNewline(fd int32) {
	f := lib.fd(fd)
	//log.Printf("xxx-newline: [%s](undefined)", f.name)
	f.Write([]byte("\n"))
}

func (lib *xlib) Xreset(length, ptr int32) int32 {
	var (
		buf   = lib.u8slice(ptr, length)
		fname = strings.TrimRight(string(buf), "\x000")
	)
	fname = lib.fname(fname)
	//log.Printf("xxx-reset: [%s](%d, %d)", fname, ptr, length)

	switch fname {
	case "tex.pool":
		h, err := wrap.FS.Open(fname)
		if err != nil {
			panic(err)
		}
		raw, err := fs.ReadFile(wrap.FS, fname)
		if err != nil {
			panic(err)
		}
		lib.files = append(lib.files, fdescr{
			name: fname,
			r:    h,
			buf:  raw,
		})
		return int32(len(lib.files) - 1)
	case "TTY:":
		lib.files = append(lib.files, fdescr{
			name:  "stdin",
			stdin: true,
			buf:   lib.input,
		})
		return int32(len(lib.files) - 1)
	}

	path, err := lib.find(fname)
	//if fname == "texsys.aux" || fname == "./texsys.aux" {
	//log.Printf("→ path=[%s]", path)
	//}
	if err != nil {
		lib.files = append(lib.files, fdescr{
			name:   fname,
			erstat: 1,
		})
		return int32(len(lib.files) - 1)
	}

	h, err := lib.ktx.Open(path)
	if err != nil {
		panic(err)
	}
	var raw []byte
	{
		hh, err := lib.ktx.Open(path)
		if err != nil {
			panic(err)
		}
		defer hh.Close()
		raw, err = io.ReadAll(hh)
		if err != nil {
			panic(err)
		}
	}
	lib.files = append(lib.files, fdescr{
		name: fname,
		r:    h,
		buf:  raw,
	})
	return int32(len(lib.files) - 1)
}

func (lib *xlib) Xgetfilesize(length, ptr int32) int32 {
	var (
		buf   = lib.u8slice(ptr, length)
		fname = lib.fname(string(buf))
	)
	//log.Printf("xxx-getfilesize: [%s](%d, %d)", fname, ptr, length)

	path, err := lib.find(fname)
	if err != nil {
		return 0
	}

	f, err := lib.ktx.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}
	sz := fi.Size()
	if sz > math.MaxInt32 {
		panic(fmt.Errorf("file %q too big (%d > %d)", fname, sz, math.MaxInt32))
	}
	return int32(sz)
}

func (lib *xlib) Xinputln(fd, bypassEOLN, bufferp, firstp, lastp, maxBufStackPtr, bufsz int32) int32 {
	var (
		f = lib.fd(fd)

		buf   = lib.u8slice(bufferp, bufsz)
		first = lib.u32slice(firstp, 1)
		last  = lib.u32slice(lastp, 1)
	)
	//log.Printf("xxx-inputln: [%s](%d+%d, %d, %d, %d)[%v, %v, %d]", lib.fd(fd).name, bufferp, bufsz, firstp, lastp, bypassEOLN, f.eof, f.eoln, f.pos2)

	// FIXME(sbinet): this should not be ignored
	//var (
	//	last_nonblank = 0 // |last| with trailing blanks removed
	//max_buf_stack = unsafe.Slice((*uint32)(unsafe.Pointer(&lib.mem[maxBufStackPtr])), 1)
	//)

	// cf.\ Matthew 19\thinspace:\thinspace30
	last[0] = first[0]

	// input the first character of the line into |f^|
	if bypassEOLN != 0 {
		if !f.eof && f.eoln {
			f.pos2++
		}
	}

	var (
		eol = -1
		sub []byte
	)
	if f.pos2 <= len(f.buf) {
		sub = f.buf[f.pos2:]
		eol = bytes.IndexByte(sub, '\n')
	}

	switch {
	case eol < 0:
		eol = len(f.buf)
		sub = f.buf[f.pos2:]
	default:
		eol = f.pos2 + eol
		sub = f.buf[f.pos2:eol]
	}

	switch {
	case f.pos2 >= len(f.buf):
		if f.stdin {
			//lib.callback()
			panic(errInitex)
		}
		f.eof = true
		return 0

	default:
		n := copy(buf[first[0]:], sub)
		last[0] = first[0] + uint32(n)
		for buf[last[0]-1] == ' ' {
			last[0]--
		}
		f.pos2 = int(eol)
		f.eoln = true
	}

	return 1
}

func (lib *xlib) Xrewrite(length, ptr int32) int32 {
	var (
		buf   = lib.u8slice(ptr, length)
		fname = strings.TrimRight(string(buf), " ")
	)
	//log.Printf("xxx-rewrite: [%s](%d, %d)", fname, ptr, length)

	switch fname {
	case "TTY:":
		lib.files = append(lib.files, fdescr{
			name:   "stdout",
			stdout: true,
			w:      os.Stdout,
		})
	default:
		f, err := os.Create(fname)
		if err != nil {
			panic(err)
		}
		lib.files = append(lib.files, fdescr{
			name:    fname,
			writing: true,
			w:       f,
		})
	}
	return int32(len(lib.files) - 1)
}

func (lib *xlib) Xget(fd, ptr, length int32) {
	var (
		f   = lib.fd(fd)
		buf = lib.mem
	)
	//log.Printf("xxx-get: [%s](%d, %d)", f.name, ptr, length)

	switch {
	case f.stdin:
		switch {
		case f.pos >= len(lib.input):
			buf[ptr] = '\r'
		default:
			buf[ptr] = lib.input[f.pos]
		}
	default:
		switch {
		case f.r != nil || f.w != nil:
			var (
				end = min(f.pos+int(length), len(f.buf))
				n   = copy(buf[ptr:], f.buf[f.pos:end])
			)
			//log.Printf("→ copy: n=%d, ptr=%d, pos=%d, eoc=%d", n, ptr, f.pos, end)
			if n == 0 {
				buf[ptr] = 0
				f.eof = true
				f.eoln = true
				return
			}
		default:
			//log.Printf("→ no file descriptor")
			f.eof = true
			f.eoln = true
			return
		}
	}

	f.eoln = false
	switch buf[ptr] {
	case '\n', '\r':
		f.eoln = true
	}
	f.pos += int(length)
}

func (lib *xlib) Xput(fd, ptr, length int32) {
	var (
		f   = lib.fd(fd)
		buf = lib.u8slice(ptr, length)
	)
	//log.Printf("xxx-put: [%s](%d, %d, writing=%v)", f.name, ptr, length, f.writing)
	if f.writing {
		f.out = append(f.out, buf...)
	}
}

func (lib *xlib) Xeof(fd int32) int32 {
	var (
		f = lib.fd(fd)
	)
	//log.Printf("xxx-eof: [%s] → [%v]", f.name, f.eof)
	switch {
	case f.eof:
		return 1
	default:
		return 0
	}
}

func (lib *xlib) Xeoln(fd int32) int32 {
	var (
		f = lib.fd(fd)
	)
	//log.Printf("xxx-eoln: [%s] → [%v]", f.name, f.eoln)
	switch {
	case f.eoln:
		return 1
	default:
		return 0
	}
}

func (lib *xlib) Xerstat(fd int32) int32 {
	var (
		f = lib.fd(fd)
	)
	//log.Printf("xxx-erstat: [%s] → [%d]", f.name, f.erstat)
	return int32(f.erstat)
}

func (lib *xlib) Xclose(fd int32) {
	var (
		f = lib.fd(fd)
	)
	//log.Printf("xxx-close: [%s]", f.name)

	if f.writing {
		f.Write(f.out)
	}

	err := f.Close()
	if err != nil {
		panic(err)
	}
}

func (lib *xlib) XgetCurrentMinutes() int32 {
	now := lib.now()
	return int32(60*now.Hour() + now.Minute())
}

func (lib *xlib) XgetCurrentDay() int32 {
	now := lib.now()
	return int32(now.Day())
}

func (lib *xlib) XgetCurrentMonth() int32 {
	now := lib.now()
	return int32(now.Month())
}

func (lib *xlib) XgetCurrentYear() int32 {
	now := lib.now()
	return int32(now.Year())
}

func (lib *xlib) Xtex_final_end() {
	lib.XprintNewline(-1)
}
