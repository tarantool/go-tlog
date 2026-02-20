package tlog

import (
	"io"
	"log/slog"

	slogcustom "github.com/tarantool/go-tlog/internal/slog"
)

// HandlerOptions are options for NewJSONHandler and NewTextHandler.
type HandlerOptions struct {
	// Level sets minimum level for the handler.
	Level Level
	// StacktraceLevel overrides the default stacktrace threshold.
	// If set, stacktraces will be attached starting from this level,
	// regardless of the main log Level.
	StacktraceLevel Level
	// ReplaceAttr is called to rewrite each non-group Attr before
	// it is logged. See slog.HandlerOptions for details.
	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr
	// AddSource causes the handler to compute the source code position
	// of the log statement and add a SourceKey attribute to the output.
	AddSource bool
}

func (o *HandlerOptions) slogHandlerOptions() slog.HandlerOptions {
	opts := slog.HandlerOptions{
		AddSource: o.AddSource,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			a = replaceAttr(groups, a)
			if o.ReplaceAttr != nil {
				a = o.ReplaceAttr(groups, a)
			}
			return a
		},
	}

	switch o.Level {
	case LevelTrace, LevelDebug:
		opts.Level = slog.LevelDebug
	case LevelDefault, LevelInfo:
		opts.Level = slog.LevelInfo
	case LevelWarn:
		opts.Level = slog.LevelWarn
	case LevelError:
		opts.Level = slog.LevelError
	}

	return opts
}

func (o *HandlerOptions) slogTraceLevel() slog.Level {
	level := o.Level
	if o.StacktraceLevel != 0 {
		level = o.StacktraceLevel
	}

	switch level {
	case LevelTrace, LevelDebug:
		return slog.LevelDebug
	case LevelDefault, LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

// NewJSONHandler creates a slog.Handler that writes JSON output to w.
// If opts is nil, default options are used.
func NewJSONHandler(w io.Writer, opts *HandlerOptions) slog.Handler {
	if opts == nil {
		opts = &HandlerOptions{}
	}

	handlerOpts := opts.slogHandlerOptions()
	base := slog.NewJSONHandler(w, &handlerOpts)

	return newStacktraceHandler(base, opts.slogTraceLevel())
}

// NewTextHandler creates a slog.Handler that writes text output to w.
// If opts is nil, default options are used.
func NewTextHandler(w io.Writer, opts *HandlerOptions) slog.Handler {
	if opts == nil {
		opts = &HandlerOptions{}
	}

	handlerOpts := opts.slogHandlerOptions()
	base := slogcustom.NewTextHandler(w, &slogcustom.HandlerOptions{
		HandlerOptions:  handlerOpts,
		OmitBuiltinKeys: true,
	})

	return newStacktraceHandler(base, opts.slogTraceLevel())
}
