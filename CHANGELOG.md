# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Roadmap

- Improve documentation:
    - Add `doc.go` for all packages.
    - Update `README.md` to point to these new `doc.go`.

## [1.3.3] - 2021-08-14
## Changed
- Improved linebreak detection and restoration.

## [1.3.2] - 2021-08-13
## Added
- Adds `PrintMessagerPerOutput` which allows you to concurrently print messages, each one, at the specified level and to the specified output. If the named output doesn't exits, the message will not be printed.
    - Cover with test.

## Changed
- Adds `output` field to `Text` and `JSON` formatters.

## [1.3.1] - 2021-08-11
## Added
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
