// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// nolint
package slog

import (
	"log/slog"
)

// attrIsEmpty reports whether a is the zero Attr, which ReplaceAttr returns to
// drop an attribute from the output.
//
// Upstream this is Attr.isEmpty, which reads the unexported Value fields
// directly (a.Value.num == 0 && a.Value.any == nil). Those are unreachable
// from here, and Value.Uint64 is not a substitute: it panics unless the value
// is of kind Uint64, and the zero Value has kind Any. Comparing against the
// zero Value through Value.Equal is kind-safe and gives the same answer.
func attrIsEmpty(a slog.Attr) bool {
	return a.Key == "" && a.Value.Equal(slog.Value{})
}
