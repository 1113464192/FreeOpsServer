package consts

import "FreeOps/internal/model"

// 拥有软删除字段的model
var SoftDeleteModelList = []interface{}{
	&model.User{},
	&model.Role{},
	&model.Menu{},
	&model.Api{},
	&model.JwtBlacklist{},
	&model.UserRecord{},
	&model.Host{},
	&model.Project{},
	&model.Game{},
	&model.OpsTemplate{},
	&model.OpsTask{},
	&model.OpsTaskLog{},
}

// MysqlTableName
const (
	MysqlTableNameUser        = "user"
	MysqlTableNameRole        = "role"
	MysqlTableNameMenu        = "menu"
	MysqlTableNameApi         = "api"
	MysqlTableNameButton      = "button"
	MysqlTableNameProject     = "project"
	MysqlTableNameJwt         = "jwt_blacklist"
	MysqlTableNameUserRecord  = "act_record"
	MysqlTableNameOpsTemplate = "ops_template"
	MysqlTableNameOpsParam    = "ops_param"
	MysqlTableNameOpsTask     = "ops_task"
	MysqlTableNameOpsTaskLog  = "ops_task_logs"

	// GORM默认Bool类型的True是1，False是0
	MysqlGormBoolIsTrue  = 1
	MysqlGormBoolIsFalse = 0
)

// 角色常量
const (
	RoleModelAdminCode = "ADMIN"
)

// 菜单常量
const (
	MenuModelMenuTypeIsDirectory = 1
	MenuModelMenuTypeIsMenu      = 2
	MenuModeIconTypeIsIconify    = 1
	MenuModeIconTypeIsLocal      = 2
	MenuModelPropsIsTrue         = "true"
	ManualComponentMenuPath      = "/document/"
)

// 用户常量
const (
	UserModelStatusEnabled      = 1
	UserModelUserGenderIsMale   = 1
	UserModelUserGenderIsFemale = 2
)

// 游戏常量
const (
	GameModeTypeIsGame      = 1
	GameModelTypeIsCross    = 2
	GameModelTypeIsCommon   = 3
	GameModelStatusIsMerged = 3
	ActionTypeIsCreate      = 1
	ActionTypeIsUpdate      = 2
)

// 运维任务常量
const (
	OpsTaskStatusIsWaiting = iota
	OpsTaskStatusIsRunning
	OpsTaskStatusIsSuccess
	OpsTaskStatusIsFailed
	OpsTaskStatusIsRejected
)
