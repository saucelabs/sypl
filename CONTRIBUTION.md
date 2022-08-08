# Contributing to _SYPL_

**Thank you for your interest in _SYPL_. Your contributions are highly welcome.**

There are multiple ways of getting involved:

- [Report a bug](#report-a-bug)
- [Suggest a feature](#suggest-a-feature)
- [Contribute code](#contribute-code)
- [Release](#release)

## Report a bug
Reporting bugs is one of the best ways to contribute. Before creating a bug report, please check that an [issue](/issues) reporting the same problem does not already exist. If there is such an issue, you may add your information as a comment.

To report a new bug you should open an issue that summarizes the bug and set the label to "bug".

If you want to provide a fix along with your bug report: That is great! In this case please send us a pull request as described in section [Contribute Code](#contribute-code).

## Suggest a feature
To request a new feature you should open an [issue](../../issues/new) and summarize the desired functionality and its use case. Set the issue label to "feature".

## Contribute code

The following guidelines aim to point to a direction that should drive the codebase to increased quality.

- Each package should have a `doc.go` file.
- Think before you make changes, design, then code. Design patterns, and well-established techniques are your friend. They allow to reduce code duplication and complexity and increase reusability and performance.
- Documentation is essential! Relevant comments should be added focusing on the **why**, not in the **what**. _Pay attention to the punctuation and casing patterns_
- Pay attention to how the code is vertically spaced and positioned, also sorted (always ascending) for example, the content of a struct, `vars` and `const` declarations, and etc.
- If you use VSCode IDE, the Go extension is installed, **_and properly setup_**, it should obey the configuration file ([.golangci.yml](.golangci.yml)) for the linter (`golangci`) and show problems the right way, otherwise, just run `$ make lint`. The same thing applies to test. If you open a test file (`*_test.go`), modify and save it, it should automatically run tests and shows coverage; otherwise, just run `$ make test`
- Always run ` $ make coverage lint` before you commit your code; it will save you time!
- If you spotted a problem or something that needs to be modified/improved, do that right way; otherwise, that with `// TODO:`
- Update the [`CHANGELOG.md`](CHANGELOG.md) but not copying your commits messages - that's not its purpose. Use that to plan changes too.

## Release

1. create a branch, commit update and push
1. once all test pass and PR is approved, merge
1. Update [the changelog](CHANGELOG)
1. make a new release by creating a tag that matches the new Sauce Connect version:
   ```sh
   $ git checkout master
   # fetch latest code
   $ git pull origin master
   $ git tag -a "vX.X.X"
   ```
1. push tag
   ```sh
   $ git push origin master --tags
   ```
1. create a "GitHub release" from tag

**Have fun, and happy hacking!**
