package crontab

import (
	"github.com/robfig/cron"
)

func Cron() {
	c := cron.New()
	c.AddFunc("0 0 0 * * *", CronMysqlLogRename)
	// 定时每日清理Token黑名单中过期的Token
	c.AddFunc("0 30 5 * * * ", CronTokenBlacklistClean)
	// 定时每日删除所有model中软删除大于30天的数据
	c.AddFunc("0 30 4 * * *", CronModelSoftDeleteClean)
	c.Start()
}
