# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Changed

### Fixed

## [1.1.0] - 2026-08-04

This release adds standard `slog.Handler` constructors — `NewTextHandler` and
`NewJSONHandler` over any `io.Writer` — so an application can build its own
`slog.Logger` with the go-tlog format, and support for reopening file outputs
so a logger can survive logrotate. It also fixes lost stacktraces on derived
loggers and a panic in the text format.

### Added

- `NewTextHandler` and `NewJSONHandler` construct a standard `slog.Handler`
  over any `io.Writer`, so an application can build its own `slog.Logger` with
  the go-tlog format and stacktrace behavior (#3).
- `Logger.Reopen` and `Logger.ReopenOnSignal` reopen file outputs, so a logger
  can survive logrotate (#11).

### Changed

- Logger outputs are safe to close concurrently with writes (#11).

### Fixed

- Fixed loggers derived with `Logger.With` or `Logger.WithGroup` silently
  losing their stacktraces (#12).
- Fixed a panic in the text format when a `ReplaceAttr` function dropped an
  attribute by returning the zero `slog.Attr` (#12).

## [1.0.0] - 2026-01-13

### Added

- Structured logging for Go on top of `log/slog`, configured through a single
  `Opts` value.
- Log levels `Trace`, `Debug`, `Info`, `Warn` and `Error`.
- Text and JSON output formats.
- Output to `stdout`, `stderr`, files, or several destinations at once through
  a comma-separated `Path`.
- Timestamp, source file and line number attached to every record.
- Stacktraces attached automatically from the log level upwards, with
  `StacktraceLevel` to set that threshold independently of the log level
  (#2).

[Unreleased]: https://github.com/tarantool/go-tlog/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/tarantool/go-tlog/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/tarantool/go-tlog/releases/tag/v1.0.0
