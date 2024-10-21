package api

type UpdateProjectReq struct {
	ID              uint    `form:"id" json:"id"`                                              // 修改才需要传，没有传算新增
	Name            string  `form:"name" json:"name" binding:"required"`                       // 项目名
	BackendVersion  string  `form:"backendVersion" json:"backendVersion" binding:"required"`   // 后端版本
	FrontendVersion string  `form:"frontendVersion" json:"frontendVersion" binding:"required"` // 前端版本
	GameRange       string  `form:"gameRange" json:"gameRange" binding:"required"`             // 游服ID范围
	CrossRange      string  `form:"crossRange" json:"crossRange" binding:"required"`           // 跨服ID范围
	CommonRange     string  `form:"commonRange" json:"commonRange"`                            // 公共服ID范围
	OneGameMem      float32 `form:"oneGameMem" json:"oneGameMem" binding:"required"`           // 单个游服进程加数据库最大占用内存(G)
	OneCrossMem     float32 `form:"oneCrossMem" json:"oneCrossMem" binding:"required"`         // 单个跨服进程加数据库最大占用内存(G)
	OneCommonMem    float32 `form:"oneCommonMem" json:"oneCommonMem"`                          // 单个公共服进程加数据库最大占用内存(G)
	CloudPlatform   string  `form:"cloudPlatform" json:"cloudPlatform" binding:"required"`     // 云平台
}

type GetProjectReq struct {
	ID              uint    `form:"id" json:"id"` // 修改才需要传，没有传算新增
	Name            string  `form:"name" json:"name"`
	BackendVersion  string  `form:"backendVersion" json:"backendVersion"`
	FrontendVersion string  `form:"frontendVersion" json:"frontendVersion"`
	GameRange       string  `form:"gameRange" json:"gameRange"`
	CrossRange      string  `form:"crossRange" json:"crossRange"`
	CommonRange     string  `form:"commonRange" json:"commonRange"`
	OneGameMem      float32 `form:"oneGameMem" json:"oneGameMem"`
	OneCrossMem     float32 `form:"oneCrossMem" json:"oneCrossMem"`
	OneCommonMem    float32 `form:"oneCommonMem" json:"oneCommonMem"`
	CloudPlatform   string  `form:"cloudPlatform" json:"cloudPlatform"`
}

type GetProjectsReq struct {
	ID            uint   `form:"id" json:"id"`
	Name          string `form:"name" json:"name"`
	CloudPlatform string `form:"cloudPlatform" json:"cloudPlatform"`
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
