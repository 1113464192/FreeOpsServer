package model

import "gorm.io/gorm"

type OpsTaskLogs struct {
	gorm.Model
	Name            string `json:"name" gorm:"type: varchar(10);comment: 任务名;index"`
	TemplateIds     string `json:"templateIds" gorm:"type: text;comment: 任务模板关联顺序(比如轮流执行模板1、2、3)"`
	StepStatus      string `json:"stepStatus" gorm:"type: text;comment: 任务执行顺序的每一步的状态(等待中、执行中、执行成功、执行失败)"`
	Status          uint8  `json:"status" gorm:"comment: 整个任务的状态(等待中、执行中、执行成功、执行失败);index"`
	Auditors        string `json:"auditors" gorm:"type: text;comment: 审核员"`
	PendingAuditors string `json:"pendingAuditors" gorm:"type: text;comment: 未审批的审核员"`
	RejectAuditor   uint   `json:"rejectAuditor" gorm:"comment: 拒绝执行的审核员"`
	TaskId          uint   `json:"taskId" gorm:"comment: 任务ID;index"`
	ProjectId       uint   `json:"projectId" gorm:"comment: 项目ID;index"`
}
