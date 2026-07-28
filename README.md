<a href="http://tarantool.org">
	<img src="https://avatars2.githubusercontent.com/u/2344919?v=2&s=250" align="right">
</a>

[![Go Reference][godoc-badge]][godoc-url]
[![Actions Status][actions-badge]][actions-url]
[![Coverage Status][coverage-badge]][coverage-url]
[![Telegram EN][telegram-badge]][telegram-en-url]
[![Telegram RU][telegram-badge]][telegram-ru-url]

## Table of contents

- [go-tlog](#go-tlog)
  - [Features](#features)
  - [Requirements](#requirements)
  - [Installation](#installation)
  - [Quick start](#quick-start)
  - [Documentation](#documentation)
  - [Configuration](#configuration)
    - [type Opts](#type-opts)
    - [Main API](#main-api)
  - [Log levels](#log-levels)
  - [Stacktraces](#stacktraces)
    - [Overriding stacktrace level](#overriding-stacktrace-level)
  - [Output formats](#output-formats)
  - [Output destinations](#output-destinations)
  - [Using the handlers directly](#using-the-handlers-directly)
  - [Examples](#examples)
  - [Testing](#testing)
  - [Contributing](#contributing)
  - [License](#license)

# go-tlog

`go-tlog` is a lightweight and configurable logging library for Go applications.  
It provides structured logging with multiple output destinations, flexible formatting,
and fine-grained log-level control.

---

## Features

- Simple setup via configuration struct
- Text or JSON output formats
- Multiple output targets: **stdout**, **stderr**, **files**
- Log levels: `Trace`, `Debug`, `Info`, `Warn`, `Error`
- Automatic timestamp, source file, and line number
- Automatic stacktraces based on log level

---

## Requirements

| Requirement | Version                                        |
|-------------|------------------------------------------------|
| Go          | 1.24 or newer                                  |
| OS          | Linux, macOS, Windows — any platform Go targets |

The module depends only on the standard library at runtime; `testify` is used
in tests.

---

## Installation

```bash
go get github.com/tarantool/go-tlog@latest
```

Then import:

```go
import "github.com/tarantool/go-tlog"
```

---

## Quick start

```go
package main

import (
	"log/slog"

	"github.com/tarantool/go-tlog"
)

func main() {
	log, err := tlog.New(tlog.Opts{
		Level:  tlog.LevelInfo,
		Format: tlog.FormatText,
		Path:   "stdout",
	})
	if err != nil {
		panic(err)
	}
	defer log.Close()

	logger := log.Logger().With(slog.String("component", "demo"))
	logger.Info("service started", slog.Int("port", 8080))
	logger.Error("failed to connect", slog.String("err", "timeout"))
}
```

Output — a timestamp, the level, the call site, the quoted message, then the
attributes:

```
2025-11-10T13:30:01+05:00 INFO /app/main.go:19 "service started" component=demo port=8080 stacktrace="..."
2025-11-10T13:30:01+05:00 ERROR /app/main.go:20 "failed to connect" component=demo err=timeout stacktrace="..."
```

Both lines carry a stacktrace because the default threshold follows the log
level — see [Stacktraces](#stacktraces) for how to raise it.

---

## Documentation

The API reference is the package documentation on
[pkg.go.dev][godoc-url]. It is generated from the source, so it never goes
stale, and the `ExampleXxx` functions shown there are compiled and executed by
CI.

Locally:

```bash
go doc -http=:6060
```

---

## Configuration

### `type Opts`

```go
type Opts struct {
    Level  Level  // minimal log level
    Format Format // FormatText or FormatJSON
    Path   string // comma-separated outputs: "stdout,/var/log/app.log"
}
```

### Main API

| Function         | Description                              |
|------------------|------------------------------------------|
| `tlog.New(opts)` | Create a new logger                      |
| `Logger()`       | Return the underlying logger for use     |
| `Close()`        | Flush buffers and close file descriptors |

---

## Log levels

| Level   | When to use                                 |
|---------|---------------------------------------------|
| `Trace` | Low-level tracing                           |
| `Debug` | Debugging information                       |
| `Info`  | Normal operational messages                 |
| `Warn`  | Non-fatal warnings                          |
| `Error` | Errors and exceptions (includes stacktrace) |

---

## Stacktraces

`go-tlog` can automatically attach stacktraces to log records.

By default, the stacktrace threshold is the same as the configured log level.
This means that stacktraces are added starting from the current log level
and for all higher-severity messages.

The default behavior is:

| Log level | Stacktrace is added for          |
|-----------|----------------------------------|
| `Trace`   | `DEBUG`, `INFO`, `WARN`, `ERROR` |
| `Debug`   | `DEBUG`, `INFO`, `WARN`, `ERROR` |
| `Info`    | `INFO`, `WARN`, `ERROR`          |
| `Warn`    | `WARN`, `ERROR`                  |
| `Error`   | `ERROR`                          |

You can override this behavior using `StacktraceLevel` to control
the stacktrace threshold independently of the log level.

### Overriding stacktrace level

```go
log, err := tlog.New(tlog.Opts{
    Level:           tlog.LevelInfo,
    StacktraceLevel: tlog.LevelError,
    Format:          tlog.FormatText,
    Path:            "stdout",
})
```

---

## Output formats

| Format       | Example                                                             |
|--------------|---------------------------------------------------------------------|
| `FormatText` | `2025-11-10T13:31:45+05:00 INFO /app/main.go:19 "message" key=value` |
| `FormatJSON` | `{"time":"...","level":"INFO","msg":"message","key":"value"}`        |

In text output the built-in fields are printed without their keys, and the
message is quoted when it contains spaces or other characters that would need
escaping.

---

## Output destinations

You can specify multiple targets separated by commas:

```go
Path: "stdout,/tmp/app.log"
```

Supported targets:

- `stdout`
- `stderr`
- File paths (created automatically if not present)

---

## Using the handlers directly

If your application already builds its own `slog.Logger`, use the handler
constructors instead of `tlog.New`. They give the same formatting and
stacktrace behavior over any `io.Writer`, without managing outputs for you:

```go
handler := tlog.NewTextHandler(os.Stdout, &tlog.HandlerOptions{
    Level:           tlog.LevelDebug,
    StacktraceLevel: tlog.LevelError,
})

logger := slog.New(handler)
```

`tlog.NewJSONHandler` is the JSON counterpart. Passing `nil` options uses the
defaults.

---

## Examples

Included examples:

- **ExampleNew_text** — basic text logger writing to stdout
- **ExampleNew_json** — JSON logging
- **ExampleNew_multi** — logging to multiple destinations (`stdout,/tmp/...`)
- **ExampleNewTextHandler** — building an `slog.Logger` over a text handler
- **ExampleNewJSONHandler** — the same with JSON output

Each example demonstrates different combinations of Path, Format, and Level,
including how to log to multiple outputs at the same time.

---

## Testing

```bash
make test           # run the test suite
make test-race      # run with the race detector
make test-coverage  # write coverage.out and print the total
make lint           # golangci-lint with the project config
```

---

## Contributing

Bug reports and pull requests are welcome — see
[CONTRIBUTING.md](CONTRIBUTING.md) for the commit message format, the review
rules, and the checks a PR is expected to pass. Issues labeled
[`good first issue`][good-first-issue] are a reasonable place to start.

---

## License

BSD 2-Clause License — see [LICENSE](LICENSE).

Parts of this module are derived from the Go standard library and are covered
by a separate BSD 3-Clause license — see [NOTICE](NOTICE).

[actions-badge]: https://github.com/tarantool/go-tlog/actions/workflows/testing.yml/badge.svg
[actions-url]: https://github.com/tarantool/go-tlog/actions/workflows/testing.yml
[coverage-badge]: https://coveralls.io/repos/github/tarantool/go-tlog/badge.svg?branch=master
[coverage-url]: https://coveralls.io/github/tarantool/go-tlog?branch=master
[good-first-issue]: https://github.com/tarantool/go-tlog/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22
[godoc-badge]: https://pkg.go.dev/badge/github.com/tarantool/go-tlog.svg
[godoc-url]: https://pkg.go.dev/github.com/tarantool/go-tlog
[telegram-badge]: https://img.shields.io/badge/Telegram-join%20chat-blue.svg
[telegram-en-url]: http://telegram.me/tarantool
[telegram-ru-url]: http://telegram.me/tarantoolru
