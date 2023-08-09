package job

import (
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"vizaduum/config"
	"vizaduum/logging"
	"vizaduum/pkg/alert"

	"github.com/go-kit/log/level"
)

var (
	lock          sync.RWMutex
	AbormalResult []string
)

const p = "pause"

var (
	Status int //0: start 1:pause
	wg     sync.WaitGroup
)

func RunCheck() {
	defer func() {
		AbormalResult = []string{}
	}()

	level.Info(logging.Logger).Log("msg", "run vizaduum task")
	if Status == 1 {
		level.Info(logging.Logger).Log("msg", "维护中,不报警")
		return
	}
	for _, service := range config.GConfig.GameConfig.Service {
		wg.Add(1)
		if strings.ToLower(service.Method) == "port" {
			portCheck(config.GConfig.GameConfig.Program.Instance, service.Ports, service.Uniqname, &wg)
		} else if strings.ToLower(service.Method) == "default" {
			defaultCheck(config.GConfig.GameConfig.Program.Instance, service.Uniqname, service.Count, &wg)
		} else {
			go func(processName string) {
				defer wg.Done()
				level.Error(logging.Logger).Log("msg", "method error", "uniqname", processName)
			}(service.Uniqname)
		}
	}
	wg.Wait()
	dataByte, err := ioutil.ReadFile(config.GConfig.GameConfig.Program.File)
	if err != nil {
		level.Error(logging.Logger).Log("err", err)
		return
	}
	// if strings.Contains(string(dataByte), p) {
	// 	level.Warn(logging.Logger).Log("msg", "维护时间段，不报警")
	// 	return
	// }
	// 当没有报警的时候每次都把报警文件清空，以防干扰
	if len(AbormalResult) == 0 {
		os.Truncate(config.GConfig.GameConfig.Program.File, 0)
		return
	}
	// 当有报警时且报警记录文件内容为空时，写入最新的报警到文件中并报警，否则要现比较报警文件和报警信息
	if string(dataByte) == "" {
		write(config.GConfig.GameConfig.Program.File, strings.Join(AbormalResult, " "))
		if config.GConfig.GameConfig.Program.Channel == "rabbitmq" {
			alert.SendToRabbitMq(strings.Join(AbormalResult, " "))
		} else if config.GConfig.GameConfig.Program.Channel == "alertmanager" {
			alert.SendToAlertmanager(strings.Join(AbormalResult, " "))
		} else {
			return
		}
	} else {
		data := strings.TrimSuffix(string(dataByte), "\n")
		if ok := compare(strings.Join(AbormalResult, " "), data); !ok {
			write(config.GConfig.GameConfig.Program.File, strings.Join(AbormalResult, " "))
			if config.GConfig.GameConfig.Program.Channel == "rabbitmq" {
				alert.SendToRabbitMq(strings.Join(AbormalResult, " "))
			} else if config.GConfig.GameConfig.Program.Channel == "alertmanager" {
				alert.SendToAlertmanager(strings.Join(AbormalResult, " "))
			} else {
				return
			}
		} else {
			level.Warn(logging.Logger).Log("msg", strings.Join(AbormalResult, " "), "notice", "宕机资源重复，不进行报警")
		}
	}
}

//清空并写入记录
func write(fname string, abormalResult string) {
	lock.Lock()
	// p, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	// path := p + "/" + fname
	os.Truncate(fname, 0)
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		level.Error(logging.Logger).Log("err", err)
		return
	}
	defer lock.Unlock()
	defer f.Close()
	f.WriteString(abormalResult)
}

// 比较文件内容和内存字符串, true代表"相等",不许报警，false需要报警
func compare(strNew, strOld string) bool {
	var newSlice []string
	//真实进程信息
	sliceNew := strings.Split(strNew, " ")
	//配置文件进程信息
	sliceOld := strings.Split(strOld, " ")

	for _, value := range sliceNew {
		if strings.Contains(strOld, value) {
			newSlice = append(newSlice, value)
		} else {
			return false
		}
	}
	if len(newSlice) == len(sliceOld) {
		return true
	} else {
		return false
	}
}
