<a href="http://tarantool.org">
	<img src="https://avatars2.githubusercontent.com/u/2344919?v=2&s=250" align="right">
</a>

[![Actions Status][actions-badge]][actions-url]
[![Telegram EN][telegram-badge]][telegram-en-url]
[![Telegram RU][telegram-badge]][telegram-ru-url]

## Table of contents

- [go-tlog](#go-tlog)
  - [Features](#features)
  - [Installation](#installation)
  - [Quick start](#quick-start)
  - [Configuration](#configuration)
    - [type Opts](#type-opts)
    - [Main API](#main-api)
  - [Log levels](#log-levels)
  - [Stacktraces](#stacktraces)
    - [Overriding stacktrace level](#overriding-stacktrace-level)
  - [Output formats](#output-formats)
  - [Output destinations](#output-destinations)
  - [Examples](#examples)
  - [Testing](#testing)
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

import "github.com/tarantool/go-tlog"

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
	logger.Info("service started", "port", 8080)
	logger.Error("failed to connect", "err", "timeout")
}
```

Output:

```
2025-11-10T13:30:01+05:00 INFO service started component=demo port=8080
2025-11-10T13:30:01+05:00 ERROR failed to connect err=timeout component=demo stacktrace="..."
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

| Format       | Example                                                       |
|--------------|---------------------------------------------------------------|
| `FormatText` | `2025-11-10T13:31:45+05:00 INFO message key=value`            |
| `FormatJSON` | `{"time":"...","level":"INFO","msg":"message","key":"value"}` |

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

## Examples

Included examples:

- **ExampleNew_text** — basic text logger writing to stdout  
- **ExampleNew_json** — JSON logging  
- **ExampleNew_multi** — logging to multiple destinations (`stdout,/tmp/...`)

Each example demonstrates different combinations of Path, Format, and Level,
including how to log to multiple outputs at the same time.

---

## Testing

```bash
make test
```

---

## License

BSD 2-Clause License — see [LICENSE](LICENSE)

[actions-badge]: https://github.com/tarantool/go-tlog/actions/workflows/testing.yml/badge.svg
[actions-url]: https://github.com/tarantool/go-tlog/actions/workflows/testing.yml
[telegram-badge]: https://img.shields.io/badge/Telegram-join%20chat-blue.svg
[telegram-en-url]: http://telegram.me/tarantool
[telegram-ru-url]: http://telegram.me/tarantoolru
