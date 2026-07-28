package tlog

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/tarantool/go-tlog/internal/outputs"
	slogcustom "github.com/tarantool/go-tlog/internal/slog"
)

// Logger contains slog.Logger for several outputs
// and methods to reopen or close these outputs.
type Logger struct {
	outputs *outputs.Outputs
	logger  *slog.Logger
}

// Opts are New options.
type Opts struct {
	// Level sets minimum level for the logger.
	Level Level
	// Format sets log format.
	Format Format
	// Path is comma-separated list of log outputs.
	// Use "stdout" and "stderr" for os streams and file paths for files.
	// Default is "stderr".
	Path string
	// StacktraceLevel overrides the default stacktrace threshold.
	// If set, stacktraces will be attached starting from this level,
	// regardless of the main log Level.
	StacktraceLevel Level
}

// New creates a new Logger with the given options.
// It configures level, format and output destinations and returns
// a ready-to-use logger instance.
func New(opts Opts) (*Logger, error) {
	var (
		logLevel   slog.Level
		traceLevel slog.Level
	)

	switch opts.Level {
	case LevelTrace:
		fallthrough
	case LevelDebug:
		traceLevel = slog.LevelDebug
		logLevel = slog.LevelDebug
	case LevelDefault:
		fallthrough
	case LevelInfo:
		traceLevel = slog.LevelInfo
		logLevel = slog.LevelInfo
	case LevelWarn:
		traceLevel = slog.LevelWarn
		logLevel = slog.LevelWarn
	case LevelError:
		traceLevel = slog.LevelError
		logLevel = slog.LevelError
	}

	// Override stacktrace level if explicitly configured.
	if opts.StacktraceLevel != 0 {
		switch opts.StacktraceLevel {
		case LevelTrace, LevelDebug:
			traceLevel = slog.LevelDebug
		case LevelInfo:
			traceLevel = slog.LevelInfo
		case LevelWarn:
			traceLevel = slog.LevelWarn
		case LevelError:
			traceLevel = slog.LevelError
		default:
			traceLevel = slog.LevelDebug
		}
	}

	if opts.Path == "" {
		// https://github.com/uber-go/zap/blob/6d482535bdd97f4d97b2f9573ac308f1cf9b574e/config.go#L167C31-L167C37
		opts.Path = "stderr"
	}

	outs, err := outputs.New(opts.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to create outputs: %w", err)
	}

	handlerOpts := slog.HandlerOptions{
		Level:       logLevel,
		ReplaceAttr: replaceAttr,
		AddSource:   true,
	}

	var baseHandler slog.Handler

	switch opts.Format {
	case FormatDefault:
		fallthrough
	case FormatText:
		baseHandler = slogcustom.NewTextHandler(outs, &slogcustom.HandlerOptions{
			HandlerOptions:  handlerOpts,
			OmitBuiltinKeys: true,
		})
	case FormatJSON:
		baseHandler = slog.NewJSONHandler(outs, &handlerOpts)
	}

	handler := newStacktraceHandler(baseHandler, traceLevel)
	l := slog.New(handler)

	return &Logger{
		outputs: outs,
		logger:  l,
	}, nil
}

func replaceAttr(group []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.TimeKey:
		return replaceTime(group, a)
	default:
		return a
	}
}

func replaceTime(_ []string, a slog.Attr) slog.Attr {
	t := a.Value.Time()

	a.Value = slog.StringValue(t.Format(time.RFC3339))

	return a
}

// Logger returns the underlying slog.Logger instance.
// It can be used directly to log messages with additional attributes.
func (l *Logger) Logger() *slog.Logger {
	return l.logger
}

// Reopen closes current file outputs and opens them once again.
// Can be used for logrotate.
func (l *Logger) Reopen() error {
	return l.outputs.Reopen()
}

// ReopenOnSignal calls Reopen each time one of sigs arrives and returns when
// ctx is done. It is blocking, intended to run in a goroutine:
//
//	go log.ReopenOnSignal(ctx, func(err error) {
//	        fmt.Fprintf(os.Stderr, "log reopen failed: %v\n", err)
//	}, syscall.SIGHUP)
//
// onErr must not write through this Logger.
func (l *Logger) ReopenOnSignal(ctx context.Context, onErr func(error), sigs ...os.Signal) {
	if len(sigs) == 0 {
		return
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sigs...)
	defer signal.Stop(ch)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ch:
			if err := l.Reopen(); err != nil && onErr != nil {
				onErr(err)
			}
		}
	}
}

// Close flushes all pending log entries and closes all opened outputs.
func (l *Logger) Close() error {
	return l.outputs.Close()
}
