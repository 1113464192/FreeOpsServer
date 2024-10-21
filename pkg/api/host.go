package api

type UpdateHostReq struct {
	ID        uint   `form:"id" json:"id"`                        // 修改才需要传，没有传算新增
	Name      string `form:"name" json:"name" binding:"required"` // 服务器名, 不允许修改
	Ipv4      string `form:"ipv4" json:"ipv4" binding:"required"`
	Ipv6      string `form:"ipv6" json:"ipv6"`
	Vip       string `form:"vip" json:"vip" binding:"required"`
	Zone      string `form:"zone" json:"zone" binding:"required"`
	Cloud     string `form:"cloud" json:"cloud" binding:"required"`
	System    string `form:"system" json:"system" binding:"required"`
	Cores     uint16 `form:"cores" json:"cores" binding:"required"`
	DataDisk  uint32 `form:"dataDisk" json:"dataDisk"`          // 数据盘,单位: G
	Mem       uint64 `form:"mem" json:"mem" binding:"required"` // 内存,单位: G
	ProjectId uint   `form:"projectId" json:"projectId" binding:"required"`
}

type GetHostReq struct {
	ID          uint   `form:"id" json:"id"` // 修改才需要传，没有传算新增
	Name        string `form:"name" json:"name"`
	Ipv4        string `form:"ipv4" json:"ipv4"`
	Ipv6        string `form:"ipv6" json:"ipv6"`
	Vip         string `form:"vip" json:"vip"`
	Zone        string `form:"zone" json:"zone"`
	Cloud       string `form:"cloud" json:"cloud"`
	System      string `form:"system" json:"system"`
	ProjectName string `form:"projectName" json:"projectName"`
}

type GetHostsReq struct {
	GetHostReq
	PageInfo
}

type GetHostRes struct {
	ID          uint   `form:"id" json:"id"`     // 修改才需要传，没有传算新增
	Name        string `form:"name" json:"name"` // 服务器名, 不允许修改
	Ipv4        string `form:"ipv4" json:"ipv4"`
	Ipv6        string `form:"ipv6" json:"ipv6"`
	Vip         string `form:"vip" json:"vip"`
	Zone        string `form:"zone" json:"zone"`
	Cloud       string `form:"cloud" json:"cloud"`
	System      string `form:"system" json:"system"`
	Cores       uint16 `form:"cores" json:"cores"`
	DataDisk    uint32 `form:"dataDisk" json:"dataDisk"`
	Mem         uint64 `form:"mem" json:"mem"`
	ProjectName string `form:"projectName" json:"projectName"`
	ProjectId   uint   `form:"projectId" json:"projectId"`
	GetHostGameInfoRes
}

type GetHostsRes struct {
	Records  []GetHostRes `json:"records" form:"records"`
	Page     int          `json:"current" form:"current"` // 页码
	PageSize int          `json:"size" form:"size"`       // 每页大小
	Total    int64        `json:"total"`
}

type GetHostGameInfoRes struct {
	GameTotal   int64 `json:"gameTotal" form:"gameTotal"`     // 游服总数
	CrossTotal  int64 `json:"crossTotal" form:"crossTotal"`   // 跨服总数
	CommonTotal int64 `json:"commonTotal" form:"commonTotal"` // 公共服总数
}

type GetHostListRes struct {
	ID   uint   `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
	Ipv4 string `json:"ipv4" form:"ipv4"`
	Ipv6 string `json:"ipv6,omitempty" form:"ipv6"`
}
