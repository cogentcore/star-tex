// Copyright ©2026 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wtex

import (
	"errors"
	"fmt"
	"os"

	wrap "star-tex.org/x/tex/internal/wtex/internal/wrap"
)

const (
	npages = 2500
)

var errInitex = errors.New("initex-dump")

func IniTeX(texmf string) (core []byte, err error) {
	lib := newLib(nil, []byte("*latex.ltx\n\\dump\n\\bye"), texmf)
	lib.initex = true
	defer os.RemoveAll(lib.tmp)

	tex := &Engine{
		m: wrap.New(lib, lib),
	}

	err = tex.Run()
	if err != nil {
		return nil, fmt.Errorf("could not process latex.ltx: %w", err)
	}
	_, _ = os.Stdout.Write([]byte("\n"))

	lib = newLib(lib.mem, []byte("\n&latex\n\\end\n"), texmf)
	tex.m = wrap.New(lib, lib)
	func() {
		defer func() {
			ee := recover()
			switch ee := ee.(type) {
			case error:
				if errors.Is(ee, errInitex) {
					err = nil
					return
				}
			}
			panic(ee)
		}()
		err = tex.Run()
	}()
	if err != nil {
		return nil, fmt.Errorf("could not create latex.fmt: %w", err)
	}

	return lib.mem, nil
}
