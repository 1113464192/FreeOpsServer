package model

import "gorm.io/gorm"

type OpsTask struct {
	gorm.Model
	Name        string `json:"name" gorm:"type: varchar(10);comment: 任务名;index"`
	TemplateIds string `json:"templateIds" gorm:"type: text;comment: 任务模板关联顺序(比如轮流执行模板1、2、3)"`
	Auditors    string `json:"userIds" gorm:"type: text;comment: 任务关联的用户(审批，都审批完了就能执行)"`
	HostId      uint   `json:"hostId" gorm:"comment: 运维管理机ID;index"`
	ProjectId   uint   `json:"projectId" gorm:"comment: 项目ID;index"`
}
