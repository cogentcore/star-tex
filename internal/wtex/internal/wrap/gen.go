// Copyright ©2026 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("wrap: ")

	web2pas()
	wasmOpt()
	wasmDis()
	wasm2Go()
}

func web2pas() {
	tie := exec.Command(
		"tie", "-m", "tex.web",
		"texlive/texk/tex.web",
		"texlive/etexdir/etex.ch",
		"changes/date.ch",
		"changes/tex-final-end.ch",
		"changes/ord-chr.ch",
		"changes/tokens.ch",
		"changes/inputln.ch",
		"changes/codes.ch",
		"changes/expanded.ch",
		"changes/strcmp.ch",
		"changes/creationdate.ch",
		"changes/filesize.ch",
		"changes/shellescape.ch",
		"changes/wordsize.ch",
		"changes/expand-depth-count.ch",
		"changes/memory.ch",
		"changes/star-latex.ch",
	)
	tie.Stdout = os.Stdout
	tie.Stderr = os.Stderr
	err := tie.Run()
	if err != nil {
		log.Fatalf("could not run tie: %v", err)
	}

	tangle := exec.Command(
		"tangle", "-underline", "tex.web",
	)
	tangle.Stdout = os.Stdout
	tangle.Stderr = os.Stderr
	err = tangle.Run()
	if err != nil {
		log.Fatalf("could not run tangle: %v", err)
	}

	bld := exec.Command(
		"podman", "build", "--network=host", "-t", "web2wasm", "./web2wasm",
	)
	bld.Stdout = os.Stdout
	bld.Stderr = os.Stderr
	err = bld.Run()
	if err != nil {
		log.Fatalf("could not build web2wasm container: %v", err)
	}

	web2wasm := exec.Command(
		"go", "run", "./web2wasm",
		"-pas", "./tex.p",
		"-pool", "./tex.pool",
		"-o", "./tex.wasm",
	)
	web2wasm.Stdout = os.Stdout
	web2wasm.Stderr = os.Stderr
	err = web2wasm.Run()
	if err != nil {
		log.Fatalf("could not run web2wasm: %v", err)
	}

	_ = os.Remove("tex.web")
	_ = os.Remove("tex.p")
}

func wasmOpt() {
	cmd := exec.Command("wasm-opt",
		"-g",
		"--code-folding",
		"--coalesce-locals-learning",
		"--precompute-propagate",
		"--code-pushing",
		"--simplify-locals",
		"--flatten",
		"--rereloop",
		"--dfo",
		"--rereloop",
		"--rereloop",
		"--ssa-nomerge",
		"--local-cse",
		"--asyncify",
		"--pass-arg=asyncify-ignore-indirect",
		"--licm",
		"--flatten",
		"--rereloop",
		"--merge-locals",
		"--merge-blocks",
		"--remove-unused-brs",
		"--remove-unused-names",
		"--dae-optimizing",
		"--inlining-optimizing",
		"--generate-stack-ir",
		"--optimize-stack-ir",
		"--optimize-instructions",
		"--vacuum",
		"-O4",
		"-o", "wtex-opt.wasm",
		"tex.wasm",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("could not run wasm-opt: %v", err)
	}
}

func wasmDis() {
	cmd := exec.Command(
		"wasm-dis", "-o", "wtex-opt.wat", "wtex-opt.wasm",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("could not run wasm-dis: %v", err)
	}
}

func wasm2Go() {
	cmd := exec.Command(
		"go", "tool", "github.com/ncruces/wasm2go",
		"-tags", "!staticcheck",
		"-o", "wtex_engine.go",
		"wtex-opt.wasm",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("could not run wasm2go: %v", err)
	}
}
