// Copyright ©2021 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"image/color"
	"os"

	"star-tex.org/x/tex/dvi"
	"star-tex.org/x/tex/dvi/dvipdf"
	"star-tex.org/x/tex/kpath"
)

var (
	_ dvi.Renderer = (*pdfRenderer)(nil)
)

type pdfRenderer struct {
	pdf *dvipdf.Renderer

	out *os.File
	err error
}

func newPDF(ctx kpath.Context, oname string) (*pdfRenderer, error) {
	f, err := os.Create(oname)
	if err != nil {
		return nil, err
	}

	return &pdfRenderer{
		pdf: dvipdf.New(ctx, f),
		out: f,
	}, nil
}

func (pr *pdfRenderer) setErr(err error) {
	if pr.pdf.Err() != nil {
		return
	}
	pr.err = err
}

func (pr *pdfRenderer) Err() error {
	err := pr.pdf.Err()
	if err != nil {
		return err
	}
	return pr.err
}

func (pr *pdfRenderer) Close() error {
	defer pr.out.Close()

	if pr.Err() != nil {
		return pr.Err()
	}

	err := pr.pdf.Close()
	if err != nil {
		err = fmt.Errorf("could not close output PDF document: %w", err)
		pr.setErr(err)
		return err
	}

	err = pr.out.Close()
	if err != nil {
		err = fmt.Errorf("could not close output PDF file: %w", err)
		pr.setErr(err)
		return err
	}

	return nil
}

func (pr *pdfRenderer) Init(pre *dvi.CmdPre, post *dvi.CmdPost) {
	pr.pdf.Init(pre, post)
}

func (pr *pdfRenderer) BOP(bop *dvi.CmdBOP) {
	pr.pdf.BOP(bop)
}

func (pr *pdfRenderer) EOP() {
	pr.pdf.EOP()
}

func (pr *pdfRenderer) DrawGlyph(x, y int32, font dvi.Font, glyph rune, c color.Color) {
	pr.pdf.DrawGlyph(x, y, font, glyph, c)
}

func (pr *pdfRenderer) DrawRule(x, y, w, h int32, c color.Color) {
	pr.pdf.DrawRule(x, y, w, h, c)
}
