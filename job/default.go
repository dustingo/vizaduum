// using ps to check process
package job

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"vizaduum/logging"

	"github.com/go-kit/log/level"
)

// 默认ps检查
func defaultCheck(host string, procressName string, expectation int, wg *sync.WaitGroup) {
	defer func() {
		if err := recover(); err != nil {
			level.Error(logging.Logger).Log("err", err)
		}
		wg.Done()
	}()
	command := fmt.Sprintf("ps aux|egrep -v \"grep|SCREEN\" | grep \"%s\"|wc -l", procressName)
	cmd := exec.Command("sh", "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		level.Error(logging.Logger).Log("msg", "exec command error", "err", err.Error())
		return
	}
	realServiceCount, err := strconv.Atoi(strings.TrimSuffix(string(out), "\n"))
	if err != nil {
		level.Error(logging.Logger).Log("err", err)
		return
	}
	if diffenceValue := expectation - realServiceCount; diffenceValue == 0 {
		level.Info(logging.Logger).Log("msg", "service online", "service", procressName)
	} else if diffenceValue > 0 {
		defer lock.Unlock()
		lock.Lock()
		// 实际进程数量与预期不一致，有所减少，可能进程宕机
		AbormalResult = append(AbormalResult, host+" "+"进程:"+procressName+"[offline]", "count:"+strconv.Itoa(diffenceValue))
		level.Info(logging.Logger).Log("msg", "service offline", "service", procressName, "count", strconv.Itoa(diffenceValue))
	} else {
		defer lock.Unlock()
		lock.Lock()
		// 实际进程数量与预期不一致，有所增多，可能由于热更或者起多了进程导致
		AbormalResult = append(AbormalResult, host+" "+"进程:"+procressName+"[abormal]!", "count:"+strconv.Itoa(diffenceValue))
		level.Info(logging.Logger).Log("msg", "service abormal", "service", procressName, "count", strconv.Itoa(diffenceValue))
	}
}
