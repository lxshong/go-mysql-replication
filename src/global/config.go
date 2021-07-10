package global

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var _config *config

type config struct {
	Target string `yaml:"target"` // 目标类型，支持redis、mongodb

	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"pass"`
	Charset  string `yaml:"charset"`
	Flavor   string `yaml:"flavor"`
	SlaveID  uint32 `yaml:"slave_id"`

	PosType string `yaml:"pos_type"`
	PosFile string `yaml:"pos_file"`
	Rules   rules  `yaml:"rules"`
	Tables  Tables
}

const (
	_POS_TYPE_FILE = "file"
)

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

	_config.Tables = map[string]map[string]bool{}
	if err := _config.Tables.Change(c.Rules); err != nil {
		return err
	}

	return nil
}

func Cfg() *config {
	return _config
}

func PosIsFile() bool {
	if _config.PosType == _POS_TYPE_FILE {
		return true
	}
	return false
}
