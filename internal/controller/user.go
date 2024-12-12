package controller

import (
	"FreeOps/internal/consts"
	"FreeOps/internal/model"
	"FreeOps/internal/service"
	"FreeOps/pkg/api"
	"FreeOps/pkg/logger"
	"FreeOps/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

// UserLogin
// @Tags 公共相关
// @title 用户登录
// @description 用户名长度不少于4位，密码不少于6位
// @Summary 用户登录
// @Produce  application/json
// @Param data formData api.AuthLoginReq true "用户名, 密码"
// @Success 200 {object} api.Response "{"code": "0000", msg: "string", data: "string"}"
// @Failure 500 {object} api.Response "{"code": "", msg: "", data: ""}"
// @Router /api/auth/login [post]
func UserLogin(c *gin.Context) {
	var loginReq api.AuthLoginReq
	if err := c.ShouldBind(&loginReq); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	u := &model.User{Username: loginReq.Username, Password: loginReq.Password}
	// 创建账号数量的变量
	var count int64
	response, err := service.UserServiceApp().Login(u)
	if err != nil {
		if err2 := model.DB.Model(&model.User{}).Where("username = ?", u.Username).Count(&count).Error; err2 != nil {
			logger.Log().Error("user", "查询用户账号失败", err2)
			c.JSON(500, util.ServerErrorResponse("获取用户数量失败", err2))
			return
		}
		var errMessage string
		if count == 0 {
			errMessage = "没有这个账号"
		} else {
			errMessage = fmt.Sprintf("登陆失败, %v", err)
		}
		c.JSON(500, util.ServerErrorResponse(errMessage, err))
		return
	}
	logger.Log().Info("user", "登录成功", "用户名: "+loginReq.Username)
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "登录成功",
		Data: response,
	})
	return
}

// RefreshToken
// @Tags 公共相关
// @title 刷新Token
// @description refreshToken放在data请求
// @Summary 刷新Token
// @Produce  application/json
// @Param refreshToken formData string true "refreshToken"
// @Success 200 {object} api.Response "{"code": "0000", msg: "string", data: "string"}"
// @Failure 500 {object} api.Response "{"code": "", msg: "", data: ""}"
// @Router /api/auth/refreshToken [post]
func RefreshToken(c *gin.Context) {
	var requestBody map[string]string
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}

	refreshToken, exists := requestBody["refreshToken"]
	if !exists {
		c.JSON(500, api.Response{
			Code: "500",
			Msg:  "refreshToken不存在",
			Data: nil,
		})
		return
	}
	res, code, err := service.UserServiceApp().RefreshToken(refreshToken)
	var msg string
	var data any
	if err == nil {
		msg = "刷新Token成功"
		data = res
	} else {
		msg = fmt.Sprint("刷新Token失败: ", err)
		data = map[string]error{
			"错误信息": err,
		}
		logger.Log().Error("user", "刷新Token失败", err)
	}
	logger.Log().Info("user", "刷新Token成功")
	c.JSON(200, api.Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

// UserLogout
// @Tags 用户相关
// @title 用户登出
// @description 登出 - 把JWT拉入黑名单
// @Summary 登出
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Success 200 {object} api.Response "{"code": "0000", msg: "string", data: "string"}"
// @Failure 500 {object} api.Response "{"code": "", msg: "", data: ""}"
// @Router /api/users/logout [post]
func UserLogout(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	token := parts[1]
	jwt := &model.JwtBlacklist{Jwt: token}
	if err := service.JwtServiceApp().JwtAddBlacklist(jwt); err != nil {
		logger.Log().Error("user", "jwt没有拉入黑名单", err)
		c.JSON(500, util.ServerErrorResponse("jwt没有成功拉入黑名单", err))
		return
	}
	logger.Log().Info("user", "jwt拉入黑名单成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "用户登出成功",
	})
}

// UpdateUser
// @Tags 用户相关
// @title 新增/修改用户信息
// @description 新增不用传用户ID，修改才传用户ID，返回用户密码
// @Summary 新增/修改用户信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData api.UpdateUserReq true "传新增或者修改用户的所需参数"
// @Success 200 {object} api.Response "{"code": "0000", msg: "string", data: "string"}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"code": "", msg: "", data: ""}"
// @Router /api/users [post]
func UpdateUser(c *gin.Context) {
	var (
		userReq api.UpdateUserReq
		err     error
	)
	if err = c.ShouldBind(&userReq); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	passwd, err := service.UserServiceApp().UpdateUser(&userReq)
	if err != nil {
		logger.Log().Error("user", "添加/修改用户失败", err)
		if err.Error() == "用户密码bcrypt加密失败" {
			c.JSON(500, util.ServerErrorResponse("用户密码bcrypt加密失败", err))
			return
		}
		c.JSON(500, util.ServerErrorResponse("添加/修改用户失败", err))
		return
	}

	logger.Log().Info("user", "添加/修改用户成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: passwd,
	})
}

// GetUsers
// @Tags 用户相关
// @title 用户列表
// @description 获取用户列表
// @Summary 获取用户列表
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetUsersReq true "所需参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/users [get]
func GetUsers(c *gin.Context) {
	var params api.GetUsersReq
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}

	res, err := service.UserServiceApp().GetUsers(params)
	if err != nil {
		logger.Log().Error("user", "获取用户列表失败", err)
		c.JSON(500, util.ServerErrorResponse("获取用户列表失败", err))
		return
	}

	logger.Log().Info("user", "获取用户列表成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}

// GetUserPrivilege
// @Tags 用户相关
// @title 用户权限
// @description 获取用户权限
// @Summary 用户权限
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/users/privilege [get]
func GetUserPrivilege(c *gin.Context) {
	user, err := util.GetClaimsUser(c)
	if err != nil {
		c.JSON(200, api.Response{
			Code: consts.SERVICE_MODAL_LOGOUT_CODE,
			Msg:  err.Error(),
		})
		return
	}
	roles, err := util.GetClaimsRole(c)
	if err != nil {
		c.JSON(200, api.Response{
			Code: consts.SERVICE_MODAL_LOGOUT_CODE,
			Msg:  err.Error(),
		})
		return
	}

	res, err := service.UserServiceApp().GetUserPrivilege(user, roles)
	if err != nil {
		logger.Log().Error("user", "获取用户权限失败", err)
		c.JSON(500, util.ServerErrorResponse("获取用户权限失败", err))
		return
	}

	logger.Log().Info("user", "获取用户权限成功", fmt.Sprintf("用户ID: %d", user.ID))
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}

// DeleteUsers
// @Tags 用户相关
// @title 删除用户
// @description 删除指定用户
// @Summary 删除用户
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param ids body api.IdsReq true "用户ID"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/users/ [delete]
func DeleteUsers(c *gin.Context) {
	var (
		param api.IdsReq
		err   error
	)
	if err = c.ShouldBind(&param); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}

	user, err := util.GetClaimsUser(c)
	if err != nil {
		c.JSON(200, api.Response{
			Code: consts.SERVICE_MODAL_LOGOUT_CODE,
			Msg:  err.Error(),
		})
		return
	}

	for _, id := range param.Ids {
		if user.ID == id {
			c.JSON(500, api.Response{
				Code: consts.SERVICE_ERROR_CODE,
				Msg:  "不能删除自己",
			})
			return
		}

	}

	if err = service.UserServiceApp().DeleteUsers(param.Ids); err != nil {
		logger.Log().Error("user", "删除用户失败", err)
		c.JSON(500, util.ServerErrorResponse("删除用户失败", err))
		return
	}

	logger.Log().Info("user", "删除用户成功", fmt.Sprintf("用户ID: %d", param.Ids))
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
	})
}

// ChangeUserPassword
// @Tags 用户相关
// @title 修改密码
// @description 修改指定用户密码
// @Summary 修改密码
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData api.ChangeUserPasswordReq true "传删除用户的所需参数,管理员操作可不传旧密码"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/users/password [patch]
func ChangeUserPassword(c *gin.Context) {
	var (
		param   api.ChangeUserPasswordReq
		err     error
		isAllow bool
	)
	if err = c.ShouldBind(&param); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	// 判断是否管理员操作
	if isAllow, err = util.IsSelfAdmin(c); err != nil {
		c.JSON(200, api.Response{
			Code: consts.SERVICE_MODAL_LOGOUT_CODE,
			Msg:  err.Error(),
		})
		return
	}

	// 非管理员则判断是否本用户自己操作
	if !isAllow {
		if isAllow, err = util.IsSelf(c, param.ID); err != nil {
			c.JSON(200, api.Response{
				Code: consts.SERVICE_MODAL_LOGOUT_CODE,
				Msg:  err.Error(),
			})
			return
		}
	}

	if !isAllow {
		c.JSON(500, api.Response{
			Code: consts.SERVICE_ERROR_CODE,
			Msg:  "没有权限修改他人密码",
		})
		return
	}

	if err = service.UserServiceApp().ChangeUserPassword(param); err != nil {
		logger.Log().Error("user", "修改用户密码失败", err)
		c.JSON(500, util.ServerErrorResponse("修改用户密码失败", err))
		return
	}

	logger.Log().Info("user", "修改用户密码成功", fmt.Sprintf("用户ID: %d", param.ID))
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
	})
}

// GetUserRecordDate
// @Tags 用户相关
// @title 存在记录的月份
// @description 获取存在记录的月份
// @Summary 存在记录的月份
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/users/history-month-exist [get]
func GetUserRecordDate(c *gin.Context) {
	dates, err := service.UserRecordApp().GetUserRecordDate()
	if err != nil {
		logger.Log().Error("user", "获取存在记录的月份失败", err)
		c.JSON(500, util.ServerErrorResponse("获取存在记录的月份失败", err))
		return
	}
	logger.Log().Info("user", "获取存在记录的月份成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: map[string][]string{
			"dates": dates,
		},
	})
}

// GetUserRecordLogs
// @Tags 用户相关
// @title 月份操作记录
// @description 查询月份操作记录
// @Summary 月份操作记录
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetUserRecordLogsReq true "所需参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/users/history-action [get]
func GetUserRecordLogs(c *gin.Context) {
	var params api.GetUserRecordLogsReq
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	logs, total, err := service.UserRecordApp().GetUserRecordLogs(params)
	if err != nil {
		logger.Log().Error("user", "查询月份操作记录失败", err)
		c.JSON(500, util.ServerErrorResponse("查询月份操作记录失败", err))
		return
	}
	res := api.GetUserRecordLogsRes{
		Logs:     logs,
		Total:    total,
		Page:     params.PageInfo.Page,
		PageSize: params.PageInfo.PageSize,
	}
	logger.Log().Info("user", "查询月份操作记录成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}

// BindUserRoles
// @Tags 用户相关
// @title 用户绑定角色
// @description 用户绑定角色
// @Summary 用户绑定角色
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.BindUserRolesReq true "传用户ID与roleCodes"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"修改权限成功，刷新Token"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/users/bind-roles [put]
func BindUserRoles(c *gin.Context) {
	var params api.BindUserRolesReq
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	if err := service.UserServiceApp().BindUserRoles(params.UserId, params.RoleIds); err != nil {
		logger.Log().Error("user", "绑定用户角色失败", err)
		c.JSON(500, util.ServerErrorResponse("绑定用户角色失败", err))
		return
	}

	logger.Log().Info("user", "绑定用户角色成功", fmt.Sprintf("用户ID: %d————角色IDs: %d", params.UserId, params.RoleIds))
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "修改权限成功，刷新Token",
	})
}

// GetUserRoles
// @Tags 用户相关
// @title 查询用户的角色
// @description 查询用户的角色
// @Summary 查询用户的角色
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param uid query uint true "用户ID"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/users/roles [get]
func GetUserRoles(c *gin.Context) {
	userId, err := strconv.ParseUint(c.Query("uid"), 10, 0)
	if err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}

	roles, err := service.UserServiceApp().GetUserRoles(uint(userId))
	if err != nil {
		logger.Log().Error("user", "获取用户角色失败", err)
		c.JSON(500, util.ServerErrorResponse("获取用户角色失败", err))
		return
	}

	res, err := service.RoleServiceApp().GetResults(&roles)
	if err != nil {
		logger.Log().Error("user", "获取用户角色失败", err)
		c.JSON(500, util.ServerErrorResponse("获取用户角色失败", err))
		return
	}

	logger.Log().Info("user", "获取用户角色成功", fmt.Sprintf("用户ID: %d", userId))
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: map[string]any{
			"records": res,
		},
	})
}

// GetUserProjectOptions
// @Tags 用户相关
// @title 自身权限的项目选项
// @description 返回权限内可以选择的项目的选项
// @Summary 自身权限的项目选型
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/users/project-options [get]
func GetUserProjectOptions(c *gin.Context) {
	res, err := service.UserServiceApp().GetUserProjectOptions(c)
	if err != nil {
		logger.Log().Error("user", "获取用户的项目选项失败", err)
		c.JSON(500, util.ServerErrorResponse("获取用户的项目选项失败", err))
		return
	}

	logger.Log().Info("user", "获取用户的项目选项成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}
