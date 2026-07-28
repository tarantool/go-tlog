// Package tlog provides structured logging for Tarantool Go applications on
// top of the standard library's log/slog.
//
// It exists to give the ecosystem one consistent log format without forcing a
// third-party logger on consumers: everything it returns is an ordinary
// *slog.Logger or slog.Handler, so code that already speaks slog needs no
// changes.
//
// On top of slog it adds three things: a Tarantool-style text format, output
// fan-out to several destinations at once, and stacktraces attached
// automatically from a configurable level upwards.
//
// # Creating a logger
//
// [New] takes an [Opts] value and returns a [Logger] that owns its outputs.
// Call [Logger.Close] to flush and release them, and [Logger.Logger] to get
// the underlying *slog.Logger:
//
//	log, err := tlog.New(tlog.Opts{
//		Level:  tlog.LevelInfo,
//		Format: tlog.FormatText,
//		Path:   "stdout",
//	})
//	if err != nil {
//		return err
//	}
//	defer log.Close()
//
//	log.Logger().Info("service started", slog.Int("port", 8080))
//
// Write attributes structurally, as in the example above, rather than
// formatting them into the message string.
//
// # Levels and formats
//
// Levels are [LevelTrace], [LevelDebug], [LevelInfo], [LevelWarn] and
// [LevelError]; [LevelDefault] means [LevelInfo]. Formats are [FormatText]
// (the default) and [FormatJSON].
//
// # Outputs
//
// Opts.Path is a comma-separated list of destinations. The names "stdout" and
// "stderr" select the corresponding OS streams; anything else is treated as a
// file path and created if missing. Writing to several destinations at once is
// a matter of listing them:
//
//	Path: "stdout,/var/log/app.log"
//
// The default is "stderr".
//
// # Stacktraces
//
// By default a stacktrace is attached to every record at or above the
// configured Level. Opts.StacktraceLevel decouples the two thresholds, which
// is the usual choice for a service that logs at Info but only wants traces
// for errors:
//
//	tlog.Opts{
//		Level:           tlog.LevelInfo,
//		StacktraceLevel: tlog.LevelError,
//	}
//
// # Using the handlers directly
//
// When an application already builds its own *slog.Logger, [NewTextHandler]
// and [NewJSONHandler] provide the same formatting and stacktrace behavior
// for any io.Writer, with no output management attached:
//
//	handler := tlog.NewTextHandler(os.Stdout, &tlog.HandlerOptions{
//		Level: tlog.LevelDebug,
//	})
//	logger := slog.New(handler)
//
// # Library authors
//
// A library should not impose a logger on its callers. Accept an
// *slog.Logger as an option and fall back to slog.Default; the application
// decides whether that default is backed by tlog.
package tlog
