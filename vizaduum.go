package main

import (
	"flag"
	"fmt"
	"vizaduum/config"
	"vizaduum/flags"
	"vizaduum/job"
	"vizaduum/logging"
	"vizaduum/pkg/register"
	"vizaduum/router"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/log/level"
	"github.com/robfig/cron/v3"
)

func init() {
	logging.Setup()
}

func main() {
	flag.Parse()
	config.Setup(flags.ConsulAddr, flags.ConsulToken, flags.GameKey, flags.SysKey)
	job.FileExist()
	go config.GameWatcher()
	go register.ServiceRegister()
	c := cron.New(cron.WithSeconds(), cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))
	_, err := c.AddFunc(fmt.Sprintf("@every %dm", config.GConfig.GameConfig.Program.Interval), job.RunCheck)
	if err != nil {
		level.Error(logging.Logger).Log("err", err.Error())
		return
	}
	c.Start()
	gin.SetMode(gin.ReleaseMode)
	gin.ForceConsoleColor()
	router := router.Setup()
	router.Run(config.GConfig.GameConfig.Program.Port)
}
