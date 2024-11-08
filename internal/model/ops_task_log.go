package model

import "gorm.io/gorm"

type OpsTaskLog struct {
	gorm.Model
	Name            string `json:"name" gorm:"type: varchar(10);comment: 任务名;index"`
	Commands        string `json:"commands" gorm:"type: text;comment: [echo '1', echo '2']"`
	StepStatus      string `json:"stepStatus" gorm:"type: text;comment: commands执行顺序的每一步的状态(等待中、执行中、执行成功、执行失败), 键值切片对应Status/Res/Command/SSHStatus"`
	Status          uint8  `json:"status" gorm:"comment: 整个任务的状态(0: 等待中、1: 执行中、2: 执行成功、3: 执行失败、4: 拒绝执行);index"`
	Auditors        string `json:"auditors" gorm:"type: text;comment: 审核员"`
	PendingAuditors string `json:"pendingAuditors" gorm:"type: text;comment: 未审批的审核员"`
	RejectAuditor   uint   `json:"rejectAuditor" gorm:"comment: 拒绝执行的审核员"`
	TaskId          uint   `json:"taskId" gorm:"comment: 任务ID;index"`
	ProjectId       uint   `json:"projectId" gorm:"comment: 项目ID;index"`
	Submitter       uint   `json:"submitter" gorm:"comment: 提交的用户"`
}