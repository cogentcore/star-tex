// Copyright ©2026 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

package main

import (
	"log"
	"os/exec"
)

func main() {
	wasmOpt()
	wasmDis()
	wasm2Go()
}

func wasmOpt() {
	cmd := exec.Command("wasm-opt",
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
	err := cmd.Run()
	if err != nil {
		log.Fatalf("could not run wasm-opt: %v", err)
	}
}

func wasmDis() {
	cmd := exec.Command(
		"wasm-dis", "-o", "wtex-opt.wat", "wtex-opt.wasm",
	)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("could not run wasm-dis: %v", err)
	}
}

func wasm2Go() {
	cmd := exec.Command(
		"wasm2go",
		"-o", "wtex_engine.go",
		"wtex-opt.wasm",
	)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("could not run wasm2go: %v", err)
	}
}
