// Copyright ©2021 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"image/color"
	"image/png"
	"os"
	"strconv"

	"star-tex.org/x/tex/dvi"
	"star-tex.org/x/tex/dvi/dvimg"
	"star-tex.org/x/tex/kpath"
)

var (
	_ dvi.Renderer = (*pngRenderer)(nil)
)

type pngRenderer struct {
	name string
	dvi  *dvimg.Renderer
	err  error
}

func newPNG(ctx kpath.Context, name string) *pngRenderer {
	return &pngRenderer{
		name: name,
		dvi:  dvimg.New(ctx),
	}
}

func (pr *pngRenderer) setErr(err error) {
	if pr.Err() != nil {
		return
	}
	pr.err = err
}

func (pr *pngRenderer) Err() error {
	err := pr.dvi.Err()
	if err != nil {
		return err
	}
	return pr.err
}

func (pr *pngRenderer) Close() error { return pr.Err() }

func (pr *pngRenderer) Init(pre *dvi.CmdPre, post *dvi.CmdPost) {
	pr.dvi.Init(pre, post)
}

func (pr *pngRenderer) BOP(bop *dvi.CmdBOP) {
	pr.dvi.BOP(bop)
}

func (pr *pngRenderer) EOP() {
	pr.dvi.EOP()
	if pr.Err() != nil {
		return
	}

	name := pr.name[:len(pr.name)-len(".png")] + "_" + strconv.Itoa(pr.dvi.Page()) + ".png"
	f, err := os.Create(name)
	if err != nil {
		pr.setErr(fmt.Errorf("could not create output PNG file: %w", err))
		return
	}
	defer f.Close()

	err = png.Encode(f, pr.dvi.Image())
	if err != nil {
		pr.setErr(fmt.Errorf("could not encode PNG image: %w", err))
		return
	}

	err = f.Close()
	if err != nil {
		pr.setErr(fmt.Errorf("could not close output PNG file %q: %w", name, err))
		return
	}
}

func (pr *pngRenderer) DrawGlyph(x, y int32, font dvi.Font, glyph rune, c color.Color) {
	pr.dvi.DrawGlyph(x, y, font, glyph, c)
}

func (pr *pngRenderer) DrawRule(x, y, w, h int32, c color.Color) {
	pr.dvi.DrawRule(x, y, w, h, c)
}
