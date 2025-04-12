// Copyright ©2021 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main // import "star-tex.org/x/tex/cmd/afm-dump"

import (
	"flag"
	"log"
	"os"

	"star-tex.org/x/tex/font/afm"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("afm-dump: ")

	flag.Parse()

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fnt, err := afm.Parse(f)
	if err != nil {
		log.Fatal(err)
	}

	met := fnt.Metrics()
	log.Printf("met: %d", met)
}
