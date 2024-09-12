package service

import (
	"FreeOps/internal/consts"
	"FreeOps/internal/model"
	"FreeOps/pkg/api"
	"errors"
	"fmt"
	"strings"
	"sync"
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
func (s *UserRecord) checkRecordTableExists() (tableName string, err error) {
	currTime := time.Now().Local()
	nowDate := currTime.Format("2006_01")
	tableName = fmt.Sprintf("%s_%s", consts.MysqlTableNameUserRecord, nowDate)
	if !model.DB.Migrator().HasTable(tableName) {
		var mutex sync.Mutex
		mutex.Lock()
		defer mutex.Unlock()
		if !model.DB.Migrator().HasTable(tableName) {
			if err = model.DB.Table(tableName).Migrator().CreateTable(&model.UserRecord{}); err != nil {
				return tableName, fmt.Errorf("创建当月记录表失败: %v", err)
			}
		}
	}
	return tableName, err
}

// RecordCreate 插入日志
func (s *UserRecord) CreateRecord(log *model.UserRecord) (err error) {
	var tableName string
	if tableName, err = s.checkRecordTableExists(); err != nil {
		return err
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
func (s *UserRecord) GetUserRecordLogs(param api.GetUserRecordLogsReq) (logs *[]api.GetUserRecordLogsServiceRes, total int64, err error) {
	tableName := fmt.Sprintf("%s_%s", consts.MysqlTableNameUserRecord, param.Date)
	if !model.DB.Migrator().HasTable(tableName) {
		return nil, 0, errors.New("没有这个日期的行为记录表存在: " + tableName)
	}

	getDB := model.DB.Table(tableName)
	if param.Username != "" {
		sqlUsername := "%" + strings.ToUpper(param.Username) + "%"
		getDB = getDB.Where("UPPER(username) LIKE ?", sqlUsername)
	}

	if param.Method != "" {
		getDB = getDB.Where("method = ?", param.Method)
	}

	if param.Ip != "" {
		getDB = getDB.Where("ip = ?", param.Ip)
	}

	if param.Status != "" {
		getDB = getDB.Where("status LIKE ?", param.Status)
	}

	// 先取出total
	if err := getDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, nil
	}

	userRecord := []model.UserRecord{}

	if err := getDB.Order("id desc").
		Offset((param.PageInfo.Page - 1) * param.PageInfo.PageSize).Limit(param.PageInfo.PageSize).Find(&userRecord).Error; err != nil {
		return nil, 0, err
	}

	logs = &[]api.GetUserRecordLogsServiceRes{}
	for _, v := range userRecord {
		*logs = append(*logs, api.GetUserRecordLogsServiceRes{
			Id:        v.ID,
			CreatedAt: v.CreatedAt,
			Ip:        v.Ip,
			Method:    v.Method,
			Path:      v.Path,
			Agent:     v.Agent,
			Body:      v.Body,
			UserId:    v.UserID,
			Username:  v.Username,
			Status:    v.Status,
			Latency:   v.Latency,
			Resp:      v.Resp,
		})
	}
	return logs, total, err
}
