// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package message

import (
	"strings"
	"time"

	"github.com/emirpasic/gods/sets/treeset"
	"github.com/saucelabs/sypl/content"
	"github.com/saucelabs/sypl/debug"
	"github.com/saucelabs/sypl/fields"
	"github.com/saucelabs/sypl/flag"
	"github.com/saucelabs/sypl/level"
	"github.com/saucelabs/sypl/options"
	"github.com/saucelabs/sypl/status"
)

// LineBreaker defines if `Content` needs linebreak strip/restoration.
//
// Context: When a message enters the pipeline, line breakers are
// removed (see: `output.Write`). Content is then processed and only, at the
// final stage - before printing, the line break is restored, if needed.
// (see: `output.write`).
type lineBreaker struct {
	// ControlChars accumulates stripped control chars.
	ControlChars []string

	// KnownLineBreakers is a list known linebreakers.
	KnownLineBreakers []string

	// Status indicates whether message had control chars stripped, or not.
	Status status.Status
}

// newLineBreaker is the lineBreaker factory.
func newLineBreaker(knownLineBreakers ...string) *lineBreaker {
	return &lineBreaker{
		ControlChars:      []string{},
		KnownLineBreakers: knownLineBreakers,
		Status:            status.Enabled,
	}
}

// Message envelops the content and contains meta-information about it.
//
// Note: Changes in the `Message` or `Options` data structure may trigger
// changes in the `New`, `Copy`, `mergeOptions` (from `Sypl`), or `New` (from
// `Options`) methods.
type message struct {
	*options.Options

	// Name of the component logging the message.
	componentName string

	// Message's linebreaker. See `lineBreaker` for more information.
	lineBreaker *lineBreaker `json:"-"`

	// tags are indicators consumed by `Output`s and `Processor`s.
	tags *treeset.Set

	// Debug capabilities.
	debug *debug.Debug

	// Content that should be written to `Output`.
	Content content.IContent

	// A randomly generated UUIDv4 that uniquely identifies the message.
	ID string

	// Message's level.
	Level level.Level

	// Output in use.
	OutputName string `json:"-"`

	// Processor in use.
	ProcessorName string `json:"-"`

	// The point in time when the message was created.
	Timestamp time.Time
}

// String interface implementation.
func (m message) String() string {
	return m.Content.GetProcessed()
}

//////
// ITag interface implementation.
//////

// AddTags adds one or more tags.
func (m *message) AddTags(tags ...string) {
	for _, tag := range tags {
		m.tags.Add(tag)
	}
}

// ContainTag verifies if tags contains the specified tag.
func (m *message) ContainTag(tag string) bool {
	return m.tags.Contains(tag)
}

// DeleteTag deletes a tag.
func (m *message) DeleteTag(tag string) {
	m.tags.Remove(tag)
}

// GetTags retrieves tags.
func (m *message) GetTags() []string {
	tags := []string{}

	m.tags.Each(func(index int, value interface{}) {
		tags = append(tags, value.(string))
	})

	return tags
}

//////
// ILineBreaker interface implementation.
//////

// getLineBreaker returns linebreaker.
func (m *message) getLineBreaker() *lineBreaker {
	return m.lineBreaker
}

// setLineBreaker sets the line break status.
func (m *message) setLineBreaker(lB *lineBreaker) IMessage {
	m.lineBreaker = lB

	return m
}

// Restore known linebreaks.
func (m *message) Restore() {
	if m.getLineBreaker().Status == status.Enabled {
		for _, controlChar := range m.getLineBreaker().ControlChars {
			m.GetContent().SetProcessed(m.GetContent().GetProcessed() + controlChar)
		}
	}
}

// Detects (cross-OS) and removes any newline/line-break, at the end of the
// content, ensuring text processing is done properly (e.g.: suffix).
func (m *message) Strip() {
	if m.getLineBreaker().Status == status.Enabled {
		for _, knownLineBreaker := range m.getLineBreaker().KnownLineBreakers {
			if strings.HasSuffix(m.GetContent().GetProcessed(), knownLineBreaker) {
				m.GetContent().SetProcessed(
					strings.TrimSuffix(m.GetContent().GetProcessed(), knownLineBreaker),
				)

				m.getLineBreaker().ControlChars = append(m.getLineBreaker().ControlChars, knownLineBreaker)

				m.Strip()
			}
		}
	}
}

//////
// IMessage interface implementation.
//////

// GetComponentName returns the component name.
func (m *message) GetComponentName() string {
	return m.componentName
}

// SetComponentName sets the component name.
func (m *message) SetComponentName(name string) IMessage {
	m.componentName = name

	return m
}

// GetContent returns the content.
func (m *message) GetContent() content.IContent {
	return m.Content
}

// SetContent sets the content.
func (m *message) SetContent(c content.IContent) IMessage {
	m.Content = c

	return m
}

// GetDebugEnvVarRegexeses returns the Debug env var regexes matchers.
func (m *message) GetDebugEnvVarRegexes() *debug.Debug {
	return m.debug
}

// SetDebugEnvVarRegexeses sets the Debug env var regexes matchers.
func (m *message) SetDebugEnvVarRegexes(d *debug.Debug) *message {
	m.debug = d

	return m
}

// GetFields returns the structured fields.
func (m *message) GetFields() fields.Fields {
	return m.Fields
}

// SetFields sets the structured fields.
func (m *message) SetFields(fields fields.Fields) IMessage {
	m.Fields = fields

	return m
}

// GetFlag returns the flag.
func (m *message) GetFlag() flag.Flag {
	return m.Flag
}

// SetFlag sets the flag.
func (m *message) SetFlag(flag flag.Flag) IMessage {
	m.Flag = flag

	return m
}

// GetID returns the id.
func (m *message) GetID() string {
	return m.ID
}

// SetID sets the id.
func (m *message) SetID(id string) {
	m.ID = id
}

// GetLevel returns the level.
func (m *message) GetLevel() level.Level {
	return m.Level
}

// SetLevel sets the level.
func (m *message) SetLevel(l level.Level) IMessage {
	m.Level = l

	return m
}

// GetMessage (low-level) returns the message.
func (m *message) GetMessage() *message {
	return m
}

// GetOutputName returns the name of the output in use.
func (m *message) GetOutputName() string {
	return m.OutputName
}

// SetOutputName sets the name of the output in use.
func (m *message) SetOutputName(outputName string) IMessage {
	m.OutputName = outputName

	return m
}

// GetOutputsNames returns the outputs names that should be used.
func (m *message) GetOutputsNames() []string {
	return m.OutputsNames
}

// SetOutputsNames sets the outputs names that should be used.
func (m *message) SetOutputsNames(outputsNames []string) IMessage {
	m.OutputsNames = outputsNames

	return m
}

// GetProcessorName returns the name of the processor in use.
func (m *message) GetProcessorName() string {
	return m.ProcessorName
}

// SetProcessorName sets the name of the processor in use.
func (m *message) SetProcessorName(processorName string) IMessage {
	m.ProcessorName = processorName

	return m
}

// GetProcessorsNames returns the processors names that should be used.
func (m *message) GetProcessorsNames() []string {
	return m.ProcessorsNames
}

// SetProcessorsNames sets the processors names that should be used.
func (m *message) SetProcessorsNames(processorsNames []string) IMessage {
	m.ProcessorsNames = processorsNames

	return m
}

// GetTimestamp returns the timestamp.
func (m *message) GetTimestamp() time.Time {
	return m.Timestamp
}

// SetTimestamp sets the timestamp.
func (m *message) SetTimestamp(timestamp time.Time) IMessage {
	m.Timestamp = timestamp

	return m
}

//////
// Helpers.
//////

// Copy message.
//
// Notes:
// - Changes in the `Message` or `Options` data structure may reflects here.
// - Could use something like the `Copier` package, but that's going to cause a
// data race, because `Output`s are processed concurrently.
//
// TODO: This can be improved.
func Copy(m IMessage) IMessage {
	msg := New(m.GetLevel(), m.GetContent().GetOriginal())

	// Copy `options.Tags`.
	msg.GetMessage().Tags = m.GetMessage().Tags

	// Adds tags to `message.tags`.
	msg.AddTags(m.GetTags()...)

	msg.SetComponentName(m.GetComponentName())
	msg.SetDebugEnvVarRegexes(m.GetDebugEnvVarRegexes())

	msg.SetFields(m.GetFields())
	msg.SetFlag(m.GetFlag())
	msg.SetID(m.GetID())

	gLB := *m.getLineBreaker()
	msg.setLineBreaker(&gLB)

	msg.SetOutputName(m.GetOutputName())
	msg.SetOutputsNames(m.GetOutputsNames())
	msg.SetProcessorName(m.GetProcessorName())
	msg.SetProcessorsNames(m.GetProcessorsNames())
	msg.SetTimestamp(m.GetTimestamp())

	return msg
}

//////
// Factory.
//////

// New is the Message factory.
//
// Note: Changes in the `Message` or `Options` data structure may reflects here.
func New(l level.Level, ct string) IMessage {
	return &message{
		Options: options.New(),

		Content:     content.New(ct),
		ID:          generateUUID(),
		Level:       l,
		lineBreaker: newLineBreaker("\n", "\r"),
		tags:        treeset.NewWithStringComparator(),
		Timestamp:   time.Now(),
	}
}
