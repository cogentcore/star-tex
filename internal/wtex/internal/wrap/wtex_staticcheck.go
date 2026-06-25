// Copyright ©2026 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build staticcheck

package wasm2go

// minimal stubs to ignore the automatically generated WASM→Go code that
// makes staticcheck use too much memory+CPU on sr.ht builds.

type Memory = interface {
	Slice() *[]byte
	Grow(delta, max int64) int64
}

type Xenv = interface {
	Xmemory() Memory
}

type Xlibrary = interface {
	XprintInteger(v0, v1 int32)
	XprintChar(v0, v1 int32)
	XprintString(v0, v1 int32)
	XprintNewline(v0 int32)
	Xreset(v0, v1 int32) int32
	Xgetfilesize(v0, v1 int32) int32
	Xinputln(v0, v1, v2, v3, v4, v5, v6 int32) int32
	Xrewrite(v0, v1 int32) int32
	Xget(v0, v1, v2 int32)
	Xput(v0, v1, v2 int32)
	Xeof(v0 int32) int32
	Xeoln(v0 int32) int32
	Xerstat(v0 int32) int32
	Xclose(v0 int32)
	XgetCurrentMinutes() int32
	XgetCurrentDay() int32
	XgetCurrentMonth() int32
	XgetCurrentYear() int32
	Xtex_final_end()
}

type Module struct {
	Memory Memory

	// Has unexported fields.
}

func New(v0 Xlibrary, v1 Xenv) *Module {
	panic("not implemented")
}

func (m *Module) Xmain() { panic("not implemented") }
