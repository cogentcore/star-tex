// Copyright ©2026 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wtex // import "star-tex.org/x/tex/internal/wtex"

import (
	"io"
	"os"

	wrap "star-tex.org/x/tex/internal/wtex/internal/wrap"
)

func New(input []byte, texmf string, texfmt []byte) *Engine {
	lib := newLib(texfmt, input, texmf)
	return &Engine{
		m:   wrap.New(lib, lib),
		tmp: lib.tmp,
	}
}

type Engine struct {
	m   *wrap.Module
	tmp string
}

func (vm *Engine) Run() error {
	defer os.RemoveAll(vm.tmp)

	vm.m.Xmain()
	return nil
}

type fdescr struct {
	name    string
	stdin   bool
	stdout  bool
	writing bool
	pos     int
	pos2    int
	erstat  int
	eof     bool
	eoln    bool
	buf     []byte
	out     []byte

	r io.ReadCloser
	w io.WriteCloser
}

func (f *fdescr) Write(p []byte) (int, error) {
	return f.w.Write(p)
}

func (f *fdescr) Read(p []byte) (int, error) {
	return f.r.Read(p)
}

func (f *fdescr) Close() (err error) {
	switch {
	case f.r != nil:
		err = f.r.Close()
		f.r = nil
	case f.w != nil:
		err = f.w.Close()
		f.w = nil
	}
	return err
}
