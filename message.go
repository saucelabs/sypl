// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sypl

// Message envelops the content storing references of the `Logger`, `Output` and
// used `Processor`s. The original content is also stored, and can be used - but
// no changed by `Processor`s.
type Message struct {
	ContentOriginal  string
	ContentProcessed string
	Level            Level

	sypl *Sypl

	// If set, the message will be printed regardless of whether it is muted,
	// or the "Maximum Level". Note: `Force` has precedence over `Mute`.
	force bool

	// If set, message will not be written to output. Note: `Force` has
	// precedence over `Mute`.
	mute bool

	output    *Output
	processor *Processor
}
