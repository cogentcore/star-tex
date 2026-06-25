// Copyright ©2025 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dvipdf_test

import (
	"bytes"
	"image/color"
	"os"
	"strings"
	"testing"

	"git.sr.ht/~sbinet/cmpimg"
	tex "github.com/cogentcore/star-tex"
	"github.com/cogentcore/star-tex/dvi"
	"github.com/cogentcore/star-tex/dvi/dvipdf"
	"github.com/cogentcore/star-tex/kpath"
)

func TestRenderer(t *testing.T) {
	for _, name := range []string{
		"testdata/hello.tex",
		"testdata/with_rule.tex",
		"testdata/full_page.tex",
		"testdata/math.tex",
	} {
		t.Run(name, func(t *testing.T) {
			src, err := os.Open(name)
			if err != nil {
				t.Fatal(err)
			}
			defer src.Close()

			var (
				stdlog = new(bytes.Buffer)
			)

			buf := new(bytes.Buffer)
			err = tex.ToDVI(buf, src)
			if err != nil {
				t.Fatalf(
					"could not process TeX document: %+v",
					err,
				)
			}

			var (
				ctx      = kpath.New()
				gotname  = strings.Replace(name, ".tex", ".pdf", 1)
				wantname = strings.Replace(name, ".tex", "_golden.pdf", 1)
			)

			pdf, err := os.Create(gotname)
			if err != nil {
				t.Fatalf("could not create output PDF: %v", err)
			}
			defer pdf.Close()

			rnd := dvipdf.New(
				ctx, pdf,
				dvipdf.WithBackground(color.White),
				dvipdf.WithEmbedFonts(true),
			)
			vm := dvi.NewMachine(
				dvi.WithContext(ctx),
				dvi.WithLogOutput(stdlog),
				dvi.WithRenderer(rnd),
				dvi.WithHandlers(dvi.NewColorHandler(ctx)),
				dvi.WithOffsetX(0),
				dvi.WithOffsetY(0),
			)

			prog, err := dvi.Compile(buf.Bytes())
			if err != nil {
				t.Fatalf("could not compile DVI program: %v", err)
			}

			err = vm.Run(prog)
			if err != nil {
				t.Fatalf("could not run DVI program: %v", err)
			}

			err = rnd.Close()
			if err != nil {
				t.Fatalf("could not render DVI program: %v", err)
			}

			err = pdf.Close()
			if err != nil {
				t.Fatalf("could not save output PDF: %v", err)
			}

			got, err := os.ReadFile(gotname)
			if err != nil {
				t.Fatalf("could not read generated PDF: %v", err)
			}

			if *cmpimg.GenerateTestData {
				err = os.WriteFile(wantname, got, 0644)
				if err != nil {
					t.Fatalf("could not regenerate reference file for %q: %v", wantname, err)
				}
			}

			want, err := os.ReadFile(wantname)
			if err != nil {
				t.Fatalf("could not read reference PDF: %v", err)
			}

			ok, err := cmpimg.Equal("pdf", got, want)
			if err != nil {
				t.Fatalf("could not compare PDFs: %v", err)
			}

			if !ok {
				// store DVI.
				dviname := strings.Replace(name, ".tex", ".dvi", 1)
				_ = os.WriteFile(dviname, buf.Bytes(), 0644)
				t.Fatalf("invalid PDFs for %q", name)
			}
			_ = os.Remove(gotname)
		})
	}
}
