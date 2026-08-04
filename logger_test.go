package tlog_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/tarantool/go-tlog"
)

func Test_Logger(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		// Test assertions expect single file opts.Path.
		opts   tlog.Opts
		log    func(l *slog.Logger)
		assert func(require *require.Assertions, logs string)
	}{
		{
			name: "InfoMessage_DebugTextLogger",
			opts: tlog.Opts{
				Level:           tlog.LevelDebug,
				StacktraceLevel: tlog.LevelError,
				Format:          tlog.FormatText,
				Path:            "InfoMessage_DebugPlainLogger.log",
			},
			log: func(l *slog.Logger) {
				l.Info("my info message")
				// Example:
				// 2025-02-19T13:51:31+03:00 INFO tlog_test.go:<line> "my info message"
			},
			assert: func(require *require.Assertions, logs string) {
				require.Contains(logs, "my info message")
				require.NotContains(logs, "stacktrace=")
			},
		},
		{
			name: "WarnMessage_InfoTextLogger_WithWarnStacktraceLevel",
			opts: tlog.Opts{
				Level:           tlog.LevelInfo,
				StacktraceLevel: tlog.LevelWarn,
				Format:          tlog.FormatText,
				Path:            "WarnMessage_InfoTextLogger_WithWarnStacktraceLevel.log",
			},
			log: func(l *slog.Logger) {
				l.Warn("my warn message")
			},
			assert: func(require *require.Assertions, logs string) {
				require.Contains(logs, "my warn message")
				require.Contains(logs, "stacktrace=")
				require.Contains(logs, "tlog_test.Test_Logger")
			},
		},
		{
			name: "ErrorMessage_DebugTextLogger",
			opts: tlog.Opts{
				Level:  tlog.LevelDebug,
				Format: tlog.FormatText,
				Path:   "ErrorMessage_DebugPlainLogger.log",
			},
			log: func(l *slog.Logger) {
				l.Error("my error message")
				// 2025-02-19T13:52:11+03:00 ERROR logger_test.go:46 "my error message"
				// stacktrace="github.com/tarantool/go-tlog_test.Test_Logger.funcY
				//     logger_test.go:<line>
				// logger_test.Test_Logger.funcX
				//     logger_test.go:<line>
				// testing.tRunner
				//     /usr/local/go/src/testing/testing.go:<line>
				// runtime.goexit
				//     /usr/local/go/src/runtime/asm_amd64.s:<line>"
			},
			assert: func(require *require.Assertions, logs string) {
				require.Contains(logs, "my error message")
				require.Contains(logs, "stacktrace=")
				require.Contains(logs, "tlog_test.Test_Logger")
			},
		},
		{
			name: "InfoMessage_ErrorTextLogger",
			opts: tlog.Opts{
				Level:  tlog.LevelError,
				Format: tlog.FormatText,
				Path:   "InfoMessage_ErrorPlainLogger.log",
			},
			log: func(l *slog.Logger) {
				l.Info("my info message")
			},
			assert: func(require *require.Assertions, logs string) {
				require.NotContains(logs, "my info message")
			},
		},
		{
			name: "InfoMessage_TraceTextLogger",
			opts: tlog.Opts{
				Level:  tlog.LevelTrace,
				Format: tlog.FormatText,
				Path:   "InfoMessage_TraceTextLogger.log",
			},
			log: func(l *slog.Logger) {
				l.Info("my info message")
				// Example (shortened):
				// 2025-02-19T13:54:00+03:00 INFO tlog_test.go:<line> "my info message"
				// stacktrace="github.com/tarantool/go-tlog_test.Test_Logger.funcY
				//     logger_test.go:<line>
				// logger_test.Test_Logger.funcX
				//     logger_test.go:<line>
				// testing.tRunner
				// 		/usr/local/go/src/testing/testing.go:<line>
				// runtime.goexit
				// 		/usr/local/go/src/runtime/asm_amd64.s:<line>"
			},
			assert: func(require *require.Assertions, logs string) {
				require.Contains(logs, "my info message")
				require.Contains(logs, "stacktrace=")
				require.Contains(logs, "tlog_test.Test_Logger")
			},
		},
		{
			name: "InfoMessage_TraceTextLogger_WithErrorStacktraceLevel",
			opts: tlog.Opts{
				Level:           tlog.LevelTrace,
				StacktraceLevel: tlog.LevelError,
				Format:          tlog.FormatText,
				Path:            "InfoMessage_TraceTextLogger_WithErrorStacktraceLevel.log",
			},
			log: func(l *slog.Logger) {
				l.Info("my info message")
			},
			assert: func(require *require.Assertions, logs string) {
				require.Contains(logs, "my info message")
				require.NotContains(logs, "stacktrace=")
			},
		},
		{
			name: "InfoMessage_DefaultLogger",
			opts: tlog.Opts{
				// Level and Format will be defaulted by New.
				Path: "InfoMessage_DefaultLogger.log",
			},
			log: func(l *slog.Logger) {
				l.Warn("my info message")
			},
			assert: func(require *require.Assertions, logs string) {
				require.Contains(logs, "my info message")
			},
		},
		{
			name: "InfoMessage_DebugJSONLogger",
			opts: tlog.Opts{
				Level:           tlog.LevelDebug,
				StacktraceLevel: tlog.LevelError,
				Format:          tlog.FormatJSON,
				Path:            "InfoMessage_DebugPlainLogger.json",
			},
			log: func(l *slog.Logger) {
				l.Info("my info message")
				// Example (shortened):
				// {
				//   "time":"2025-02-19T13:55:16+03:00",
				//   "level":"INFO",
				//   "source":{
				//  	"function":"github.com/tarantool/go-tlog_test.Test_Logger.funcZ",
				// 		"file":"tlog_test.go",
				//		"line":<line>
				//	},
				//   "msg":"my info message"
				// }
			},
			assert: func(require *require.Assertions, logs string) {
				require.Contains(logs, `"msg":"my info message"`)
				require.NotContains(logs, `"stacktrace"`)
			},
		},
		{
			name: "ErrorMessage_DebugJSONLogger",
			opts: tlog.Opts{
				Level:  tlog.LevelDebug,
				Format: tlog.FormatJSON,
				Path:   "ErrorMessage_DebugJSONLogger.json",
			},
			log: func(l *slog.Logger) {
				l.Error("my error message")
				// Example (shortened):
				// {
				//   "time":"2025-02-19T13:56:56+03:00",
				//   "level":"ERROR",
				//   "source":{
				//   	"function":"github.com/tarantool/go-tlog_test.Test_Logger.funcN",
				//   	"file":"tlog_test.go",
				//   	"line":<line>
				//   },
				//   "msg":"my error message",
				//   "stacktrace":"github.com/tarantool/go-tlog_test.Test_Logger.funcN\n
				//                \ttlog_test.go:<line>\n
				//                testing.tRunner\n
				//                \truntime.goexit"
				// }
			},
			assert: func(require *require.Assertions, logs string) {
				require.Contains(logs, `"msg":"my error message"`)
				require.Contains(logs, `"stacktrace":"`)
				require.Contains(logs, "tlog_test.Test_Logger")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := require.New(t)

			tmpDir := t.TempDir()
			tc.opts.Path = filepath.Join(tmpDir, tc.opts.Path)

			l, err := tlog.New(tc.opts)
			r.NoError(err)

			defer func() {
				_ = l.Close()
			}()

			tc.log(l.Logger())

			logs, err := os.ReadFile(tc.opts.Path)
			r.NoError(err)

			// If there are stacktraces, there are no redundant internal frames.
			r.NotContains(string(logs), "slog.(*Logger).")

			tc.assert(r, string(logs))
		})
	}
}

func Test_LoggerReopen(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	dir := t.TempDir()

	filename := filepath.Join(dir, "Test_LoggerReopen.log")

	l, err := tlog.New(tlog.Opts{Path: filename})
	r.NoError(err)

	defer func() {
		_ = l.Close()
	}()

	l.Logger().Info("log_message1")

	rotated := filename + ".1"

	err = os.Rename(filename, rotated)
	r.NoError(err)

	err = l.Reopen()
	r.NoError(err)

	l.Logger().Info("log_message2")

	outOld, err := os.ReadFile(rotated)
	r.NoError(err)
	r.Contains(string(outOld), "log_message1")
	r.NotContains(string(outOld), "log_message2")

	out, err := os.ReadFile(filename)
	r.NoError(err)
	r.NotContains(string(out), "log_message1")
	r.Contains(string(out), "log_message2")
}

func Test_LoggerReopenOnSignal(t *testing.T) {
	// Do not run in parallel, sends signals to self process.
	r := require.New(t)

	dir := t.TempDir()

	filename := filepath.Join(dir, "Test_LoggerReopenOnSignal.log")

	l, err := tlog.New(tlog.Opts{Path: filename})
	r.NoError(err)

	defer func() {
		_ = l.Close()
	}()

	ctx, cancel := context.WithCancel(t.Context())
	go l.ReopenOnSignal(ctx, func(err error) {
		fmt.Fprintf(os.Stderr, "log reopen failed: %v\n", err)
	}, syscall.SIGUSR1)

	oldMessage := "log_message"
	l.Logger().Info(oldMessage)

	rotated := filename + ".1"

	err = os.Rename(filename, rotated)
	r.NoError(err)

	// We cannot guarantee that signal has been already processed, so let's
	// retry. The signal may also be lost if the handler goroutine has not
	// registered with signal.Notify yet, so we re-send it on each iteration
	// until the reopen takes effect.
	var newMessage string

	for i := range 10 {
		err = syscall.Kill(os.Getpid(), syscall.SIGUSR1)
		r.NoError(err)

		newMessage = fmt.Sprintf("log_newmessage%d", i)

		l.Logger().Info(newMessage)

		out, err := os.ReadFile(rotated)
		r.NoError(err)

		if !strings.Contains(string(out), newMessage) {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}

	cancel()

	outOld, err := os.ReadFile(rotated)
	r.NoError(err)
	r.Contains(string(outOld), oldMessage)
	r.NotContains(string(outOld), newMessage)

	out, err := os.ReadFile(filename)
	r.NoError(err)
	r.NotContains(string(out), oldMessage)
	r.Contains(string(out), newMessage)
}

func Test_LoggerConcurrency(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	dir := t.TempDir()

	filename := filepath.Join(dir, "Test_LoggerConcurrency.log")

	l, err := tlog.New(tlog.Opts{Path: filename})
	r.NoError(err)

	defer func() {
		_ = l.Close()
	}()

	const goroutines = 100

	var wg sync.WaitGroup

	for i := range goroutines {
		wg.Add(3)

		// Concurrent Logger().Info() calls.
		go func(n int) {
			defer wg.Done()

			l.Logger().Info("concurrent info message", slog.Int("n", n))
		}(i)

		// Concurrent Reopen() calls.
		go func() {
			defer wg.Done()

			_ = l.Reopen()
		}()

		// Concurrent Close() calls.
		go func() {
			defer wg.Done()

			_ = l.Close()
		}()
	}

	wg.Wait()
}
