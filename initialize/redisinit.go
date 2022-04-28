package initialize

// 报错没有做处理

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/go-redis/redis/v7"
)

// 声明全局rdb变量
var Rdb *redis.Client

func redisinit() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       1,
	})
	_, err := Rdb.Ping().Result()
	
	if err != nil {
		logs.Error("redis链接失败", err)
		return
	}
}
