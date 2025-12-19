package slog_test

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	slogcustom "github.com/tarantool/go-tlog/internal/slog"
)

func Test_TextHandler_OmitBuiltinKeys(t *testing.T) {
	require := require.New(t)

	var b bytes.Buffer

	handler := slogcustom.NewTextHandler(&b, &slogcustom.HandlerOptions{
		HandlerOptions:  slog.HandlerOptions{AddSource: true},
		OmitBuiltinKeys: true,
	})
	l := slog.New(handler)

	l.Info("my message")

	out := strings.TrimSpace(b.String())

	require.NotContains(out, "time=")
	require.Regexp(`^20\d\d\-\d\d\-\d\dT\d\d:\d\d:\d\d.*$`, out)

	require.NotContains(out, "level=")
	require.Contains(out, "INFO")

	require.NotContains(out, "source=")
	require.Contains(out, "internal/slog/text_handler_test.go:25")

	require.NotContains(out, "msg=")
	require.Contains(out, "my message")
}
