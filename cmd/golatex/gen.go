// Copyright ©2026 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

package main

import (
	"archive/zip"
	"flag"
	"log"
	"os"

	"star-tex.org/x/tex/internal/wtex"
)

func main() {
	var (
		texmf = flag.String("texmf", "", "path to TEXMF distribution")
	)
	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatal("missing path to output format name")
	}

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

	oname := flag.Arg(0)

	bin, err := wtex.IniTeX(*texmf)
	if err != nil {
		log.Fatal(err)
	}

	out, err := os.Create(oname + ".zip")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	zout := zip.NewWriter(out)

	w, err := zout.Create(oname)
	if err != nil {
		log.Fatal(err)
	}

	_, err = w.Write(bin)
	if err != nil {
		log.Fatalf("could not write to %q: %v", oname, err)
	}
	err = zout.Close()
	if err != nil {
		log.Fatalf("could not save %q: %v", oname, err)
	}

	err = out.Close()
	if err != nil {
		log.Fatalf("could not save %q: %v", oname, err)
	}

	_ = os.Remove(oname)

	return
}
