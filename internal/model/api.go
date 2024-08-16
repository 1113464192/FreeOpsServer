package model

import "gorm.io/gorm"

type Api struct {
	gorm.Model
	Path        string `json:"path" gorm:"type:varchar(100);comment:api路径"`          // api路径
	Method      string `json:"method" gorm:"type:varchar(10);comment:方法"`            // 方法:创建/更新POST(默认)|查看GET|删除DELETE
	ApiGroup    string `json:"apiGroup" gorm:"type:varchar(10);comment:api组"`        // api组
	Description string `json:"description" gorm:"type:varchar(255);comment:api中文描述"` // api中文描述
}
