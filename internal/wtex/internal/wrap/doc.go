// Copyright ©2026 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wasm2go // import "github.com/cogentcore/star-tex/internal/wtex/internal/wrap"

import "embed"

//go:generate go run ./gen.go

//go:embed tex.pool
var FS embed.FS
