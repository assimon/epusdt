package bootstrap

import (
	"github.com/assimon/luuu/command"
	"github.com/assimon/luuu/config"
	"github.com/assimon/luuu/model/dao"
	"github.com/assimon/luuu/mq"
	"github.com/assimon/luuu/task"
	"github.com/assimon/luuu/telegram"
	"github.com/assimon/luuu/util/log"
)

//Start 服务启动
func Start() {
	// 配置加载
	config.Init()
	// 日志加载
	log.Init()
	// Mysql启动
	dao.MysqlInit()
	// redis启动
	dao.RedisInit()
	// 队列启动
	mq.Start()
	// telegram机器人启动
	go telegram.BotStart()
	// 定时任务
	go task.Start()
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}
