// Copyright ©2026 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"star-tex.org/x/tex/internal/wtex"
)

func main() {
	log.SetPrefix("golatex: ")
	log.SetPrefix("")
	log.SetFlags(0)

	var (
		initex = flag.Bool("ini", false, "run INITEX")
		texfmt = flag.String("fmt", "", "path to preloaded fmt file")
		texmf  = flag.String("texmf", "", "path to TEXMF distribution")
	)
	flag.Parse()

	if *texmf == "" {
	loop:
		for _, env := range []string{"TEXMF", "TEXMFROOT", "TEXMFDIST", "TEXMFHOME"} {
			v := os.Getenv(env)
			if v != "" {
				*texmf = v
				break loop
			}
		}

		if *texmf == "" {
			// try /usr/share/texmf-dist
			*texmf = "/usr/share/texmf-dist"
		}
	}

	if *initex {
		fmt, err := wtex.IniTeX(*texmf)
		if err != nil {
			log.Fatal(err)
		}
		err = os.WriteFile(*texfmt, fmt, 0644)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if flag.NArg() <= 0 {
		flag.Usage()
		log.Fatal("missing path to input (La)TeX file")
	}

	for _, fname := range flag.Args() {
		err := process(fname, *texmf, *texfmt)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func process(fname, texmf, texfmtName string) error {
	doc, err := os.ReadFile(fname)
	if err != nil {
		return fmt.Errorf("could not read %q: %w", fname, err)
	}

	var texfmt []byte
	if texfmtName != "" {
		raw, err := os.ReadFile(texfmtName)
		if err != nil {
			return fmt.Errorf("could not read TeX preloaded format %q: %w", texfmtName, err)
		}
		texfmt = raw
	}

	tex := wtex.New(doc, texmf, texfmt)
	err = tex.Run()
	if err != nil {
		return fmt.Errorf("could not run tex: %w", err)
	}

	return nil
}
