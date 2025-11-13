package tlog_test

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

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
				Level:  tlog.LevelDebug,
				Format: tlog.FormatText,
				Path:   "InfoMessage_DebugPlainLogger.log",
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
				Level:  tlog.LevelDebug,
				Format: tlog.FormatJSON,
				Path:   "InfoMessage_DebugPlainLogger.json",
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
		tc := tc

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
