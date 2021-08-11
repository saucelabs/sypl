// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package processor

import (
	"github.com/saucelabs/sypl/message"
	"github.com/saucelabs/sypl/meta"
)

// IProcessor specifies what a processor does.
type IProcessor interface {
	meta.IMeta

	// String interface.
	String() string

	// Run the processor, if enabled.
	Run(m message.IMessage) error
}
