package model

type OpsParam struct {
	// id、提取变量内容的关键字、变量内容、映射到模板的变量名
	ID       uint   `json:"id" form:"id" gorm:"primarykey"`
	Keyword  string `json:"keyword" form:"keyword" gorm:"type: varchar(30);comment: 提取变量内容的关键字"`
	Variable string `json:"variable" form:"variable" gorm:"type: varchar(30);comment: 映射到模板的变量名"`
}
