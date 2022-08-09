// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package processor

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/saucelabs/sypl/color"
	"github.com/saucelabs/sypl/flag"
	"github.com/saucelabs/sypl/level"
	"github.com/saucelabs/sypl/message"
)

// Casing definition, e.g.: Upper, Lower, Title, etc.
type Casing string

const (
	// Lowercase casing.
	Lowercase Casing = "lowercase"

	// Uppercase casing.
	Uppercase Casing = "uppercase"
)

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

// ChangeFirstCharCase changes message content's first char to the specified
// case.
//
// Notes:
// - `casing` because `case` is a reserved word.
// - Order matters! If this comes after another processor like the Prefixer, it
// will change the case of the first char of the Prefix mask, not the message
// content!
func ChangeFirstCharCase(casing Casing) IProcessor {
	return New("ChangeFirstCharCase", func(m message.IMessage) error {
		firstChar := string(m.GetContent().GetProcessed()[0])
		contentWithoutFirstChar := m.GetContent().GetProcessed()[1:len(m.GetContent().GetProcessed())]

		switch casing {
		case Uppercase:
			m.GetContent().SetProcessed(strings.ToUpper(firstChar) + contentWithoutFirstChar)
		case Lowercase:
			m.GetContent().SetProcessed(strings.ToLower(firstChar) + contentWithoutFirstChar)
		}

		return nil
	})
}

// ColorizeBasedOnLevel colorize messages based on the specified levels.
func ColorizeBasedOnLevel(levelColorMap map[level.Level]color.Color) IProcessor {
	return New("ColorizeBasedOnLevel", func(m message.IMessage) error {
		for level, color := range levelColorMap {
			if m.GetLevel() == level {
				m.GetContent().SetProcessed(color(m.GetContent().GetProcessed()))
			}
		}

		return nil
	})
}

// ColorizeBasedOnWord colorize a messages with the specified colors if a
// message contains a specific word.
func ColorizeBasedOnWord(wordColorMap map[string]color.Color) IProcessor {
	return New("ColorizeBasedOnWord", func(m message.IMessage) error {
		for word, color := range wordColorMap {
			if strings.Contains(m.GetContent().GetProcessed(), word) {
				m.GetContent().SetProcessed(color(m.GetContent().GetProcessed()))
			}
		}

		return nil
	})
}

// Decolourizer removes any colour.
func Decolourizer() IProcessor {
	return New("Decolourizer", func(m message.IMessage) error {
		m.GetContent().SetProcessed(stripansi.Strip(m.GetContent().GetProcessed()))

		return nil
	})
}

// ErrorSimulator simulates an error in the pipeline.
//
//nolint:goerr113
func ErrorSimulator(msg string) IProcessor {
	return New("ErrorSimulator", func(m message.IMessage) error {
		return errors.New(msg)
	})
}

// ForceBasedOnLevel force messages to be printed based on the specified levels.
func ForceBasedOnLevel(levels ...level.Level) IProcessor {
	return New("ForceBasedOnLevel", func(m message.IMessage) error {
		concatenatedLevels := level.LevelsToString(levels)

		if strings.Contains(concatenatedLevels, m.GetLevel().String()) {
			m.SetFlag(flag.Force)
		}

		return nil
	})
}

// MuteBasedOnLevel mute messages based on the specified levels.
func MuteBasedOnLevel(levels ...level.Level) IProcessor {
	return New("MuteBasedOnLevel", func(m message.IMessage) error {
		concatenatedLevels := level.LevelsToString(levels)

		if strings.Contains(concatenatedLevels, m.GetLevel().String()) {
			m.SetFlag(flag.Mute)
		}

		return nil
	})
}

// PrefixBasedOnMask prefixes messages with the predefined mask.
//
// Example: 2021-06-22 12:51:46.089 [80819] [CLI] [Info].
func PrefixBasedOnMask(timestampFormat string) IProcessor {
	return New("PrefixBasedOnMask", func(m message.IMessage) error {
		m.GetContent().SetProcessed(generateDefaultPrefix(
			m.GetTimestamp().Format(timestampFormat),
			m.GetComponentName(),
			m.GetLevel(),
		) + m.GetContent().GetProcessed())

		return nil
	})
}

// PrefixBasedOnMaskExceptForLevels is a specialized version of the
// `PrefixBasedOnMask`. It prefixes all messages, except for the specified
// levels.
func PrefixBasedOnMaskExceptForLevels(timestampFormat string, levels ...level.Level) IProcessor {
	return New("PrefixBasedOnMaskExceptForLevels", func(m message.IMessage) error {
		concatenatedLevels := level.LevelsToString(levels)

		if !strings.Contains(concatenatedLevels, m.GetLevel().String()) {
			m.GetContent().SetProcessed(generateDefaultPrefix(
				m.GetTimestamp().Format(timestampFormat),
				m.GetComponentName(),
				m.GetLevel(),
			) + m.GetContent().GetProcessed())
		}

		return nil
	})
}

// Prefixer prefixes a message with the specified `prefix`.
func Prefixer(prefix string) IProcessor {
	return New("Prefixer", func(m message.IMessage) error {
		m.GetContent().SetProcessed(prefix + m.GetContent().GetProcessed())

		return nil
	})
}

// PrintOnlyAtLevel prints only if message is at the specified level.
func PrintOnlyAtLevel(levels ...level.Level) IProcessor {
	return New("PrintOnlyAtLevel", func(m message.IMessage) error {
		concatenatedLevels := level.LevelsToString(levels)

		if !strings.Contains(concatenatedLevels, m.GetLevel().String()) {
			m.SetFlag(flag.Mute)
		}

		return nil
	})
}

// PrintOnlyIfTagged prints only if message contains the specified tag.
func PrintOnlyIfTagged(tag string) IProcessor {
	return New("PrintOnlyIfTagged", func(m message.IMessage) error {
		if !m.ContainTag(tag) {
			m.SetFlag(flag.Mute)
		}

		return nil
	})
}

// Suffixer suffixes a message with the specified `suffix`.
func Suffixer(suffix string) IProcessor {
	return New("Suffixer", func(m message.IMessage) error {
		m.GetContent().SetProcessed(m.GetContent().GetProcessed() + suffix)

		return nil
	})
}
