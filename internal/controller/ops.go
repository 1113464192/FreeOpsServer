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
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"time"
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
// @Param data query api.GetOpsTemplatesReq true "传新增或者修改操作模板的所需参数"
// @Success 200 {object} api.Response "{"code": "0000", msg: "string", data: "string"}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"code": "", msg: "", data: ""}"
// @Router /ops/template [get]
func GetOpsTemplate(c *gin.Context) {
	var (
		temReq         api.GetOpsTemplatesReq
		bindProjectIds []uint
		err            error
	)
	if err = c.ShouldBind(&temReq); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	if bindProjectIds, err = service.UserServiceApp().GetUserProjectIDs(c); err != nil {
		c.JSON(500, util.ServerErrorResponse("获取用户的项目IDs失败", err))
		return
	}
	result, err := service.OpsServiceApp().GetOpsTemplate(&temReq, bindProjectIds)
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

// DeleteOpsTemplate
// @Tags 运维操作相关
// @title 删除操作模板
// @description 删除操作模板
// @Summary 删除操作模板
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param ids body api.IdsReq true "模板IDs"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /ops/template [delete]
func DeleteOpsTemplate(c *gin.Context) {
	var param api.IdsReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	if err := service.OpsServiceApp().DeleteOpsTemplate(param.Ids); err != nil {
		logger.Log().Error("ops", "删除运维操作模板失败", err)
		c.JSON(500, util.ServerErrorResponse("删除运维操作模板失败", err))
		return
	}

	logger.Log().Info("ops", "删除运维操作模板成功", fmt.Sprintf("ID: %v", param.Ids))
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
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
// @Param data query api.GetOpsParamsTemplatesReq true "传新增或者修改操作模板的所需参数"
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

// DeleteOpsParamsTemplate
// @Tags 运维操作相关
// @title 删除运维操作的参数模板
// @description 删除运维操作的参数模板
// @Summary 删除运维操作的参数模板
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param ids body api.IdsReq true "模板IDs"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /ops/param-template [delete]
func DeleteOpsParamsTemplate(c *gin.Context) {
	var param api.IdsReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	if err := service.OpsServiceApp().DeleteOpsParamsTemplate(param.Ids); err != nil {
		logger.Log().Error("ops", "删除运维操作的参数模板失败", err)
		c.JSON(500, util.ServerErrorResponse("删除运维操作的参数模板失败", err))
		return
	}

	logger.Log().Info("ops", "删除运维操作的参数模板成功", fmt.Sprintf("ID: %v", param.Ids))
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
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
// @Param data body api.UpdateOpsTaskReq true "传新增或者修改操作模板的所需参数"
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
// @Summary 查询运维操作任务信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetOpsTaskReq true "传新增或者修改操作模板的所需参数"
// @Success 200 {object} api.Response "{"code": "0000", msg: "string", data: "string"}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"code": "", msg: "", data: ""}"
// @Router /ops/task [get]
func GetOpsTask(c *gin.Context) {
	var (
		taskReq        api.GetOpsTaskReq
		bindProjectIds []uint
		err            error
	)
	if err = c.ShouldBind(&taskReq); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	if bindProjectIds, err = service.UserServiceApp().GetUserProjectIDs(c); err != nil {
		c.JSON(500, util.ServerErrorResponse("获取用户的项目IDs失败", err))
		return
	}
	result, err := service.OpsServiceApp().GetOpsTask(taskReq, bindProjectIds)
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

// DeleteOpsTask
// @Tags 运维操作相关
// @title 删除运维操作任务信息
// @description 删除运维操作任务信息
// @Summary 删除运维操作任务信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param ids body api.IdsReq true "任务IDs"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /ops/task [delete]
func DeleteOpsTask(c *gin.Context) {
	var param api.IdsReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	if err := service.OpsServiceApp().DeleteOpsTask(param.Ids); err != nil {
		logger.Log().Error("ops", "删除运维操作任务信息失败", err)
		c.JSON(500, util.ServerErrorResponse("删除运维操作任务信息失败", err))
		return
	}

	logger.Log().Info("ops", "删除运维操作任务信息成功", fmt.Sprintf("ID: %v", param.Ids))
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
	})
}

// RunOpsTaskCheckScript
// @Tags 运维操作相关
// @title 执行一个阻塞的任务并返回结果
// @description 目前主要是为了执行运维的检查脚本，返回给运营审批时阅览
// @Summary 执行一个阻塞的任务并返回结果
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.RunOpsTaskCheckScriptReq true "请输入需要的参数"
// @Success 200 {object} api.Response "{"code": "0000", msg: "string", data: "string"}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"code": "", msg: "", data: ""}"
// @Router /ops/run-task-check-script [post]
func RunOpsTaskCheckScript(c *gin.Context) {
	var (
		taskReq api.RunOpsTaskCheckScriptReq
		err     error
		result  *[]api.SSHResultRes
	)
	if err = c.ShouldBind(&taskReq); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	if result, err = service.OpsServiceApp().RunOpsTaskCheckScript(taskReq); err != nil {
		logger.Log().Error("ops", "执行单个运维任务失败", err)
		c.JSON(500, util.ServerErrorResponse("执行单个运维任务失败", err))
		return
	}

	logger.Log().Info("ops", "执行单个运维任务成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: result,
	})
}

// GetOpsTaskTmpCommands
// @Tags 运维操作相关
// @title 查看根据参数会生成的命令
// @description 查看根据参数会生成的命令
// @Summary 查看根据参数会生成的命令
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.GetOpsTaskTmpCommandsReq true "请输入需要的参数"
// @Success 200 {object} api.Response "{"code": "0000", msg: "string", data: "string"}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"code": "", msg: "", data: ""}"
// @Router /ops/commands [post]
func GetOpsTaskTmpCommands(c *gin.Context) {
	var (
		params api.GetOpsTaskTmpCommandsReq
		err    error
	)
	if err = c.ShouldBind(&params); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}

	res, err := service.OpsServiceApp().GetOpsTaskTmpCommands(params)
	if err != nil {
		logger.Log().Error("ops", "查看根据参数会生成的命令失败", err)
		c.JSON(500, util.ServerErrorResponse("查看根据参数会生成的命令失败", err))
		return
	}

	logger.Log().Info("ops", "查看根据参数会生成的命令成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}

// SubmitOpsTask
// @Tags 运维操作相关
// @title 提交运维操作任务
// @description 提交运维操作任务
// @Summary 提交运维操作任务
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.SubmitOpsTaskReq true "请输入需要的参数"
// @Success 200 {object} api.Response "{"code": "0000", msg: "string", data: "string"}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"code": "", msg: "", data: ""}"
// @Router /ops/submit-task [post]
func SubmitOpsTask(c *gin.Context) {
	var (
		taskReq api.SubmitOpsTaskReq
		err     error
	)
	if err = c.ShouldBind(&taskReq); err != nil {
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

	if err = service.OpsServiceApp().SubmitOpsTask(taskReq, user.ID); err != nil {
		logger.Log().Error("ops", "提交运维操作任务失败", err)
		c.JSON(500, util.ServerErrorResponse("提交运维操作任务失败", err))
		return
	}

	logger.Log().Info("ops", "提交运维操作任务成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
	})
}

// ApproveOpsTask
// @Tags 运维操作相关
// @title 用户审批任务
// @description 用户审批任务
// @Summary 用户审批任务
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param tid formData api.ApproveOpsTaskReq true "请输入需要的参数"
// @Success 200 {object} api.Response "{"code": "0000", msg: "string", data: "string"}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"code": "", msg: "", data: ""}"
// @Router /ops/approve-task [put]
func ApproveOpsTask(c *gin.Context) {
	var (
		params api.ApproveOpsTaskReq
		err    error
	)
	if err = c.ShouldBind(&params); err != nil {
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

	if err = service.OpsServiceApp().ApproveOpsTask(params, user.ID); err != nil {
		logger.Log().Error("ops", "用户审批任务失败", err)
		c.JSON(500, util.ServerErrorResponse("用户审批任务失败", err))
		return
	}

	logger.Log().Info("ops", "用户审批任务成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
	})
}

// GetTaskPendingApprovers
// @Tags 运维操作相关
// @title 获取等待中的任务与尚未审批的用户
// @description 获取等待中的任务与对应尚未审批的用户
// @Summary 获取等待中的任务与对应尚未审批的用户
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Success 200 {object} api.Response "{"code": "0000", msg: "string", data: "string"}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"code": "", msg: "", data: ""}"
// @Router /ops/task-pending-approvers [get]
func GetTaskPendingApprovers(c *gin.Context) {
	res, err := service.OpsServiceApp().GetTaskPendingApprovers()
	if err != nil {
		logger.Log().Error("ops", "获取等待中的任务与尚未审批的用户失败", err)
		c.JSON(500, util.ServerErrorResponse("获取等待中的任务与尚未审批的用户失败", err))
		return
	}

	logger.Log().Info("ops", "获取等待中的任务与尚未审批的用户成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}

// GetOpsTaskLog
// @Tags 运维操作相关
// @title 查询运维操作任务日志
// @description 要获取具体信息直接传ID, 不获取commands等，只批量获取name等基础数据不用传ID
// @Summary 查询运维操作任务日志
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetOpsTaskLogReq true "传新增或者修改操作模板的所需参数"
// @Success 200 {object} api.Response "{"code": "0000", msg: "string", data: "string"}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"code": "", msg: "", data: ""}"
// @Router /ops/task-log [get]
func GetOpsTaskLog(c *gin.Context) {
	var (
		taskReq        api.GetOpsTaskLogReq
		bindProjectIds []uint
		err            error
	)
	if err = c.ShouldBind(&taskReq); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	if bindProjectIds, err = service.UserServiceApp().GetUserProjectIDs(c); err != nil {
		c.JSON(500, util.ServerErrorResponse("获取用户的项目IDs失败", err))
		return
	}
	result, err := service.OpsServiceApp().GetOpsTaskLog(taskReq, bindProjectIds)
	if err != nil {
		logger.Log().Error("ops", "查询运维操作任务日志失败", err)
		c.JSON(500, util.ServerErrorResponse("查询运维操作任务日志失败", err))
		return
	}
	logger.Log().Info("ops", "查询运维操作任务日志成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: result,
	})
}

// GetOpsTaskRunningWS
// @Tags 运维操作相关
// @title 实时同步执行中的任务状态
// @description websocket实时同步权限内的项目执行中的任务状态
// @Summary 实时同步执行中的任务状态
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /ops/task-running-ws [get]
func GetOpsTaskRunningWS(c *gin.Context) {
	var (
		roleIds        []uint
		bindProjectIds []uint
	)
	upgrader := websocket.Upgrader{
		HandshakeTimeout: 2 * time.Second,                            //握手超时时间
		ReadBufferSize:   1024,                                       //读缓冲大小
		WriteBufferSize:  1024,                                       //写缓冲大小
		CheckOrigin:      func(r *http.Request) bool { return true }, // 可以规定跨域域名
		// 这个函数在升级 HTTP 连接到 WebSocket 连接时(握手)发生错误时被调用。例如可以在这个函数中向 HTTP 响应中写入一个错误消息
		Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
			// 写入状态码
			w.WriteHeader(status)
			// 写入错误信息
			w.Write([]byte("WebSocket upgrade failed: " + reason.Error()))
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(500, util.ServerErrorResponse("websocket连接创建失败", err))
		logger.Log().Error("ops", "websocket连接创建失败", err)
		return
	}
	// 获取角色对应的项目ID
	if roleIds, err = service.RoleServiceApp().GetSelfRoleIDs(c); err != nil {
		c.JSON(500, util.ServerErrorResponse("获取角色IDs失败", err))
		logger.Log().Error("ops", "获取角色IDs失败", err)
		return
	}
	if bindProjectIds, err = service.RoleServiceApp().GetRoleProjects(roleIds); err != nil {
		logger.Log().Error("ops", "获取角色对应的项目ID失败", err)
		c.JSON(500, util.ServerErrorResponse("获取角色对应的项目ID失败", err))
		return
	}
	if err = service.OpsServiceApp().GetOpsTaskRunningWS(conn, bindProjectIds); err != nil {
		logger.Log().Error("ops", "实时同步执行中的任务状态失败", err)
		c.JSON(500, util.ServerErrorResponse("实时同步执行中的任务状态失败", err))
		return
	}
}
