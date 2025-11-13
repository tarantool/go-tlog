package tlog

// Level represents logger level.
type Level int

const (
	// LevelDefault is the default level. Logger uses LevelInfo as a default one.
	LevelDefault Level = iota
	// LevelTrace prints messages up to Debug. Messages up to Debug have stacktraces.
	LevelTrace
	// LevelDebug prints messages up to Debug. Messages up to Error have stacktraces.
	LevelDebug
	// LevelInfo prints messages up to Info. Messages up to Error have stacktraces.
	LevelInfo
	// LevelWarn prints messages up to Warn. Messages up to Error have stacktraces.
	LevelWarn
	// LevelError prints messages up to Error. Messages up to Error have stacktraces.
	LevelError
)
