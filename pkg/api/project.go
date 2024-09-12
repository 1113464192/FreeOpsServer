package api

type UpdateProjectReq struct {
	ID   uint   `form:"id" json:"id"`                        // 修改才需要传，没有传算新增
	Name string `form:"name" json:"name" binding:"required"` // 项目名
}

type GetProjectReq struct {
	ID   uint   `form:"id" json:"id"` // 修改才需要传，没有传算新增
	Name string `form:"name" json:"name"`
}

type GetProjectsReq struct {
	GetProjectReq
	PageInfo
}

type GetProjectsRes struct {
	Records  []GetProjectReq `json:"records" form:"records"`
	Page     int             `json:"current" form:"current"` // 页码
	PageSize int             `json:"size" form:"size"`       // 每页大小
	Total    int64           `json:"total"`
}
