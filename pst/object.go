// Copyright ©2021 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pst

// Object describes a PostScript object.
type Object interface {
	Kind() Kind
	Attr() Attr
}

type Kind uint8

// Attr describes the attributes of a PostScript object.
// Attributes affect the behavior of a PostScript object when it is executed or
// when certain operations are performed on it.
// Attributes do not affect the behavior of PostScript objects when they are
// treated strictly as data.
type Attr uint8

// Attributes of a PostScript object.
const (
	AttrInvalid    Attr = 0x0
	AttrLiteral    Attr = 0x1
	AttrExecutable Attr = 0x2

	AttrUnlimited Attr = 0x10
	AttrReadOnly  Attr = 0x20
	AttrExecOnly  Attr = 0x40
	AttrNone      Attr = 0x80
)
