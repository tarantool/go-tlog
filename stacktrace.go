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

	return h.Handler.Handle(ctx, record)
}
