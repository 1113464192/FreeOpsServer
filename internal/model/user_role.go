package model

type UserRole struct {
	UserId uint `json:"userId" gorm:"index;comment:用户ID"`
	RoleId uint `json:"roleId" gorm:"index;comment:角色ID"`
}
