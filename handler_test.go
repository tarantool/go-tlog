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

// Test_Handler_StacktraceSurvivesDerivation is a regression test: the
// stacktrace handler embeds slog.Handler, so without explicit WithAttrs and
// WithGroup methods the embedded handler answered for it and every logger
// derived with With or WithGroup lost its stacktraces.
func Test_Handler_StacktraceSurvivesDerivation(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name   string
		derive func(*slog.Logger) *slog.Logger
	}{
		{
			name:   "With",
			derive: func(l *slog.Logger) *slog.Logger { return l.With(slog.String("k", "v")) },
		},
		{
			name:   "WithGroup",
			derive: func(l *slog.Logger) *slog.Logger { return l.WithGroup("g") },
		},
		{
			name: "With after WithGroup",
			derive: func(l *slog.Logger) *slog.Logger {
				return l.WithGroup("g").With(slog.String("k", "v"))
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer

			handler := tlog.NewJSONHandler(&buf, &tlog.HandlerOptions{
				Level:           tlog.LevelDebug,
				StacktraceLevel: tlog.LevelError,
			})

			tc.derive(slog.New(handler)).Error("fail")

			require.Contains(t, buf.String(), `"stacktrace"`)
		})
	}
}

// Test_Handler_ReplaceAttrCanDropAttrs is a regression test: dropping an
// attribute by returning the zero slog.Attr from ReplaceAttr — the way the
// standard library documents — used to panic in the text handler.
func Test_Handler_ReplaceAttrCanDropAttrs(t *testing.T) {
	t.Parallel()

	drop := func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey || a.Key == "secret" {
			return slog.Attr{}
		}

		return a
	}

	for _, tc := range []struct {
		name    string
		timeKey string
		new     func(*bytes.Buffer, *tlog.HandlerOptions) slog.Handler
	}{
		{name: "text", timeKey: "time=", new: func(b *bytes.Buffer, o *tlog.HandlerOptions) slog.Handler {
			return tlog.NewTextHandler(b, o)
		}},
		{name: "json", timeKey: `"time":`, new: func(b *bytes.Buffer, o *tlog.HandlerOptions) slog.Handler {
			return tlog.NewJSONHandler(b, o)
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer

			handler := tc.new(&buf, &tlog.HandlerOptions{
				Level:       tlog.LevelInfo,
				ReplaceAttr: drop,
			})

			require.NotPanics(t, func() {
				slog.New(handler).Info("hello", slog.String("secret", "s3cr3t"), slog.Int("port", 8080))
			})

			r := require.New(t)
			r.Contains(buf.String(), "hello")
			r.Contains(buf.String(), "8080")
			r.NotContains(buf.String(), "s3cr3t")
			r.NotContains(buf.String(), tc.timeKey)
		})
	}
}
