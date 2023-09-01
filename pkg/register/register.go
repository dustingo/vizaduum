package register

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"vizaduum/config"
	"vizaduum/flags"
	"vizaduum/logging"

	"github.com/go-kit/log/level"
	capi "github.com/hashicorp/consul/api"
)

func ServiceRegister() {
	client, _ := capi.NewClient(&capi.Config{
		Address: flags.ConsulAddr,
		Token:   flags.ConsulToken,
	})
	p, _ := strconv.Atoi(strings.Split(config.GConfig.GameConfig.Program.Port, ":")[1])
	regService := capi.AgentServiceRegistration{
		ID:      md5sum(config.GConfig.GameConfig.Program.Name),
		Name:    config.GConfig.GameConfig.Program.Name,
		Address: config.GConfig.GameConfig.Program.Instance,
		Port:    p,
		Check: &capi.AgentServiceCheck{
			CheckID:  md5sum(config.GConfig.GameConfig.Program.Instance),
			HTTP:     fmt.Sprintf("http://%s:%d/api/v1/health", config.GConfig.GameConfig.Program.Instance, p),
			Timeout:  "1s",
			Interval: "10s",
		},
	}
	err := client.Agent().ServiceRegister(&regService)
	if err != nil {
		level.Error(logging.Logger).Log("msg", "register to consul failed", "err", err.Error())
		os.Exit(1)
	}
}
func md5sum(id string) string {
	m := md5.New()
	m.Write([]byte(id))
	cipherString := m.Sum(nil)
	return hex.EncodeToString(cipherString)
}
