// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// nolint
package slog

import "sync"

// buffer is a byte buffer.
//
// This implementation is adapted from the unexported type buffer
// in go/src/fmt/print.go.
type buffer []byte

// Having an initial size gives a dramatic speedup.
var bufPool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, 1024)
		return (*buffer)(&b)
	},
}

func newBuffer() *buffer {
	return bufPool.Get().(*buffer)
}

// Free releases the buffer back to the pool.
func (b *buffer) Free() {
	// To reduce peak allocation, return only smaller buffers to the pool.
	const maxBufferSize = 16 << 10
	if cap(*b) <= maxBufferSize {
		*b = (*b)[:0]
		bufPool.Put(b)
	}
}

// Reset clears the buffer contents by setting its length to zero.
func (b *buffer) Reset() {
	b.SetLen(0)
}

// Write appends the provided bytes to the buffer.
func (b *buffer) Write(p []byte) (int, error) {
	*b = append(*b, p...)
	return len(p), nil
}

// WriteString appends the provided string to the buffer.
func (b *buffer) WriteString(s string) (int, error) {
	*b = append(*b, s...)
	return len(s), nil
}

// WriteByte appends a single byte to the buffer.
func (b *buffer) WriteByte(c byte) error {
	*b = append(*b, c)
	return nil
}

// String returns the contents of the buffer as a string.
func (b *buffer) String() string {
	return string(*b)
}

// Len returns the number of bytes currently stored in the buffer.
func (b *buffer) Len() int {
	return len(*b)
}

// SetLen reslices the buffer to the specified length.
func (b *buffer) SetLen(n int) {
	*b = (*b)[:n]
}
