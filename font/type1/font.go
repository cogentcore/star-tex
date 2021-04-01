// Copyright ©2021 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package type1

import (
	"fmt"
	"io"

	"golang.org/x/image/font"
	"star-tex.org/x/tex/font/afm"
)

// Font is a Type 1 font.
type Font struct {
	blks []block

	afm afm.Font
}

// Parse parses a (binary) Type 1 font, together with its Adobe File metrics (AFM)
// font file.
func Parse(r, metrics io.Reader) (*Font, error) {
	var (
		f   Font
		err error
	)

	f.afm, err = afm.Parse(metrics)
	if err != nil {
		return nil, fmt.Errorf("could not parse AFM: %w", err)
	}

	f.blks, err = parsePFB(r)
	if err != nil {
		return nil, fmt.Errorf("could not read PFB: %w", err)
	}
	return &f, nil
}

// Metrics returns the metrics of this Type 1 font.
func (f *Font) Metrics() font.Metrics {
	// FIXME(sbinet): rescale?
	return f.afm.Metrics()
}
