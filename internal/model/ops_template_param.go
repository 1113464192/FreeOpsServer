package model

type OpsTemplateParam struct {
	TemplateId uint `gorm:"comment: 模板ID;index"`
	ParamId    uint `gorm:"comment: 参数ID;index"`
}
