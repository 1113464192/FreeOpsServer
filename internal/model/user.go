package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID         uint `gorm:"primarykey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time      `gorm:"index"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Status     uint           `json:"status" gorm:"default:1;index;comment:状态: 1(enabled),2(disabled);type:tinyint"`
	Username   string         `json:"username" gorm:"type:varchar(20);uniqueIndex;comment:用户名"`
	Password   string         `json:"password" gorm:"type:varchar(60);comment:密码"`
	UserGender string         `json:"userGender" gorm:"type:varchar(2);comment: 1(male),2(female)"`
	Nickname   string         `json:"nickname" gorm:"type:varchar(30);comment:昵称"`
	UserPhone  string         `json:"userPhone" gorm:"type:varchar(50);comment:手机号"`
	UserEmail  string         `json:"userEmail" gorm:"type:varchar(50);comment:邮箱"`
}
