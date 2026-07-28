package tlog

import (
	"context"
	"log/slog"

	"github.com/tarantool/go-tlog/internal/stacktrace"
)

type stacktraceHandler struct {
	slog.Handler

	fromLevel slog.Level
}

func newStacktraceHandler(h slog.Handler, fromLevel slog.Level) stacktraceHandler {
	return stacktraceHandler{
		Handler:   h,
		fromLevel: fromLevel,
	}
}

// Strip stacktraceHandler.Handle, slog.(*Logger).log and
// slog.(*Logger).<Level>.
var internalsStripLevel = 3

func (h stacktraceHandler) Handle(ctx context.Context, record slog.Record) error {
	if record.Level >= h.fromLevel {
		record.Add("stacktrace", stacktrace.Get(internalsStripLevel))
	}

	//nolint:wrapcheck // This handler only decorates the record; the error
	// belongs to the wrapped handler and is passed through unchanged.
	return h.Handler.Handle(ctx, record)
}

// WithAttrs re-wraps the result of the underlying handler.
//
// Without this, the embedded slog.Handler would answer for us and return a
// bare handler, so slog.Logger.With would silently switch stacktraces off for
// the derived logger.
func (h stacktraceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return newStacktraceHandler(h.Handler.WithAttrs(attrs), h.fromLevel)
}

// WithGroup re-wraps the result of the underlying handler, for the same reason
// as [stacktraceHandler.WithAttrs].
func (h stacktraceHandler) WithGroup(name string) slog.Handler {
	return newStacktraceHandler(h.Handler.WithGroup(name), h.fromLevel)
}
