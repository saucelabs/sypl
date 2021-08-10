// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// The goal of Color isn't to expose all possible colors - its up to the user
// to do that, but to provide some common used colors in logging.

package color

import "github.com/fatih/color"

// Color specification.
type Color func(a ...interface{}) string

// Built-in available colors.
var (
	Red        = color.New(color.FgRed).SprintFunc()
	BoldRed    = color.New(color.FgRed, color.Bold).SprintFunc()
	Green      = color.New(color.FgGreen).SprintFunc()
	BoldGreen  = color.New(color.FgGreen, color.Bold).SprintFunc()
	Yellow     = color.New(color.FgYellow).SprintFunc()
	BoldYellow = color.New(color.FgYellow, color.Bold).SprintFunc()
)
