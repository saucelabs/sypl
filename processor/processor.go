// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package processor

import (
	"github.com/saucelabs/sypl/message"
	"github.com/saucelabs/sypl/status"
)

// RunFunc specifies the function used to process a message.
type RunFunc func(m message.IMessage) error

// Processor is a self-contained algorithms that run in isolation.
type processor struct {
	// Function used to process a message.
	f RunFunc

	// Name of the processor.
	name string

	// Status of the processor.
	status status.Status
}

// String interface implementation.
func (p processor) String() string {
	return p.name
}

//////
// IMeta interface implementation.
//////

// GetName returns the processor name.
func (p *processor) GetName() string {
	return p.name
}

// SetName sets the processor name.
func (p *processor) SetName(name string) {
	p.name = name
}

// GetStatus returns the processor status.
func (p *processor) GetStatus() status.Status {
	return p.status
}

// SetStatus sets the processor status.
func (p *processor) SetStatus(s status.Status) {
	p.status = s
}

//////
// IProcessor interface implementation.
//////

// Run the processor, if enabled.
func (p *processor) Run(m message.IMessage) error {
	if p.GetStatus() != status.Enabled {
		return nil
	}

	return p.f(m)
}

//////
// Factory.
//////

// New is the Processor factory.
func New(name string, f RunFunc) IProcessor {
	return &processor{
		f:      f,
		name:   name,
		status: status.Enabled,
	}
}
