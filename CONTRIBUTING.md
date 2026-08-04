# Contributing to go-tlog

Thanks for taking the time to contribute. This document describes how to report
issues, how to prepare a pull request, and what we check during review.

Everything in this project is in English — code, identifiers, comments,
documentation, commit messages and review discussions.

## Reporting issues

Open an issue at https://github.com/tarantool/go-tlog/issues.

For a bug report, include:

- the version of `go-tlog` (module version or commit) and the Go version;
- a minimal reproducer — ideally a failing test or a short `main.go`;
- what you expected to happen and what happened instead.

For a feature request, describe the problem you are solving before proposing an
API. Issues labeled `good first issue` are a reasonable place to start if you
are new to the project.

## Requirements

- Go 1.24 or newer.
- `golangci-lint` v2.11.4 — the version is pinned in `Makefile` and in the CI
  workflow so that local and CI results match. Install it with
  `make install-lint`.
- `codespell` for the spell-check target (optional locally, enforced in CI).

## Development workflow

```bash
make test           # run the test suite
make test-race      # run with the race detector
make test-coverage  # write coverage.out and print the total
make lint           # golangci-lint with the project config
make fmt            # gofmt/gofumpt/goimports via golangci-lint
make tidy           # go mod tidy
```

Run `make test lint` before pushing. CI runs the same targets.

## Branches

`master` is protected: no force-push, no direct commits, review required.

Work in a branch that names the issue it closes:

- `<username>/gh-<issue>-<topic>` — for a GitHub issue, e.g.
  `bigbes/gh-12-reopen-race`;
- `<username>/tntp-<issue>-<topic>` — for a Jira ticket;
- `<username>/gh-no-<topic>` — for deliberately ticketless work. Allowed, but
  discouraged: prefer filing an issue first.

Release branches are `<username>/release-X.Y.Z`.

The ticket in the branch name and the references in the commit bodies must
agree — a `gh-no-` branch carries no references at all.

## Commit messages

A commit message has three parts: a subject, a body, and references to related
issues.

```
handler: add separate json/text handlers

Add NewJSONHandler and NewTextHandler constructors that accept any
io.Writer and return a standard slog.Handler. This makes the library
conformant with the slog handler guide and allows usage with
slog.New(handler) directly.

Closes #3
```

Rules:

- Subject is `scope: summary`, where the scope is the package or area touched,
  not a change type — `handler: add …`, never `feat: add …`. Scopes used in
  this repository: `handler`, `logger`, `outputs`, `stacktrace`, `slog`,
  `api`, `ci`, `test`, `doc`.
- Keep the summary within ~50 characters, lowercase, no trailing period, in the
  imperative mood — it must complete the sentence
  "If applied, this commit will ...".
- No issue numbers in the subject line. References go in the body.
- Separate the body from the subject with a blank line; wrap body lines at 72
  characters.
- The body explains *what* and *why*, not *how*.
- Issue references go on the last lines of the body. Use `Part of #NNN` for
  intermediate commits of a multi-commit change and `Closes #NNN` in the final
  one.
- Record co-authorship with a `Co-authored-by:` trailer.
- Use your real name and a working email address.
- Every commit must be atomic and logically complete.

## Pull requests

Every patch needs a test. A regression test for a bug fix must fail without the
fix. Exceptions are rare and discussed case by case.

Before requesting review, self-review your own diff:

- all intended changes are in the PR, and nothing else;
- the diff is under 500–700 lines, excluding file deletions;
- you can explain every line;
- `make test lint` passes;
- new functionality is documented — doc comments, `README.md`, and an
  `ExampleXxx` where it helps;
- user-visible changes have a `CHANGELOG.md` entry under `## [Unreleased]`.

A PR needs **two approvals** to merge. Review threads are resolved by the
reviewer who opened them, not by the author. An unresolved disagreement is
decided by a maintainer.

Reviewers check readability, whether the change matches the issue, code style,
tests, and the absence of references to non-public resources.

## Code style

We follow the [Google Go Style Guide][style] and
[Go Code Review Comments][review]. Formatting and linting are automated —
`gofmt`, `gofumpt` and `goimports` run as part of `golangci-lint`, and the lint
config enables all linters by default with individually justified exceptions.

Style disagreements are settled by the style guide, not by argument in the PR.

The `depguard` linter restricts what may be imported: the standard library and
the explicitly allowed modules in `.golangci.yaml`. This keeps non-public
dependencies out of a public module. If you genuinely need a new dependency,
raise it in the issue first.

Code under `internal/` is not part of the public API and may change without
notice. `internal/slog/` is a vendored fork of the standard library's
`log/slog`; keep changes there minimal and preserve the upstream copyright
headers.

## Changelog

`CHANGELOG.md` follows [Keep a Changelog 1.1.0][kac] and the project follows
[Semantic Versioning][semver].

Write entries for users, not for yourself. Do not log changes to internal APIs,
tests, or CI. Merge scattered notes about one feature into a single entry.
Describe fixes in the past tense and put the issue reference at the end:

```markdown
### Fixed

- Fixed a panic when Path pointed at an unwritable directory (#42).
```

## License

By contributing you agree that your contribution is licensed under the BSD
2-Clause License, the same license that covers this project — see `LICENSE`.

[style]: https://google.github.io/styleguide/go/
[review]: https://go.dev/wiki/CodeReviewComments
[kac]: https://keepachangelog.com/en/1.1.0/
[semver]: https://semver.org/spec/v2.0.0.html
