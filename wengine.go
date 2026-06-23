// Copyright ©2021 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tex

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"fmt"
	"image/color"
	"io"
	"os"

	"star-tex.org/x/tex/dvi"
	"star-tex.org/x/tex/dvi/dvipdf"
	"star-tex.org/x/tex/internal/wtex"
	"star-tex.org/x/tex/kpath"
)

//go:embed golatex.fmt.zip
var latexFmt []byte

// cached version of format file
var cachedFmt []byte

// LaTeXToDVI reads the provided LaTeX document from r and compiles it to
// the provided writer as a ToDVI document.
func LaTeXToDVI(w io.Writer, r io.Reader) error {
	return NewLaTeX().Process(w, r)
}

// LaTeXToPDF reads the provided LaTeX document from r and compiles it to
// the provided writer as a PDF document.
func LaTeXToPDF(ctx kpath.Context, w io.Writer, r io.Reader) error {
	buf := new(bytes.Buffer)
	err := LaTeXToDVI(buf, r)
	if err != nil {
		return fmt.Errorf("could not compile LaTeX document to DVI: %w", err)
	}

	log := new(bytes.Buffer)
	pdf := dvipdf.New(
		ctx, w,
		dvipdf.WithBackground(color.White),
		dvipdf.WithEmbedFonts(true),
	)
	vm := dvi.NewMachine(
		dvi.WithContext(ctx),
		dvi.WithLogOutput(log),
		dvi.WithRenderer(pdf),
		dvi.WithHandlers(dvi.NewColorHandler(ctx)),
		dvi.WithOffsetX(0),
		dvi.WithOffsetY(0),
	)

	prog, err := dvi.Compile(buf.Bytes())
	if err != nil {
		return fmt.Errorf("could not compile DVI program: %w", err)
	}

	err = vm.Run(prog)
	if err != nil {
		return fmt.Errorf("could not run DVI program: %w", err)
	}

	err = pdf.Close()
	if err != nil {
		return fmt.Errorf("could not render DVI program to PDF: %w", err)
	}

	return nil
}

// LaTeXEngine is a LaTeX engine.
type LaTeXEngine struct {
	// Stdin specifies the TeX engine's standard input.
	//
	// If Stdin is nil, the engine reads from the null device (os.DevNull).
	Stdin io.Reader

	// Stdout and Stderr specify the TeX engine's standard output and error.
	//
	// If either is nil, Process connects the corresponding file descriptor
	// to the null device (os.DevNull).
	Stdout io.Writer
	Stderr io.Writer

	// Jobname used for TeX output.
	// Default is "output".
	Jobname string

	// Stdlog collects TeX logging messages.
	//
	// If Stdlog is nil, Process connects Stdlog to the null device (os.DevNull).
	Stdlog io.Writer
}

// NewLaTeX creates a new LaTeX engine.
func NewLaTeX() *LaTeXEngine {
	return &LaTeXEngine{Jobname: defaultJobname}
}

// Process reads the provided LaTeX document and
// compiles it to the provided writer.
func (engine *LaTeXEngine) Process(w io.Writer, r io.Reader) error {
	jobname := engine.Jobname
	if jobname == "" {
		jobname = defaultJobname
	}
	if jobname == "plain" {
		jobname = "plain_"
	}

	var (
	// stdin io.Reader = strings.NewReader(`\input plain \input ` + jobname)
	// stdout           = new(bytes.Buffer)
	// stderr           = new(bytes.Buffer)
	// stdlog = io.Discard
	)

	// if engine.Stdin != nil {
	// 	stdin = io.MultiReader(stdin, engine.Stdin)
	// }
	// if engine.Stdlog != nil {
	// 	stdlog = engine.Stdlog
	// }

	texfmt := cachedFmt
	if texfmt == nil {
		zr, err := zip.NewReader(bytes.NewReader(latexFmt), int64(len(latexFmt)))
		if err != nil {
			return fmt.Errorf("could not read compressed golatex.fmt: %w", err)
		}
		f, err := zr.Open("golatex.fmt")
		if err != nil {
			return fmt.Errorf("could not open embedded golatex.fmt file: %w", err)
		}
		defer f.Close()

		texfmt, err = io.ReadAll(f)
		if err != nil {
			return fmt.Errorf("could not read embedded golatex.fmt file: %w", err)
		}
		cachedFmt = texfmt
	}

	texmf := ""
	for _, env := range []string{"TEXMF", "TEXMFROOT", "TEXMFDIST", "TEXMFHOME"} {
		v := os.Getenv(env)
		if v != "" {
			texmf = v
			break
		}
	}
	if texmf == "" {
		// try /usr/local/share/texmf-dist
		texmf = "/usr/local/share/texmf-dist"
	}

	doc, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	wt := wtex.New(doc, texmf, texfmt)
	err = wt.Run()
	if err != nil {
		return fmt.Errorf("could not run tex: %w", err)
	}

	// output is now at texput.dvi -- read that into stdout
	out, err := os.Open("texput.dvi")
	if err != nil {
		return fmt.Errorf("texput.dvi output file could not be opened: %w", err)
	}
	defer out.Close()
	io.Copy(w, out)

	return nil
}
