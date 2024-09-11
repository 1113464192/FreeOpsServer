package model

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	Name string `json:"name" gorm:"type:varchar(20);uniqueIndex;comment:用户名"`
}
