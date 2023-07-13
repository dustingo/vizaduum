package job

import (
	"os"
	"vizaduum/config"
	"vizaduum/logging"

	"github.com/go-kit/log/level"
)

func FileExist() {
	_, err := os.Stat(config.GConfig.GameConfig.Program.File)
	if err == nil {
		level.Debug(logging.Logger).Log("msg", "file exist", "file", config.GConfig.GameConfig.Program.File)
		return
	}
	if os.IsNotExist((err)) {
		level.Debug(logging.Logger).Log("msg", "file not exist,then create it", "file", config.GConfig.GameConfig.Program.File)
		f, err := os.OpenFile(config.GConfig.GameConfig.Program.File, os.O_CREATE|os.O_RDWR, 0755)
		if err != nil {
			level.Error(logging.Logger).Log("msg", "create file failed", "file", config.GConfig.GameConfig.Program.File)
			panic(err)
		}
		defer f.Close()
	}
}
