package model

type RoleButton struct {
	RoleId   uint `json:"roleId" gorm:"index;comment:角色ID"`
	ButtonId uint `json:"buttonId" gorm:"index;comment:按钮ID"`
}
