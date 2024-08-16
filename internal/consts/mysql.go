package consts

import "FreeOps/internal/model"

// 拥有软删除字段的model
var SoftDeleteModelList = []interface{}{
	&model.User{},
	&model.Role{},
	&model.Menu{},
	&model.Api{},
	&model.JwtBlacklist{},
}

// MysqlTableName
const (
	MysqlTableNameUser       = "user"
	MysqlTableNameRole       = "role"
	MysqlTableNameMenu       = "menu"
	MysqlTableNameApi        = "api"
	MysqlTableNameButton     = "button"
	MysqlTableNameJwt        = "jwt_blacklist"
	MysqlTableNameUserRecord = "act_record"
)

// 角色常量
const (
	RoleModelAdminCode = "ADMIN"
)

// 菜单常量
const (
	MenuModelMenuTypeIsDirectory  = 1
	MenuModelMenuTypeIsMenu       = 2
	MenuModeIconTypeIsIconify     = 1
	MenuModeIconTypeIsLocal       = 2
	MenuModelHideInMenuIsYes      = 1
	MenuModelMultiTabIsYes        = 1
	MenuModelIsConstantRouteIsYes = 1
	MenuModelPropsIsTrue          = "true"
)

// 用户常量
const (
	UserModelStatusEnabled      = 1
	UserModelUserGenderIsMale   = 1
	UserModelUserGenderIsFemale = 2
)
