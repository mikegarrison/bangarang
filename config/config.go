package config

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/eliothedeman/bangarang/alarm"
)

type Configer interface {
	ConfigStruct() interface{}
	Init(interface{}) error
}

type AppConfig struct {
	Escalations      []*alarm.Escalation    `json:"-"`
	EscalationsDir   string                 `json:"escalations_dir"`
	KeepAliveAge     time.Duration          `json:"-"`
	Raw_KeepAliveAge string                 `json:"keep_alive_age"`
	TcpPort          *int                   `json:"tcp_port"`
	HttpPort         *int                   `json:"http_port"`
	Alarms           *alarm.AlarmCollection `json:"alarms"`
}

func LoadConfigFile(fileName string) (*AppConfig, error) {
	buff, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	ac := &AppConfig{
		Raw_KeepAliveAge: "25m",
	}
	ac.Alarms = &alarm.AlarmCollection{}

	err = json.Unmarshal(buff, ac)
	if err != nil {
		return nil, err
	}

	ac.KeepAliveAge, err = time.ParseDuration(ac.Raw_KeepAliveAge)
	if err != nil {
		return ac, err
	}

	paths, err := filepath.Glob(ac.EscalationsDir + "*.json")
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		e, err := loadEscalation(path)
		if err != nil {
			return ac, err
		}

		ac.Escalations = append(ac.Escalations, e)
	}

	return ac, nil
}

func loadEscalation(fileName string) (*alarm.Escalation, error) {
	if !filepath.IsAbs(fileName) {
		fileName, _ = filepath.Abs(fileName)
	}
	buff, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	e := &alarm.Escalation{}
	err = json.Unmarshal(buff, e)
	if err != nil {
		return e, err
	}
	e.Policy.Compile()

	err = e.LoadAlarms()
	return e, err
}
