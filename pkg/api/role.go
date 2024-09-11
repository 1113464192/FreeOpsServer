package api

type UpdateRoleReq struct {
	ID       uint   `form:"id" json:"id"` // 修改才需要传，没有传算新增
	RoleName string `form:"roleName" json:"roleName" binding:"required"`
	RoleCode string `form:"roleCode" json:"roleCode" binding:"required"`
	RoleDesc string `form:"roleDesc" json:"roleDesc"`
}

type GetRolesReq struct {
	ID       uint   `form:"id" json:"id"`
	RoleName string `form:"roleName" json:"roleName"`
	RoleCode string `form:"roleCode" json:"roleCode"`
	PageInfo
}

type GetRolesRes struct {
	Records  []UpdateRoleReq `json:"records" form:"records"`
	Page     int             `json:"current" form:"current"` // 页码
	PageSize int             `json:"size" form:"size"`       // 每页大小
	Total    int64           `json:"total"`
}

type GetAllRolesSummaryRes struct {
	ID       uint   `form:"id" json:"id"`
	RoleName string `form:"roleName" json:"roleName"`
	RoleCode string `form:"roleCode" json:"roleCode"`
}

type BindRoleRelationReq struct {
	RoleId          uint   `json:"roleId"  binding:"required"`                                 // 角色id
	AssociationType uint8  `form:"associationType" json:"associationType"  binding:"required"` // 1: api 2: menu 3: button 4: project
	ObjectIds       []uint `form:"objectIds" json:"objectIds"  binding:"required"`
}
