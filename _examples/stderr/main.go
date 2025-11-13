package main

import (
	"log/slog"

	"github.com/tarantool/go-tlog"
)

func main() {
	l, err := tlog.New(tlog.Opts{
		Level:  tlog.LevelDebug,
		Format: tlog.FormatJSON,
		Path:   "stderr",
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	log := l.Logger().With(slog.String("mode", "stderr"))
	log.Debug("debug message", "module", "init")
	log.Info("component loaded", "component", "api")
	log.Error("unexpected nil pointer")
}
