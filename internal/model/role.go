package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	RoleName string `json:"roleName" gorm:"type:varchar(30);comment:角色名称"`
	RoleCode string `json:"roleCode" gorm:"type:varchar(30);index;comment:角色代码"`
	RoleDesc string `json:"roleDesc" gorm:"type:varchar(255);comment:角色描述"`
}
