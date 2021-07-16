# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Roadmap

- Add support for structured logging.

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
