package api

import (
	"FreeOps/internal/model"
)

type UpdateOpsTemplateReq struct {
	ID        uint   `form:"id" json:"id"`
	Name      string `form:"name" json:"name"  binding:"required"`
	Content   string `form:"content" json:"content"  binding:"required"` // 如: ls ${path}
	ProjectId uint   `form:"projectId" json:"projectId"  binding:"required"`
}

type GetOpsTemplatesReq struct {
	ID        uint   `form:"id" json:"id"`
	Name      string `form:"name" json:"name"`
	ProjectId uint   `form:"projectId" json:"projectId"`
	PageInfo
}

type GetOpsTemplateRes struct {
	ID          uint   `form:"id" json:"id"`
	UpdatedAt   string `form:"updatedAt" json:"updatedAt"`
	Name        string `form:"name" json:"name"`
	Content     string `form:"content" json:"content,omitempty"`
	ProjectName string `form:"projectName" json:"projectName"`
	ProjectId   uint   `form:"projectId" json:"projectId"`
}

type GetOpsTemplatesRes struct {
	Records  []GetOpsTemplateRes `json:"records" form:"records"`
	Page     int                 `json:"current" form:"current"` // 页码
	PageSize int                 `json:"size" form:"size"`       // 每页大小
	Total    int64               `json:"total"`
}

type GetOpsParamsTemplatesReq struct {
	ID      uint   `form:"id" json:"id"`
	Keyword string `form:"keyword" json:"keyword"`
	PageInfo
}

type GetOpsParamsTemplatesRes struct {
	Records  []model.OpsParam `json:"records" form:"records"`
	Page     int              `json:"current" form:"current"` // 页码
	PageSize int              `json:"size" form:"size"`       // 每页大小
	Total    int64            `json:"total"`
}

type BindTemplateParamsReq struct {
	TemplateID uint   `form:"templateId" json:"templateId" binding:"required"`
	ParamIDs   []uint `form:"paramIds" json:"paramIds"`
}

type UpdateOpsTaskReq struct {
	ID              uint   `form:"id" json:"id"`
	Name            string `form:"name" json:"name"  binding:"required"`
	CheckTemplateId uint   `form:"checkTemplateId" json:"checkTemplateId"` // 默认被执行的运维检测脚本的模板，返回的打印信息是用于给运营审批查看的，不传则跳过检查
	TemplateIds     []uint `form:"templateIds" json:"templateIds" binding:"required"`
	Auditors        []uint `form:"auditors" json:"auditors"`
	HostId          uint   `form:"hostId" json:"hostId"  binding:"required"`
	IsIntranet      bool   `form:"isIntranet" json:"isIntranet"`
	IsConcurrent    bool   `form:"isConcurrent" json:"isConcurrent"`
	ProjectId       uint   `form:"projectId" json:"projectId"  binding:"required"`
}

type GetOpsTaskReq struct {
	ID        uint   `form:"id" json:"id"`
	Name      string `form:"name" json:"name"`
	ProjectId uint   `form:"projectId" json:"projectId"`
	PageInfo
}

type GetOpsTaskRes struct {
	ID              uint   `form:"id" json:"id"`
	Name            string `form:"name" json:"name"`
	CheckTemplateId uint   `form:"checkTemplateId" json:"checkTemplateId,omitempty"`
	TemplateIds     []uint `form:"templateIds" json:"templateIds,omitempty"`
	Auditors        []uint `form:"userIds" json:"userIds,omitempty"`
	HostId          uint   `form:"hostId" json:"hostId"`
	IsIntranet      bool   `form:"isIntranet" json:"isIntranet"`
	IsConcurrent    bool   `form:"isConcurrent" json:"isConcurrent"`
	ProjectId       uint   `form:"projectId" json:"projectId"`
	ProjectName     string `form:"projectName" json:"projectName"`
}

type GetOpsTasksRes struct {
	Records  []GetOpsTaskRes `json:"records" form:"records"`
	Page     int             `json:"current" form:"current"` // 页码
	PageSize int             `json:"size" form:"size"`       // 每页大小
	Total    int64           `json:"total"`
}

type RunOpsTaskCheckScriptReq struct {
	ExecContext string `form:"execContent" json:"execContent" binding:"required"` // 运营的执行内容文案，从中根据Params提取参数放入模板中执行
	TaskId      uint   `form:"taskId" json:"taskId" binding:"required"`
}

type OpsTaskLogtepStatus struct {
	Command           string `json:"command"`
	Status            int    `json:"status"`
	Response          string `json:"response"`
	SSHResponseStatus int    `json:"sshResponseStatus"`
}

type SubmitOpsTaskReq struct {
	TaskId      uint   `form:"taskId" json:"taskId" binding:"required"`
	ExecContext string `form:"execContent" json:"execContent" binding:"required"` // 运营的执行内容文案，从中根据Params提取参数放入各个模板中执行
	TemplateIds []uint `form:"templateIds" json:"templateIds" binding:"required"` // 模板是可勾选的，因此不一定完全执行taskId的所有模板，所以需要单独传。按顺序如:1,2,3
	Auditors    []uint `form:"auditors" json:"auditors"`
	Submitter   uint   `form:"submitter" json:"submitter" binding:"required"` // 提交者
}

type ApproveOpsTaskReq struct {
	TaskId  uint `form:"taskId" json:"taskId" binding:"required"` // 任务日志ID
	IsAllow bool `form:"isAllow" json:"isAllow"`
}

type GetTaskPendingApproversRes struct {
	TaskName     string   `json:"taskName"`
	Submitter    string   `json:"submitter"`
	PendingUsers []string `json:"pendingUsers"`
}

type GetOpsTaskLogReq struct {
	ID        uint   `form:"id" json:"id"`
	Name      string `form:"name" json:"name"`
	Status    uint8  `form:"status" json:"status"`
	ProjectId uint   `form:"projectId" json:"projectId"`
	Username  string `form:"username" json:"username"`
	PageInfo
}

type GetOpsTaskLogRes struct {
	ID              uint                  `json:"id"`
	Name            string                `json:"name"`
	Commands        []string              `json:"commands,omitempty"`
	StepStatus      []OpsTaskLogtepStatus `json:"stepStatus,omitempty"`
	Status          uint8                 `json:"status"`
	Auditors        []uint                `json:"auditors"`
	PendingAuditors []uint                `json:"pendingAuditors,omitempty"`
	RejectAuditor   uint                  `json:"rejectAuditor"`
	ProjectName     string                `json:"projectName"`
	ProjectId       uint                  `json:"projectId"`
	Submitter       uint                  `json:"submitter"`
}

type GetOpsTaskLogsRes struct {
	Records  []GetOpsTaskLogRes `json:"records" form:"records"`
	Page     int                `json:"current" form:"current"` // 页码
	PageSize int                `json:"size" form:"size"`       // 每页大小
	Total    int64              `json:"total"`
}

type GetOpsTaskRunningWSRes struct {
	Name         string `json:"name"`
	Command      string `json:"command"`
	Status       uint8  `json:"status"`
	ProjectName  string `json:"projectName"`
	Submitter    string `json:"submitter"`
	IsIntranet   bool   `json:"isIntranet"`
	IsConcurrent bool   `json:"isConcurrent"`
}
