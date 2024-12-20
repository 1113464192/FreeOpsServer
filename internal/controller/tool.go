package controller

import (
	"FreeOps/global"
	"FreeOps/internal/consts"
	"FreeOps/internal/model"
	"FreeOps/internal/service"
	"FreeOps/internal/service/tool"
	"FreeOps/pkg/api"
	"FreeOps/pkg/logger"
	"FreeOps/pkg/util"
	"encoding/json"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"strconv"
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
		roles  *[]model.Role
		err    error
	)
	obj := c.Request.URL.RequestURI()
	act := c.Request.Method
	if err = global.IncreaseWebSSHConn(); err != nil {
		c.JSON(200, util.ServerErrorResponse("已达到最大webSSH数量", err))
		return
	}
	defer func() {
		global.ReduceWebSSHConn()
	}()

	if wsConn, user, roles, err = util.UpgraderWebSocket(c, true); err != nil {
		c.JSON(200, util.ServerErrorResponse("升级websocket失败", err))
		return
	}
	defer func() {
		wsConn.WriteMessage(websocket.CloseMessage, []byte("websocket连接关闭"))
		wsConn.Close()
	}()
	var (
		sub string
		e   *casbin.SyncedEnforcer
		// 记录casbin权限审核成功次数
		success int
		isAdmin bool
	)
	// 判断是否管理员操作
	for _, role := range *roles {
		if role.RoleCode == consts.RoleModelAdminCode {
			isAdmin = true
		}
	}
	if !isAdmin {
		for _, role := range *roles {
			sub = strconv.FormatUint(uint64(role.ID), 10)
			e = service.CasbinServiceApp().Casbin()
			if s, _ := e.Enforce(sub, obj, act); s {
				success++
			}
		}
		if success == 0 {
			tool.Tool().WebSSHSendErr(wsConn, "无权限")
			return
		}
	}
	_, message, err := wsConn.ReadMessage()
	if err != nil {
		tool.Tool().WebSSHSendErr(wsConn, "读取websocket消息失败")
		logger.Log().Warning("tool", "读取websocket消息失败", err)
		return
	}

	if err = json.Unmarshal(message, &param); err != nil {
		tool.Tool().WebSSHSendErr(wsConn, "解析参数失败")
		logger.Log().Error("tool", "解析参数失败", err)
		return
	}

	wsRes, err := tool.Tool().WebSSHConn(wsConn, user, param)
	if err != nil {
		tool.Tool().WebSSHSendErr(wsConn, fmt.Sprintf("连接Webssh失败: %s", err.Error()))
		logger.Log().Warning("tool", wsRes+"连接Webssh失败", err)
		return
	}
}
