package config

import (
	"fmt"
	"time"
	"vizaduum/flags"
	"vizaduum/logging"
	"vizaduum/pkg/util"

	"github.com/go-kit/log/level"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"gopkg.in/yaml.v3"
)

func GameWatcher() {
	var gameConfig util.GameConfig
	param := map[string]interface{}{
		"type": "key",
		"key":  flags.GameKey,
	}
	watcher, err := watch.Parse(param)
	if err != nil {
		level.Error(logging.Logger).Log("msg", "failed to parse watcher", "err", err.Error())
		return
	}
	watcher.Token = flags.ConsulToken
	watcher.Handler = func(u uint64, i interface{}) {
		kv, ok := i.(*api.KVPair)
		if !ok {
			level.Error(logging.Logger).Log("msg", "invalid data type received")
			return
		}
		yaml.Unmarshal(kv.Value, &gameConfig)
		GConfig.GameConfig = gameConfig
	}
	go func() {
		err := watcher.Run(flags.ConsulAddr)
		if err != nil {
			level.Error(logging.Logger).Log("msg", "failed to run watcher", "err", err.Error())
			return
		}
	}()
	for {
		time.Sleep(30 * time.Second)
		level.Info(logging.Logger).Log("msg", "sync consul config every 30 seconds")
		level.Debug(logging.Logger).Log("msg", fmt.Sprint(GConfig.GameConfig))
	}
}
