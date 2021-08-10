// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package status

// Status definition.
type Status int

const (
	// Disabled means that a component is inactive, disabled.
	Disabled Status = iota

	// Enabled means that a component is active, enabled.
	Enabled
)

var names = [...]string{"Disabled", "Enabled"}

// String translates enum Status to string.
func (f Status) String() string {
	if f < Disabled || f > Enabled {
		return "Unknown"
	}

	return names[f]
}
