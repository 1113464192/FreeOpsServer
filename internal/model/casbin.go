package model

type CasbinRule struct {
	RoleId string `gorm:"column:v0" json:"roleId"`
	Path   string `gorm:"column:v1" json:"path"`
	Method string `gorm:"column:v2" json:"method"`
}

func (c CasbinRule) TableName() string {
	return "casbin_rule"
}
