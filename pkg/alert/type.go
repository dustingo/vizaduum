package alert

import "time"

/*----------rbbitmq-----------*/
type Annotations struct {
	Description string `json:"description"`
	Summary     string `json:"summary,omitempty"`
	Call        string `json:"call,omitempty"` // 是否需要电话通知
}
type Labels struct {
	Alertname string `json:"alertname"`
	Altype    string `json:"altype"`
	App       string `json:"app,omitempty"`
	Cluster   string `json:"cluster,omitempty"`
	Instance  string `json:"instance,omitempty"`
	Job       string `json:"job,omitempty"`
	Name      string `json:"name,omitempty"`
	Region    string `json:"region,omitempty"`
	Team      string `json:"team"`
}

type Alerts struct {
	Status       string      `json:"status"`
	Labels       Labels      `json:"labels"`
	Annotations  Annotations `json:"annotations"`
	StartsAt     time.Time   `json:"startsAt,omitempty"`
	EndsAt       time.Time   `json:"endsAt,omitempty"`
	GeneratorURL string      `json:"generatorURL,omitempty"`
	Fingerprint  string      `json:"fingerprint,omitempty"`
}

type GroupLabels struct {
	Alertname string `json:"alertname"`
}
type CommonLabels struct {
	Alertname string `json:"alertname"`
	Altype    string `json:"altype"`
	App       string `json:"app,omitempty"`
	Cluster   string `json:"cluster,omitempty"`
	Instance  string `json:"instance,omitempty"`
	Job       string `json:"job,omitempty"`
	Name      string `json:"name,omitempty"`
	Region    string `json:"region,omitempty"`
	Team      string `json:"team"`
}
type CommonAnnotations struct {
	Description string `json:"description"`
	Summary     string `json:"summary,omitempty"`
	Call        string `json:"call,omitempty"` // 是否需要电话通知
}

type AlertmanagerData struct {
	Receiver          string            `json:"receiver"`
	Status            string            `json:"status"`
	Alerts            []Alerts          `json:"alerts"`
	GroupLabels       GroupLabels       `json:"groupLabels"`
	CommonLabels      CommonLabels      `json:"commonLabels"`
	CommonAnnotations CommonAnnotations `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
	TruncatedAlerts   int               `json:"truncatedAlerts"`
}

/*----alertmanager-----*/
type NotifyAlert struct {
	// Label value pairs for purpose of aggregation, matching, and disposition
	// dispatching. This must minimally include an "alertname" label.
	Labels Label `json:"labels"`

	// Extra key/value information which does not define alert identity.
	Annotations Anno `json:"annotations"`

	// The known time range for this alert. Both ends are optional.
	StartsAt     time.Time `json:"startsAt,omitempty"`
	EndsAt       time.Time `json:"endsAt,omitempty"`
	GeneratorURL string    `json:"generatorURL,omitempty"`
}

// Label is a key/value pair of strings.
type Label struct {
	Alertname string `json:"alertname"`
	Name      string `json:"name"`
	Cluster   string `json:"cluster"`
	Instance  string `json:"instance"`
	Altype    string `json:"altype"`
	Team      string `json:"team"`
	App       string `json:"app"`
}
type Anno struct {
	Summary     string `json:"summary"`
	Call        string `json:"call"`
	Description string `json:"description"`
}
