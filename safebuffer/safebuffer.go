// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package safebuffer

import (
	"bytes"
	"sync"
)

// Buffer is a concurrent-safe bytes.Buffer.
type Buffer struct {
	buffer bytes.Buffer
	mutex  sync.Mutex
}

// String returns the contents of the unread portion of the buffer as a string.
// If the Buffer is a nil pointer, it returns "<nil>".
//
// String interface implementation.
func (s *Buffer) String() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.buffer.String()
}

// Write appends the contents of p to the buffer, growing the buffer as needed.
// It returns the number of bytes written.
//
// io.Writer interface implementation.
func (s *Buffer) Write(p []byte) (int, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.buffer.Write(p)
}

// Reset resets the buffer to be empty, but it retains the underlying storage
// for use by future writes. Reset is the same as Truncate(0).
//
// Proxify buffer.Reset.
func (s *Buffer) Reset() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.buffer.Reset()
}
