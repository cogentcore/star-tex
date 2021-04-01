// Copyright ©2021 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package type1

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	markerASCII  = 0x01
	markerBinary = 0x02
	markerDone   = 0x03
	markerStart  = 0x80

	recordHeaderSize = 6

	t1Header = `%!FontType1-`
	psHeader = `%!PS-AdobeFont-`
)

var (
	errInvalidPFB    = errors.New("invalid PFB file")
	errNoStartMarker = errors.New("no START marker")
	errNoMarker      = errors.New("no marker")
	errInvalidRecord = errors.New("invalid record size")
	errInvalidHeader = errors.New("invalid Type1 header")
)

type block struct {
	kind byte
	data []byte
}

func parsePFB(r io.Reader) (blks []block, err error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("could not read PFB: %w", err)
	}

	var (
		n int
		c int
	)
loop:
	for {
		v := raw[c]
		switch v {
		case markerStart, markerASCII, markerBinary, markerDone:
			// ok.
		default:
			return nil, errInvalidPFB
		}
		kind := raw[c+1]
		if kind == markerDone {
			break loop
		}
		n = int(binary.LittleEndian.Uint32(raw[c+2:]))
		if n > len(raw[c+4:]) {
			return nil, errInvalidRecord
		}
		beg := c + recordHeaderSize
		end := beg + n
		blks = append(blks, block{
			kind: kind,
			data: raw[beg:end],
		})
		c = end
	}

	switch {
	case bytes.HasPrefix(blks[0].data, []byte(psHeader)):
		// ok.
	case bytes.HasPrefix(blks[0].data, []byte(t1Header)):
		// ok.
	default:
		return nil, errInvalidHeader
	}

	return blks, nil
}

func readRecord(p []byte, marker uint8) ([]byte, error) {
	if p[0] != markerStart {
		return nil, errNoStartMarker
	}

	if p[1] != marker {
		return nil, errNoMarker
	}

	n := int(binary.LittleEndian.Uint32(p[2:]))
	if n > len(p)-recordHeaderSize {
		return nil, errInvalidRecord
	}

	beg := recordHeaderSize
	end := beg + n
	return p[beg:end], nil
}
