// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package sypl provides a Simple Yet Powerful Logger built on top of the Golang
// logger. A sypl logger can have many `Output`s, and each `Output` is
// responsible for writing to a specified destination. Each Output can have
// multiple `Processor`s, which run in isolation manipulating the log message.
// The order of execution is according to the registering order. The above
// features allow sypl to fit into many different logging flows and needs.
package sypl
