package api

import "FreeOps/internal/model"

type UpdateApiReq struct {
	ID          uint   `json:"id" form:"id"`
	Path        string `json:"path" form:"path" binding:"required"`     // 如果API为/user就/user，不要写成/user/
	Method      string `json:"method" form:"method" binding:"required"` // 必须大写
	ApiGroup    string `json:"apiGroup" form:"apiGroup" binding:"required"`
	Description string `json:"description" form:"description"`
}

type GetApiReq struct {
	ID       uint   `form:"id" json:"id"`
	ApiGroup string `json:"apiGroup" form:"apiGroup"`
	PageInfo
}

type GetApiRes struct {
	Records  []model.Api `json:"records" form:"records"`
	Page     int         `json:"current" form:"current"` // 页码
	PageSize int         `json:"size" form:"size"`       // 每页大小
	Total    int64       `json:"total"`
}
