package endpoint

import "go-mysql-replication/src/global"

type Endpoint interface {
	Consume(event global.Event) error
}

func GetEndpoint() Endpoint {
	switch global.Cfg().Endpoint {
	case global.END_POINT_REDIS:
		return NewRedisEndPoint()
	case global.END_POINT_STDIO:
		return NewStdioEndpoint()
	}
	return nil
}
