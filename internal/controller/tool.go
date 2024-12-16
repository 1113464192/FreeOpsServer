package controller

import (
	"FreeOps/global"
	"FreeOps/internal/model"
	"FreeOps/internal/service/tool"
	"FreeOps/pkg/api"
	"FreeOps/pkg/logger"
	"FreeOps/pkg/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSSHConn
// @Tags 工具相关
// @title webSSH连接Linux
// @description webSSH连接Linux,这里用jumpserver举例，默认使用配置文件用户与密钥。可以改为自动获取当前用户，防止冒用其它user
// @Summary webSSH连接Linux
// @Produce  application/json
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/tools/webSSH [get]
func WebSSHConn(c *gin.Context) {
	var (
		param  api.WebSSHConnReq
		wsConn *websocket.Conn
		user   *model.User
		err    error
	)
	if err = global.IncreaseWebSSHConn(); err != nil {
		c.JSON(500, util.ServerErrorResponse("已达到最大webSSH数量", err))
		return
	}
	defer func() {
		global.ReduceWebSSHConn()
	}()

	if wsConn, user, _, err = util.UpgraderWebSocket(c, true); err != nil {
		c.JSON(500, util.ServerErrorResponse("升级websocket失败", err))
		return
	}
	defer func() {
		wsConn.WriteMessage(websocket.CloseMessage, []byte("websocket连接关闭"))
		wsConn.Close()
	}()
	_, message, err := wsConn.ReadMessage()
	if err != nil {
		c.JSON(500, util.ServerErrorResponse("读取websocket消息失败", err))
		return
	}

	if err = json.Unmarshal(message, &param); err != nil {
		c.JSON(500, util.ServerErrorResponse("解析参数失败", err))
		return
	}

	wsRes, err := tool.Tool().WebSSHConn(wsConn, user, param)
	if err != nil {
		c.JSON(500, util.ServerErrorResponse(wsRes+"连接Webssh失败", err))
		logger.Log().Error("tool", wsRes+"连接Webssh失败", err)
		return
	}
}
