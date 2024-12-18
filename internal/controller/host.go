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

// UpdateHost
// @Tags 服务器相关
// @title 新增/修改服务器信息
// @description 新增不用传服务器ID，修改才传服务器ID
// @Summary 新增/修改服务器信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData api.UpdateHostReq true "传新增或者修改服务器的所需参数"
// @Success 200 {object} api.Response "{"code": "0000", msg: "string", data: "string"}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"code": "", msg: "", data: ""}"
// @Router /api/hosts [post]
func UpdateHost(c *gin.Context) {
	var (
		hostReq api.UpdateHostReq
		err     error
	)
	if err = c.ShouldBind(&hostReq); err != nil {
		c.JSON(200, util.BindErrorResponse(err))
		return
	}
	if err = service.HostServiceApp().UpdateHost(&hostReq); err != nil {
		logger.Log().Error("host", "创建/修改服务器失败", err)
		c.JSON(200, util.ServerErrorResponse("创建/修改服务器失败", err))
		return
	}
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
	})
}

// GetHosts
// @Tags 服务器相关
// @title 查询服务器信息
// @description 查询所有/指定条件服务器的信息
// @Summary 查询服务器信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetHostsReq true "传对应条件的参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/hosts [get]
func GetHosts(c *gin.Context) {
	var (
		params         api.GetHostsReq
		bindProjectIds []uint
	)
	err := c.ShouldBind(&params)
	if err != nil {
		c.JSON(200, util.BindErrorResponse(err))
		return
	}
	if bindProjectIds, err = service.UserServiceApp().GetUserProjectIDs(c); err != nil {
		c.JSON(200, util.ServerErrorResponse("查询用户项目ID失败", err))
		return
	}
	res, err := service.HostServiceApp().GetHosts(&params, bindProjectIds)
	if err != nil {
		logger.Log().Error("host", "查询服务器信息失败", err)
		c.JSON(200, util.ServerErrorResponse("查询服务器信息失败", err))
		return
	}

	logger.Log().Info("host", "查询服务器信息成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}

// DeleteHosts
// @Tags 服务器相关
// @title 删除服务器
// @description 删除指定服务器
// @Summary 删除服务器
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param ids body api.IdsReq true "服务器ID"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/hosts [delete]
func DeleteHosts(c *gin.Context) {
	var param api.IdsReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(200, util.BindErrorResponse(err))
		return
	}
	if err := service.HostServiceApp().DeleteHosts(param.Ids); err != nil {
		logger.Log().Error("host", "删除服务器失败", err)
		c.JSON(200, util.ServerErrorResponse("删除服务器失败", err))
		return
	}

	logger.Log().Info("host", "删除服务器成功", fmt.Sprintf("ID: %v", param.Ids))
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
	})
}

// GetHostList
// @Tags 服务器相关
// @title 获取服务器列表
// @description 查询服务器列表
// @Summary 获取服务器列表
// @Produce   application/json
// @Param Authorization header string true "格式为：Bearer 登录返回的用户令牌"
// @Param id query uint true "项目ID"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/hosts/summary [get]
func GetHostList(c *gin.Context) {
	pidStr := c.Query("id")
	pid, err := strconv.ParseUint(pidStr, 10, 64)
	if err != nil {
		c.JSON(200, util.BindErrorResponse(err))
		return
	}
	var res []api.GetHostListRes
	if res, err = service.HostServiceApp().GetHostList(uint(pid)); err != nil {
		logger.Log().Error("host", "获取服务器列表失败", err)
		c.JSON(200, util.ServerErrorResponse("获取服务器列表失败", err))
		return
	}

	logger.Log().Info("host", "获取服务器列表成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}

// GetHostGameInfo
// @Tags 服务器相关
// @title 获取服务器各业务信息总数
// @description 获取服务器各业务信息总数
// @Summary 获取服务器各业务信息总数
// @Produce   application/json
// @Param Authorization header string true "格式为：Bearer 登录返回的用户令牌"
// @Param id query uint true "传服务器ID"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/hosts/game-info [get]
func GetHostGameInfo(c *gin.Context) {
	var (
		id  uint64
		err error
	)
	if id, err = strconv.ParseUint(c.Query("id"), 10, 32); err != nil {
		c.JSON(200, util.BindErrorResponse(err))
		return
	}
	res, err := service.HostServiceApp().GetHostGameInfo(uint(id))
	if err != nil {
		logger.Log().Error("host", "获取服务器各资产总数失败", err)
		c.JSON(200, util.ServerErrorResponse("获取服务器各资产总数失败", err))
		return
	}

	logger.Log().Info("host", "获取服务器各资产总数成功", fmt.Sprintf("服务器ID: %d", id))
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}
