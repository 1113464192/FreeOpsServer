package model

import "gorm.io/gorm"

type OpsTemplate struct {
	gorm.Model
	Name      string `json:"name" gorm:"type: varchar(10);comment: 模板名;index"`
	Content   string `json:"content,omitempty" gorm:"type: text;comment: 模板内容"`
	ProjectId uint   `json:"projectId" gorm:"comment: 项目ID;index"`
}
