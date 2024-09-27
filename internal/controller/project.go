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

// UpdateProject
// @Tags 项目相关
// @title 新增/修改项目信息
// @description 新增不用传项目ID，修改才传项目ID
// @Summary 新增/修改项目信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData api.UpdateProjectReq true "传新增或者修改项目的所需参数"
// @Success 200 {object} api.Response "{"code": "0000", msg: "string", data: "string"}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"code": "", msg: "", data: ""}"
// @Router /projects [post]
func UpdateProject(c *gin.Context) {
	var (
		projectReq api.UpdateProjectReq
		err        error
	)
	if err = c.ShouldBind(&projectReq); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	if err = service.ProjectServiceApp().UpdateProject(&projectReq); err != nil {
		logger.Log().Error("project", "创建/修改项目失败", err)
		c.JSON(500, util.ServerErrorResponse("创建/修改项目失败", err))
		return
	}
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
	})
}

// GetProjects
// @Tags 项目相关
// @title 查询项目信息
// @description 查询所有/指定条件项目的信息
// @Summary 查询项目信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetProjectsReq true "传对应条件的参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /projects [get]
func GetProjects(c *gin.Context) {
	var params api.GetProjectsReq
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	res, err := service.ProjectServiceApp().GetProjects(&params)
	if err != nil {
		logger.Log().Error("project", "查询项目信息失败", err)
		c.JSON(500, util.ServerErrorResponse("查询项目信息失败", err))
		return
	}

	logger.Log().Info("project", "查询项目信息成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}

// GetProjectList
// @Tags 项目相关
// @title 获取项目列表
// @description 查询项目列表
// @Summary 获取项目列表
// @Produce   application/json
// @Param Authorization header string true "格式为：Bearer 登录返回的用户令牌"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /projects/all-summary [get]
func GetProjectList(c *gin.Context) {
	var err error
	res, err := service.ProjectServiceApp().GetProjectList()
	if err != nil {
		logger.Log().Error("project", "获取项目列表失败", err)
		c.JSON(500, util.ServerErrorResponse("获取项目列表失败", err))
		return
	}

	logger.Log().Info("project", "获取项目列表成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}

// DeleteProjects
// @Tags 项目相关
// @title 删除项目
// @description 删除指定项目
// @Summary 删除项目
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param ids body api.IdsReq true "项目ID"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /projects [delete]
func DeleteProjects(c *gin.Context) {
	var param api.IdsReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	if err := service.ProjectServiceApp().DeleteProjects(param.Ids); err != nil {
		logger.Log().Error("project", "删除项目失败", err)
		c.JSON(500, util.ServerErrorResponse("删除项目失败", err))
		return
	}

	logger.Log().Info("project", "删除项目成功", fmt.Sprintf("ID: %v", param.Ids))
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
	})
}

// GetProjectHosts
// @Tags 项目相关
// @title 获取项目的服务器列表
// @description 由于swagger本身的限制，get请求的切片会报错，并非接口本身问题，请换个方式，如http://127.0.0.1:9081/api/v1/group/apis?ids=3&ids=4
// @Summary 获取项目的服务器列表
// @Produce   application/json
// @Param Authorization header string true "格式为：Bearer 登录返回的用户令牌"
// @Param data query api.IdsReq true "传项目ID切片"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /projects/hosts [get]
func GetProjectHosts(c *gin.Context) {
	var params api.IdsReq
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	res, err := service.ProjectServiceApp().GetProjectHosts(params)
	if err != nil {
		logger.Log().Error("project", "获取项目服务器失败", err)
		c.JSON(500, util.ServerErrorResponse("获取项目服务器失败", err))
		return
	}

	logger.Log().Info("project", "获取项目服务器成功", fmt.Sprintf("项目ID: %d", params.Ids))
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}

// GetProjectGames
// @Tags 项目相关
// @title 获取项目的游戏服列表
// @description 由于swagger本身的限制，get请求的切片会报错，并非接口本身问题，请换个方式，如http://127.0.0.1:9081/api/v1/group/apis?ids=3&ids=4
// @Summary 获取项目的游戏服列表
// @Produce   application/json
// @Param Authorization header string true "格式为：Bearer 登录返回的用户令牌"
// @Param data query api.IdsReq true "传项目ID切片"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /projects/games [get]
func GetProjectGames(c *gin.Context) {
	var params api.IdsReq
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	res, err := service.ProjectServiceApp().GetProjectHosts(params)
	if err != nil {
		logger.Log().Error("project", "获取项目游戏服失败", err)
		c.JSON(500, util.ServerErrorResponse("获取项目游戏服失败", err))
		return
	}

	logger.Log().Info("project", "获取项目游戏服成功", fmt.Sprintf("项目ID: %d", params.Ids))
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}

// GetProjectAssetsTotal
// @Tags 项目相关
// @title 获取项目各资产总数
// @description 查询项目各资产总数
// @Summary 获取项目各资产总数
// @Produce   application/json
// @Param Authorization header string true "格式为：Bearer 登录返回的用户令牌"
// @Param id query uint true "传项目ID"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /projects/assets-total [get]
func GetProjectAssetsTotal(c *gin.Context) {
	var (
		id  uint64
		err error
	)
	if id, err = strconv.ParseUint(c.Query("id"), 10, 32); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	res, err := service.ProjectServiceApp().GetProjectAssetsTotal(uint(id))
	if err != nil {
		logger.Log().Error("project", "获取项目各资产总数失败", err)
		c.JSON(500, util.ServerErrorResponse("获取项目各资产总数失败", err))
		return
	}

	logger.Log().Info("project", "获取项目各资产总数成功", fmt.Sprintf("项目ID: %d", id))
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}
