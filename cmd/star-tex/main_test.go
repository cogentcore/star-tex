// Copyright ©2021 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"git.sr.ht/~sbinet/cmpimg"
	"star-tex.org/x/tex/kpath"
)

func TestProcess(t *testing.T) {
	for _, name := range []string{
		"../../testdata/hello.tex",
	} {
		t.Run(filepath.Base(name), func(t *testing.T) {
			r, err := os.Open(name)
			if err != nil {
				t.Fatalf("could not open TeX document: %+v", err)
			}
			defer r.Close()
			oname := strings.Replace(name, ".tex", ".pdf", 1)
			_ = os.RemoveAll(oname)

			var (
				o      = new(bytes.Buffer)
				stdout = new(bytes.Buffer)
				stderr = new(bytes.Buffer)
				ktx    = kpath.New()
			)

			err = process(ktx, o, r, stdout, stderr)
			if err != nil {
				t.Fatalf("could not process TeX document: %+v", err)
			}

			want, err := os.ReadFile(strings.Replace(name, ".tex", "_golden.pdf", 1))
			if err != nil {
				t.Fatalf("could not read reference PDF file: %+v", err)
			}

			got := o.Bytes()
			ok, err := cmpimg.Equal("pdf", got, want)
			if err != nil {
				t.Fatalf("could not compare PDFs: %+v", err)
			}
			if !ok {
				_ = os.WriteFile(oname, got, 0644)
				t.Fatalf("PDF files compare different.")
			}
		})
	}
}
