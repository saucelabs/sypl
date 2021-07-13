// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sypl

import (
	"fmt"
	"os"
	"strings"

	"github.com/saucelabs/sypl/flag"
	"github.com/saucelabs/sypl/level"
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

// GetName returns the processor's name.
func (p *Processor) GetName() string {
	return p.name
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

// ProcessorsNames extract the names of the given processors.
func ProcessorsNames(processors []*Processor) string {
	processorsNames := []string{}

	for _, processor := range processors {
		processorsNames = append(processorsNames, processor.name)
	}

	return strings.Join(processorsNames, ",")
}

//////
// Helpers
//////

// generateDefaultPrefix generates prefix for the `PrefixBasedOnMask` processor.
func generateDefaultPrefix(timestamp, component string, level level.Level) string {
	return fmt.Sprintf("%s [%d] [%s] [%s] ",
		// Timestamp.
		timestamp,

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
		message.SetProcessedContent(prefix + message.GetProcessedContent())
	})
}

// Suffixer suffixes messages with the given string.
func Suffixer(suffix string) *Processor {
	return NewProcessor("Suffixer", func(message *Message) {
		message.SetProcessedContent(message.GetProcessedContent() + suffix)
	})
}

// PrefixBasedOnMask prefixes messages with the predefined mask.
//
// Example: 2021-06-22 12:51:46.089 [80819] [CLI] [Info].
func PrefixBasedOnMask(timestampFormat string) *Processor {
	return NewProcessor("PrefixBasedOnMask", func(message *Message) {
		message.SetProcessedContent(generateDefaultPrefix(
			message.GetTimestamp().Format(timestampFormat),
			message.GetSypl().GetName(),
			message.GetLevel(),
		) + message.GetProcessedContent())
	})
}

// PrefixBasedOnMaskExceptForLevels is a specialized version of the
// `PrefixBasedOnMask`. It prefixes all messages, except for the specified
// levels.
func PrefixBasedOnMaskExceptForLevels(timestampFormat string, levels ...level.Level) *Processor {
	return NewProcessor("PrefixBasedOnMaskExceptForLevels", func(message *Message) {
		concatenatedLevels := level.LevelsToString(levels)

		if !strings.Contains(concatenatedLevels, message.GetLevel().String()) {
			message.SetProcessedContent(generateDefaultPrefix(
				message.GetTimestamp().Format(timestampFormat),
				message.GetSypl().GetName(),
				message.GetLevel(),
			) + message.GetProcessedContent())
		}
	})
}

// ColorizeBasedOnLevel colorize messages based on the specified levels.
func ColorizeBasedOnLevel(levelColorMap map[level.Level]Color) *Processor {
	return NewProcessor("ColorizeBasedOnLevel", func(message *Message) {
		for level, color := range levelColorMap {
			if message.GetLevel() == level {
				message.SetProcessedContent(color(message.GetProcessedContent()))
			}
		}
	})
}

// ColorizeBasedOnWord colorize a messages with the specified colors if a
// message contains a specific word.
func ColorizeBasedOnWord(wordColorMap map[string]Color) *Processor {
	return NewProcessor("ColorizeBasedOnWord", func(message *Message) {
		for word, color := range wordColorMap {
			if strings.Contains(message.GetProcessedContent(), word) {
				message.SetProcessedContent(color(message.GetProcessedContent()))
			}
		}
	})
}

// MuteBasedOnLevel mute messages based on the specified levels.
func MuteBasedOnLevel(levels ...level.Level) *Processor {
	return NewProcessor("MuteBasedOnLevel", func(message *Message) {
		concatenatedLevels := level.LevelsToString(levels)

		if strings.Contains(concatenatedLevels, message.GetLevel().String()) {
			message.SetFlag(flag.Mute)
		}
	})
}

// ForceBasedOnLevel force messages to be printed based on the specified levels.
func ForceBasedOnLevel(levels ...level.Level) *Processor {
	return NewProcessor("ForceBasedOnLevel", func(message *Message) {
		concatenatedLevels := level.LevelsToString(levels)

		if strings.Contains(concatenatedLevels, message.GetLevel().String()) {
			message.SetFlag(flag.Force)
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

		for i, processor := range message.GetOutput().GetProcessors() {
			if strings.Contains(concatenatedNames, processor.GetName()) {
				message.GetOutput().GetProcessor(i).SetStatus(status)
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

		for i, output := range message.GetSypl().GetOutputs() {
			if strings.Contains(concatenatedNames, output.GetName()) {
				message.GetSypl().GetOutput(i).SetStatus(status)
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
		firstChar := string(message.GetProcessedContent()[0])
		contentWithoutFirstChar := message.GetProcessedContent()[1:len(message.GetProcessedContent())]

		switch casing {
		case Uppercase:
			message.SetProcessedContent(strings.ToUpper(firstChar) + contentWithoutFirstChar)
		case Lowercase:
			message.SetProcessedContent(strings.ToLower(firstChar) + contentWithoutFirstChar)
		}
	})
}
