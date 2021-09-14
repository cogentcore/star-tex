// Copyright ©2021 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package afm

import (
	"image"
	"math"
	"os"
	"testing"

	"golang.org/x/image/font"
	"star-tex.org/x/tex/kpath"
)

func TestParseCM(t *testing.T) {
	ktx := kpath.New()
	for _, tc := range []struct {
		name string
		want font.Metrics
	}{
		{
			"cmr10.afm",
			font.Metrics{
				Height:     64000,
				Ascent:     44416,
				Descent:    12416,
				XHeight:    27584,
				CapHeight:  43712,
				CaretSlope: image.Point{X: 0, Y: 1},
			},
		},
		{
			"cmitt10.afm",
			font.Metrics{
				Height:     59456,
				Ascent:     39104,
				Descent:    14208,
				XHeight:    27968,
				CapHeight:  39104,
				CaretSlope: image.Point{X: 250069, Y: 1000000},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			name, err := ktx.Find(tc.name)
			if err != nil {
				t.Fatalf("could not find %q file: %+v", tc.name, err)
			}

			f, err := ktx.Open(name)
			if err != nil {
				t.Fatalf("could not open afm: %+v", err)
			}
			defer f.Close()

			fnt, err := Parse(f)
			if err != nil {
				t.Fatalf("could not parse afm: %+v", err)
			}

			got := fnt.Metrics()
			if got != tc.want {
				t.Fatalf("invalid metrics:\ngot= %d\nwant=%d", got, tc.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	for _, tc := range []string{
		"testdata/fake-vertical.afm",
		"testdata/times-with-composites.afm",
	} {
		t.Run(tc, func(t *testing.T) {
			f, err := os.Open(tc)
			if err != nil {
				t.Fatalf("could not open AFM test file: %+v", err)
			}
			defer f.Close()

			_, err = Parse(f)
			if err != nil {
				t.Fatalf("could not parse AFM test file: %+v", err)
			}
		})
	}
}

func TestInt16_16(t *testing.T) {
	const tol = 1e-5
	for _, tc := range []struct {
		str  string
		want float64
	}{
		{"0", 0},
		{"1.0", 1},
		{"1.2", 1.2},
		{"+1.2", +1.2},
		{"-1.2", -1.2},
	} {
		t.Run("", func(t *testing.T) {
			v := fixedFrom(tc.str)
			got := v.Float64()
			if diff := math.Abs(got - tc.want); diff > tol {
				t.Fatalf("invalid 16:16 value: got=%v, want=%v (diff=%e)", got, tc.want, diff)
			}
		})
	}
}
