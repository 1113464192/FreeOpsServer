package controller

import (
	"FreeOps/internal/consts"
	"FreeOps/internal/service"
	"FreeOps/pkg/api"
	"FreeOps/pkg/logger"
	"FreeOps/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

// UpdateApi
// @Tags API相关
// @title 新增/修改API信息
// @description 新增不用传ID，修改才传ID
// @Summary 新增/修改API信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData api.UpdateApiReq true "传新增/修改API的参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/apis [post]
func UpdateApi(c *gin.Context) {
	var (
		params api.UpdateApiReq
		err    error
	)
	if err = c.ShouldBind(&params); err != nil {
		c.JSON(200, util.BindErrorResponse(err))
		return
	}

	if strings.ToUpper(params.Method) != params.Method {
		c.JSON(200, util.BindErrorResponse(fmt.Errorf("method必须大写,当前值: %s", params.Method)))
	}

	if err = service.ApiServiceApp().UpdateApi(&params); err != nil {
		logger.Log().Error("api", "创建/修改Api失败", err)
		c.JSON(200, util.ServerErrorResponse("创建/修改Api失败", err))
		return
	}

	logger.Log().Info("api", "创建/修改Api成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
	})
}

// GetApis
// @Tags API相关
// @title 获取apis信息
// @description 查询apis信息
// @Summary 获取apis信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetApiReq true "传新增/修改API的参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/apis [get]
func GetApis(c *gin.Context) {
	var params api.GetApiReq
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(200, util.BindErrorResponse(err))
		return
	}
	res, err := service.ApiServiceApp().GetApis(&params)
	if err != nil {
		logger.Log().Error("api", "创建/修改Api失败", err)
		c.JSON(200, util.ServerErrorResponse("创建/修改Api失败", err))
		return
	}

	logger.Log().Info("api", "查询Api成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}

// DeleteApi
// @Tags API相关
// @title 删除apis
// @description 删除apis
// @Summary 删除apis
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.IdsReq true "传新增/修改API的参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/apis [delete]
func DeleteApi(c *gin.Context) {
	var params api.IdsReq
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(200, util.BindErrorResponse(err))
		return
	}
	if err := service.ApiServiceApp().DeleteApi(params); err != nil {
		logger.Log().Error("api", "删除Api失败", err)
		c.JSON(200, util.ServerErrorResponse("删除Api失败", err))
		return
	}

	logger.Log().Info("api", "删除Api成功", fmt.Sprintf("ApiID: %v", params.Ids))
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
	})
}

// GetApiGroup
// @Tags API相关
// @title 获取API所有组
// @description 获取API所有组
// @Summary 获取API所有组
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/apis/group [get]
func GetApiGroup(c *gin.Context) {
	res, err := service.ApiServiceApp().GetApiGroup()
	if err != nil {
		logger.Log().Error("api", "查询Api组群失败", err)
		c.JSON(200, util.ServerErrorResponse("查询Api组群失败", err))
		return
	}

	logger.Log().Info("api", "查询Api组群成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}

// GetApiTree
// @Tags API相关
// @title 获取API树
// @description 获取API树
// @Summary 获取API树
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/apis/tree [get]
func GetApiTree(c *gin.Context) {
	res, err := service.ApiServiceApp().GetApiTree()
	if err != nil {
		logger.Log().Error("api", "获取API树失败", err)
		c.JSON(200, util.ServerErrorResponse("获取API树失败", err))
		return
	}

	logger.Log().Info("api", "获取API树成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}
