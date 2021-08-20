// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package flag

// Flag defines behaviours.
type Flag int

const (
	None Flag = iota

	// Force message will be processed, and printed independent of `Level`
	// restrictions.
	Force

	// Mute message will be processed, but not printed.
	Mute

	// Skip message will not be processed, neither formatted, but printed.
	Skip

	// SkipAndForce message will not be processed, , neither formatted, but will
	// be printed independent of `Level` restrictions.
	SkipAndForce

	// SkipAndMute message will not be processed, neither printed.
	SkipAndMute
)

var names = [...]string{"None", "Force", "Mute", "Skip", "SkipAndForce", "SkipAndMute"}

// String interface implementation.
func (f Flag) String() string {
	if f < None || f > SkipAndMute {
		return "Unknown"
	}

	return names[f]
}
