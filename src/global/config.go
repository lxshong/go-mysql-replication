package global

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var _config *config

type config struct {
	Target string `yaml:"target"` // 目标类型，支持redis、mongodb

	Addr     string `yaml:"addr"`
	User     string `yaml:"user"`
	Password string `yaml:"pass"`
	Charset  string `yaml:"charset"`

	SlaveID uint32 `yaml:"slave_id"`
}

func InitConfig(fileName string) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	var c config

	if err := yaml.Unmarshal(data, &c); err != nil {
		return err
	}

	_config = &c

	return nil
}

func Cfg() *config {
	return _config
}