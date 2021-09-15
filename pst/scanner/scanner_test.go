// Copyright ©2021 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scanner

import (
	"reflect"
	"testing"

	"star-tex.org/x/tex/pst/token"
)

func TestScanner(t *testing.T) {
	for _, tc := range []struct {
		src  string
		want []Token
	}{
		{
			" % comment\n",
			[]Token{{token.COMMENT, " comment"}},
		},
		{
			"%comment (with) %",
			[]Token{{token.COMMENT, "comment (with) %"}},
		},
		{
			" name",
			[]Token{{token.NAME, "name"}},
		},
		{
			"name$123",
			[]Token{{token.NAME, "name$123"}},
		},
		{
			"name ",
			[]Token{{token.NAME, "name"}},
		},
		{
			"name\n",
			[]Token{{token.NAME, "name"}},
		},
		{
			"name%comment",
			[]Token{
				{token.NAME, "name"},
				{token.COMMENT, "comment"},
			},
		},
		{
			"a.b",
			[]Token{{token.NAME, "a.b"}},
		},
		{
			"a-b",
			[]Token{{token.NAME, "a-b"}},
		},
		{
			"23A",
			[]Token{{token.NAME, "23A"}},
		},
		{
			"12-34",
			[]Token{{token.NAME, "12-34"}},
		},
		{
			"$foo",
			[]Token{{token.NAME, "$foo"}},
		},
		{
			"$$",
			[]Token{{token.NAME, "$$"}},
		},
		{
			"@pattern",
			[]Token{{token.NAME, "@pattern"}},
		},
		{
			"/",
			[]Token{{token.SLASH, ""}},
		},
		{
			"/ ",
			[]Token{{token.SLASH, ""}},
		},
		{
			"/name",
			[]Token{
				{token.SLASH, ""},
				{token.NAME, "name"},
			},
		},
		{
			"//name",
			[]Token{
				{token.SLASHSLASH, ""},
				{token.NAME, "name"},
			},
		},
		{
			" (a (sub) string)",
			[]Token{{token.STRING, "a (sub) string"}},
		},
		{
			`(a \(sub string)`,
			[]Token{{token.STRING, `a \(sub string`}},
		},
		{
			"(a \r\n string)",
			[]Token{{token.STRING, "a \r\n string"}},
		},
		{
			// FIXME(sbinet): check with GS whether that's the correct interpretation.
			"<>",
			[]Token{
				{token.LT, ""},
				{token.GT, ""},
			},
		},
		{
			// FIXME(sbinet): check with GS whether that's the correct interpretation.
			"< >",
			[]Token{
				{token.LT, ""},
				{token.GT, ""},
			},
		},
		{
			"<901fa3>",
			[]Token{{token.STRING, string([]byte{0x90, 0x1f, 0xa3})}},
		},
		{
			"<901fa>",
			[]Token{{token.STRING, string([]byte{0x90, 0x1f, 0xa0})}},
		},
		{
			// FIXME(sbinet): check with GS whether that's the correct interpretation.
			"< 901fa >",
			[]Token{
				{token.LT, ""},
				{token.NAME, "901fa"},
				{token.GT, ""},
			},
		},
		{
			"<cafe babe\ndead beef>\n",
			[]Token{{token.STRING, string([]byte{
				0xca, 0xfe, 0xba, 0xbe,
				0xde, 0xad, 0xbe, 0xef,
			})}},
		},
		{
			"<~BOu!rD]j7BEbo7~>",
			[]Token{{token.STRING, "hello world"}},
		},
		{
			"1 2 3",
			[]Token{
				{token.INT, "1"},
				{token.INT, "2"},
				{token.INT, "3"},
			},
		},
		{
			"1 a 3",
			[]Token{
				{token.INT, "1"},
				{token.NAME, "a"},
				{token.INT, "3"},
			},
		},
		{
			"1.2 2e3 16#901fa",
			[]Token{
				{token.FLOAT, "1.2"},
				{token.FLOAT, "2e3"},
				{token.INT, "16#901fa"},
			},
		},
		{
			"[ 123 /abc (xyz) ]",
			[]Token{
				{token.LBRACK, ""},
				{token.INT, "123"},
				{token.SLASH, ""},
				{token.NAME, "abc"},
				{token.STRING, "xyz"},
				{token.RBRACK, ""},
			},
		},
		{
			"[123 /abc (xyz)]",
			[]Token{
				{token.LBRACK, ""},
				{token.INT, "123"},
				{token.SLASH, ""},
				{token.NAME, "abc"},
				{token.STRING, "xyz"},
				{token.RBRACK, ""},
			},
		},
		{
			"[]",
			[]Token{
				{token.LBRACK, ""},
				{token.RBRACK, ""},
			},
		},
		{
			"{ add 2 div }",
			[]Token{
				{token.LBRACE, ""},
				{token.NAME, "add"},
				{token.INT, "2"},
				{token.NAME, "div"},
				{token.RBRACE, ""},
			},
		},
		{
			"{add 2 div}",
			[]Token{
				{token.LBRACE, ""},
				{token.NAME, "add"},
				{token.INT, "2"},
				{token.NAME, "div"},
				{token.RBRACE, ""},
			},
		},
		{
			"{}",
			[]Token{
				{token.LBRACE, ""},
				{token.RBRACE, ""},
			},
		},
		{
			"<<>>",
			[]Token{
				{token.LSHIFT, ""},
				{token.RSHIFT, ""},
			},
		},
		{
			"<<k1 v1 k2 v2>>",
			[]Token{
				{token.LSHIFT, ""},
				{token.NAME, "k1"},
				{token.NAME, "v1"},
				{token.NAME, "k2"},
				{token.NAME, "v2"},
				{token.RSHIFT, ""},
			},
		},
		{
			"<<\nk1 v1\nk2 v2\n>>",
			[]Token{
				{token.LSHIFT, ""},
				{token.NAME, "k1"},
				{token.NAME, "v1"},
				{token.NAME, "k2"},
				{token.NAME, "v2"},
				{token.RSHIFT, ""},
			},
		},
	} {
		t.Run("", func(t *testing.T) {
			got, err := Scan([]byte(tc.src))
			if err != nil {
				t.Fatalf("could not scan %q: %+v", tc.src, err)
			}
			if len(got) != len(tc.want) {
				t.Errorf("invalid scan length for %q: got=%d, want=%d", tc.src, len(got), len(tc.want))
			}
			n := len(got)
			if n > len(tc.want) {
				n = len(tc.want)
			}
			for i, got := range got[:n] {
				want := tc.want[i]
				if !reflect.DeepEqual(got, want) {
					t.Errorf(
						"invalid scan Token for %q:\ngot[%d]= {%+v, %q}\nwant[%d]={%+v, %q}",
						tc.src,
						i, got.Kind, got.Value,
						i, want.Kind, want.Value,
					)
				}
			}
		})
	}
}

func TestScan(t *testing.T) {
	got, err := Scan([]byte("123 abc (xyz)"))
	if err != nil {
		t.Fatalf("could not scan: %+v", err)
	}
	want := []Token{
		{token.INT, "123"},
		{token.NAME, "abc"},
		{token.STRING, "xyz"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid scan:\ngot= %+v\nwant=%+v", got, want)
	}
}

func TestScannerPS(t *testing.T) {
	const src = `%!PS-AdobeFont-1.0: CMR10 003.002
%%Title: CMR10
%Version: 003.002
%%CreationDate: Mon Jul 13 16:17:00 2009
%%Creator: David M. Jones
%Copyright: Copyright (c) 1997, 2009 American Mathematical Society
%Copyright:  (<http://www.ams.org>), with Reserved Font Name CMR10.
% This Font Software is licensed under the SIL Open Font License, Version 1.1.
% This license is in the accompanying file OFL.txt, and is also
% available with a FAQ at: http://scripts.sil.org/OFL.
%%EndComments

FontDirectory/CMR10 known{/CMR10 findfont dup/UniqueID known{dup
/UniqueID get 5000793 eq exch/FontType get 1 eq and}{pop false}ifelse
{save true}{false}ifelse}{false}ifelse
11 dict begin
/FontType 1 def
/FontMatrix [0.001 0 0 0.001 0 0 ]readonly def
/FontName /CMR10 def
/FontBBox {-40 -250 1009 750 }readonly def
/UniqueID 5000793 def
/PaintType 0 def
/FontInfo 9 dict dup begin
 /version (003.002) readonly def
 /Notice (Copyright \050c\051 1997, 2009 American Mathematical Society \050<http://www.ams.org>\051, with Reserved Font Name CMR10.) readonly def
 /FullName (CMR10) readonly def
 /FamilyName (Computer Modern) readonly def
 /Weight (Medium) readonly def
 /ItalicAngle 0 def
 /isFixedPitch false def
 /UnderlinePosition -100 def
 /UnderlineThickness 50 def
end readonly def
`

	got, err := Scan([]byte(src))
	if err != nil {
		t.Fatalf("could not scan: %+v", err)
	}

	want := []Token{
		{token.COMMENT, "!PS-AdobeFont-1.0: CMR10 003.002"},
		{token.COMMENT, "%Title: CMR10"},
		{token.COMMENT, "Version: 003.002"},
		{token.COMMENT, "%CreationDate: Mon Jul 13 16:17:00 2009"},
		{token.COMMENT, "%Creator: David M. Jones"},
		{token.COMMENT, "Copyright: Copyright (c) 1997, 2009 American Mathematical Society"},
		{token.COMMENT, "Copyright:  (<http://www.ams.org>), with Reserved Font Name CMR10."},
		{token.COMMENT, " This Font Software is licensed under the SIL Open Font License, Version 1.1."},
		{token.COMMENT, " This license is in the accompanying file OFL.txt, and is also"},
		{token.COMMENT, " available with a FAQ at: http://scripts.sil.org/OFL."},
		{token.COMMENT, "%EndComments"},
		{token.NAME, "FontDirectory"},
		{token.SLASH, ""},
		{token.NAME, "CMR10"},
		{token.NAME, "known"},
		{token.LBRACE, ""},
		{token.SLASH, ""},
		{token.NAME, "CMR10"},
		{token.NAME, "findfont"},
		{token.NAME, "dup"},
		{token.SLASH, ""},
		{token.NAME, "UniqueID"},
		{token.NAME, "known"},
		{token.LBRACE, ""},
		{token.NAME, "dup"},
		{token.SLASH, ""},
		{token.NAME, "UniqueID"},
		{token.NAME, "get"},
		{token.INT, "5000793"},
		{token.NAME, "eq"},
		{token.NAME, "exch"},
		{token.SLASH, ""},
		{token.NAME, "FontType"},
		{token.NAME, "get"},
		{token.INT, "1"},
		{token.NAME, "eq"},
		{token.NAME, "and"},
		{token.RBRACE, ""},
		{token.LBRACE, ""},
		{token.NAME, "pop"},
		{token.NAME, "false"},
		{token.RBRACE, ""},
		{token.NAME, "ifelse"},
		{token.LBRACE, ""},
		{token.NAME, "save"},
		{token.NAME, "true"},
		{token.RBRACE, ""},
		{token.LBRACE, ""},
		{token.NAME, "false"},
		{token.RBRACE, ""},
		{token.NAME, "ifelse"},
		{token.RBRACE, ""},
		{token.LBRACE, ""},
		{token.NAME, "false"},
		{token.RBRACE, ""},
		{token.NAME, "ifelse"},
		{token.INT, "11"},
		{token.NAME, "dict"},
		{token.NAME, "begin"},
		{token.SLASH, ""},
		{token.NAME, "FontType"},
		{token.INT, "1"},
		{token.NAME, "def"},
		{token.SLASH, ""},
		{token.NAME, "FontMatrix"},
		{token.LBRACK, ""},
		{token.FLOAT, "0.001"},
		{token.INT, "0"},
		{token.INT, "0"},
		{token.FLOAT, "0.001"},
		{token.INT, "0"},
		{token.INT, "0"},
		{token.RBRACK, ""},
		{token.NAME, "readonly"},
		{token.NAME, "def"},
		{token.SLASH, ""},
		{token.NAME, "FontName"},
		{token.SLASH, ""},
		{token.NAME, "CMR10"},
		{token.NAME, "def"},
		{token.SLASH, ""},
		{token.NAME, "FontBBox"},
		{token.LBRACE, ""},
		{token.INT, "-40"},
		{token.INT, "-250"},
		{token.INT, "1009"},
		{token.INT, "750"},
		{token.RBRACE, ""},
		{token.NAME, "readonly"},
		{token.NAME, "def"},
		{token.SLASH, ""},
		{token.NAME, "UniqueID"},
		{token.INT, "5000793"},
		{token.NAME, "def"},
		{token.SLASH, ""},
		{token.NAME, "PaintType"},
		{token.INT, "0"},
		{token.NAME, "def"},
		{token.SLASH, ""},
		{token.NAME, "FontInfo"},
		{token.INT, "9"},
		{token.NAME, "dict"},
		{token.NAME, "dup"},
		{token.NAME, "begin"},
		{token.SLASH, ""},
		{token.NAME, "version"},
		{token.STRING, "003.002"},
		{token.NAME, "readonly"},
		{token.NAME, "def"},
		{token.SLASH, ""},
		{token.NAME, "Notice"},
		{token.STRING, "Copyright \\050c\\051 1997, 2009 American Mathematical Society \\050<http://www.ams.org>\\051, with Reserved Font Name CMR10."},
		{token.NAME, "readonly"},
		{token.NAME, "def"},
		{token.SLASH, ""},
		{token.NAME, "FullName"},
		{token.STRING, "CMR10"},
		{token.NAME, "readonly"},
		{token.NAME, "def"},
		{token.SLASH, ""},
		{token.NAME, "FamilyName"},
		{token.STRING, "Computer Modern"},
		{token.NAME, "readonly"},
		{token.NAME, "def"},
		{token.SLASH, ""},
		{token.NAME, "Weight"},
		{token.STRING, "Medium"},
		{token.NAME, "readonly"},
		{token.NAME, "def"},
		{token.SLASH, ""},
		{token.NAME, "ItalicAngle"},
		{token.INT, "0"},
		{token.NAME, "def"},
		{token.SLASH, ""},
		{token.NAME, "isFixedPitch"},
		{token.NAME, "false"},
		{token.NAME, "def"},
		{token.SLASH, ""},
		{token.NAME, "UnderlinePosition"},
		{token.INT, "-100"},
		{token.NAME, "def"},
		{token.SLASH, ""},
		{token.NAME, "UnderlineThickness"},
		{token.INT, "50"},
		{token.NAME, "def"},
		{token.NAME, "end"},
		{token.NAME, "readonly"},
		{token.NAME, "def"},
	}

	if got, want := len(got), len(want); got != want {
		t.Errorf("invalid number of tokens: got=%d, want=%d", got, want)
	}
	n := len(got)
	if n > len(want) {
		n = len(want)
	}
	for i, got := range got[:n] {
		want := want[i]
		if !reflect.DeepEqual(got, want) {
			t.Errorf(
				"invalid scan Token for %q:\ngot[%d]= {%+v, %q}\nwant[%d]={%+v, %q}",
				src,
				i, got.Kind, got.Value,
				i, want.Kind, want.Value,
			)
		}
	}
}

func TestErrScan(t *testing.T) {
	for _, tc := range []struct {
		src  string
		want []Token
	}{
		{
			"(str) (invalid",
			[]Token{
				{token.STRING, "str"},
				{token.ILLEGAL, "(invalid"},
			},
		},
		{
			"(str) <0102",
			[]Token{
				{token.STRING, "str"},
				{token.ILLEGAL, "<0102"},
			},
		},
		{
			"(str) <0a02x>",
			[]Token{
				{token.STRING, "str"},
				{token.ILLEGAL, "<0a02x>"},
			},
		},
		{
			"(str) <~0102",
			[]Token{
				{token.STRING, "str"},
				{token.ILLEGAL, "<~0102"},
			},
		},
		{
			"(str) <~xxxx~>",
			[]Token{
				{token.STRING, "str"},
				{token.ILLEGAL, "<~xxxx~>"},
			},
		},
	} {
		t.Run("", func(t *testing.T) {
			got, err := Scan([]byte(tc.src))
			if err == nil {
				t.Fatalf("expected an error")
			}

			if got, want := len(got), len(tc.want); got != want {
				t.Errorf("invalid number of tokens: got=%d, want=%d", got, want)
			}
			n := len(got)
			if n > len(tc.want) {
				n = len(tc.want)
			}
			for i, got := range got[:n] {
				want := tc.want[i]
				if !reflect.DeepEqual(got, want) {
					t.Errorf(
						"invalid scan Token for %q:\ngot[%d]= {%+v, %q}\nwant[%d]={%+v, %q}",
						tc.src,
						i, got.Kind, got.Value,
						i, want.Kind, want.Value,
					)
				}
			}

		})
	}
}
