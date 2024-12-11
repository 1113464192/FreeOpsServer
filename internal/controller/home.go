package controller

import (
	"FreeOps/internal/consts"
	"FreeOps/internal/service"
	"FreeOps/pkg/api"
	"FreeOps/pkg/logger"
	"FreeOps/pkg/util"
	"github.com/gin-gonic/gin"
)

// GetHomeInfo
// @Tags 首页相关
// @title 首页展示信息
// @description 获取展示基本信息
// @Summary 首页展示信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /home/info [get]
func GetHomeInfo(c *gin.Context) {
	res, err := service.HomeServiceApp().GetHomeInfo()
	if err != nil {
		logger.Log().Error("home", "获取首页展示信息失败", err)
		c.JSON(500, util.ServerErrorResponse("获取首页展示信息失败", err))
		return
	}

	logger.Log().Info("home", "获取首页展示信息成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}
