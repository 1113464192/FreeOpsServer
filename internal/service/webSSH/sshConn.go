package webSSH

import (
	"FreeOps/global"
	"FreeOps/internal/consts"
	"FreeOps/internal/model"
	"FreeOps/pkg/api"
	"FreeOps/pkg/logger"
	"bytes"
	"fmt"
	"io"
	"sync"

	"golang.org/x/crypto/ssh"
)

// Webssh的连接配置
type SSHConnect struct {
	Session       *ssh.Session
	Client        *ssh.Client
	CombineOutput *WebsshBufferWriter
	StdinPipe     io.WriteCloser
	Once          sync.Once
	Logger        *logger.Logger
	User          *model.User
	Host          *model.Host
}

type WebsshBufferWriter struct {
	Buffer bytes.Buffer
	Mu     sync.Mutex
}

// SSHNewConnect 取消支持agent forward，连接太不安全
func SSHNewConnect(client *ssh.Client, session *ssh.Session, windowSize api.WindowSize, user *model.User, host *model.Host) (conn *SSHConnect, err error) {
	//if err = agent.ForwardToAgent(client, *sshAgentPointer); err != nil {
	//	return nil, fmt.Errorf("启动agent.ForwardToAgent失败: %v", err)
	//}
	//if err = agent.RequestAgentForwarding(session); err != nil {
	//	return nil, fmt.Errorf("启用 SSH agent forwarding 的请求失败: %v", err)
	//}
	// 生成输入管道
	stdinPipe, err := session.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("生成session.StdinPipe失败: %v", err)
	}

	modes := ssh.TerminalModes{
		// 设置控制是否在终端上显示输入。设置为0表示不显示输入，这在输入密码等敏感信息时很有用
		// "回显"是指当你在终端输入字符时，这些字符会被显示在终端上(开启了回显，那么用户的输入会被服务器返回发送给前端)
		ssh.ECHO: global.Conf.WebSSH.SshEcho,
		// 设置控制终端的输入输出的速度。设置为14400表示输入速度为14400bps。
		// 如果你设置的值超过了你的网络带宽，可能会导致数据传输不稳定，可能会出现丢包、延迟增大等问题
		ssh.TTY_OP_ISPEED: global.Conf.WebSSH.SshTtyOpISpeed,
		ssh.TTY_OP_OSPEED: global.Conf.WebSSH.SshTtyOpOSpeed,
	}
	// 伪终端屏幕高、宽大小(单位为字符)
	// 基本的操作是所有终端类型都支持的。不同的终端类型可能会支持不同的特性，比如颜色、鼠标事件、窗口大小改变事件等。然而，许多基本的终端操作，比如输入和输出文本，是所有终端类型都支持的。
	if err = session.RequestPty(consts.WebSSHXTerminal, windowSize.Height, windowSize.Weight, modes); err != nil {
		return nil, fmt.Errorf("生成Session.RequestPty失败: %v", err)
	}

	// 直接使用session.Run方法。这个方法会在远程服务器上启动一个新的shell，执行你的命令，然后关闭shell。这个过程对用户是透明的，你不需要手动启动和关闭shell
	// 在远程服务器上启动一个命令行界面(新的 shell 会话)。这个命令行界面可以接收和处理命令，返回命令的结果(允许在 SSH 会话中执行命令并与远程服务器进行交互)。只有启动了shell，才能在SSH会话中执行命令
	if err = session.Shell(); err != nil {
		return nil, fmt.Errorf("生成Session.Shell失败: %v", err)
	}
	comboWriter := new(WebsshBufferWriter)
	setStderr(comboWriter, session)
	setStdout(comboWriter, session)

	return &SSHConnect{
		Session:       session,
		Client:        client,
		CombineOutput: comboWriter,
		StdinPipe:     stdinPipe,
		Once:          sync.Once{},
		Logger:        logger.Log(),
		User:          user,
		Host:          host,
	}, err
}

// 实现io.Writer接口的多态条件
func (s *WebsshBufferWriter) Write(p []byte) (int, error) {
	// 防止多个goroutine同时写入接收池，导致多个webssh间数据错乱
	s.Mu.Lock()
	defer s.Mu.Unlock()
	return s.Buffer.Write(p)
}

func setStderr(stderr io.Writer, session *ssh.Session) {
	session.Stderr = stderr
}

func setStdout(stdout io.Writer, session *ssh.Session) {
	session.Stdout = stdout
}
