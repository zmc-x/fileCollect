package initialize

import (
	"context"
	"fileCollect/global"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func InitReids(sc *global.ServerConfig) {
	cli := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", sc.RedisAddr, sc.RedisPort),
		Password: sc.RedisPasswd,
		DB: int(sc.RedisDb),
	})
	// test connection
	Pong, err := cli.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("initialize/redis.go InitRedis:" + err.Error())
	} else {
		log.Println("initialize/redis.go InitReids: redis connect ping response :" + Pong)
		global.Rdb = cli
	}
}