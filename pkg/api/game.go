package api

type UpdateGameReq struct {
	ID             uint   `form:"id" json:"id"` // 修改才需要传，没有传算新增
	Name           string `form:"name" json:"name" binding:"required"`
	ServerId       uint   `form:"serverId" json:"serverId" binding:"required"`
	Type           uint8  `form:"type" json:"type" binding:"required"`
	Status         uint8  `form:"status" json:"status" binding:"required"`
	LbName         string `form:"lbName" json:"lbName" binding:"required"`
	LbListenerPort uint   `form:"lbListenerPort" json:"lbListenerPort" binding:"required"`
	ServerPort     uint   `form:"serverPort" json:"serverPort" binding:"required"`
	ProjectId      uint   `form:"projectId" json:"projectId" binding:"required"`
	HostId         uint   `form:"hostId" json:"hostId" binding:"required"`
	CrossId        uint   `form:"crossId" json:"crossId"`
	GlobalId       uint   `form:"globalId" json:"globalId"`
}

type GetGameReq struct {
	ID          uint   `form:"id" json:"id"`
	Name        string `form:"name" json:"name"`
	ServerId    uint   `form:"serverId" json:"serverId"`
	Type        uint8  `form:"type" json:"type"`
	Status      uint8  `form:"status" json:"status"`
	ProjectName string `form:"projectName" json:"projectName"`
	HostName    string `form:"hostName" json:"hostName"`
	CrossId     uint   `form:"crossId" json:"crossId"`
	GlobalId    uint   `form:"globalId" json:"globalId"`
}

type GetGamesReq struct {
	GetGameReq
	PageInfo
}

type GetGameRes struct {
	ID             uint   `form:"id" json:"id"`
	Name           string `form:"name" json:"name"`
	ServerId       uint   `form:"serverId" json:"serverId"`
	Type           uint8  `form:"type" json:"type"`
	Status         uint8  `form:"status" json:"status"`
	LbName         string `form:"lbName" json:"lbName"`
	LbListenerPort uint   `form:"lbListenerPort" json:"lbListenerPort"`
	ServerPort     uint   `form:"serverPort" json:"serverPort"`
	ProjectName    string `form:"projectName" json:"projectName"`
	HostName       string `form:"hostName" json:"hostName"`
	CrossId        uint   `form:"crossId" json:"crossId"`
	GlobalId       uint   `form:"globalId" json:"globalId"`
}

type GetGamesRes struct {
	Records  []GetGameRes `json:"records" form:"records"`
	Page     int          `json:"current" form:"current"` // 页码
	PageSize int          `json:"size" form:"size"`       // 每页大小
	Total    int64        `json:"total"`
}
