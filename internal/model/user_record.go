package model

import (
	"gorm.io/gorm"
	"time"
)

type UserRecord struct {
	gorm.Model
	Ip       string        `json:"ip" form:"ip" gorm:"type:varchar(50);column:ip;comment:请求ip"`                      // 请求ip
	Method   string        `json:"method" form:"method" gorm:"type:varchar(10);column:method;comment:请求方法"`          // 请求方法
	Path     string        `json:"path" form:"path" gorm:"type:varchar(100);column:path;comment:请求路径"`               // 请求路径
	Agent    string        `json:"agent" form:"agent" gorm:"type:varchar(255);column:agent;comment:代理"`              // 代理
	Body     string        `json:"body" form:"body" gorm:"type:text;column:body;comment:请求Body"`                     // 请求Body
	UserID   uint          `json:"userId" form:"userId" gorm:"column:user_id;comment:用户id;index"`                    // 用户id
	Username string        `json:"username" form:"username" gorm:"type:varchar(30);column:username;comment:用户账号"`    // 用户账号
	Status   int           `json:"status" form:"status" gorm:"column:status;comment:返回状态"`                           // 请求状态
	Latency  time.Duration `json:"latency" form:"latency" gorm:"column:latency;comment:延迟(纳秒)" swaggertype:"string"` // 延迟
	Resp     string        `json:"resp" form:"resp" gorm:"type:text;column:resp;comment:响应Body"`                     // 响应Body
}
