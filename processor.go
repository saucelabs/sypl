// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sypl

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Casing definition, e.g.: Upper, Lower, Title, etc.
type Casing string

const (
	// Lowercase casing.
	Lowercase Casing = "lowercase"

	// Uppercase casing.
	Uppercase Casing = "uppercase"
)

// ProcessorFunc is the processor's `do` specification.
type ProcessorFunc func(message *Message)

// Processor processes messages. `Processor`s are self-contained algorithms that
// run in isolation. Any error, should be properly handled, within the processor
// context itself, and not bubbled up. Don't need to handle cases where message
// has no content - it's already done, see `sypl.Process`.
type Processor struct {
	do      ProcessorFunc
	enabled bool
	name    string
}

// GetStatus returns if the processor is enabled or disabled.
func (p *Processor) GetStatus() bool {
	return p.enabled
}

// SetStatus allows to enable or disable a processor.
func (p *Processor) SetStatus(status bool) {
	p.enabled = status
}

// Process the message.
func (p *Processor) Run(message *Message) {
	if p.enabled {
		p.do(message)
	}
}

// NewProcessor creates a new `Processor`.
//
// Notes:
// - The created `Processor` is enabled by default.
// - This method is chainable.
func NewProcessor(name string, processorFunc ProcessorFunc) *Processor {
	return &Processor{
		enabled: true,
		do:      processorFunc,
		name:    name,
	}
}

//////
// Helpers
//////

// generateDefaultPrefix generates prefix for the `PrefixBasedOnMask` processor.
func generateDefaultPrefix(timestampFormat, component string, level Level) string {
	return fmt.Sprintf("%s [%d] [%s] [%s] ",
		// Timestamp.
		time.Now().Format(timestampFormat),

		// PID.
		os.Getpid(),

		// Component name.
		component,

		// Message level.
		level,
	)
}

//////
// Built-in processors.
//////

// Prefixer prefixes messages with the given string.
func Prefixer(prefix string) *Processor {
	return NewProcessor("Prefixer", func(message *Message) {
		message.ContentProcessed = prefix + message.ContentProcessed
	})
}

// Suffixer suffixes messages with the given string.
func Suffixer(suffix string) *Processor {
	return NewProcessor("Suffixer", func(message *Message) {
		message.ContentProcessed += suffix
	})
}

// PrefixBasedOnMask prefixes messages with the predefined mask.
//
// Example: 2021-06-22 12:51:46.089 [80819] [CLI] [INFO].
func PrefixBasedOnMask(timestampFormat string) *Processor {
	return NewProcessor("PrefixBasedOnMask", func(message *Message) {
		message.ContentProcessed = generateDefaultPrefix(
			timestampFormat,
			message.sypl.name,
			message.Level,
		) + message.ContentProcessed
	})
}

// PrefixBasedOnMaskExceptForLevels is a specialized version of the
// `PrefixBasedOnMask`. It prefixes all messages, except for the specified
// levels.
func PrefixBasedOnMaskExceptForLevels(timestampFormat string, levels ...Level) *Processor {
	return NewProcessor("PrefixBasedOnMaskExceptForLevels", func(message *Message) {
		concatenatedLevels := LevelsToString(levels)

		if !strings.Contains(concatenatedLevels, message.Level.String()) {
			message.ContentProcessed = generateDefaultPrefix(
				timestampFormat,
				message.sypl.name,
				message.Level,
			) + message.ContentProcessed
		}
	})
}

// ColorizeBasedOnLevel colorize messages based on the specified levels.
func ColorizeBasedOnLevel(levelColorMap map[Level]Color) *Processor {
	return NewProcessor("ColorizeBasedOnLevel", func(message *Message) {
		for level, color := range levelColorMap {
			if message.Level == level {
				message.ContentProcessed = color(message.ContentProcessed)
			}
		}
	})
}

// ColorizeBasedOnWord colorize a messages with the specified colors if a
// message contains a specified word.
func ColorizeBasedOnWord(wordColorMap map[string]Color) *Processor {
	return NewProcessor("ColorizeBasedOnWord", func(message *Message) {
		for word, color := range wordColorMap {
			if strings.Contains(message.ContentProcessed, word) {
				message.ContentProcessed = color(message.ContentProcessed)
			}
		}
	})
}

// MuteBasedOnLevel mute messages based on the specified levels.
func MuteBasedOnLevel(levels ...Level) *Processor {
	return NewProcessor("MuteBasedOnLevel", func(message *Message) {
		concatenatedLevels := LevelsToString(levels)

		if strings.Contains(concatenatedLevels, message.Level.String()) {
			message.mute = true
		}
	})
}

// ForceBasedOnLevel force messages to be printed based on the specified levels.
func ForceBasedOnLevel(levels ...Level) *Processor {
	return NewProcessor("ForceBasedOnLevel", func(message *Message) {
		concatenatedLevels := LevelsToString(levels)

		if strings.Contains(concatenatedLevels, message.Level.String()) {
			message.force = true
		}
	})
}

// EnableDisableProcessors enables or disables the specified processors.
//
// Note: Order matters! Enabling or disabling a processor that was already
// executed as no effect at all!
func EnableDisableProcessors(status bool, names ...string) *Processor {
	return NewProcessor("EnableDisableProcessors", func(message *Message) {
		concatenatedNames := strings.Join(names, ",")

		for i, processor := range message.output.processors {
			if strings.Contains(concatenatedNames, processor.name) {
				message.output.processors[i].SetStatus(status)
			}
		}
	})
}

// EnableDisableOutputs enables or disables the specified outputs.
//
// Note: Order matters! Enabling or disabling an output that was already
// executed as no effect at all!
func EnableDisableOutputs(status bool, names ...string) *Processor {
	return NewProcessor("EnableDisableOutputs", func(message *Message) {
		concatenatedNames := strings.Join(names, ",")

		for i, output := range message.sypl.outputs {
			if strings.Contains(concatenatedNames, output.name) {
				message.sypl.outputs[i].SetStatus(status)
			}
		}
	})
}

// ChangeFirstCharCase changes message content's first char to the specified
// case.
//
// Notes:
// - `casing` because `case` is a reserved word.
// - Order matters! If this comes after another processor like the Prefixer, it
// will change the case of the first char of the Prefix mask, not the message
// content!
func ChangeFirstCharCase(casing Casing) *Processor {
	return NewProcessor("ChangeFirstCharCase", func(message *Message) {
		firstChar := string(message.ContentProcessed[0])
		contentWithoutFirstChar := message.ContentProcessed[1:len(message.ContentProcessed)]

		switch casing {
		case Uppercase:
			message.ContentProcessed = strings.ToUpper(firstChar) + contentWithoutFirstChar
		case Lowercase:
			message.ContentProcessed = strings.ToLower(firstChar) + contentWithoutFirstChar
		}
	})
}
