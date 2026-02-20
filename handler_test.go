package tlog_test

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tarantool/go-tlog"
)

func Test_NewJSONHandler(t *testing.T) {
	t.Parallel()

	t.Run("writes JSON to io.Writer", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		handler := tlog.NewJSONHandler(&buf, &tlog.HandlerOptions{
			Level:           tlog.LevelDebug,
			StacktraceLevel: tlog.LevelError,
		})
		logger := slog.New(handler)

		logger.Info("hello json")

		r := require.New(t)
		r.Contains(buf.String(), `"msg":"hello json"`)
		r.Contains(buf.String(), `"level":"INFO"`)
		r.NotContains(buf.String(), `"stacktrace"`)
	})

	t.Run("includes stacktrace at error", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		handler := tlog.NewJSONHandler(&buf, &tlog.HandlerOptions{
			Level: tlog.LevelDebug,
		})
		logger := slog.New(handler)

		logger.Error("fail")

		r := require.New(t)
		r.Contains(buf.String(), `"msg":"fail"`)
		r.Contains(buf.String(), `"stacktrace"`)
	})

	t.Run("nil opts uses defaults", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		handler := tlog.NewJSONHandler(&buf, nil)
		logger := slog.New(handler)

		logger.Info("default opts")

		r := require.New(t)
		r.Contains(buf.String(), `"msg":"default opts"`)
	})

	t.Run("respects level filtering", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		handler := tlog.NewJSONHandler(&buf, &tlog.HandlerOptions{
			Level: tlog.LevelError,
		})
		logger := slog.New(handler)

		logger.Info("should be dropped")

		r := require.New(t)
		r.Empty(buf.String())
	})

	t.Run("ReplaceAttr is applied", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		handler := tlog.NewJSONHandler(&buf, &tlog.HandlerOptions{
			Level:           tlog.LevelInfo,
			StacktraceLevel: tlog.LevelError,
			ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
				if a.Key == slog.MessageKey {
					a.Value = slog.StringValue("replaced")
				}
				return a
			},
		})
		logger := slog.New(handler)

		logger.Info("original")

		r := require.New(t)
		r.Contains(buf.String(), `"msg":"replaced"`)
		r.NotContains(buf.String(), `"msg":"original"`)
	})
}

func Test_NewTextHandler(t *testing.T) {
	t.Parallel()

	t.Run("writes text to io.Writer", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		handler := tlog.NewTextHandler(&buf, &tlog.HandlerOptions{
			Level:           tlog.LevelDebug,
			StacktraceLevel: tlog.LevelError,
		})
		logger := slog.New(handler)

		logger.Info("hello text")

		r := require.New(t)
		r.Contains(buf.String(), "hello text")
		r.Contains(buf.String(), "INFO")
		r.NotContains(buf.String(), "stacktrace=")
	})

	t.Run("includes stacktrace at error", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		handler := tlog.NewTextHandler(&buf, &tlog.HandlerOptions{
			Level: tlog.LevelDebug,
		})
		logger := slog.New(handler)

		logger.Error("fail text")

		r := require.New(t)
		r.Contains(buf.String(), "fail text")
		r.Contains(buf.String(), "stacktrace=")
	})

	t.Run("nil opts uses defaults", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		handler := tlog.NewTextHandler(&buf, nil)
		logger := slog.New(handler)

		logger.Info("default opts")

		r := require.New(t)
		r.Contains(buf.String(), "default opts")
	})

	t.Run("respects level filtering", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		handler := tlog.NewTextHandler(&buf, &tlog.HandlerOptions{
			Level: tlog.LevelError,
		})
		logger := slog.New(handler)

		logger.Info("should be dropped")

		r := require.New(t)
		r.Empty(buf.String())
	})
}
