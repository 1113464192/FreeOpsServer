package tool

import (
	"FreeOps/global"
	"FreeOps/internal/model"
	"FreeOps/pkg/api"
	"FreeOps/pkg/logger"
	"FreeOps/pkg/util"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

type ToolService struct {
}

var (
	insTool = &ToolService{}
)

func Tool() *ToolService {
	return insTool
}

// 发送信息给websocket
func (s *ToolService) WebSSHSendText(wsConn *websocket.Conn, b []byte) error {
	if err := wsConn.WriteMessage(websocket.TextMessage, b); err != nil {
		return fmt.Errorf("发送信息给websocket报错: %v", err)
	}
	return nil
}

// 接收错误信息返回给前端
func (s *ToolService) WebSSHSendErr(wsConn *websocket.Conn, msg string) error {
	// 前端接收到一个json并有wsError这个key的时候，代表这个消息是发送给前端的websocket报错，而不是给用户的
	errMsg := map[string]string{
		"wsError": msg,
	}
	errMsgBytes, err := json.Marshal(errMsg)
	if err != nil {
		return err
	}

	if err = wsConn.WriteJSON(errMsgBytes); err != nil {
		return err
	}
	return nil
}

func (s *ToolService) WebSSHConn(wsConn *websocket.Conn, user *model.User, param api.WebSSHConnReq) (wsRes string, err error) {
	var (
		host    model.Host
		sshConn *SSHConnect
		client  *ssh.Client
		session *ssh.Session
	)

	if err = model.DB.First(&host, param.Hid).Error; err != nil {
		return "", fmt.Errorf("服务器 %d 查询失败: %v", param.Hid, err)
	}
	sshParam := &api.SSHRunReq{
		SSHPort:    host.SSHPort,
		HostIp:     host.Ipv4,
		Username:   global.Conf.SshConfig.OpsSSHUsername,
		Key:        global.OpsSSHKey,
		Passphrase: nil,
	}

	if global.Conf.SshConfig.OpsKeyPassphrase != "" {
		sshParam.Passphrase = []byte(global.Conf.SshConfig.OpsKeyPassphrase)
	}

	// 生成sshClient
	if client, _, _, err = util.SSHNewClient(sshParam.HostIp, sshParam.Username, sshParam.SSHPort, sshParam.Key, sshParam.Passphrase, ""); err != nil {
		if e := s.WebSSHSendErr(wsConn, "生成ssh.Client时发生错误: "+err.Error()); e != nil {
			logger.Log().Warning("tool", "发送错误信息至websocket失败", err)
		}
		return "", fmt.Errorf("生成ssh.Client时发生错误: %v", err)
	}

	// 生成sshSession
	if session, err = util.SSHNewSession(client); err != nil {
		if e := s.WebSSHSendErr(wsConn, "生成ssh.Session时发生错误: "+err.Error()); e != nil {
			logger.Log().Warning("tool", "发送错误信息至websocket失败", err)
		}
		return "", fmt.Errorf("生成ssh.Session时发生错误: %v", err)
	}

	// 生成sshConn
	if sshConn, err = SSHNewConnect(client, session, param.WindowSize, user, &host); err != nil {
		if e := s.WebSSHSendErr(wsConn, "创建ssh连接时发生错误: "+err.Error()); e != nil {
			logger.Log().Warning("tool", "发送错误信息至websocket失败", err)
		}
		return "", fmt.Errorf("创建ssh连接时发生错误: %v", err)
	}

	// 在外层关闭SSH，内层关闭恐导致提前关闭ssh连接
	defer sshConn.Client.Close()
	defer sshConn.Session.Close()

	quit := make(chan struct{}, 1)
	go sshConn.WsSend(wsConn, quit)
	go sshConn.WsRec(wsConn, quit)
	go sshConn.SessionWait(quit)
	<-quit
	return wsRes, nil
}
