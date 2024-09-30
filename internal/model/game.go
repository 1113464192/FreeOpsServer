package model

type Game struct {
	Id             uint    `gorm:"comment: 游戏服ID;index"`
	Type           uint8   `gorm:"comment: 1 游服 2 跨服 3 公共服  ...后续有需要再加;index"`
	Status         uint8   `gorm:"comment: 1 运行中 2 停服 3 已合服 4 操作中;index"`
	LbName         *string `gorm:"type: varchar(30);comment: 负载均衡名"`
	LbListenerPort *uint   `gorm:"comment: 负载均衡监听器的端口"`
	ServerPort     uint    `gorm:"comment: 游戏服务端口"`
	ProjectId      uint    `gorm:"comment: 项目ID;index"`
	HostId         uint    `gorm:"comment: 服务器ID;index"`
	CrossId        *uint   `gorm:"comment: 关联跨服ID;index"`
	CommonId       *uint   `gorm:"comment: 关联公共服ID;index"`
}
