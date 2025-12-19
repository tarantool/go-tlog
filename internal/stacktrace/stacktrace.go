package stacktrace

import (
	"runtime"
	"strconv"
	"strings"
)

// Get returns a formatted stacktrace starting from the specified number of frames to skip.
func Get(skip int) string {
	frames := getFrames(skip)

	var b strings.Builder

	for {
		frame, more := frames.Next()

		writeFrame(&b, frame)

		if !more {
			break
		}

		b.WriteByte('\n')
	}

	return b.String()
}

const (
	defaultProgramCounters = 64

	// runtime.Counters, Get and getFrames.
	baseNestingLevel = 3

	pcsExtendFactor = 2
)

func getFrames(skip int) *runtime.Frames {
	pcs := make([]uintptr, defaultProgramCounters)

	for {
		n := runtime.Callers(baseNestingLevel+skip, pcs)
		if n < cap(pcs) {
			pcs = pcs[:n]

			break
		}

		pcs = make([]uintptr, len(pcs)*pcsExtendFactor)
	}

	return runtime.CallersFrames(pcs)
}

func writeFrame(b *strings.Builder, frame runtime.Frame) {
	b.WriteString(frame.Function)
	b.WriteByte('\n')
	b.WriteByte('\t')
	b.WriteString(frame.File)
	b.WriteByte(':')
	b.WriteString(strconv.Itoa(frame.Line))
}
