package model

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	Name string `json:"name" gorm:"type:varchar(20);uniqueIndex;comment:项目名"`
	// 拆分开，方便后续服务器维护，以及不同云商操作隔离
	CloudPlatform string `json:"cloudPlatform" gorm:"type:varchar(20);comment:云平台,国语即可"`
}
