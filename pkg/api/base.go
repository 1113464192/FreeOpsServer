package api

type Response struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// PageInfo 分页请求
type PageInfo struct {
	Page     int `json:"current" form:"current"` // 页码
	PageSize int `json:"size" form:"size"`       // 每页大小
}

type IdsReq struct {
	Ids []uint `json:"ids" form:"ids" binding:"required"`
}

type IdPageReq struct {
	Id uint `json:"id" form:"id"`
	PageInfo
}

// 定义通用关联返回的struct
type RelationPageRes struct {
	Records  []uint `json:"records" form:"records"`
	Page     int    `json:"current" form:"current"` // 页码
	PageSize int    `json:"size" form:"size"`       // 每页大小
	Total    int64  `json:"total"`
}

type CustomErrorReq struct {
	Code string `json:"code" form:"code" binding:"required"`
	Msg  string `json:"msg" form:"msg" binding:"required"`
}

type GetIdAndNameRes struct {
	ID   uint   `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
}
