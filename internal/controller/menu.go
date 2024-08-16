package controller

import (
	"FreeOps/internal/consts"
	"FreeOps/internal/service"
	"FreeOps/pkg/api"
	"FreeOps/pkg/logger"
	"FreeOps/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

// UpdateMenu
// @Tags 菜单相关
// @title 新增/修改菜单信息
// @description 新增不用传ID，修改才传ID
// @Summary 新增/修改菜单信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.UpdateMenuReq true "传新增/修改菜单的参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /menus [post]
func UpdateMenu(c *gin.Context) {
	var (
		params api.UpdateMenuReq
		err    error
	)
	if err = c.ShouldBind(&params); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	if err = service.MenuServiceApp().UpdateMenu(&params); err != nil {
		logger.Log().Error("menu", "创建/修改菜单失败", err)
		c.JSON(500, util.ServerErrorResponse("创建/修改菜单失败", err))
		return
	}

	logger.Log().Info("menu", "创建/修改菜单成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
	})
}

// GetMenus
// @Tags 菜单相关
// @title 查询菜单信息
// @description 查询菜单信息
// @Summary 查询菜单信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.IdPageReq false "传查询菜单的参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /menus [get]
func GetMenus(c *gin.Context) {
	var param api.IdPageReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	res, err := service.MenuServiceApp().GetMenus(&param)
	if err != nil {
		logger.Log().Error("menu", "查询菜单失败", err)
		c.JSON(500, util.ServerErrorResponse("查询菜单失败", err))
		return
	}

	logger.Log().Info("menu", "查询菜单成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}

// DeleteMenu
// @Tags 菜单相关
// @title 删除菜单
// @description 删除菜单信息
// @Summary 删除菜单
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.IdsReq true "删除菜单的参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /menus [delete]
func DeleteMenu(c *gin.Context) {
	var param api.IdsReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	err := service.MenuServiceApp().DeleteMenus(param.Ids)
	if err != nil {
		logger.Log().Error("menu", "查询菜单失败", err)
		c.JSON(500, util.ServerErrorResponse("查询菜单失败", err))
		return
	}

	logger.Log().Info("menu", "删除菜单成功", fmt.Sprintf("ID:%v", param.Ids))
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
	})
}

// GetMenuButtons
// @Tags 菜单相关
// @title 获取菜单的关联按钮列表
// @description 获取菜单的关联按钮列表
// @Summary 获取菜单的关联按钮列表
// @Produce   application/json
// @Param Authorization header string true "格式为：Bearer 登录返回的用户令牌"
// @Param id query uint true "菜单ID"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /menus/buttons [get]
func GetMenuButtons(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 0)
	if err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	res, err := service.MenuServiceApp().GetMenuButtons(uint(id))
	if err != nil {
		logger.Log().Error("menu", "获取菜单按钮失败", err)
		c.JSON(500, util.ServerErrorResponse("获取菜单按钮失败", err))
		return
	}

	logger.Log().Info("menu", "获取菜单按钮成功", fmt.Sprintf("菜单ID: %d", id))
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}

// GetAllPages
// @Tags 菜单相关
// @title 获取所有页面
// @description 获取所有页面
// @Summary 获取所有页面
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /menus/all-pages [get]
func GetAllPages(c *gin.Context) {
	res, err := service.MenuServiceApp().GetAllPages()
	if err != nil {
		logger.Log().Error("menu", "获取所有页面失败", err)
		c.JSON(500, util.ServerErrorResponse("获取所有页面失败", err))
		return
	}

	logger.Log().Info("menu", "获取所有页面成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}

// GetConstantRoutes
// @Tags 菜单相关
// @title 获取所有常量路由
// @description 获取所有常量路由
// @Summary 获取所有常量路由
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /menus/constant-routes [get]
func GetConstantRoutes(c *gin.Context) {
	res, err := service.MenuServiceApp().GetConstantRoutes()
	if err != nil {
		logger.Log().Error("menu", "获取所有页面失败", err)
		c.JSON(500, util.ServerErrorResponse("获取所有页面失败", err))
		return
	}

	logger.Log().Info("menu", "获取所有页面成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}

// GetMenuTree
// @Tags 菜单相关
// @title 获取菜单树
// @description 获取菜单树
// @Summary 获取菜单树
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /menus/tree [get]
func GetMenuTree(c *gin.Context) {
	res, err := service.MenuServiceApp().GetMenuTree()
	if err != nil {
		logger.Log().Error("menu", "获取菜单树失败", err)
		c.JSON(500, util.ServerErrorResponse("获取菜单树失败", err))
		return
	}

	logger.Log().Info("menu", "获取菜单树成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}

// IsRouteExist
// @Tags 菜单相关
// @title 判断路由是否存在
// @description 判断路由是否存在
// @Summary 判断路由是否存在
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param routeName query string true "路由名称"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /menus/is-route-exist [get]
func IsRouteExist(c *gin.Context) {
	routeName := c.Query("routeName")
	mBool, err := service.MenuServiceApp().IsRouteExist(routeName)
	if err != nil {
		logger.Log().Error("menu", "判断路由是否存在失败", err)
		c.JSON(500, util.ServerErrorResponse("判断路由是否存在失败", err))
		return
	}

	logger.Log().Info("menu", "判断路由是否存在成功", fmt.Sprintf("路由名称: %s", routeName))
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: mBool,
	})
}
