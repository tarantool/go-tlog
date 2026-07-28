# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Text and JSON handlers compatible with slog.
- Support for reopening log files.

### Changed
- Logger outputs close is thread-safe now.

### Fixed


## [v1.0.0] - 2026-01-13

### Added

- Core structured logging library for Go.
- Support for log levels: Trace, Debug, Info, Warn, Error.
- Text and JSON output formats.
- Multiple output destinations: stdout, stderr, file paths, multi-target.
- Automatic timestamp, source file and line number.
- Stacktrace support with configurable stacktrace level.
- Test suite for core functionality.
- Idiomatic Go examples (testable examples).
- Makefile, GitHub Actions CI workflow, README, LICENSE, lint configuration.
