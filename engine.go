// Copyright ©2021 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tex

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"modernc.org/knuth/tex"
)

const (
	defaultJobname = "output"
)

// Engine is a TeX engine.
type Engine struct {
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

// New creates a new TeX engine.
func New() *Engine {
	return &Engine{Jobname: defaultJobname}
}

// Run reads the provided TeX document from r and compiles it to
// the provided writer.
func Run(w io.Writer, r io.Reader) error {
	return New().Process(w, r)
}

// Process reads the provided TeX document and
// compiles it to the provided writer.
func (engine *Engine) Process(w io.Writer, r io.Reader) error {
	jobname := engine.Jobname
	if jobname == "" {
		jobname = defaultJobname
	}
	if jobname == "plain" {
		jobname = "plain_"
	}

	var (
		stdin  io.Reader = strings.NewReader(`\input plain \input ` + jobname)
		stdout           = new(bytes.Buffer)
		stderr           = new(bytes.Buffer)
		stdlog           = io.Discard
	)

	if engine.Stdin != nil {
		stdin = io.MultiReader(stdin, engine.Stdin)
	}
	if engine.Stdlog != nil {
		stdlog = engine.Stdlog
	}

	err := tex.Main(
		stdin,
		wtee(stdout, engine.Stdout),
		wtee(stderr, engine.Stderr),
		tex.WithDVIFile(w),
		tex.WithInputFile(jobname+".tex", r),
		tex.WithLogFile(stdlog),
	)
	if err != nil {
		return fmt.Errorf("could not run knuth·main:\nstdout:\n%s\nstderr:\n%s\nerror: %w", stdout, stderr, err)
	}

	return nil
}

func wtee(ws ...io.Writer) io.Writer {
	vs := make([]io.Writer, 0, len(ws))
	for _, w := range ws {
		if w == nil {
			continue
		}
		vs = append(vs, w)
	}
	return io.MultiWriter(vs...)
}
