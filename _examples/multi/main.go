package main

import (
	"log/slog"
	"os"

	"github.com/tarantool/go-tlog"
)

func main() {
	_ = os.MkdirAll("/tmp/tlog_multi", 0755)

	log, err := tlog.New(tlog.Opts{
		Level:  tlog.LevelInfo,
		Format: tlog.FormatText,
		Path:   "stdout,/tmp/tlog_multi/app.log",
	})
	if err != nil {
		panic(err)
	}
	defer log.Close()

	logger := log.Logger().With(slog.String("example", "multi-output"))
	logger.Info("message written to stdout and file")
}
