package controller

import (
	"FreeOps/global"
	"FreeOps/internal/service"
	"FreeOps/pkg/api"
	"FreeOps/pkg/logger"
	"FreeOps/pkg/util"
	"github.com/gin-gonic/gin"
)

// WebSSHConn
// @Tags 工具相关
// @title webSSH连接Linux
// @description webSSH连接Linux,这里用jumpserver举例，默认使用配置文件用户与密钥。可以改为自动获取当前用户，防止冒用其它user
// @Summary webSSH连接Linux
// @Produce  application/json
// @Param data query api.WebSSHConnReq true "传HostID、屏幕高宽"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/tools/webSSH [get]
func WebSSHConn(c *gin.Context) {
	var (
		param api.WebSSHConnReq
		err   error
	)
	if err = global.IncreaseWebSSHConn(); err != nil {
		c.JSON(500, util.ServerErrorResponse("已达到最大webSSH数量", err))
		return
	}
	defer func() {
		global.ReduceWebSSHConn()
	}()

	if err = c.ShouldBindQuery(&param); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}

	wsRes, err := service.Tool().WebSSHConn(c, param)
	if err != nil {
		logger.Log().Error("Webssh", wsRes+"连接Webssh失败", err)
		c.JSON(500, util.ServerErrorResponse(wsRes+"连接Webssh失败", err))
		return
	}
}
