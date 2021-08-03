package endpoint

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"go-mysql-replication/src/global"
)

type redisEndpoint struct {
	redisdb *redis.Client
}

func (r redisEndpoint) Consume(event global.Event) error {
	listKey := fmt.Sprintf("%s:%s",event.Schema,event.Table)
	return r.redisdb.RPush(listKey,event.ToString()).Err()
}

func NewRedisEndPoint() (Endpoint, error) {
	redisClient, err := getRedisClient()
	if err != nil {
		return nil, err
	}
	return &redisEndpoint{redisdb: redisClient}, nil
}

// 获取redis客户端
func getRedisClient() (*redis.Client, error) {
	redisConfig := global.Cfg().Redis
	//连接服务器
	redisdb := redis.NewClient(&redis.Options{
		Addr:     redisConfig.GetAddr(),
		Password: redisConfig.GetPasswd(),
		DB:       redisConfig.GetDb(),
	})

	//ping
	err := redisdb.Ping().Err()
	if err != nil {
		return nil, errors.New("redis conn err")
	}
	return redisdb, err
}
