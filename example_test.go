// Copyright ©2025 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tex_test

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	tex "github.com/cogentcore/star-tex"
	"github.com/cogentcore/star-tex/kpath"
)

func ExampleToPDF() {
	const src = `%% A simple TeX document.

Hello, world !
\hrule
Bye.
\bye
`

	ktx := kpath.New()
	pdf := new(bytes.Buffer)

	err := tex.ToPDF(ktx, pdf, strings.NewReader(src))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", pdf.Bytes()[:8])
	// Output:
	// %PDF-1.4
}
