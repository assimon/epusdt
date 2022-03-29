package dao

import (
	"context"
	"fmt"
	"github.com/assimon/luuu/util/log"
	"github.com/go-redis/redis/v8"
	"github.com/gookit/color"
	"github.com/spf13/viper"
	"time"
)

var Rdb *redis.Client

func RedisInit() {
	options := redis.Options{
		Addr: fmt.Sprintf(
			"%s:%s",
			viper.GetString("redis_host"),
			viper.GetString("redis_port")), // Redis地址
		DB:          viper.GetInt("redis_db"),                                        // Redis库
		PoolSize:    viper.GetInt("redis_pool_size"),                                 // Redis连接池大小
		MaxRetries:  viper.GetInt("redis_max_retries"),                               // 最大重试次数
		IdleTimeout: time.Second * time.Duration(viper.GetInt("redis_idle_timeout")), // 空闲链接超时时间
	}
	if viper.GetString("redis_passwd") != "" {
		options.Password = viper.GetString("redis_passwd")
	}
	Rdb = redis.NewClient(&options)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pong, err := Rdb.Ping(ctx).Result()
	if err == redis.Nil {
		log.Sugar.Debug("[store_redis] Nil reply returned by Rdb when key does not exist.")
	} else if err != nil {
		color.Red.Printf("[store_redis] redis connRdb err,err=%s", err)
		panic(err)
	} else {
		log.Sugar.Debug("[store_redis] redis connRdb success,suc=%s", pong)
	}
}
