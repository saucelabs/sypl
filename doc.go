// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package sypl provides a Simple Yet Powerful Logger built on top of the Golang
// logger. A sypl logger can have many `Output`s, and each `Output` is
// responsible for writing to a specified destination. Each Output can have
// multiple `Processor`s, which run in isolation manipulating the log message.
// The order of execution is according to the registering order. The above
// features allow sypl to fit into many different logging flows and needs.
//
// In a application with many loggers, and child loggers, sometimes more fine
// control is needed, specially when debugging applications. Sypl offers two
// powerful ways to achieve that: `SYPL_FILTER`, and `SYPL_DEBUG` env vars.
//
// `SYPL_FILTER` allows to specify the name(s) of the component(s) that should
// be logged, for example, for a given application with the following loggers:
// `svc`, `pv`, and `cm`, if a developer wants only to see `svc`, and `pv`
// logging, it's achieved just setting `SYPL_FILTER="svc,pv"`.
//
// `SYPL_DEBUG` allows to specify the max level, for example, for a given
// application with the following loggers: `svc`, `pv`, and `cm`, if a developer
// sets:
//   - `SYPL_DEBUG="debug"`: any application running using Sypl, any component,
//     any output, will log messages bellow the `debug` level
//   - `SYPL_DEBUG="console:debug"`: any application running using Sypl with an
//     output called `console`, will log messages bellow the `debug` level
//   - `SYPL_DEBUG="warn,console:debug"`: any application running using Sypl, any
//     component, any output, will log messages bellow the `warn` level, AND any
//     application running using Sypl with an output called `console`, will log
//     messages bellow the `debug` level. NOTE that `warn` is specified first.
//     Only for this case - global max level scope, it's a requirement! In this
//     case -> `SYPL_DEBUG="console:debug,warn"`, `warn` will be discarded.
//   - `SYPL_DEBUG="svc:console:debug"`: any application running using Sypl with a
//     component called `svc` with an output called `console`, will log messages
//     bellow the `debug` level
//   - `SYPL_DEBUG="file:warn,svc:console:debug"`: any application running using
//     Sypl with an output called `file` will log messages bellow the `warn`
//     level, and any application running using Sypl with a component called `svc`
//     with an output called `console` will log messages bellow the `debug`.
//
// The possibilities are endless! Checkout the [`debugAndFilter`](example_test.go)
// for more.
package sypl
