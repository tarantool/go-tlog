package tlog_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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

// ExampleLogger_Reopen shows how to reopen log file outputs.
// It is useful for integrating with logrotate: after an external
// tool moves the log file aside, call Reopen so that subsequent
// log entries are written to a freshly created file.
func ExampleLogger_Reopen() {
	path := "/tmp/tlog_example_reopen.log"

	log, err := tlog.New(tlog.Opts{
		Level:  tlog.LevelInfo,
		Format: tlog.FormatText,
		Path:   path,
	})
	if err != nil {
		panic(err)
	}
	defer func() { _ = log.Close() }()

	logger := log.Logger().With(slog.String("component", "example_reopen"))
	logger.Info("message before rotation")

	// Simulate logrotate: move the current file aside.
	if err := os.Rename(path, path+".1"); err != nil {
		panic(err)
	}

	// Reopen outputs so new entries go to a newly created file.
	if err := log.Reopen(); err != nil {
		panic(err)
	}

	logger.Info("message after rotation")
}

// ExampleLogger_ReopenOnSignal shows how to reopen log file outputs
// automatically when a signal is received. It is intended to be run in
// a separate goroutine, e.g. to reopen logs on SIGHUP sent by logrotate.
func ExampleLogger_ReopenOnSignal() {
	log, err := tlog.New(tlog.Opts{
		Level:  tlog.LevelInfo,
		Format: tlog.FormatText,
		Path:   "/tmp/tlog_example_reopen_on_signal.log",
	})
	if err != nil {
		panic(err)
	}
	defer func() { _ = log.Close() }()

	// Stop reopening on Ctrl+C / termination.
	ctx, stop := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM,
	)
	defer stop()

	// Reopen log files each time SIGHUP is received.
	go log.ReopenOnSignal(ctx, func(err error) {
		fmt.Fprintf(os.Stderr, "log reopen failed: %v\n", err)
	}, syscall.SIGHUP)

	log.Logger().Info("service started")
}
