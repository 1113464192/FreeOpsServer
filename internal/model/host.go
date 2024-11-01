package model

import (
	"gorm.io/gorm"
)

type Host struct {
	gorm.Model
	Name      string  `gorm:"type: varchar(30);index;comment: 服务器名"`
	Ipv4      string  `gorm:"type: varchar(30);uniqueIndex;comment: 如: 11.129.212.42"`
	Ipv6      *string `gorm:"type: varchar(100);uniqueIndex;comment: 如: 241d:c000:2022:601c:0:91aa:274c:e7ac/64"`
	Vip       string  `gorm:"type: varchar(30);uniqueIndex;comment: 内网IP"`
	SSHPort   uint16  `gorm:"type:smallint unsigned;comment: SSH端口"`
	Zone      string  `gorm:"type: varchar(100);comment: 服务器所在地区"`
	Cloud     string  `gorm:"comment: 云平台,自定义语言,与Project的CloudPlatform字符串一致"`
	System    string  `gorm:"type: varchar(30)"`
	Cores     uint16  `gorm:"comment: CPU核数"`
	DataDisk  uint32  `gorm:"comment: 数据盘, 单位为G"`
	Mem       uint64  `gorm:"comment: 单位为M"`
	ProjectId uint    `gorm:"comment: 项目ID;index"`
}
