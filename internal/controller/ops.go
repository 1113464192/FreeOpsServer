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
)

// UpdateOpsTemplate
// @Tags 运维操作相关
// @title 新增/修改操作模板信息
// @description 新增不用传ID，修改才传ID
// @Summary 新增/修改操作模板信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData api.UpdateOpsTemplateReq true "传新增或者修改操作模板的所需参数"
// @Success 200 {object} api.Response "{"code": "0000", msg: "string", data: "string"}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"code": "", msg: "", data: ""}"
// @Router /ops/template [post]
func UpdateOpsTemplate(c *gin.Context) {
	var (
		temReq api.UpdateOpsTemplateReq
		err    error
	)
	if err = c.ShouldBind(&temReq); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	if err = service.OpsServiceApp().UpdateOpsTemplate(&temReq); err != nil {
		logger.Log().Error("ops", "创建/修改运维操作模板失败", err)
		c.JSON(500, util.ServerErrorResponse("创建/修改运维操作模板失败", err))
		return
	}

	logger.Log().Info("ops", "创建/修改运维操作模板成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
	})
}

// GetOpsTemplate
// @Tags 运维操作相关
// @title 查询操作模板信息
// @description 要获取content直接传ID, 不获取content等, 只批量获取name等基础数据不用传ID
// @Summary 查询操作模板信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data Query api.GetOpsTemplatesReq true "传新增或者修改操作模板的所需参数"
// @Success 200 {object} api.Response "{"code": "0000", msg: "string", data: "string"}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"code": "", msg: "", data: ""}"
// @Router /ops/template [get]
func GetOpsTemplate(c *gin.Context) {
	var (
		temReq api.GetOpsTemplatesReq
		err    error
	)
	if err = c.ShouldBind(&temReq); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	result, err := service.OpsServiceApp().GetOpsTemplate(&temReq)
	if err != nil {
		logger.Log().Error("ops", "查询运维操作模板失败", err)
		c.JSON(500, util.ServerErrorResponse("查询运维操作模板失败", err))
		return
	}
	logger.Log().Info("ops", "查询运维操作模板成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: result,
	})
}

// UpdateOpsParamsTemplate
// @Tags 运维操作相关
// @title 新增/修改 运维操作的参数模板
// @description 从运营文案信息获取参数的正则模板。新增不用传ID，修改才传ID。获取结构如(keyword游服,variable为填入运维操作模板的{var},英文逗号分割到var就是数组[100_300, 400_700],可以有多个关键字): 游服: 100-300,400-700\n......
// @Summary 新增/修改 运维操作的参数模板
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData model.OpsParam true "传新增或者修改模板的所需参数"
// @Success 200 {object} api.Response "{"code": "0000", msg: "string", data: "string"}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"code": "", msg: "", data: ""}"
// @Router /ops/param-template [post]
func UpdateOpsParamsTemplate(c *gin.Context) {
	var (
		temReq model.OpsParam
		err    error
	)
	if err = c.ShouldBind(&temReq); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	if err = service.OpsServiceApp().UpdateOpsParamsTemplate(temReq); err != nil {
		logger.Log().Error("ops", "创建/修改运维操作的参数模板失败", err)
		c.JSON(500, util.ServerErrorResponse("创建/修改运维操作的参数模板失败", err))
		return
	}
	logger.Log().Info("ops", "创建/修改运维操作的参数模板成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
	})
}

// GetOpsParamsTemplate
// @Tags 运维操作相关
// @title 获取运维操作的参数模板
// @description 获取运维操作的参数模板
// @Summary 获取运维操作的参数模板
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data Query api.GetOpsParamsTemplatesReq true "传新增或者修改操作模板的所需参数"
// @Success 200 {object} api.Response "{"code": "0000", msg: "string", data: "string"}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"code": "", msg: "", data: ""}"
// @Router /ops/param-template [get]
func GetOpsParamsTemplate(c *gin.Context) {
	var (
		temReq api.GetOpsParamsTemplatesReq
		err    error
	)
	if err = c.ShouldBind(&temReq); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	result, err := service.OpsServiceApp().GetOpsParamsTemplate(temReq)
	if err != nil {
		logger.Log().Error("ops", "查询运维操作的参数模板失败", err)
		c.JSON(500, util.ServerErrorResponse("查询运维操作的参数模板失败", err))
		return
	}
	logger.Log().Info("ops", "查询运维操作的参数模板成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: result,
	})
}

// BindTemplateParams
// @Tags 运维操作相关
// @title 运维操作模板与运维操作的参数模板关系绑定
// @description 运维操作模板与运维操作的参数模板关系绑定
// @Summary 运维操作模板与运维操作的参数模板关系绑定
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.BindTemplateParamsReq true "传运维操作模板ID与运维操作的参数模板IDs"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"修改权限成功，刷新Token"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /ops/bind-template-params [put]
func BindTemplateParams(c *gin.Context) {
	var (
		params api.BindTemplateParamsReq
		err    error
	)
	if err = c.ShouldBind(&params); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	if err = service.OpsServiceApp().BindTemplateParams(params.TemplateID, params.ParamIDs); err != nil {
		logger.Log().Error("ops", "绑定 运维操作模板与运维操作的参数模板 失败", err)
		c.JSON(500, util.ServerErrorResponse("绑定 运维操作模板与运维操作的参数模板 失败", err))
		return
	}

	logger.Log().Info("ops", "绑定 运维操作模板与运维操作的参数模板 成功", fmt.Sprintf("模板ID: %d————参数模板IDs: %d", params.TemplateID, params.ParamIDs))
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
	})
}

// GetTemplateParams
// @Tags 运维操作相关
// @title 查询运维操作模板对应的参数模板
// @description 查询运维操作模板对应的参数模板
// @Summary 查询运维操作模板对应的参数模板
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param id query uint true "模板ID"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"修改权限成功，刷新Token"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /ops/template-params [get]
func GetTemplateParams(c *gin.Context) {
	var res []model.OpsParam
	id, err := strconv.ParseUint(c.Query("id"), 10, 0)
	if err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	if res, err = service.OpsServiceApp().GetTemplateParams(uint(id)); err != nil {
		logger.Log().Error("ops", "查询 运维操作模板与运维操作的参数模板 失败", err)
		c.JSON(500, util.ServerErrorResponse("查询 运维操作模板与运维操作的参数模板 失败", err))
		return
	}

	logger.Log().Info("ops", "查询 运维操作模板与运维操作的参数模板 成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}

// UpdateOpsTask
// @Tags 运维操作相关
// @title 新增/修改 运维操作任务信息
// @description 新增不用传ID，修改才传ID
// @Summary 新增/修改 运维操作任务信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData api.UpdateOpsTaskReq true "传新增或者修改操作模板的所需参数"
// @Success 200 {object} api.Response "{"code": "0000", msg: "string", data: "string"}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"code": "", msg: "", data: ""}"
// @Router /ops/task [post]
func UpdateOpsTask(c *gin.Context) {
	var (
		taskReq api.UpdateOpsTaskReq
		err     error
	)
	if err = c.ShouldBind(&taskReq); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	if err = service.OpsServiceApp().UpdateOpsTask(taskReq); err != nil {
		logger.Log().Error("ops", "创建/修改运维操作模板失败", err)
		c.JSON(500, util.ServerErrorResponse("创建/修改运维操作任务信息失败", err))
		return
	}

	logger.Log().Info("ops", "创建/修改运维操作任务信息成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
	})
}

// GetOpsTask
// @Tags 运维操作相关
// @title 查询运维操作任务信息
// @description 要获取具体信息直接传ID, 不获取content等，只批量获取name等基础数据不用传ID
// @Summary 查询操作模板信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data Query api.GetOpsTaskReq true "传新增或者修改操作模板的所需参数"
// @Success 200 {object} api.Response "{"code": "0000", msg: "string", data: "string"}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"code": "", msg: "", data: ""}"
// @Router /ops/task [get]
func GetOpsTask(c *gin.Context) {
	var (
		taskReq api.GetOpsTaskReq
		err     error
	)
	if err = c.ShouldBind(&taskReq); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	result, err := service.OpsServiceApp().GetOpsTask(taskReq)
	if err != nil {
		logger.Log().Error("ops", "查询运维操作任务信息失败", err)
		c.JSON(500, util.ServerErrorResponse("查询运维操作任务信息失败", err))
		return
	}
	logger.Log().Info("ops", "查询运维操作任务信息成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: result,
	})
}
