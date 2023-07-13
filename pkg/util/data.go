package util

/*--------game process config --------*/
type GameConfig struct {
	Program Program   `yaml:"program"`
	Service []Service `yaml:"service"`
}
type Program struct {
	App      string `yaml:"app"`
	Cluster  string `yaml:"cluster"`
	Name     string `yaml:"name"`
	Instance string `yaml:"instance"`
	Altype   string `yaml:"altype"`
	Team     string `yaml:"team"`
	Call     string `yaml:"call"`
	Channel  string `yaml:"channel"`
	Port     string `yaml:"port"`
	File     string `yaml:"file"`
	Interval int    `yaml:"interval"`
}
type Service struct {
	Method   string `yaml:"method"`
	Uniqname string `yaml:"uniqname"`
	Count    int    `yaml:"count"`
	Ports    []int  `yaml:"ports"`
}

/*--------sysconfig--------*/

type SYSConfig struct {
	Rabbitmq struct {
		User  string `json:"user"`
		Pass  string `json:"pass"`
		Url   string `json:"url"`
		Queue string `json:"queue"`
	} `yaml:"rabbitmq"`
	Alertmanager struct {
		Url string `json:"url"`
	} `yaml:"alertmanager"`
}

type ConsulConfig struct {
	GameConfig GameConfig `json:"gameConfig"`
	SYSConfig  SYSConfig  `json:"mqConfig"`
}
