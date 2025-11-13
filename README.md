[![CI](https://github.com/tarantool/go-tlog/actions/workflows/ci.yml/badge.svg)](https://github.com/tarantool/go-tlog/actions/workflows/ci.yml) •
[Telegram EN](https://t.me/tarantool) •
[Telegram RU](https://t.me/tarantoolru)

---

<p href="https://www.tarantool.io">
	<img src="https://github.com/tarantool.png" align="right" width=250>
</p>

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
- Stacktrace for errors

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

	logger := log.Logger().With(tlog.String("component", "demo"))
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

Ready-to-run examples are located in the `_examples/` directory:

```
_examples/
 ├── stdout/
 │   └── main.go
 ├── stderr/
 │   └── main.go
 ├── file/
 │   └── main.go
 └── multi/
     └── main.go
```

Run examples:

```bash
# Example 1 — log to STDOUT in text format
go run ./_examples/stdout

# Example 2 — log to STDERR in JSON format
# Redirect stderr to a file and inspect its contents
go run ./_examples/stderr 2> logs.json
cat logs.json

# Example 3 — log to a file in /tmp directory
# The file will be created automatically if it doesn’t exist
go run ./_examples/file
cat /tmp/tlog_demo/app.log

# Example 4 — log to multiple destinations (stdout + file)
# This writes the same log entry both to console and to /tmp/tlog_multi/app.log
go run ./_examples/multi
cat /tmp/tlog_multi/app.log
```

Each example demonstrates different combinations of Path, Format, and Level,
including how to log to multiple outputs at the same time.

---

## Testing

```bash
go test ./...
```

---

## License

BSD 2-Clause License — see [LICENSE](LICENSE)
