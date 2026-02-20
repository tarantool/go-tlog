package tlog_test

import (
	"log/slog"
	"os"

	"github.com/tarantool/go-tlog"
)

// ExampleNew_text shows how to create a text logger that writes to stdout.
func ExampleNew_text() {
	log, err := tlog.New(tlog.Opts{
		Level:  tlog.LevelInfo,
		Format: tlog.FormatText,
		Path:   "stdout",
	})
	if err != nil {
		panic(err)
	}
	defer func() { _ = log.Close() }()

	logger := log.Logger().With(slog.String("component", "example_text"))
	logger.Info("service started")
}

// ExampleNew_json shows how to create a JSON logger that writes to stdout.
func ExampleNew_json() {
	log, err := tlog.New(tlog.Opts{
		Level:  tlog.LevelInfo,
		Format: tlog.FormatJSON,
		Path:   "stdout",
	})
	if err != nil {
		panic(err)
	}
	defer func() { _ = log.Close() }()

	logger := log.Logger().With(
		slog.String("component", "example_json"),
		slog.String("request_id", "abc-123"),
	)
	logger.Error("request failed")
}

// ExampleNew_multi shows how to log to multiple destinations at once.
func ExampleNew_multi() {
	log, err := tlog.New(tlog.Opts{
		Level:  tlog.LevelInfo,
		Format: tlog.FormatText,
		Path:   "stdout,/tmp/tlog_example_multi.log",
	})
	if err != nil {
		panic(err)
	}
	defer func() { _ = log.Close() }()

	logger := log.Logger().With(slog.String("component", "example_multi"))
	logger.Info("message written to stdout and file")
}

// ExampleNewJSONHandler shows how to create a JSON handler
// that writes to any io.Writer.
func ExampleNewJSONHandler() {
	handler := tlog.NewJSONHandler(os.Stdout, &tlog.HandlerOptions{
		Level:           tlog.LevelInfo,
		StacktraceLevel: tlog.LevelError,
	})
	logger := slog.New(handler)

	logger.Info("service started")
}

// ExampleNewTextHandler shows how to create a text handler
// that writes to any io.Writer.
func ExampleNewTextHandler() {
	handler := tlog.NewTextHandler(os.Stdout, &tlog.HandlerOptions{
		Level:           tlog.LevelInfo,
		StacktraceLevel: tlog.LevelError,
	})
	logger := slog.New(handler)

	logger.Info("service started")
}
