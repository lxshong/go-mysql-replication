package endpoint

import "go-mysql-replication/src/global"

type redisEndpoint struct {
}

func (r redisEndpoint) Consume(event global.Event) error {
	return nil
}

func NewRedisEndPoint() Endpoint {
	return &redisEndpoint{}
}
