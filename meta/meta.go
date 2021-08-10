// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package meta

import "github.com/saucelabs/sypl/status"

// IMeta specifies how to get/set information about a component.
type IMeta interface {
	// GetName returns the component name.
	GetName() string

	// SetName sets the component name.
	SetName(name string)

	// GetStatus returns the component status.
	GetStatus() status.Status

	// SetStatus sets the component status.
	SetStatus(s status.Status)
}
