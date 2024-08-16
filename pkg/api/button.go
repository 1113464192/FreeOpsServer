package api

type UpdateButtonReq struct {
	ButtonCode string `form:"buttonCode" json:"buttonCode" binding:"required"`
	ButtonDesc string `form:"buttonDesc" json:"buttonDesc"`
	MenuId     uint   `form:"menuId" json:"menuId" binding:"required"`
}

type UpdateButtonsReq struct {
	Buttons []UpdateButtonReq `json:"buttons" binding:"required"`
}

type GetButtonReq struct {
	ID         uint   `form:"id" json:"id"` // 修改才需要传，没有传算新增
	ButtonCode string `form:"buttonCode" json:"buttonCode"`
	ButtonDesc string `form:"buttonDesc" json:"buttonDesc"`
	MenuId     uint   `form:"menuId" json:"menuId"`
}

type GetButtonsReq struct {
	GetButtonReq
	PageInfo
}

type GetButtonsRes struct {
	Records  []GetButtonReq `json:"records" form:"records"`
	Page     int            `json:"current" form:"current"` // 页码
	PageSize int            `json:"size" form:"size"`       // 每页大小
	Total    int64          `json:"total"`
}
