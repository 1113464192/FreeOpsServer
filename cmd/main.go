package main

import (
	"FreeOps/configs"
	"FreeOps/internal/controller"
	"FreeOps/internal/crontab"
	"FreeOps/internal/model"
	"FreeOps/internal/service"
	"log"
)

func main() {
	configs.Init()
	model.Database()
	crontab.Cron()
	model.AutoMigrateMysql()
	// 创建casbin_rule表
	service.CasbinServiceApp().Casbin()

	r := controller.NewRoute()

	err := r.Run(":9080")
	if err != nil {
		log.Fatalf("启动报错: \n%v", err)
		return
	}
}
