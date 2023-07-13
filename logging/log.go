package logging

import (
	"os"

	"github.com/go-kit/log"
)

var (
	Logger log.Logger
)

func Setup() {
	Logger = log.NewLogfmtLogger(os.Stdout)
	Logger = log.With(Logger, "ts", log.DefaultTimestamp)
	Logger = log.With(Logger, "caller", log.DefaultCaller)
}
