package model

import (
	"FreeOps/global"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB
var LogFile *os.File
var NewLogger logger.Interface

// 所有需要创建的字段
var autoMigrateModelList = []interface{}{
	&User{},
	&Role{},
	&UserRole{},
	&Button{},
	&Menu{},
	&MenuRole{},
	&RoleButton{},
	&Api{},
	&JwtBlacklist{},
}

func Database() {
	var err error
	if dirStat, err := os.Stat(global.RootPath + "/logs/mysql/"); err != nil || !dirStat.IsDir() {
		if err = os.MkdirAll(global.RootPath+"/logs/mysql/", 0777); err != nil {
			log.Fatalf("创建mysql目录失败: %v", err)
		}
	}

	logFileName := global.RootPath + "/logs/mysql/mysql.log"

	LogFile, err = os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("无法打开日志文件:", err)
		return
	}

	NewLogger = logger.New(
		log.New(LogFile, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Duration(global.Conf.Mysql.SlowThreshold) * time.Second, // Slow SQL threshold
			LogLevel:                  logger.LogLevel(global.Conf.Mysql.LogLevel),                  // Log level(这里记得根据需求改一下)
			IgnoreRecordNotFoundError: true,                                                         // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,                                                         // Disable color
		},
	)
	dsn := global.Conf.Mysql.Conf
	if dsn == "" {
		log.Fatalf("dsn配置为空")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger:          NewLogger,
		CreateBatchSize: global.Conf.Mysql.CreateBatchSize,
		NowFunc: func() time.Time {
			tmp := time.Now().Local().Format("2006-01-02 15:04:05")
			now, _ := time.ParseInLocation("2006-01-02 15:04:05", tmp, time.Local)
			return now
		},
	})
	if err != nil {
		log.Fatalf("生成GORM.DB失败: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取DB失败: %v", err)
	}
	// 连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(30)
	// 设置连接的最大生命周期为 0，这意味着连接在连接池中没有最大生命周期的限制，它可以一直保持打开状态
	sqlDB.SetConnMaxLifetime(0)
	DB = db
}

func AutoMigrateMysql() {
	if global.Conf.System.Mode != "product" {
		err := DB.AutoMigrate(
			autoMigrateModelList...,
		)
		if err != nil {
			log.Fatalf("自动迁移报错: \n%v", err)
			return
		}
	}
}
