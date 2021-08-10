// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package level

import (
	"strings"
)

// Level specification.
type Level int

// Available levels.
const (
	None Level = iota
	Fatal
	Error
	Info
	Warn
	Debug
	Trace
)

var names = [...]string{"None", "Fatal", "Error", "Info", "Warn", "Debug", "Trace"}

// String interface implementation.
func (l Level) String() string {
	if l < None || l > Trace {
		return "Unknown"
	}

	return names[l]
}

// FromInt returns a `Level` from a given integer.
//
// Note: Failure will return "Unknown".
func FromInt(level int) Level {
	return Level(level)
}

// FromInt returns a `Level` from a given integer.
//
// Note: Failure will return `None` level.
func FromString(level string) Level {
	for i, levelString := range names {
		if strings.EqualFold(level, levelString) {
			return Level(i)
		}
	}

	return None
}

// LevelsToString converts a slice of levels to string (concatenated).
func LevelsToString(levels []Level) string {
	names := []string{}

	for _, level := range levels {
		names = append(names, level.String())
	}

	return strings.Join(names, ",")
}
