package crontab

import (
	"FreeOps/global"
	"FreeOps/internal/consts"
	"FreeOps/internal/model"
	"FreeOps/pkg/logger"
	"FreeOps/pkg/util"
	"fmt"
	gLogger "gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func CronMysqlLogRename() {
	var err error
	now := time.Now().Local()
	previousDay := now.AddDate(0, 0, -1) // 获取前一天的日期
	logFileName := fmt.Sprintf(global.RootPath+"/logs/mysql/%s.log", previousDay.Format("20060102"))

	// 关闭之前的日志文件描述符
	if err := model.LogFile.Close(); err != nil {
		logger.Log().Error("crontab", "关闭mysql日志文件描述符失败", err)
		return
	}

	// 重命名日志文件
	if err := os.Rename(global.RootPath+"/logs/mysql/mysql.log", logFileName); err != nil {
		logger.Log().Error("crontab", "重命名mysql日志失败,原日志名:", err)
		return
	}

	logFileName = global.RootPath + "/logs/mysql/mysql.log"

	model.LogFile, err = os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Log().Error("crontab", "CronMysqlLogRename无法打开日志文件:", err)
		return
	}

	model.NewLogger = gLogger.New(
		log.New(model.LogFile, "\r\n", log.LstdFlags), // io writer
		gLogger.Config{
			SlowThreshold:             time.Duration(global.Conf.Mysql.SlowThreshold) * time.Second, // Slow SQL threshold
			LogLevel:                  gLogger.LogLevel(global.Conf.Mysql.LogLevel),                 // Log level(这里记得根据需求改一下)
			IgnoreRecordNotFoundError: true,                                                         // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,                                                         // Disable color
		},
	)
	model.DB.Logger = model.NewLogger

}

func CronTokenBlacklistClean() {
	var err error
	// 获取黑名单中的token
	var tokenList []model.JwtBlacklist
	if err = model.DB.Find(&tokenList).Error; err != nil {
		logger.Log().Error("token", "获取黑名单中的token失败", err)
		return
	}
	// 遍历黑名单中的token，判断是否过期
	var expiredTokenIds []uint
	for _, token := range tokenList {
		// 如果过期，删除
		_, err := util.ParseToken(token.Jwt)
		if err != nil && err.Error() == "invalid token" {
			expiredTokenIds = append(expiredTokenIds, token.ID)
		}
	}
	// 检查 expiredTokenIds 是否为空
	if len(expiredTokenIds) == 0 {
		logger.Log().Info("token", "没有过期的token需要删除")
		return
	}
	// 从黑名单删除token以节省空间,不使用软删除
	if err = model.DB.Delete(&model.JwtBlacklist{}, expiredTokenIds).Error; err != nil {
		logger.Log().Error("token", "删除过期token失败", err.Error())
	}
}

func CronModelSoftDeleteClean() {
	dateThreshold := time.Now().AddDate(0, 0, global.Conf.Mysql.SoftDeleteRetainDays) // 获取SoftDeleteRetainDays天前的日期
	for _, m := range consts.SoftDeleteModelList {
		// 防止报错
		if !model.DB.Migrator().HasColumn(m, "deleted_at") {
			logger.Log().Warning("crontab", fmt.Sprintf("表 %T 中不存在 deleted_at 列, 请修改consts", m))
			continue
		}
		for {
			mDB := model.DB.Where("deleted_at < ?", dateThreshold).Unscoped().Delete(m).Limit(global.Conf.Mysql.CreateBatchSize)
			if mDB.Error != nil {
				logger.Log().Error("crontab", "删除软数据失败", mDB.Error)
				break
			}
			if mDB.RowsAffected < int64(global.Conf.Mysql.CreateBatchSize) {
				break
			}
		}
	}
}
