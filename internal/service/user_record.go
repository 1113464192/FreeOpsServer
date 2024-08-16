package service

import (
	"FreeOps/internal/consts"
	"FreeOps/internal/model"
	"FreeOps/pkg/api"
	"errors"
	"fmt"
	"strings"
	"time"
)

type UserRecord struct {
}

var (
	insRecord = UserRecord{}
)

func UserRecordApp() *UserRecord {
	return &insRecord
}

// checkRecordTableExists 检查当月记录表是否存在，并返回表名
func (s *UserRecord) checkRecordTableExists() (tableName string, exist bool) {
	currTime := time.Now().Local()
	nowDate := currTime.Format("2006_01")
	tableName = fmt.Sprintf("%s_%s", consts.MysqlTableNameUserRecord, nowDate)
	// 等待表的创建，最多等待20s
	for i := 0; i < 5; i++ {
		if model.DB.Migrator().HasTable(tableName) {
			exist = true
			break
		}
		_ = model.DB.Table(tableName).Migrator().CreateTable(&model.UserRecord{})
		time.Sleep(time.Second)
	}
	return tableName, exist
}

// RecordCreate 插入日志
func (s *UserRecord) CreateRecord(log *model.UserRecord) (err error) {
	tableName, exist := s.checkRecordTableExists()
	if !exist {
		return errors.New("当月表尚未创建，请联系运维查看")
	}
	if err = model.DB.Table(tableName).Create(log).Error; err != nil {
		return fmt.Errorf("在当月记录表中创建记录失败: %v", err)
	}
	return err
}

// 查询有多少个月份表可供查询
func (s *UserRecord) GetUserRecordDate() (dates []string, err error) {
	// 构建原生 SQL 查询语句
	sql := fmt.Sprintf(`SHOW TABLES LIKE 'act_record_%%'`)
	if err = model.DB.Raw(sql).Scan(&dates).Error; err != nil {
		return nil, fmt.Errorf("获取所有记录表的表名失败: %v", err)
	}
	for i := 0; i < len(dates); i++ {
		dates[i] = strings.Replace(dates[i], "act_record_", "", -1)
	}
	return dates, err
}

// 查询月份记录
func (s *UserRecord) GetUserRecordLogs(param api.GetUserRecordLogsReq) (logs *[]model.UserRecord, total int64, err error) {
	tableName := fmt.Sprintf("%s_%s", consts.MysqlTableNameUserRecord, param.Date)
	if !model.DB.Migrator().HasTable(tableName) {
		return nil, 0, errors.New("没有这个日期的行为记录表存在: " + tableName)
	}

	// 先取出total
	if err := model.DB.Table(tableName).Where("user_id = ?", param.Id).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, nil
	}

	logs = &[]model.UserRecord{}

	if err := model.DB.Table(tableName).Where("user_id = ?", param.Id).Order("id desc").
		Offset((param.PageInfo.Page - 1) * param.PageInfo.PageSize).Limit(param.PageInfo.PageSize).Find(logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, err
}
