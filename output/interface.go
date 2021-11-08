// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package output

import (
	"io"

	"github.com/saucelabs/sypl/formatter"
	"github.com/saucelabs/sypl/internal/builtin"
	"github.com/saucelabs/sypl/level"
	"github.com/saucelabs/sypl/message"
	"github.com/saucelabs/sypl/meta"
	"github.com/saucelabs/sypl/processor"
)

// IProcessor specifies what an output does.
type IOutput interface {
	meta.IMeta

	// String interface.
	String() string

	// GetBuiltinLogger returns the Golang's builtin logger.
	GetBuiltinLogger() *builtin.Builtin

	// SetBuiltinLogger sets the Golang's builtin logger.
	SetBuiltinLogger(builtinLogger *builtin.Builtin) IOutput

	// GetFormatter returns the formatter.
	GetFormatter() formatter.IFormatter

	// SetFormatter sets the formatter.
	SetFormatter(fmtr formatter.IFormatter) IOutput

	// GetMaxLevel returns the max level.
	GetMaxLevel() level.Level

	// SetMaxLevel sets the max level.
	SetMaxLevel(l level.Level) IOutput

	// AddProcessors adds one or more processors.
	AddProcessors(processors ...processor.IProcessor) IOutput

	// GetProcessor returns the registered processor by its name. If not found, will be nil.
	GetProcessor(name string) processor.IProcessor

	// SetProcessors sets one or more processors.
	SetProcessors(processors ...processor.IProcessor) IOutput

	// GetProcessors returns registered processors.
	GetProcessors() []processor.IProcessor

	// GetProcessorsNames returns the names of the registered processors.
	GetProcessorsNames() []string

	// GetWriter returns the writer.
	GetWriter() io.Writer

	// SetWriter sets the writer.
	SetWriter(w io.Writer) IOutput

	// Write write the message to the defined output.
	Write(m message.IMessage) error
}
