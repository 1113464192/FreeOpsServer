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
	Auditors        []uint `form:"auditors" json:"auditors,omitempty"`
	HostId          uint   `form:"hostId" json:"hostId"`
	HostName        string `form:"hostName" json:"hostName"`
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

type OpsTaskLogStepStatus struct {
	Command           string `json:"command"`
	StartTime         string `json:"startTime"`
	EndTime           string `json:"endTime"`
	Status            uint8  `json:"status"`
	Response          string `json:"response"`
	SSHResponseStatus int    `json:"sshResponseStatus"`
}

type GetOpsTaskTmpCommandsReq struct {
	TemplateIds []uint `form:"templateIds" json:"templateIds" binding:"required"`
	ExecContext string `form:"execContent" json:"execContent" binding:"required"` // 运营的执行内容文案，从中根据Params提取参数放入模板中执行
}

type SubmitOpsTaskReq struct {
	TaskId        uint   `form:"taskId" json:"taskId" binding:"required"`
	ExecContext   string `form:"execContent" json:"execContent" binding:"required"` // 运营的执行内容文案，从中根据Params提取参数放入各个模板中执行
	CheckResponse string `form:"checkResponse" json:"checkResponse"`                // 运维检测脚本根据ExecContext返回的信息，存储到任务日志中
	TemplateIds   []uint `form:"templateIds" json:"templateIds" binding:"required"` // 模板是可勾选的，因此不一定完全执行taskId的所有模板，所以需要单独传。按顺序如:1,2,3
	Auditors      []uint `form:"auditors" json:"auditors"`
	Submitter     uint   `form:"submitter" json:"submitter" binding:"required"` // 提交者
	ExecTime      int64  `json:"execTime"`                                      // 指定执行时间(不选默认审批完或者没有审批人就立即执行)
}

type ApproveOpsTaskReq struct {
	ID      uint `form:"id" json:"id" binding:"required"` // 任务日志ID
	IsAllow bool `form:"isAllow" json:"isAllow"`
}

type GetUserTaskPendingReq struct {
	TaskName  string `form:"taskName" json:"taskName"`
	ProjectId uint   `form:"projectId" json:"projectId"`
	PageInfo
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
	ID                  uint                   `json:"id"`
	StartTime           string                 `json:"startTime"`
	EndTime             string                 `json:"endTime"`
	Name                string                 `json:"name"`
	HostIp              string                 `json:"hostIp"`
	ExecContext         string                 `json:"execContext,omitempty"`
	CheckResponse       string                 `json:"checkResponse,omitempty"`
	Commands            []string               `json:"commands,omitempty"`
	StepStatus          []OpsTaskLogStepStatus `json:"stepStatus,omitempty"`
	Status              uint8                  `json:"status"`
	Auditors            []uint                 `json:"auditors,omitempty"`
	AuditorNames        []string               `json:"auditorNames,omitempty"`
	PendingAuditors     []uint                 `json:"pendingAuditors,omitempty"`
	PendingAuditorNames []string               `json:"pendingAuditorNames,omitempty"`
	RejectAuditor       uint                   `json:"rejectAuditor,omitempty"`
	RejectAuditorName   string                 `json:"rejectAuditorName,omitempty"`
	ProjectName         string                 `json:"projectName"`
	ProjectId           uint                   `json:"projectId"`
	Submitter           uint                   `json:"submitter"`
	SubmitterName       string                 `json:"submitterName"`
	ExecTime            string                 `json:"execTime"`
}

type GetOpsTaskLogsRes struct {
	Records  []GetOpsTaskLogRes `json:"records" form:"records"`
	Page     int                `json:"current" form:"current"` // 页码
	PageSize int                `json:"size" form:"size"`       // 每页大小
	Total    int64              `json:"total"`
}

type GetOpsTaskRunningWSRes struct {
	Name          string                 `json:"name"`
	StartTime     string                 `json:"startTime"`
	EndTime       string                 `json:"endTime"`
	Status        uint8                  `json:"status"`
	SubmitterName string                 `json:"submitterName"`
	Children      []OpsTaskLogStepStatus `json:"children,omitempty"`
}
