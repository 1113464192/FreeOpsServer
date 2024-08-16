package controller

import (
	"FreeOps/internal/consts"
	"FreeOps/internal/service"
	"FreeOps/pkg/api"
	"FreeOps/pkg/logger"
	"FreeOps/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
)

// UpdateRole
// @Tags 角色相关
// @title 新增/修改角色信息
// @description 新增不用传ID，修改才传ID
// @Summary 新增/修改角色信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData api.UpdateRoleReq true "传新增/修改角色的参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /roles [post]
func UpdateRole(c *gin.Context) {
	var (
		params api.UpdateRoleReq
		err    error
	)
	if err = c.ShouldBind(&params); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	if err = service.RoleServiceApp().UpdateRole(&params); err != nil {
		logger.Log().Error("role", "创建/修改组失败", err)
		c.JSON(500, util.ServerErrorResponse("创建/修改组失败", err))
		return
	}

	logger.Log().Info("role", "创建/修改组成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
	})
}

// GetRoles
// @Tags 角色相关
// @title 查询角色信息
// @description 查询所有/指定条件角色的信息
// @Summary 查询角色信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetRolesReq true "传对应条件的参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /roles [get]
func GetRoles(c *gin.Context) {
	var params api.GetRolesReq
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	res, err := service.RoleServiceApp().GetRoles(&params)
	if err != nil {
		logger.Log().Error("role", "查询角色信息失败", err)
		c.JSON(500, util.ServerErrorResponse("查询角色信息失败", err))
		return
	}

	logger.Log().Info("role", "查询角色信息成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}

// GetAllRolesSummary
// @Tags 角色相关
// @title 所有角色简略信息
// @description 查询所有角色的简略信息
// @Summary 所有角色简略信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌" "
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /roles/all-summary [get]
func GetAllRolesSummary(c *gin.Context) {
	roles, err := service.RoleServiceApp().GetAllRolesSummary()
	if err != nil {
		logger.Log().Error("role", "查询所有角色的简略信息失败", err)
		c.JSON(500, util.ServerErrorResponse("查询所有角色的简略信息失败", err))
		return
	}

	logger.Log().Info("role", "查询所有角色的简略信息成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: roles,
	})
}

// DeleteRoles
// @Tags 角色相关
// @title 删除角色
// @description 删除指定角色
// @Summary 删除角色
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param ids body api.IdsReq true "角色ID"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /roles/ [delete]
func DeleteRoles(c *gin.Context) {
	var param api.IdsReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	if err := service.RoleServiceApp().DeleteRoles(param.Ids); err != nil {
		logger.Log().Error("role", "删除角色失败", err)
		c.JSON(500, util.ServerErrorResponse("删除角色失败", err))
		return
	}

	logger.Log().Info("role", "删除角色成功", fmt.Sprintf("角色ID: %d", param.Ids))
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
	})
}

// BindRoleRelation
// @Tags 角色相关
// @title 关联角色关系
// @description 1: api 2: menu 3: button
// @Summary 关联角色关系
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.BindRoleRelationReq true "associationType参数 1: api 2: menu 3: button"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /roles/bind [put]
func BindRoleRelation(c *gin.Context) {
	var param api.BindRoleRelationReq
	var err error
	if err = c.ShouldBind(&param); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	if err = service.RoleServiceApp().BindRoleRelation(param); err != nil {
		logger.Log().Error("role", "绑定角色关系失败", err)
		c.JSON(500, util.ServerErrorResponse("绑定角色关系失败", err))
		return
	}

	logger.Log().Info("role", "绑定角色关系成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
	})
}

// GetRoleMenus
// @Tags 角色相关
// @title 获取角色的关联菜单列表
// @description 获取角色的关联菜单列表
// @Summary 获取角色的关联菜单列表
// @Produce   application/json
// @Param Authorization header string true "格式为：Bearer 登录返回的用户令牌"
// @Param data query api.IdsReq true "所需参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /roles/menus [get]
func GetRoleMenus(c *gin.Context) {
	var params api.IdsReq
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	res, err := service.RoleServiceApp().GetRoleMenus(params)
	if err != nil {
		logger.Log().Error("role", "获取角色菜单失败", err)
		c.JSON(500, util.ServerErrorResponse("获取角色菜单失败", err))
		return
	}

	logger.Log().Info("role", "获取角色菜单成功", fmt.Sprintf("角色ID: %d", params.Ids))
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}

// GetRoleApis
// @Tags 角色相关
// @title 获取角色的API权限列表
// @description 由于swagger本身的限制，get请求的切片会报错，并非接口本身问题，请换个方式，如http://127.0.0.1:9081/api/v1/group/apis?ids=3&ids=4
// @Summary 获取角色的API权限列表
// @Produce   application/json
// @Param Authorization header string true "格式为：Bearer 登录返回的用户令牌"
// @Param data query api.IdsReq true "传角色ID切片"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /roles/apis [get]
func GetRoleApis(c *gin.Context) {
	var param api.IdsReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	res, err := service.CasbinServiceApp().GetPolicyPathByGroupIds(param.Ids)
	if err != nil {
		logger.Log().Error("role", "获取角色API失败", err)
		c.JSON(500, util.ServerErrorResponse("获取角色API失败", err))
		return
	}

	logger.Log().Info("role", "获取角色API成功", fmt.Sprintf("角色ID: %d", param.Ids))
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}

// GetRoleUsers
// @Tags 角色相关
// @title 获取角色的关联用户列表
// @description 获取角色的关联用户列表
// @Summary 获取角色的关联用户列表
// @Produce   application/json
// @Param Authorization header string true "格式为：Bearer 登录返回的用户令牌"
// @Param data query api.IdPageReq true "所需参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /roles/users [get]
func GetRoleUsers(c *gin.Context) {
	var params api.IdPageReq
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	res, err := service.RoleServiceApp().GetRoleUsers(params)
	if err != nil {
		logger.Log().Error("role", "获取角色用户失败", err)
		c.JSON(500, util.ServerErrorResponse("获取角色用户失败", err))
		return
	}

	logger.Log().Info("role", "获取角色用户成功", fmt.Sprintf("角色ID: %d", params.Id))
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}
