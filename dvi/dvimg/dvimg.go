// Copyright ©2021 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package dvimg provides a DVI renderer that renders to an image.
package dvimg // import "star-tex.org/x/tex/dvi/dvimg"

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"

	"golang.org/x/image/font"
	xfix "golang.org/x/image/math/fixed"

	"star-tex.org/x/tex/dvi"
	"star-tex.org/x/tex/font/fixed"
	"star-tex.org/x/tex/font/pkf"
	"star-tex.org/x/tex/kpath"
)

const (
	shrink = 1
)

var (
	_ dvi.Renderer = (*Renderer)(nil)
)

// Renderer renders a DVI document to an image.Image.
type Renderer struct {
	page int

	bkg color.Color

	pre   dvi.CmdPre
	post  dvi.CmdPost
	conv  float32 // converts DVI units to pixels
	tconv float32 // converts unmagnified DVI units to pixels
	dpi   float32

	ctx   kpath.Context
	faces map[fntkey]font.Face

	sub image.Rectangle // subset of image containing (only) drawn rules and glyphes

	img draw.Image
	err error
}

// Option customizes a dvimg renderer
type Option func(r *Renderer)

// WithBackground configures the renderer to use the provided color c as
// a background for the resulting image.
func WithBackground(c color.Color) Option {
	return func(r *Renderer) {
		r.bkg = c
	}
}

// New creates a new image renderer, with the provided kpathsea context.
func New(ctx kpath.Context, opts ...Option) *Renderer {
	r := &Renderer{ctx: ctx, faces: make(map[fntkey]font.Face)}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Page returns the current page number of the DVI document.
func (pr *Renderer) Page() int {
	return pr.page
}

// Image returns the current image representing the DVI document.
func (pr *Renderer) Image() image.Image {
	return pr.img
}

// Bounds returns the bounds of the current image containing only drawn rules
// and glyphes.
func (pr *Renderer) Bounds() image.Rectangle {
	return pr.sub
}

// Err returns the current error status of the renderer, if any.
func (pr *Renderer) Err() error {
	return pr.err
}

func (pr *Renderer) setErr(err error) {
	if pr.err != nil {
		return
	}
	pr.err = err
}

func (pr *Renderer) Init(pre *dvi.CmdPre, post *dvi.CmdPost) {
	pr.pre = *pre
	pr.post = *post
	if pr.dpi == 0 {
		pr.dpi = 600 // FIXME(sbinet): infer from fonts?
	}
	res := pr.dpi
	conv := float32(pr.pre.Num) / 254000.0 * (res / float32(pr.pre.Den))
	pr.tconv = conv
	pr.conv = conv * float32(pr.pre.Mag) / 1000.0

	//	conv = 1/(float32(pre.Num)/float32(pre.Den)*
	//		(float32(pre.Mag)/1000.0)*
	//		(pr.dpi*shrink/254000.0)) + 0.5

	if pr.bkg == nil {
		pr.bkg = color.White
	}
}

func (pr *Renderer) BOP(bop *dvi.CmdBOP) {
	if pr.err != nil {
		return
	}

	pr.page = int(bop.C0)
	pr.sub = image.Rectangle{}

	bnd := image.Rect(0, 0,
		int(pr.pixels(int32(pr.post.Width))),
		int(pr.pixels(int32(pr.post.Height))),
	)
	pr.img = image.NewRGBA(bnd)
	draw.Draw(pr.img, bnd, image.NewUniform(pr.bkg), image.Point{}, draw.Over)
}

func (pr *Renderer) EOP() {
	if pr.err != nil {
		return
	}
}

func (pr *Renderer) DrawGlyph(x, y int32, font dvi.Font, glyph rune, c color.Color) {
	if pr.err != nil {
		return
	}

	dot := xfix.Point26_6{
		X: xfix.I(int(pr.pixels(x))),
		Y: xfix.I(int(pr.pixels(y))),
	}

	face, ok := pr.face(font)
	if !ok {
		return
	}

	dr, mask, maskp, _, ok := face.Glyph(dot, glyph)
	if !ok {
		pr.setErr(fmt.Errorf("could not find glyph 0x%02x", glyph))
		return
	}

	draw.DrawMask(
		pr.img, dr,
		image.NewUniform(c), image.Point{},
		mask, maskp, draw.Over,
	)
	pr.sub = pr.sub.Union(dr)
}

func (pr *Renderer) DrawRule(x, y, w, h int32, c color.Color) {
	if pr.err != nil {
		return
	}

	r := image.Rect(
		int(pr.pixels(x+0)), int(pr.pixels(y+0)),
		int(pr.pixels(x+w)), int(pr.pixels(y-h)),
	)

	draw.Draw(pr.img, r, image.NewUniform(c), image.Point{}, draw.Over)
	pr.sub = pr.sub.Union(r)
}

func roundF32(v float32) int32 {
	if v > 0 {
		return int32(v + 0.5)
	}
	return int32(v - 0.5)
}

func (pr *Renderer) pixels(v int32) int32 {
	x := pr.conv * float32(v)
	return roundF32(x / shrink)
}

// func (pr *pngRenderer) rulepixels(v int32) int32 {
// 	x := int32(pr.conv * float32(v))
// 	if float32(x) < pr.conv*float32(v) {
// 		return x + 1
// 	}
// 	return x
// }

type fntkey struct {
	name string
	size fixed.Int12_20
}

func (pr *Renderer) face(fnt dvi.Font) (font.Face, bool) {
	key := fntkey{
		name: fnt.Name(),
		size: fnt.Size(),
	}
	if f, ok := pr.faces[key]; ok {
		return f, ok
	}

	fname, err := pr.ctx.Find(fnt.Name() + ".pk")
	if err != nil {
		log.Printf("could not find font face %q: %+v", fnt.Name(), err)
		name := "cmr10"
		log.Printf("replacing with %q", name)
		fname, err = pr.ctx.Find(name + ".pk")
	}
	if err != nil {
		pr.setErr(fmt.Errorf("could not find font face %q: %+v", fnt.Name(), err))
		return nil, false
	}

	f, err := pr.ctx.Open(fname)
	if err != nil {
		pr.setErr(fmt.Errorf("could not open font face %q: %+v", fnt.Name(), err))
		return nil, false
	}
	defer f.Close()

	pk, err := pkf.Parse(f)
	if err != nil {
		pr.setErr(fmt.Errorf("could not parse font face %q: %+v", fnt.Name(), err))
		return nil, false
	}

	tfm := fnt.Metrics()

	if tfm.Checksum() != pk.Checksum() {
		pr.setErr(fmt.Errorf(
			"TFM and PK checksum do not match for %q: tfm=0x%x, pk=0x%x",
			fnt.Name(),
			tfm.Checksum(),
			pk.Checksum(),
		))
		return nil, false
	}

	face := pkf.NewFace(pk, tfm, &pkf.FaceOptions{
		Size: tfm.DesignSize().Float64(),
		DPI:  float64(pr.dpi),
	})
	pr.faces[key] = face

	return face, true
}
