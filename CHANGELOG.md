# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Roadmap

- Add support for structured logging.

## [1.1.2] - 2021-07-13
### Changed
- Fix typo (`spyl`)

## [1.1.1] - 2021-07-13
# Added
- Adds `print{ln}Pretty` which allows to print data structures as JSON text.

### Changed
- Prefixes sypl errors making it easier to identify when happens.
- Fixes a bug in `level.FromString` where invalid string would call `log.Fatal`

## [1.1.0] - 2021-07-13
### Added
- Adds the ability to tag a message, see new `print{f,ln}WithOptions` example.
- Adds the ability to flag a message, see new `Skip` flag.
- Adds `print{f,ln}WithOptions` which allows to specify message's `Options` such as a list of `Output`s and `Processor`s to be used.
- Functional approach: no direct-access to data structure properties.
- Adds more examples.
- Adds more tests.
- Adds more documentation.
- Extracted `Flag`, `Content` and `Level` to packages.

## [1.0.0] - 2021-07-08
### Added
- First release.
