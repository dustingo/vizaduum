package alert

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"vizaduum/config"
	"vizaduum/logging"

	"github.com/go-kit/log/level"
	"github.com/rabbitmq/amqp091-go"
)

func SendToRabbitMq(message string) {
	anno := Annotations{
		Description: message,
		Call:        config.GConfig.GameConfig.Program.Call,
		Summary:     fmt.Sprintf("%s 进程报警", config.GConfig.GameConfig.Program.Name),
	}
	labels := Labels{
		Alertname: "进程报警",
		Altype:    config.GConfig.GameConfig.Program.Altype,
		App:       config.GConfig.GameConfig.Program.App,
		Cluster:   config.GConfig.GameConfig.Program.Cluster,
		Name:      config.GConfig.GameConfig.Program.Name,
		Instance:  config.GConfig.GameConfig.Program.Instance,
		Team:      config.GConfig.GameConfig.Program.Team,
	}
	msg := Alerts{
		Status:      "firing",
		Labels:      labels,
		Annotations: anno,
		StartsAt:    time.Now(),
	}
	conn, err := amqp091.Dial(fmt.Sprintf("amqp://%s:%s@%s/", config.GConfig.SYSConfig.Rabbitmq.User,
		config.GConfig.SYSConfig.Rabbitmq.Pass,
		config.GConfig.SYSConfig.Rabbitmq.Url))
	if err != nil {
		level.Error(logging.Logger).Log("msg", "failed to connect rabbitmq", "err", err)
		return
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		level.Error(logging.Logger).Log("msg", "failed to open a rabbit channel")
		return
	}
	defer ch.Close()
	alertInfo, _ := json.Marshal(&msg)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = ch.PublishWithContext(
		ctx,
		"",
		config.GConfig.SYSConfig.Rabbitmq.Queue,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        alertInfo,
		},
	)
	if err != nil {
		level.Error(logging.Logger).Log("msg", "failed to to publish alert to rabbitmq")
		return
	}
	level.Debug(logging.Logger).Log("msg", "send to rabbitmq succeed", "alertmessage", string(alertInfo))
}

func SendToAlertmanager(message string) {
	labels := Label{
		Alertname: "进程报警",
		Name:      config.GConfig.GameConfig.Program.Name,
		Cluster:   config.GConfig.GameConfig.Program.Cluster,
		Instance:  config.GConfig.GameConfig.Program.Instance,
		Altype:    config.GConfig.GameConfig.Program.Altype,
		Team:      config.GConfig.GameConfig.Program.Team,
		App:       config.GConfig.GameConfig.Program.App,
	}
	anno := Anno{
		Summary:     fmt.Sprintf("%s 进程报警", config.GConfig.GameConfig.Program.Name),
		Call:        config.GConfig.GameConfig.Program.Call,
		Description: message,
	}
	a := &NotifyAlert{
		StartsAt:    time.Now(),
		Labels:      labels,
		Annotations: anno,
	}
	adata := []NotifyAlert{}
	adata = append(adata, *a)
	client := http.Client{
		Timeout: 5 * time.Second, // FIXME: timeout
	}
	data, err := json.Marshal(adata)
	if err != nil {
		level.Error(logging.Logger).Log("err", err.Error())
		return
	}
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/alerts", config.GConfig.SYSConfig.Alertmanager.Url), bytes.NewReader(data))
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		level.Error(logging.Logger).Log("msg", "send to alertmanager error", "err", err.Error())
		return
	}
	if resp.StatusCode == 200 {
		level.Info(logging.Logger).Log("msg", "send to alertmanager succeed")
	} else {
		d, _ := ioutil.ReadAll(resp.Body)
		level.Error(logging.Logger).Log("msg", "send to alertmanager failed", "err", string(d))
	}
}
