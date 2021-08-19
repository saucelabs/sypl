// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package options

import "github.com/saucelabs/sypl/flag"

// Fields allows to add structured fields to a message.
type Fields map[string]interface{}

// Options extends printer's capabilities.
//
// Note: Changes in the `Message` or `Options` data structure may trigger
// changes in the `Copy`, `mergeOptions`, `New`, or `NewOptions` methods.
type Options struct {
	// Structured fields.
	Fields Fields

	// Flags define behaviors.
	Flag flag.Flag

	// OutputsNames are the names of the outputs to be used.
	OutputsNames []string

	// ProcessorsNames are the names of the processors to be used.
	ProcessorsNames []string

	// Tags are indicators consumed by `Output`s and `Processor`s.
	Tags []string
}

//////
// Factory.
//////

// New is the `Options` factory.
//
// Note: Changes in the `Message` or `Options` data structure may reflects here.
func New() *Options {
	return &Options{
		Fields:          Fields{},
		Flag:            flag.None,
		OutputsNames:    []string{},
		ProcessorsNames: []string{},
		Tags:            []string{},
	}
}
