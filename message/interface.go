// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package message

import (
	"time"

	"github.com/saucelabs/sypl/content"
	"github.com/saucelabs/sypl/fields"
	"github.com/saucelabs/sypl/flag"
	"github.com/saucelabs/sypl/level"
)

// ILineBreaker specifies what a LineBreaker does.
type ILineBreaker interface {
	// Restore known linebreaks.
	Restore()

	// Detects (cross-OS) and removes any newline/line-break, at the end of the
	// content, ensuring text processing is done properly (e.g.: suffix).
	Strip()
}

// IMessage specifies what a message does.
type IMessage interface {
	ILineBreaker
	ITag

	// String interface.
	String() string

	// GetComponentName returns the component name.
	GetComponentName() string

	// SetComponentName sets the component name.
	SetComponentName(name string)

	// GetContent returns the content.
	GetContent() content.IContent

	// SetContent sets the content.
	SetContent(c content.IContent)

	// GetFields returns the structured fields.
	GetFields() fields.Fields

	// SetFields sets the structured fields.
	SetFields(fields fields.Fields)

	// GetFlag returns the flag.
	GetFlag() flag.Flag

	// SetFlag sets the flag.
	SetFlag(flag flag.Flag)

	// GetID returns the id.
	GetID() string

	// SetID sets the id.
	SetID(id string)

	// GetLevel returns the level.
	GetLevel() level.Level

	// SetLevel sets the level.
	SetLevel(l level.Level)

	// getLineBreaker returns linebreaker.
	getLineBreaker() *lineBreaker

	// setLineBreaker sets the linebreaker.
	setLineBreaker(lB *lineBreaker)

	// GetMessage (low-level) returns the message.
	GetMessage() *message

	// GetOutputName returns the name of the output in use.
	GetOutputName() string

	// SetOutputName sets the name of the output in use.
	SetOutputName(outputName string)

	// GetOutputsNames returns the outputs names that should be used.
	GetOutputsNames() []string

	// SetOutputsNames sets the outputs names that should be used.
	SetOutputsNames(outputsNames []string)

	// GetProcessorName returns the name of the processor in use.
	GetProcessorName() string

	// SetProcessorName sets the name of the processor in use.
	SetProcessorName(processorName string)

	// GetProcessorsNames returns the processors names that should be used.
	GetProcessorsNames() []string

	// SetProcessorsNames sets the processors names that should be used.
	SetProcessorsNames(processorsNames []string)

	// GetTimestamp returns the timestamp.
	GetTimestamp() time.Time

	// SetTimestamp sets the timestamp.
	SetTimestamp(timestamp time.Time)
}

// ITag specifies what a Tag does.
type ITag interface {
	// AddTags adds one or more tags.
	AddTags(tags ...string)

	// ContainTag verifies if tags contains the specified tag.
	ContainTag(tag string) bool

	// DeleteTag deletes a tag.
	DeleteTag(tag string)

	// GetTags retrieves tags.
	GetTags() []string
}
