// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sypl

import (
	"log"
	"time"

	"github.com/emirpasic/gods/sets/treeset"
	"github.com/google/uuid"
	"github.com/saucelabs/sypl/content"
	"github.com/saucelabs/sypl/flag"
	"github.com/saucelabs/sypl/level"
)

// generateUUID generates UUIDv4 for message ID.
func generateUUID() string {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Println("[sypl] [Error] generateUUID: Failed to generate UUID for message", err)
	}

	return id.String()
}

// Message envelops the content storing references of the `Logger`, `Output` and
// used `Processor`s. The original content is also stored, and can be used - but
// no changed by `Processor`s.
type Message struct {
	// Content that should be written to `Output`.
	content *content.Content

	// Flags define behaviors.
	flag flag.Flag

	// A randomly generated UUIDv4 that uniquely identifies the message.
	id string

	// Message's level.
	level level.Level

	// A reference to the `Output` in-use.
	output *Output

	// A reference to the `Processor` in-use.
	processor *Processor

	// A reference to the `sypl` in-use.
	sypl *Sypl

	// Tags are indicators consumed by `Output`s and `Processor`s.
	tags *treeset.Set

	// The point in time when the message was created.
	timestamp time.Time
}

// GetProcessedContent returns the content to be processed.
func (m *Message) GetProcessedContent() string {
	return m.content.GetProcessed()
}

// SetProcessedContent sets the processed content.
func (m *Message) SetProcessedContent(content string) {
	m.content.SetProcessed(content)
}

// GetOriginalContent returns the original, non-modified content.
func (m *Message) GetOriginalContent() string {
	return m.content.GetOriginal()
}

// GetFlag returns message's `Flag`.
func (m *Message) GetFlag() flag.Flag {
	return m.flag
}

// SetFlag flags message.
func (m *Message) SetFlag(flag flag.Flag) {
	m.flag = flag
}

// GetID returns the message's `id`.
func (m *Message) GetID() string {
	return m.id
}

// GetLevel returns the message's `Level`.
func (m *Message) GetLevel() level.Level {
	return m.level
}

// GetProcessor returns the message's `Processor`.
func (m *Message) GetProcessor() *Processor {
	return m.processor
}

// SetProcessor sets message's `Processor`.
func (m *Message) SetProcessor(processor *Processor) {
	m.processor = processor
}

// GetOutput returns the message's `Output`.
func (m *Message) GetOutput() *Output {
	return m.output
}

// SetOutput sets message's `Output`.
func (m *Message) SetOutput(output *Output) {
	m.output = output
}

// GetSypl returns the message's `Sypl`.
func (m *Message) GetSypl() *Sypl {
	return m.sypl
}

// SetSypl sets message's `Sypl`.
func (m *Message) SetSypl(sypl *Sypl) {
	m.sypl = sypl
}

// AddTags adds one or more tags.
func (m *Message) AddTags(tags ...string) {
	for _, tag := range tags {
		m.tags.Add(tag)
	}
}

// GetTags retrieves tags.
func (m *Message) GetTags() []string {
	tags := []string{}

	m.tags.Each(func(index int, value interface{}) {
		tags = append(tags, value.(string))
	})

	return tags
}

// DeleteTag deletes a tag.
func (m *Message) DeleteTag(tag string) {
	m.tags.Remove(tag)
}

// ContainTag verifies if `tags` contains the specified `tag`.
func (m *Message) ContainTag(tag string) bool {
	return m.tags.Contains(tag)
}

// GetTimestamp returns the message's `timestamp`.
func (m *Message) GetTimestamp() time.Time {
	return m.timestamp
}

// NewMessage creates a new message.
func NewMessage(sypl *Sypl, output *Output, processor *Processor, level level.Level, ct string) *Message {
	return &Message{
		content:   content.NewContent(ct),
		id:        generateUUID(),
		flag:      flag.None,
		level:     level,
		output:    output,
		processor: processor,
		sypl:      sypl,
		tags:      treeset.NewWithStringComparator(),
		timestamp: time.Now(),
	}
}
