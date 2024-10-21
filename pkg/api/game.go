package api

type UpdateGameReq struct {
	Id             uint   `form:"id" json:"id" binding:"required"`
	Type           uint8  `form:"type" json:"type" binding:"required"`
	Status         uint8  `form:"status" json:"status" binding:"required"`
	LbName         string `form:"lbName" json:"lbName"`
	LbListenerPort uint   `form:"lbListenerPort" json:"lbListenerPort"`
	ServerPort     uint   `form:"serverPort" json:"serverPort" binding:"required"`
	ProjectId      uint   `form:"projectId" json:"projectId" binding:"required"`
	HostId         uint   `form:"hostId" json:"hostId" binding:"required"`
	CrossId        uint   `form:"crossId" json:"crossId"`
	CommonId       uint   `form:"commonId" json:"commonId"`
	ActionType     uint8  `form:"actionType" json:"actionType" binding:"required"` // 1: 创建 2: 更新
}

type UpdateGameStatusReq struct {
	Id     uint  `form:"id" json:"id" binding:"required"`
	Status uint8 `form:"status" json:"status" binding:"required"`
}

type GetGameReq struct {
	Id          uint   `form:"id" json:"id"`
	Type        uint8  `form:"type" json:"type"`
	Status      uint8  `form:"status" json:"status"`
	ProjectName string `form:"projectName" json:"projectName"`
	HostName    string `form:"hostName" json:"hostName"`
	Ipv4        string `form:"ipv4" json:"ipv4"`
	CrossId     uint   `form:"crossId" json:"crossId"`
	CommonId    uint   `form:"commonId" json:"commonId"`
}

type GetGamesReq struct {
	GetGameReq
	PageInfo
}

type GetGameRes struct {
	Id             uint   `form:"id" json:"id"`
	Type           uint8  `form:"type" json:"type"`
	Status         uint8  `form:"status" json:"status"`
	LbName         string `form:"lbName" json:"lbName"`
	LbListenerPort uint   `form:"lbListenerPort" json:"lbListenerPort"`
	ServerPort     uint   `form:"serverPort" json:"serverPort"`
	ProjectName    string `form:"projectName" json:"projectName"`
	HostName       string `form:"hostName" json:"hostName"`
	Ipv4           string `form:"ipv4" json:"ipv4"`
	ProjectId      uint   `form:"projectId" json:"projectId"`
	HostId         uint   `form:"hostId" json:"hostId"`
	CrossId        uint   `form:"crossId" json:"crossId"`
	CommonId       uint   `form:"commonId" json:"commonId"`
}

type GetGamesRes struct {
	Records  []GetGameRes `json:"records" form:"records"`
	Page     int          `json:"current" form:"current"` // 页码
	PageSize int          `json:"size" form:"size"`       // 每页大小
	Total    int64        `json:"total"`
}
