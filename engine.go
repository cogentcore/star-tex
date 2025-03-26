// Copyright ©2021 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tex

import (
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
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer

	// Jobname used for TeX output.
	// Default is "output".
	Jobname string

	// Stdlog collects TeX logging messages.
	Stdlog io.Writer
}

// NewEngine creates a new TeX engine connected to the provided
// stdin and stdout file descriptors.
func NewEngine(stdout, stderr io.Writer, stdin io.Reader) *Engine {
	return &Engine{
		stdin:   stdin,
		stdout:  stdout,
		stderr:  stderr,
		Jobname: defaultJobname,
		Stdlog:  io.Discard,
	}
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
		stdout io.Writer
		stderr io.Writer
	)

	if engine.stdin != nil {
		stdin = io.MultiReader(stdin, engine.stdin)
	}
	if engine.stdout != nil {
		stdout = engine.stdout
	}
	if engine.stderr != nil {
		stderr = engine.stderr
	}

	err := tex.Main(
		stdin, stdout, stderr,
		tex.WithDVIFile(w),
		tex.WithInputFile(jobname+".tex", r),
		tex.WithLogFile(engine.Stdlog),
	)
	if err != nil {
		return fmt.Errorf("could not run knuth·main: %w", err)
	}

	return nil
}

func writerCloser(w io.Writer) io.WriteCloser {
	if w, ok := w.(io.WriteCloser); ok {
		return w
	}
	return nopWriteCloser{w}
}

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error { return nil }
