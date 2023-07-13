package config

import (
	"os"
	"vizaduum/logging"
	"vizaduum/pkg/util"

	"github.com/go-kit/log/level"
	capi "github.com/hashicorp/consul/api"
	"gopkg.in/yaml.v3"
)

var (
	GConfig util.ConsulConfig
)

// 初始化rabbitmq配置和游戏进程列表配置
func Setup(address string, token string, gkv, mkv string) {
	defer func() {
		if err := recover(); err != nil {
			level.Error(logging.Logger).Log("msg", "load consul game config failed", "err", err)
		}
	}()
	var gameConfig util.GameConfig
	var sysConfig util.SYSConfig
	var err error
	client, err := capi.NewClient(&capi.Config{
		Address: address,
		Token:   token,
	})
	if err != nil {
		level.Error(logging.Logger).Log("msg", "connect to consul failed", "err", err.Error())
		os.Exit(1)
	}
	mqPair, _, err := client.KV().Get(mkv, nil)
	if err != nil {
		level.Error(logging.Logger).Log("msg", "failed to get kv value", "err", err.Error(), "mqKey", mkv)
		return
	}
	err = yaml.Unmarshal(mqPair.Value, &sysConfig)
	if err != nil {
		level.Error(logging.Logger).Log("msg", "unmarshal config failed", "err", err.Error())
		return
	}
	gamesPair, _, err := client.KV().Get(gkv, nil)
	if err != nil {
		level.Error(logging.Logger).Log("msg", "failed to get kv value", "err", err.Error(), "gameKey", gkv)
		return
	}
	err = yaml.Unmarshal(gamesPair.Value, &gameConfig)
	if err != nil {
		level.Error(logging.Logger).Log("msg", "unmarshal config failed", "err", err.Error())
		return
	}
	GConfig = util.ConsulConfig{
		GameConfig: gameConfig,
		SYSConfig:  sysConfig,
	}
	level.Info(logging.Logger).Log("msg", "setup consul gameconfig")
}
