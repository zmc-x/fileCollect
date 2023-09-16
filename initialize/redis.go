package initialize

import (
	"context"
	"fileCollect/global"
	"fileCollect/utils/zaplog"
	"fmt"

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
		zaplog.GetLogLevel(zaplog.FATAL, err.Error())
	} else {
		zaplog.GetLogLevel(zaplog.INFO, "redis connect ping response :" + Pong)
		global.Rdb = cli
	}
}