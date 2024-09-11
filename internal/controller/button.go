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

// UpdateButton
// @Tags 按钮相关
// @title 新增/修改按钮信息
// @description 新增不用传ID，修改才传ID
// @Summary 新增/修改按钮信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data body api.UpdateButtonsReq true "传新增/修改按钮的参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /buttons [post]
func UpdateButton(c *gin.Context) {
	var params api.UpdateButtonsReq
	var err error
	if err = c.ShouldBind(&params); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	if err = service.ButtonServiceApp().UpdateButtons(&params); err != nil {
		logger.Log().Error("button", "创建/修改按钮失败", err)
		c.JSON(500, util.ServerErrorResponse("创建/修改按钮失败", err))
		return
	}

	logger.Log().Info("button", "创建/修改按钮成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
	})
}

// GetButtons
// @Tags 按钮相关
// @title 查询按钮信息
// @description 查询所有/指定条件按钮的信息
// @Summary 查询按钮信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetButtonsReq true "传对应条件的参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /buttons [get]
func GetButtons(c *gin.Context) {
	var params api.GetButtonsReq
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	res, err := service.ButtonServiceApp().GetButtons(&params)
	if err != nil {
		logger.Log().Error("button", "查询按钮信息失败", err)
		c.JSON(500, util.ServerErrorResponse("查询按钮信息失败", err))
		return
	}

	logger.Log().Info("button", "查询按钮信息成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}

// DeleteButtons
// @Tags 按钮相关
// @title 删除按钮
// @description 删除指定按钮
// @Summary 删除按钮
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param ids body api.IdsReq true "按钮ID"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /buttons [delete]
func DeleteButtons(c *gin.Context) {
	var param api.IdsReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	if err := service.ButtonServiceApp().DeleteButtons(param.Ids); err != nil {
		logger.Log().Error("button", "删除按钮失败", err)
		c.JSON(500, util.ServerErrorResponse("删除按钮失败", err))
		return
	}

	logger.Log().Info("button", "删除按钮成功", fmt.Sprintf("ID: %v", param.Ids))
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
	})
}
