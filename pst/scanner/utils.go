// Copyright ©2021 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scanner

var delimiters = string([]byte{
	0, '\t', '\n', '\f', '\r', ' ',
})

func isWhiteSpace(c byte) bool {
	switch c {
	case 0, '\t', '\n', '\f', '\r', ' ':
		return true
	}
	return false
}

func isSpecial(c byte) bool {
	switch c {
	case '(', ')', '<', '>', '[', ']', '{', '}', '/', '%':
		return true
	}
	return false
}

func isRegular(c byte) bool {
	switch {
	case isWhiteSpace(c):
		return false
	case isSpecial(c):
		return false
	}
	return true
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// func isASCIILetter(c byte) bool {
// 	switch {
// 	case 'a' <= c && c <= 'z':
// 		return true
// 	case 'A' <= c && c <= 'Z':
// 		return true
// 	}
// 	return false
// }

func isHexChar(c byte) bool {
	switch {
	case isDigit(c):
		return true
	case 'a' <= c && c <= 'f':
		return true
	case 'A' <= c && c <= 'F':
		return true
	}
	return false
}
