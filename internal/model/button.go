package model

type Button struct {
	ID         uint   `gorm:"primarykey"`
	ButtonCode string `json:"buttonCode" gorm:"type:varchar(30);comment:按钮代码;index"`
	ButtonDesc string `json:"buttonDesc" gorm:"type:varchar(255);comment:按钮描述"`
	MenuId     uint   `json:"menuId" gorm:"comment:菜单ID"`
}
