// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package message

import (
	"time"

	"github.com/emirpasic/gods/sets/treeset"
	"github.com/saucelabs/sypl/content"
	"github.com/saucelabs/sypl/flag"
	"github.com/saucelabs/sypl/level"
	"github.com/saucelabs/sypl/options"
)

// Message envelops the content and contains meta-information about it.
//
// Note: Changes in the `Message` or `Options` data structure may trigger
// changes in the `Copy`, `mergeOptions`, `NewMessage`, or `NewOptions` methods.
type message struct {
	*options.Options

	// Name of the component logging the message.
	componentName string

	// Content that should be written to `Output`.
	Content content.IContent

	// A randomly generated UUIDv4 that uniquely identifies the message.
	ID string

	// Message's level.
	Level level.Level

	// Output in use.
	OutputName string

	// Processor in use.
	ProcessorName string

	// If set to true, a new line break will be added before printing.
	//
	// Context: When a message enters the pipeline, the last line break is
	// removed - if any (see: `Println`). Content is then processed and only, at
	// the final stage - before printing, the line break is restored, if needed.
	// (see: `write`). This flag controls/indicates if the line break should be
	// restored or not.
	RestoreLineBreak bool

	// tags are indicators consumed by `Output`s and `Processor`s.
	tags *treeset.Set

	// The point in time when the message was created.
	Timestamp time.Time
}

// String interface implementation.
func (m message) String() string {
	return m.Content.GetProcessed()
}

//////
// IMessage interface implementation.
//////

// GetComponentName returns the component name.
func (m *message) GetComponentName() string {
	return m.componentName
}

// SetComponentName sets the component name.
func (m *message) SetComponentName(name string) {
	m.componentName = name
}

// GetRestoreLineBreak returns the line break status.
func (m *message) GetRestoreLineBreak() bool {
	return m.RestoreLineBreak
}

// SetRestoreLineBreak sets the line break status.
func (m *message) SetRestoreLineBreak(s bool) {
	m.RestoreLineBreak = s
}

// GetContent returns the content.
func (m *message) GetContent() content.IContent {
	return m.Content
}

// SetContent sets the content.
func (m *message) SetContent(c content.IContent) {
	m.Content = c
}

// GetFields returns the structured fields.
func (m *message) GetFields() options.Fields {
	return m.Fields
}

// SetFields sets the structured fields.
func (m *message) SetFields(fields options.Fields) {
	m.Fields = fields
}

// GetFlag returns the flag.
func (m *message) GetFlag() flag.Flag {
	return m.Flag
}

// SetFlag sets the flag.
func (m *message) SetFlag(flag flag.Flag) {
	m.Flag = flag
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
func (m *message) SetLevel(l level.Level) {
	m.Level = l
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
func (m *message) SetOutputName(outputName string) {
	m.OutputName = outputName
}

// GetOutputsNames returns the outputs names that should be used.
func (m *message) GetOutputsNames() []string {
	return m.OutputsNames
}

// SetOutputsNames sets the outputs names that should be used.
func (m *message) SetOutputsNames(outputsNames []string) {
	m.OutputsNames = outputsNames
}

// GetProcessorName returns the name of the processor in use.
func (m *message) GetProcessorName() string {
	return m.ProcessorName
}

// SetProcessorName sets the name of the processor in use.
func (m *message) SetProcessorName(processorName string) {
	m.ProcessorName = processorName
}

// GetProcessorsNames returns the processors names that should be used.
func (m *message) GetProcessorsNames() []string {
	return m.ProcessorsNames
}

// SetProcessorsNames sets the processors names that should be used.
func (m *message) SetProcessorsNames(processorsNames []string) {
	m.ProcessorsNames = processorsNames
}

// AddTags adds one or more tags.
func (m *message) AddTags(tags ...string) {
	for _, tag := range tags {
		m.tags.Add(tag)
	}
}

// GetTags retrieves tags.
func (m *message) GetTags() []string {
	tags := []string{}

	m.tags.Each(func(index int, value interface{}) {
		tags = append(tags, value.(string))
	})

	return tags
}

// DeleteTag deletes a tag.
func (m *message) DeleteTag(tag string) {
	m.tags.Remove(tag)
}

// ContainTag verifies if tags contains the specified tag.
func (m *message) ContainTag(tag string) bool {
	return m.tags.Contains(tag)
}

// GetTimestamp returns the timestamp.
func (m *message) GetTimestamp() time.Time {
	return m.Timestamp
}

// SetTimestamp sets the timestamp.
func (m *message) SetTimestamp(timestamp time.Time) {
	m.Timestamp = timestamp
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
func Copy(m IMessage) IMessage {
	msg := NewMessage(m.GetLevel(), m.GetContent().GetOriginal())

	// Copy `options.Tags`.
	msg.GetMessage().Tags = m.GetMessage().Tags

	// Adds tags to `message.tags`.
	msg.AddTags(m.GetTags()...)

	msg.SetComponentName(m.GetComponentName())
	msg.SetFields(m.GetFields())
	msg.SetFlag(m.GetFlag())
	msg.SetID(m.GetID())
	msg.SetOutputName(m.GetOutputName())
	msg.SetOutputsNames(m.GetOutputsNames())
	msg.SetProcessorName(m.GetProcessorName())
	msg.SetProcessorsNames(m.GetProcessorsNames())
	msg.SetRestoreLineBreak(m.GetRestoreLineBreak())
	msg.SetTimestamp(m.GetTimestamp())

	return msg
}

//////
// Factory.
//////

// NewMessage is the Message factory.
//
// Note: Changes in the `Message` or `Options` data structure may reflects here.
func NewMessage(l level.Level, ct string) IMessage {
	return &message{
		Options: options.NewDefaultOptions(),

		Content:          content.NewContent(ct),
		ID:               generateUUID(),
		Level:            l,
		RestoreLineBreak: false,
		tags:             treeset.NewWithStringComparator(),
		Timestamp:        time.Now(),
	}
}
