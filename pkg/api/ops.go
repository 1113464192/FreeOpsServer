package api

import "FreeOps/internal/model"

type UpdateOpsTemplateReq struct {
	ID        uint   `form:"id" json:"id"`
	Name      string `form:"name" json:"name"  binding:"required"`
	Content   string `form:"content" json:"content"  binding:"required"`
	ProjectId uint   `form:"projectId" json:"projectId"  binding:"required"`
}

type GetOpsTemplatesReq struct {
	ID        uint   `form:"id" json:"id"`
	Name      string `form:"name" json:"name"`
	ProjectId uint   `form:"projectId" json:"projectId"`
	PageInfo
}

type GetOpsTemplateRes struct {
	ID        uint   `form:"id" json:"id"`
	UpdatedAt string `form:"updatedAt" json:"updatedAt"`
	Name      string `form:"name" json:"name"`
	Content   string `form:"content" json:"content,omitempty"`
	ProjectId uint   `form:"projectId" json:"projectId"`
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
	ParamIDs   []uint `form:"paramIds" json:"paramIds" binding:"required"`
}

type UpdateOpsTaskReq struct {
	ID          uint   `form:"id" json:"id"`
	Name        string `form:"name" json:"name"  binding:"required"`
	TemplateIds string `form:"templateIds" json:"templateIds" binding:"required"`
	Auditors    string `form:"auditors" json:"auditors"`
	HostId      uint   `form:"hostId" json:"hostId"  binding:"required"`
	IsIntranet  bool   `form:"isIntranet" json:"isIntranet"`
	ProjectId   uint   `form:"projectId" json:"projectId"  binding:"required"`
}

type GetOpsTaskReq struct {
	ID        uint   `form:"id" json:"id"`
	Name      string `form:"name" json:"name"`
	ProjectId uint   `form:"projectId" json:"projectId"`
	PageInfo
}

type GetOpsTaskRes struct {
	ID          uint   `form:"id" json:"id"`
	Name        string `form:"name" json:"name"`
	TemplateIds []uint `form:"templateIds" json:"templateIds"`
	Auditors    []uint `form:"userIds" json:"userIds,omitempty"`
	HostId      uint   `form:"hostId" json:"hostId"`
	ProjectId   uint   `form:"projectId" json:"projectId"`
}

type GetOpsTasksRes struct {
	Records  []GetOpsTaskRes `json:"records" form:"records"`
	Page     int             `json:"current" form:"current"` // 页码
	PageSize int             `json:"size" form:"size"`       // 每页大小
	Total    int64           `json:"total"`
}
