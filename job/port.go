// using port to check process
package job

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
	"vizaduum/logging"

	"github.com/go-kit/log/level"
)

// 端口检查
func portCheck(host string, ports []int, procressName string, wg *sync.WaitGroup) {
	defer func() {
		if err := recover(); err != nil {
			level.Error(logging.Logger).Log("err", err)
		}
		wg.Done()
	}()
	for _, port := range ports {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 5*time.Second)
		if err != nil {
			level.Error(logging.Logger).Log("msg", "端口无法链接，可能已宕机", "port", port, "service", procressName)
			AbormalResult = append(AbormalResult, host+" "+"port:"+strconv.Itoa(port)+"[offline]")
			continue
		}
		level.Info(logging.Logger).Log("msg", "端口链接正常", "port", port, "service", procressName)
		defer conn.Close()
	}
}
