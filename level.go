// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sypl

import (
	"log"
	"strings"
)

// Level specification.
type Level int

// Available levels.
const (
	NONE Level = iota
	FATAL
	ERROR
	INFO
	WARN
	DEBUG
	TRACE
)

var levelsString = []string{"", "FATAL", "ERROR", "INFO", "WARN", "DEBUG", "TRACE"}

// String translates enum levels to string. Invalid translation will result in
// Fatal error. The following would causa a fatal error:
//  fmt.Println(Level(15))
// `15` is out of the range of `levelsString`.
func (l Level) String() string {
	// Handles non-valid levels.
	MustBeValidFromInt(int(l))

	return levelsString[l]
}

// LevelsToString converts a slice of levels to string (concatenated).
func LevelsToString(levels []Level) string {
	levelsString := []string{}

	for _, level := range levels {
		levelsString = append(levelsString, level.String())
	}

	return strings.Join(levelsString, ",")
}

// IsValidLevelFromInt validates if a given level - integer format, is valid.
func IsValidLevelFromInt(level int) bool {
	return !(level < 0 || level > len(levelsString)-1)
}

// MustBeValidFromInt validates the given `level`. It'll exit if it's invalid.
// If `level` is valid, it returns the respective `Level`.
func MustBeValidFromInt(level int) Level {
	if !IsValidLevelFromInt(level) {
		log.Fatalf("Invalid specified log level %d", level)
	}

	return Level(level)
}

// LevelFromInt return a `Level` from a given integer.
//
// Notes:
// - It'll be automatically validated.
// - It'll exit if it's invalid.
func LevelFromInt(level int) Level {
	return MustBeValidFromInt(level)
}

// LevelFromInt return a `Level` from a given integer.
//
// Note: It'll exit if it's invalid.
func LevelFromString(level string) Level {
	for i, levelString := range levelsString {
		if strings.EqualFold(level, levelString) {
			return Level(i)
		}
	}

	log.Fatalf("Invalid specified log level %s", level)

	return NONE
}
