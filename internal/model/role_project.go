package model

type RoleProject struct {
	RoleId    uint `json:"roleId" gorm:"index;comment:角色ID"`
	ProjectId uint `json:"projectId" gorm:"index;comment:项目ID"`
}
