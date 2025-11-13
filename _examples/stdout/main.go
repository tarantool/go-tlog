package main

import (
	"log/slog"

	"github.com/tarantool/go-tlog"
)

func main() {
	l, err := tlog.New(tlog.Opts{
		Level:  tlog.LevelInfo,
		Format: tlog.FormatText,
		Path:   "stdout",
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	log := l.Logger().With(slog.String("mode", "stdout"))
	log.Info("service started")
	log.Warn("cache warming", "duration", "1.3s")
	log.Error("failed to connect", "host", "db1")
}
