package controller

import (
	"FreeOps/internal/consts"
	"FreeOps/pkg/api"
	"github.com/gin-gonic/gin"
)

// Test
// @Tags 测试相关
// @title 测试Gin能否正常访问
// @description 无设置权限，返回"Hello world!~~(无权限版)"
// @Summary 测试Gin能否正常访问
// @Produce  application/json
// @Success 200 {} api.Response "{"code": "0000", msg: "string", data: "string"}"
// @Router /ping [get]
func Test(c *gin.Context) {
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: "Hellp world!~~(无权限版)",
	})
}

// Test2
// @Tags 测试相关
// @title 测试Gin能否正常访问
// @description 设置权限，返回"Hello world!~~(验证权限版)"
// @Summary 测试Gin能否正常访问
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Success 200 {} api.Response "{"code": "0000", msg: "string", data: "string"}"
// @Router /api/ping2 [get]
func Test2(c *gin.Context) {
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: "Hello world!~~(验证权限版)",
	})
}
