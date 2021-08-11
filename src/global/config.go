package global

import (
	"fmt"
	"github.com/go-mysql-org/go-mysql/schema"
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
	PosZk []string `yaml:"pos_zk"`
	Rules   rules  `yaml:"rules"`
	Tables  Tables
	// 端点
	Endpoint string `yaml:"endpoint"`
	Redis    redis  `yaml:"redis"`
}

const (
	_POS_TYPE_ZK  = "zk"
	_POS_TYPE_FILE  = "file"
	END_POINT_REDIS = "redis"
	END_POINT_STDIO = "stdio"
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

	_config.Tables = map[string]map[string]*schema.Table{}
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

func PosIsZk() bool {
	if _config.PosType == _POS_TYPE_ZK {
		return true
	}
	return false
}

// redis
type redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Db       int    `yaml:"db"`
	Password string `yaml:"passwd"`
}

func (receiver redis) GetAddr() string {
	return fmt.Sprintf("%s:%d", receiver.Host, receiver.Port)
}

func (receiver redis) GetDb() int {
	return receiver.Db
}

func (receiver redis) GetPasswd() string {
	return receiver.Password
}
