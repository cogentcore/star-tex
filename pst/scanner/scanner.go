// Copyright ©2021 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package scanner implements a scanner for PostScript source text.
// It takes a []byte as source which can then be tokenized through repeated
// calls to the Scan method.
package scanner // import "star-tex.org/x/tex/pst/scanner"

import (
	"bytes"
	"encoding/ascii85"
	"errors"
	"strconv"
	"strings"

	"star-tex.org/x/tex/pst/token"
)

var errInvalid = errors.New("invalid PostScript program")

// Token is a PostScript token.
type Token struct {
	Kind  token.Token
	Value string
}

// Scan scans the provided PostScript source and returns the list of Tokens.
func Scan(src []byte) ([]Token, error) {
	var toks []Token
	sc := NewScanner(src)
	for {
		tok := sc.Scan()
		switch tok.Kind {
		case token.ILLEGAL:
			if tok.Value == "" {
				tok.Value = string(sc.src[sc.cur:])
			}
			toks = append(toks, tok)
			return toks, errInvalid
		case token.EOF:
			return toks, nil
		}
		toks = append(toks, tok)
	}
	return toks, nil
}

// Scanner scans PostScript source text.
type Scanner struct {
	src []byte // src is the source this Scanner is scanning
	cur int    // cur is the current offset into src
}

// NewScanner creates a Scanner scanning src.
func NewScanner(src []byte) Scanner {
	return Scanner{src: src}
}

// Scan scans the next token and returns the token and its literal string
// if applicable.
// The source end is indicated by token.EOF.
func (sc *Scanner) Scan() Token {
	c, ok := sc.scan()
	if !ok {
		return Token{token.EOF, ""}
	}

	for ok && isWhiteSpace(c) {
		c, ok = sc.scan()
	}
	if !ok {
		return Token{token.EOF, ""}
	}

	switch c {
	case '[':
		return Token{token.LBRACK, ""}
	case ']':
		return Token{token.RBRACK, ""}
	case '{':
		return Token{token.LBRACE, ""}
	case '}':
		return Token{token.RBRACE, ""}
	case '%':
		return sc.readComment()
	case '<':
		c, ok = sc.peek()
		if ok && c == '<' {
			sc.cur++
			return Token{token.LSHIFT, ""}
		}
		if ok && c == '~' {
			sc.cur++
			return sc.readA85String()
		}
		if ok && isHexChar(c) {
			return sc.readHexString()
		}
		return Token{token.LT, ""}

	case '>':
		c, ok = sc.peek()
		if ok && c == '>' {
			sc.cur++
			return Token{token.RSHIFT, ""}
		}
		return Token{token.GT, ""}

	case '(':
		return sc.readLitString()
	case '/':
		c, ok = sc.peek()
		if ok && c == '/' {
			sc.cur++
			return Token{token.SLASHSLASH, ""}
		}
		return Token{token.SLASH, ""}

	default:
		sc.cur--
		pos := sc.cur
		tok, ok := sc.readNumber()
		if ok {
			return tok
		}

		sc.cur = pos // rewind
		return sc.readName()
	}

	panic("impossible")
}

func (sc *Scanner) peek() (byte, bool) {
	i := sc.cur
	if i >= len(sc.src) {
		return 0, false
	}
	return sc.src[i], true
}

func (sc *Scanner) scan() (byte, bool) {
	if sc.cur >= len(sc.src) {
		return 0, false
	}
	v := sc.src[sc.cur]
	sc.cur++
	return v, true
}

func (sc *Scanner) readComment() Token {
	var (
		beg = sc.cur
		end = bytes.Index(sc.src[sc.cur:], []byte("\n")) // FIXME(sbinet): handle \r\n\f
	)
	switch {
	case end < 0:
		end = len(sc.src)
	default:
		end += sc.cur
	}
	sc.cur = end + 1
	return Token{token.COMMENT, string(sc.src[beg:end])}
}

func (sc *Scanner) readLitString() Token {
	var (
		beg  = sc.cur
		nest = 0
		esc  = false
	)
	for nest >= 0 {
		c, ok := sc.scan()
		if !ok {
			sc.cur = beg - 1 // rewind to first '('
			return Token{token.ILLEGAL, ""}
		}
		switch c {
		case '(':
			if !esc {
				nest++
			}
		case ')':
			if !esc {
				nest--
			}
		case '\\':
			esc = true
			continue
		}
		esc = false
	}
	end := sc.cur - 1 // points at last ')'
	return Token{token.STRING, string(sc.src[beg:end])}
}

func (sc *Scanner) readHexString() Token {
	var (
		beg = sc.cur
		end = bytes.Index(sc.src[beg:], []byte(">"))
	)
	if end < 0 {
		sc.cur--
		return Token{token.ILLEGAL, ""}
	}
	end += beg

	src := make([]byte, 0, end-beg)
	for _, v := range sc.src[beg:end] {
		switch v {
		case ' ', '\t', '\n', '\r':
			continue
		}
		src = append(src, v)
	}

	if n := len(src); n%2 != 0 {
		ref := src
		src = make([]byte, n+1)
		copy(src, ref)
	}
	dst := make([]byte, 0, len(src))
	hex := func(c byte) (uint8, bool) {
		switch {
		case c == '\x00':
			return 0, true
		case '0' <= c && c <= '9':
			return c - '0', true
		case 'a' <= c && c <= 'f':
			return c - 'a' + 10, true
		case 'A' <= c && c <= 'F':
			return c - 'A' + 10, true
		}
		return c, false
	}

	for i := 0; i < len(src); i += 2 {
		v0, ok0 := hex(src[i+0])
		v1, ok1 := hex(src[i+1])
		if !ok0 || !ok1 {
			sc.cur--
			return Token{token.ILLEGAL, ""}
		}
		dst = append(dst, v0<<4+v1)
	}
	sc.cur = end + len(">")
	return Token{token.STRING, string(dst)}
}

func (sc *Scanner) readA85String() Token {
	var (
		beg = sc.cur
		end = bytes.Index(sc.src[beg:], []byte("~>"))
	)
	if end < 0 {
		sc.cur -= len("<~")
		return Token{token.ILLEGAL, ""}
	}
	end += beg
	src := sc.src[beg:end]
	dst := make([]byte, len(src))
	n, _, err := ascii85.Decode(dst, src, true)
	if err != nil {
		sc.cur -= len("<~")
		return Token{token.ILLEGAL, ""}
	}
	dst = dst[:n]
	sc.cur = end + len("~>")
	return Token{token.STRING, string(dst)}
}

func (sc *Scanner) readName() Token {
	beg := sc.cur
loop:
	for sc.cur < len(sc.src) {
		c, ok := sc.scan()
		if !ok {
			sc.cur = beg
			return Token{token.ILLEGAL, ""}
		}
		if !isRegular(c) {
			sc.cur--
			break loop
		}
	}
	end := sc.cur
	return Token{token.NAME, string(sc.src[beg:end])}
}

func (sc *Scanner) readNumber() (Token, bool) {
	var (
		beg = sc.cur
		end = bytes.IndexAny(sc.src[beg:], delimiters)
	)
	if end < 0 {
		end = len(sc.src)
	} else {
		end += sc.cur
	}

	var (
		src = sc.src[beg:end]
		tok token.Token
	)
	switch {
	case bytes.Contains(src, []byte("#")):
		toks := strings.Split(string(src), "#")
		if len(toks) != 2 {
			return Token{token.ILLEGAL, ""}, false
		}
		base, err := strconv.Atoi(toks[0])
		if err != nil {
			return Token{token.ILLEGAL, ""}, false
		}
		_, err = strconv.ParseInt(toks[1], base, 64)
		if err != nil {
			return Token{token.ILLEGAL, ""}, false
		}
		tok = token.INT

	case bytes.ContainsAny(src, "eE."):
		_, err := strconv.ParseFloat(string(src), 64)
		if err != nil {
			return Token{token.ILLEGAL, ""}, false
		}
		tok = token.FLOAT

	default:
		_, err := strconv.Atoi(string(src))
		if err != nil {
			return Token{token.ILLEGAL, ""}, false
		}
		tok = token.INT
	}

	sc.cur = end
	return Token{tok, string(src)}, true
}
