// Copyright ©2026 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("web2wasm: ")

	var (
		fname = flag.String("pas", "tex.p", "path to PASCAL file")
		pname = flag.String("pool", "tex.pool", "path to POOL file")
		oname = flag.String("o", "tex.wasm", "path to output WASM file")
	)

	flag.Parse()

	if *fname == "" {
		flag.Usage()
		log.Fatal("missing path to input PASCAL file")
	}

	err := compile(*oname, *fname, *pname)
	if err != nil {
		log.Fatal(err)
	}
}

func compile(oname, fname, pname string) error {
	tmp, err := os.MkdirTemp("", "web2wasm-")
	if err != nil {
		return fmt.Errorf("could not create tmp dir: %w", err)
	}
	defer os.RemoveAll(tmp)

	pas, err := os.ReadFile(fname)
	if err != nil {
		return fmt.Errorf("could not open PASCAL file: %w", err)
	}

	err = os.WriteFile(filepath.Join(tmp, "tex.p"), pas, 0644)
	if err != nil {
		return fmt.Errorf("could not install PASCAL file into podman volume: %w", err)
	}

	pool, err := os.ReadFile(pname)
	if err != nil {
		return fmt.Errorf("could not open POOL file: %w", err)
	}

	err = os.WriteFile(filepath.Join(tmp, "tex.pool"), pool, 0644)
	if err != nil {
		return fmt.Errorf("could not install PASCAL file into podman volume: %w", err)
	}

	cmd := exec.Command(
		"podman", "run", "--rm",
		"-v", tmp+":/opt/web2wasm/share",
		"web2wasm", "/bin/bash", "/opt/web2wasm/run",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("could not run web2wasm: %w", err)
	}

	wasm, err := os.ReadFile(filepath.Join(tmp, "tex.wasm"))
	if err != nil {
		return fmt.Errorf("could not read WASM file output: %w", err)
	}

	err = os.WriteFile(oname, wasm, 0644)
	if err != nil {
		return fmt.Errorf("could not write WASM file output: %w", err)
	}

	return nil
}
