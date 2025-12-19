// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// nolint
package slog

import (
	"log/slog"
)

func attrIsEmpty(a slog.Attr) bool {
	return a.Key == "" && a.Value.Uint64() == 0 && a.Value.Any() == nil
}
