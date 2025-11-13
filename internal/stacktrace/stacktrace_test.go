package stacktrace_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tarantool/go-tlog/internal/stacktrace"
)

func funcNested() string {
	// Skip funcNested, include funcWrapper.
	return stacktrace.Get(1)
}

func funcWrapper() string {
	return funcNested()
}

func Test_Get(t *testing.T) {
	require := require.New(t)

	stack := funcWrapper()

	require.Contains(stack, "internal/stacktrace/stacktrace_test.go:17")
	require.Contains(stack, "funcWrapper")

	require.NotContains(stack, "internal/stacktrace/stacktrace_test.go:12")
	require.NotContains(stack, "funcNested")
}
