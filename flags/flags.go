package flags

import (
	"flag"
)

var (
	ConsulAddr  string
	ConsulToken string
	GameKey     string
	SysKey      string
)

func init() {
	flag.StringVar(&ConsulAddr, "consul-address", "127.0.0.1:8500", "consul  address")
	flag.StringVar(&ConsulToken, "consul-token", "", "consul secretID")
	flag.StringVar(&GameKey, "game-key", "", "key of consul's game config")
	flag.StringVar(&SysKey, "sys-key", "", "key of consul's sys config")
}
