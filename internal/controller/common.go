package controller

import (
	"FreeOps/pkg/api"
	"FreeOps/pkg/logger"
	"FreeOps/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
)

// CustomError
// @Tags 公共相关
// @title 自定义错误返回
// @description 自定义错误返回
// @Summary 自定义错误返回
// @Produce  application/json
// @Param data query api.CustomErrorReq true "必要参数"
// @Success 200 {object} api.Response "{"code": "0000", msg: "string", data: "string"}"
// @Failure 500 {object} api.Response "{"code": "", msg: "", data: ""}"
// @Router /api/auth/error [get]
func CustomError(c *gin.Context) {
	var param api.CustomErrorReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(200, util.BindErrorResponse(err))
		return
	}
	logInfo := fmt.Sprintf("code: %s, msg: %s", param.Code, param.Msg)
	logger.Log().Info("common", "自定义错误成功", logInfo)
	c.JSON(200, api.Response{
		Code: param.Code,
		Msg:  param.Msg,
	})
}
