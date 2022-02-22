# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Roadmap

- Fix possible race condition on `sypl.SetMaxLevel`
- Add badges to README.md:
  - GoDoc
  - Go Report Card
  - Build Status
  - Coverage Status
  - License
  - Release
- Flag should be an array.
- Move `Options`' methods from `Message` to `Options`
- Improve documentation:
  - Add `doc.go` for all packages

Refs. for badges:

- http://github.com/wayneashleyberry/terminal-dimensions
- https://github.com/golangci/golangci-lint

## [1.5.9] - 2022-02-21
### Added
- For convenience, conforms with `io.Writer` interface. Default level: `error`. `SetIoWriterLevel` changes the default level.

### Changed
- Fixed chained example.
- Lowercased all levels.
- It now warns when an application tries to write to a closed writer.

## [1.5.8] - 2021-11-08
### Changed
- All `SetXYZ` methods returns its proper interface allowing method chaining.
- `Breakpoint` is now variadic.
- Properly handle cases where sypl writes to a piped output, but it's broken.

## [1.5.7] - 2021-11-02
### Changed
- Fixed `ExampleNew_globalFields` test.

## [1.5.6] - 2021-11-02
### Added
Added the ability to set breakpoints. If a `Breakpoint` is set it'll stop execution waiting the user press `/n` (**"enter"**) to continue. It helps users doing quick, and effective log-to-console debug. A message with the breakpoint `name`, and `PID` of the process will be printed using the `debug` level. Arbitrary `data` can optionally be set - if set, it'll be printed. Errors are printed using the standard `error` level. Set logging level to `trace` for more.

Previously, flow would look like:
- Log markers are set, e.g.: `logger.Debugln("Here 1", whatever)`
- Application runs
- Scan visually `output`, or a `file` - via `grep` for the markers.

Now:
- Named `Breakpoint`s are set
- Application runs
- Breakpoint is hit. Information about it is printed.
- Runtime is paused, allowing analysis of `data` - if any, right way. Additionally, an external and more advanced debugger can be attached.
- Dev controls the flow, pressing `enter` at any time, continue.

## [1.5.5] - 2021-10-29
### Changed
- Exported `sypl.Name` to deal with https://github.com/golang/go/issues/5819.

## [1.5.4] - 2021-10-13
### Added
In a application with many loggers, and child loggers, sometimes more fine control is needed, specially when debugging applications. Sypl offers two powerful ways to achieve that: `SYPL_FILTER`, and `SYPL_DEBUG` env vars.

`SYPL_FILTER` allows to specify the name(s) of the component(s) that should be logged, for example, for a given application with the following loggers: `svc`, `pv`, and `cm`, if a developer wants only to see `svc`, and `pv` logging, it's achieved just setting `SYPL_FILTER="svc,pv"`.

`SYPL_DEBUG` allows to specify the max level, for example, for a given application with the following loggers: `svc`, `pv`, and `cm`, if a developer sets:

- `SYPL_DEBUG="debug"`: any application running using Sypl, any component, any output, will log messages bellow the `debug` level
- `SYPL_DEBUG="console:debug"`: any application running using Sypl with an output called `console`, will log messages bellow the `debug` level
- `SYPL_DEBUG="warn,console:debug"`: any application running using Sypl, any component, any output, will log messages bellow the `warn` level, AND any application running using Sypl with an output called `console`, will log messages bellow the `debug` level.

_NOTE: `warn` is specified first. Only for this case - **global scope**, it's a requirement.
`SYPL_DEBUG="console:debug,warn"`: In this case `warn` will be **discarded!**._

- `SYPL_DEBUG="svc:console:debug"`: any application running using Sypl with a component called `svc` with an output called `console`, will log messages bellow the `debug` level
- `SYPL_DEBUG="file:warn,svc:console:debug"`: any application running using Sypl with an output called `file` will log messages bellow the `warn` level, and any application running using Sypl with a component called `svc` with an output called `console` will log messages bellow the `debug`.

Possible scopes:

- `{componentName:outputName:level}`: Component, and output scoped.
- `{outputName:level}`: Output scoped.
- `{level}`: Global scope.

The possibilities are endless! Checkout the [`debugAndFilter`](example_test.go) example for more.
### Changed
- Renamed logging component filtering env var from `SYPL_DEBUG` to `SYPL_FILTER`.

## [1.5.3] - 2021-09-21
### Changed
- Fix bug where setting fields for a message would set globally too.

## [1.5.2] - 2021-09-21
### Changed
- Level `FromString`, and `MustFromString` methods validates if `level` param is empty.

## [1.5.1] - 2021-09-10
### Changed
- Sypl `SetFields` is chainable.

## [1.5.0] - 2021-09-10
### Added
- Adds the ability to set global Fields.

## [1.4.6] - 2021-08-30
### Changed
- `FromString` error now prints also available levels.
- `LevelsNames` returns lower-cased levels.

## [1.4.5] - 2021-08-30
### Changed
- `StdErr` now only prints `Error` AND `Fatal` instead of only `Error`.
- `Console` now ignores `Error` AND `Fatal` instead of only `Error`.
- `PrintOnlyAtLevel` now handle multiples levels.
- `FromString` now returns the level, and error instead of level, and bool (ok).
- Internal `sypl.process` is now validated. In case of failure it throws `ErrSyplNotInitialized`.
- All `error.go` files were renamed to `errors.go`, following Go standards.

## [1.4.4] - 2021-08-20
### Added
- Adds `PrintNewLine`.

### Changed
- `Skip` and `SkipAndForce` flags now skips formatters too.

## [1.4.3] - 2021-08-20
### Changed
- Removes unused entries from `Makefile`.
- `sypl.New` now returns `*Sypl`.

## [1.4.2] - 2021-08-19
### Added
- Adds `PrintMessagesToOutputsWithOptions`.

## [1.4.1] - 2021-08-19
### Changed
- Allows to specify the name of `dashHandler` output.
    - Now, when `-` is specified as a path, `dashHandler` is named after the original output.

## [1.4.0] - 2021-08-18
### Changed
- Fixed names of the factories, so it doesn't stutters.

Note: Breaking change.

## [1.3.11] - 2021-08-18
### Added
- Adds `LevelsNames`.

## [1.3.10] - 2021-08-18
### Changed
- Improved `FromString`.

## [1.3.9] - 2021-08-18
### Added
- Adds `MustFromString`.

## [1.3.8] - 2021-08-18
### Added
- Adds the ability to get and set outputs' max level.

## [1.3.7] - 2021-08-17
### Changed
- `NewDefault` only prints errors to `stderr`.

## [1.3.6] - 2021-08-17
### Added
- Adds `PrintOnlyIfTagged` built-in processor.

### Changed
- Renames `PrintOnlyLevel` to `PrintOnlyAtLevel`.

## [1.3.5] - 2021-08-17
### Changed
- `StdErr` only prints @ `Error` `Level`.

## [1.3.4] - 2021-08-17
### Added
- Creates `StdErr` built-in `Output`.

### Changed
- Removes `path` (unused) from `FileBased` `Output`.

## [1.3.3] - 2021-08-14
### Changed
- Improved linebreak detection and restoration.

## [1.3.2] - 2021-08-13
### Added
- Adds `PrintMessagerPerOutput` which allows you to concurrently print messages, each one, at the specified level and to the specified output. If the named output doesn't exits, the message will not be printed.
    - Cover with test.

### Changed
- Adds `output` field to `Text` and `JSON` formatters.

## [1.3.1] - 2021-08-11
### Added
- Adds the ability to create child loggers (`New`). The child logger is an accurate, and efficient shallow copy of the parent logger. Changes to internals, such as the state of outputs, and processors, are reflected cross all other loggers.
- Adds `Text`, and `JSON` formatters. It also process fields. See `example_test.go/ExampleNew_textFormatter` and `example_test.go/ExampleNew_jsonFormatter` for examples. Both formatters automatically adds:
    - Component name
    - Level
    - Timestamp (RFC3339).
- Add more tests. Covered `ErrorSimulator` processor.
- Adds ability to filter logging message. See `example_test.go/ExampleNew_childLoggers` for example. Having many loggers can be, sometimes, noisy. Also, sometimes - for debugging reason, you may want to see only `componentA`, and `componentC`. Now, it's possible. Just specify the name of the components (comma-separated list) in the `SYPL_DEBUG` env var.

## [1.3.0] - 2021-08-10
### Added
- Adds support for structured logging.
    - See `example_test.go/ExampleNew_fieldsProcessing`.
- Components are interface(behaviour)-driven (design-pattern).
- Components are Factory built (design-pattern).
- Adds `Buffer` built-in `output`, it's a concurrent-safe buffer.
- Refactored code, components are packaged.

## [1.2.5] - 2021-07-22
### Added
- Adds the `Decolourizer` processor.

## [1.2.4] - 2021-07-16
### Changed
- Go mod checksum.

## [1.2.3] - 2021-07-15
### Added
- Adds `Sprint{f|lnf|ln}`, and `{Level}{f|lnf|ln}` Convenient methods. It's your `Sprint`, or `Sinfo` (example) but also returning the non-processed content.

Before:

```go
// ...
var errMsg := "Some error"

logger.Errorln(errMsg)

return errors.New(errMsg)
```

Now:

```go
// ...
return logger.Serrorln("Some error") // Prints and process like `Errorln`, and returns an error.
```

## [1.2.2] - 2021-07-15
### Changed
- Fixes `Flag`s processing logic.
- Covers `Flag`s with test.

## [1.2.1] - 2021-07-15
### Changed
- Fixes `prettify` not printing the error if it fails.

## [1.2.0] - 2021-07-14
### Added
- Finer-control on message's behaviour with two new `Flags`: `SkipAndForce` and `SkipAndMute`.
- Adds `Printlnf`, and `{Level}{lnf}` Convenient methods. It's your `Printf`, or `Infof` (example) without the need to add `"\n"` to the format - less annoying repetition.

Before:

```go
// ...
exampleContent := "example"
logger.Printf("Something %s\n", exampleContent)
```

Now:

```go
// ...
exampleContent := "example"
logger.Printlnf("Something %s", exampleContent)
```

### Changed
- Improves testability, and maintainability: All "Convenient methods" are based on "Base methods" that are based on the implementation of the interface.
    - Testability: You mock the interface, and have full control over how it works.
    - Maintainability: You change the interface implementation, you change how everything works.

## [1.1.2] - 2021-07-13
### Changed
- Fix typo (`spyl`).

## [1.1.1] - 2021-07-13
# Added
- Adds `Print{ln}Pretty` which allows to print data structures as JSON text.

Now:

```go
// ...
logger.PrintlnPretty(&SomeStruct{
    nonExportedKey: "Value1",
    SomeExportedKey: "Value2",
})

// Prints:
// {
//     "SomeExportedKey": "Value2"
// }
```

### Changed
- Prefixes sypl errors making it easier to identify when happens.
- Fixes a bug in `level.FromString` where invalid string would call `log.Fatal`.

## [1.1.0] - 2021-07-13
### Added
- Adds the ability to tag a message, see new `Print{f,ln}WithOptions` example.
- Adds the ability to flag a message, see new `Skip` flag.
- Adds `Print{f,ln}WithOptions` which allows to specify message's `Options` such as a list of `Output`s and `Processor`s to be used.
- Functional approach: no direct-access to data structure properties.
- Adds more examples.
- Adds more tests.
- Adds more documentation.
- Extracted `Flag`, `Content` and `Level` to packages.

## [1.0.0] - 2021-07-08
### Added
- First release.
