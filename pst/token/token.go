// Copyright ©2021 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package token defines constants representing the lexical tokens of the
// PostScript programming language and basic operations on tokens (printing,
// predicates).
package token // import "star-tex.org/x/tex/pst/token"

//go:generate stringer -type Token

// Token is the set of lexical tokens of the PostScript programming language.
type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	COMMENT // %

	// Identifiers and basic type literals
	// (these tokens stand for classes of literals)
	NAME   // main
	INT    // 12345 16#1234
	FLOAT  // 123.45
	STRING // (abc) <0f21a3> <~...~>

	LPAREN // (
	LBRACK // [
	LBRACE // {

	RPAREN // )
	RBRACK // ]
	RBRACE // }

	LT     // <
	GT     // >
	LSHIFT // <<
	RSHIFT // >>

	SLASH      // /
	SLASHSLASH // //
)
