package main

import (
	"log/slog"
	"os"

	"github.com/tarantool/go-tlog"
)

func main() {
	_ = os.MkdirAll("/tmp/tlog_demo", 0755)

	l, err := tlog.New(tlog.Opts{
		Level:  tlog.LevelInfo,
		Format: tlog.FormatText,
		Path:   "/tmp/tlog_demo/app.log",
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	log := l.Logger().With(slog.String("mode", "file"))
	log.Info("logging to file", "path", "/tmp/tlog_demo/app.log")
	log.Warn("network delay", "ms", 250)
	log.Error("write failed", "err", "disk quota exceeded")
}
