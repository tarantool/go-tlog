package tlog

// Format represents logger format.
type Format int

const (
	// FormatDefault is the default format. Logger uses FormatText as a default one.
	FormatDefault Format = iota
	// FormatText prints messages as a human-readable text string.
	FormatText Format = iota
	// FormatJSON prints each message as a JSON object.
	FormatJSON
)
