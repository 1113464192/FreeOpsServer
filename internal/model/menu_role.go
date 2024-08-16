package model

type MenuRole struct {
	MenuId uint `json:"menuId" gorm:"index;comment:菜单ID"`
	RoleId uint `json:"roleId" gorm:"index;comment:角色ID"`
}
