# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Roadmap

- Add support for structured logging.

## [1.2.1] - 2021-07-15
### Changed
- Adds the error message for when `prettify` fails.

## [1.2.0] - 2021-07-14
### Changed
- Finer-control on message's behaviour with two new `Flags`: `SkipAndForce` and `SkipAndMute`.
- Improves testability, and maintainability: All "Convenient methods" are based on "Base methods" that are based on the implementation of the interface.
    - Testability: You mock the interface, and have full control over how it works.
    - Maintainability: You change the interface implementation, you change how everything works.
- Adds `Printlnf`, and `{Level}{lnf}` Convenient methods. It's your `Printf`, or `Infof` without the need to add `"\n"` to the format - less annoying repetition.

## [1.1.2] - 2021-07-13
### Changed
- Fix typo (`spyl`).

## [1.1.1] - 2021-07-13
# Added
- Adds `Print{ln}Pretty` which allows to print data structures as JSON text.

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
