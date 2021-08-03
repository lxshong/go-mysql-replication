package endpoint

import (
	"errors"
	"go-mysql-replication/src/global"
)

type Endpoint interface {
	Consume(event global.Event) error
}

func GetEndpoint() (Endpoint,error) {
	switch global.Cfg().Endpoint {
	case global.END_POINT_REDIS:
		return NewRedisEndPoint()
	case global.END_POINT_STDIO:
		return NewStdioEndpoint()
	}
	return nil,errors.New("不支持终端类型")
}
