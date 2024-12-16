package controller

import (
	"FreeOps/internal/consts"
	"FreeOps/internal/service"
	"FreeOps/pkg/api"
	"FreeOps/pkg/logger"
	"FreeOps/pkg/util"
	"github.com/gin-gonic/gin"
)

// CreateCloudProject
// @Tags 云平台相关
// @title 新增云项目信息
// @description 新增云项目信息
// @Summary 新增云项目信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.CloudProjectReq true "传对应参数"
// @Success 200 {object} api.Response "{"code": "0000", msg: "string", data: "string"}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"code": "", msg: "", data: ""}"
// @Router /api/clouds/create/project [post]
func CreateCloudProject(c *gin.Context) {
	var (
		param api.CloudProjectReq
		err   error
	)

	if err = c.ShouldBind(&param); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}

	if err = service.CloudServiceApp().CreateCloudProject(param.Name, param.CloudPlatform); err != nil {
		logger.Log().Error("cloud", "创建云项目失败", err)
		c.JSON(500, util.ServerErrorResponse("创建云项目失败", err))
		return
	}

	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
	})
}

// CreateCloudHost
// @Tags 云平台相关
// @title 购买服务器
// @description 购买服务器
// @Summary 购买服务器
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.CreateCloudHostReq true "传对应参数"
// @Success 200 {object} api.Response "{"code": "0000", msg: "string", data: "string"}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"code": "", msg: "", data: ""}"
// @Router /api/clouds/create/host [post]
func CreateCloudHost(c *gin.Context) {
	var (
		param api.CreateCloudHostReq
		err   error
	)

	if err = c.ShouldBind(&param); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}

	if err = service.CloudServiceApp().CreateCloudHost(param.ProjectId, param.CloudPlatform, param.HostCount); err != nil {
		logger.Log().Error("cloud", "购买云服务器失败", err)
		c.JSON(500, util.ServerErrorResponse("购买云服务器失败", err))
		return
	}

	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
	})
}

// UpdateCloudProject
// @Tags 云平台相关
// @title 修改云项目信息
// @description 修改云项目信息
// @Summary 修改云项目信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData api.CloudProjectReq true "传需要的参数"
// @Success 200 {object} api.Response "{"code": "0000", msg: "string", data: "string"}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"code": "", msg: "", data: ""}"
// @Router /api/clouds/update/project [post]
func UpdateCloudProject(c *gin.Context) {
	var (
		param api.CloudProjectReq
		err   error
	)

	if err = c.ShouldBind(&param); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}

	if err = service.CloudServiceApp().UpdateCloudProject(param.Name, param.CloudPlatform); err != nil {
		logger.Log().Error("cloud", "更新云项目失败", err)
		c.JSON(500, util.ServerErrorResponse("更新云项目失败", err))
		return
	}

	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
	})
}

// GetCloudProjectId
// @Tags 云平台相关
// @title 获取云项目ID
// @description 获取云项目ID
// @Summary 获取云项目ID
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param name query api.CloudProjectReq true "传需要的参数"
// @Success 200 {object} api.Response "{"code": "0000", msg: "string", data: "string"}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"code": "", msg: "", data: ""}"
// @Router /api/clouds/query/project [get]
func GetCloudProjectId(c *gin.Context) {
	var (
		param api.CloudProjectReq
		err   error
	)

	if err = c.ShouldBind(&param); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}

	cid, err := service.CloudServiceApp().GetCloudProjectId(param.Name, param.CloudPlatform, 0)
	if err != nil {
		logger.Log().Error("cloud", "获取云项目ID失败", err)
		c.JSON(500, util.ServerErrorResponse("获取云项目ID失败", err))
		return
	}
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: cid,
	})
}
