package model

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	Name            string  `json:"name" gorm:"type:varchar(20);uniqueIndex;comment:项目名"`
	BackendVersion  string  `json:"backendVersion" gorm:"type:varchar(50);comment:后端版本"`
	FrontendVersion string  `json:"frontendVersion" gorm:"type:varchar(50);comment:前端版本"`
	GameRange       string  `json:"gameRange" gorm:"type:varchar(50);comment:游服ID范围"`
	CrossRange      string  `json:"crossRange" gorm:"type:varchar(50);comment:跨服ID范围"`
	CommonRange     string  `json:"commonRange" gorm:"type:varchar(50);comment:公共服ID范围"`
	OneGameMem      float32 `json:"oneGameMem" gorm:"type:FLOAT;comment:单个游服进程加数据库最大占用内存(G)"`
	OneCrossMem     float32 `json:"oneCrossMem" gorm:"type:FLOAT;comment:单个跨服进程加数据库最大占用内存(G)"`
	OneCommonMem    float32 `json:"oneCommonMem" gorm:"type:FLOAT;comment:单个公共服进程加数据库最大占用内存(G)"`
	// 拆分开，方便后续服务器维护，以及不同云商操作隔离
	CloudPlatform string `json:"cloudPlatform" gorm:"type:varchar(20);comment:云平台,自定义语言,与host的Cloud字符串一致"`
}
