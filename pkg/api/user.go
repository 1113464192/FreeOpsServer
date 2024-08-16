package api

import "FreeOps/internal/model"

type AuthLoginReq struct {
	Username string `form:"username" json:"username" binding:"required,min=4,max=16"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=18"`
}

// 用户结果返回
type UserRes struct {
	ID         uint   `json:"id"`
	Status     uint   `json:"status"`
	Username   string `json:"username"`
	Password   string `json:"password,omitempty"`
	UserGender uint   `json:"userGender"`
	Nickname   string `json:"nickname"`
	UserPhone  string `json:"userPhone"`
	UserEmail  string `json:"userEmail"`
}

type GetUsersReq struct {
	ID         uint   `json:"id" form:"id"`
	Status     uint   `json:"status" form:"status"`
	Username   string `json:"username" form:"username"`
	UserGender uint   `json:"userGender" form:"userGender"`
	Nickname   string `json:"nickname" form:"nickname"`
	UserPhone  string `json:"userPhone" form:"userPhone"`
	UserEmail  string `json:"userEmail" form:"userEmail"`
	PageInfo
}

// 用户列表返回并携带页码
type GetUsersRes struct {
	Records  []UserRes `json:"records" form:"records"`
	Page     int       `json:"current" form:"current"` // 页码
	PageSize int       `json:"size" form:"size"`       // 每页大小
	Total    int64     `json:"total"`
}

// 登录返回
type AuthLoginRes struct {
	//UserRes
	//RoleCodes    []string `json:"roleCodes"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

// 查询用户操作记录
type GetUserRecordLogsReq struct {
	Id   uint   `json:"id" form:"id" binding:"required"`         // 用户ID
	Date string `json:"string" form:"string" binding:"required"` // 年月，如：2006_01
	PageInfo
}

type GetUserRecordLogsRes struct {
	Logs     *[]model.UserRecord `json:"records" form:"records"`
	Page     int                 `json:"current" form:"current"` // 页码
	PageSize int                 `json:"size" form:"size"`       // 每页大小
	Total    int64               `json:"total"`
}

type GetUserPrivilegeRes struct {
	UserId   uint     `json:"userId" form:"userId"`
	Username string   `json:"username" form:"username"`
	Roles    []string `json:"roles" form:"roles"`
	Buttons  []string `json:"buttons" form:"buttons"`
}

type UpdateUserReq struct {
	ID         uint   `form:"id" json:"id"` // 修改才需要传，没有传算新增
	Status     uint   `form:"status" json:"status" binding:"required,oneof=1 2"`
	Username   string `form:"username" json:"username" binding:"required,min=4,max=16"`
	UserGender uint   `form:"userGender" json:"userGender" binding:"required"`
	Nickname   string `form:"nickname" json:"nickname" binding:"required"`
	UserPhone  string `form:"userPhone" json:"userPhone"`
	UserEmail  string `form:"userEmail" json:"userEmail"`
}

type ChangeUserPasswordReq struct {
	ID          uint   `form:"id" json:"id" binding:"required"`
	OldPassword string `form:"oldPassword" json:"oldPassword"`
	NewPassword string `form:"newPassword" json:"newPassword" binding:"required,min=6,max=18"`
}

type BindUserRolesReq struct {
	UserId  uint   `form:"userId" json:"userId" binding:"required"`
	RoleIds []uint `form:"roleIds" json:"roleIds" binding:"required"`
}
