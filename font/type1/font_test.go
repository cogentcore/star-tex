// Copyright ©2021 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package type1

import (
	"testing"

	"star-tex.org/x/tex/kpath"
)

func TestParse(t *testing.T) {
	ctx := kpath.New()

	type blkDescr struct {
		kind byte
		size int
	}

	for _, tc := range []struct {
		name string
		blks []blkDescr
	}{
		{
			name: "cmr10",
			blks: []blkDescr{
				{kind: markerASCII, size: 4287},
				{kind: markerBinary, size: 30900},
				{kind: markerASCII, size: 545},
			},
		},
		{
			name: "cmitt10",
			blks: []blkDescr{
				{kind: markerASCII, size: 4379},
				{kind: markerBinary, size: 21113},
				{kind: markerASCII, size: 545},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			pfbName, err := ctx.Find(tc.name + ".pfb")
			if err != nil {
				t.Fatalf("could not find type1 font file: %+v", err)
			}

			pfb, err := ctx.Open(pfbName)
			if err != nil {
				t.Fatalf("could not read file %q: %+v", pfbName, err)
			}
			defer pfb.Close()

			afmName, err := ctx.Find(tc.name + ".afm")
			if err != nil {
				t.Fatalf("could not find type1 metrics font file: %+v", err)
			}

			afm, err := ctx.Open(afmName)
			if err != nil {
				t.Fatalf("could not read file %q: %+v", afmName, err)
			}
			defer afm.Close()

			fnt, err := Parse(pfb, afm)
			if err != nil {
				t.Fatalf("could not parse file %q: %+v", pfbName, err)
			}

			if got, want := len(fnt.blks), len(tc.blks); got != want {
				t.Fatalf("invalid number of blocks: got=%d, want=%d", got, want)
			}

			for i, got := range fnt.blks {
				want := tc.blks[i]
				if want.kind != got.kind {
					t.Fatalf("invalid block[%d] kind: got=0x%02x, want=0x%02x", i, got.kind, want.kind)
				}
				if want.size != len(got.data) {
					t.Fatalf("invalid block[%d] size: got=%d, want=%d", i, len(got.data), want.size)
				}
			}

			_ = fnt.Metrics()
		})
	}

}
