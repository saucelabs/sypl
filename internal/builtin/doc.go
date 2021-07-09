// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package builtin contains the Golang's built-in logger with only one change:
// Don't always add a newline. If someday they fix that, remove this internal
// package.
//
// See this https://github.com/golang/go/issues/16564 for more info.
package builtin
