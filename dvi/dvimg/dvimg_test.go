// Copyright ©2025 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dvimg_test

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
	"testing"

	"git.sr.ht/~sbinet/cmpimg"
	"star-tex.org/x/tex"
	"star-tex.org/x/tex/dvi"
	"star-tex.org/x/tex/dvi/dvimg"
	"star-tex.org/x/tex/kpath"
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
			err = tex.Run(buf, src)
			if err != nil {
				t.Fatal(err)
			}

			ctx := kpath.New()

			rnd := dvimg.New(ctx,
				dvimg.WithBackground(color.Transparent),
				dvimg.WithBackground(color.White),
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

			err = rnd.Err()
			if err != nil {
				t.Fatalf("could not render DVI program: %v", err)
			}

			var (
				gotname  = strings.Replace(name, ".tex", ".png", 1)
				wantname = strings.Replace(name, ".tex", "_golden.png", 1)
			)

			img, err := os.Create(gotname)
			if err != nil {
				t.Fatalf("could not create output image: %v", err)
			}
			defer img.Close()
			err = png.Encode(img, rnd.Image().(*image.RGBA).SubImage(rnd.Bounds()))
			if err != nil {
				t.Fatalf("could not encode output image: %v", err)
			}
			err = img.Close()
			if err != nil {
				t.Fatalf("could not save output image: %v", err)
			}

			got, err := os.ReadFile(gotname)
			if err != nil {
				t.Fatalf("could not read generated image: %v", err)
			}

			want, err := os.ReadFile(wantname)
			if err != nil {
				t.Fatalf("could not read reference image: %v", err)
			}

			ok, err := cmpimg.Equal("png", got, want)
			if err != nil {
				t.Fatalf("could not compare images: %v", err)
			}

			if !ok {
				// store DVI.
				dviname := strings.Replace(name, ".tex", ".dvi", 1)
				_ = os.WriteFile(dviname, buf.Bytes(), 0644)
				t.Fatalf("invalid images for %q", name)
			}
			_ = os.Remove(gotname)
		})
	}
}
