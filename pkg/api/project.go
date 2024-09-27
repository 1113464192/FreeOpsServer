package api

type UpdateProjectReq struct {
	ID            uint   `form:"id" json:"id"`                                          // 修改才需要传，没有传算新增
	Name          string `form:"name" json:"name" binding:"required"`                   // 项目名
	CloudPlatform string `form:"cloudPlatform" json:"cloudPlatform" binding:"required"` // 云平台
}

type GetProjectReq struct {
	ID            uint   `form:"id" json:"id"` // 修改才需要传，没有传算新增
	Name          string `form:"name" json:"name"`
	CloudPlatform string `form:"cloudPlatform" json:"cloudPlatform"`
}

type GetProjectsReq struct {
	GetProjectReq
	PageInfo
}

type GetProjectRes struct {
	GetProjectReq
	GetProjectAssetsTotalRes
}

type GetProjectsRes struct {
	Records  []GetProjectRes `json:"records" form:"records"`
	Page     int             `json:"current" form:"current"` // 页码
	PageSize int             `json:"size" form:"size"`       // 每页大小
	Total    int64           `json:"total"`
}

type GetProjectAssetsTotalRes struct {
	HostTotal   int64 `json:"hostTotal" form:"hostTotal"`     // 服务器总数
	GameTotal   int64 `json:"gameTotal" form:"gameTotal"`     // 游服总数
	CrossTotal  int64 `json:"crossTotal" form:"crossTotal"`   // 跨服总数
	CommonTotal int64 `json:"commonTotal" form:"commonTotal"` // 公共服总数
}
