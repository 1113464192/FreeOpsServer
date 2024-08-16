package main

import (
	"FreeOps/configs"
	"FreeOps/internal/controller"
	"FreeOps/internal/crontab"
	"FreeOps/internal/model"
	"log"
)

func main() {
	configs.Init()
	model.Database()
	crontab.Cron()
	model.AutoMigrateMysql()

	r := controller.NewRoute()

	err := r.Run(":9080")
	if err != nil {
		log.Fatalf("启动报错: \n%v", err)
		return
	}
}
