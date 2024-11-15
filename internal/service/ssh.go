package service

import (
	"FreeOps/global"
	"FreeOps/internal/consts"
	"FreeOps/pkg/api"
	"FreeOps/pkg/logger"
	"FreeOps/pkg/util"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"sync"
)

type SSHService struct {
}

var (
	insSSH = &SSHService{}
)

func SSH() *SSHService {
	return insSSH
}

type sshClientGroup struct {
	clientMap      map[string]*ssh.Client
	clientMapMutex sync.Mutex
}

func (s *SSHService) RunSSHCmdAsync(param *[]api.SSHRunReq) (*[]api.SSHResultRes, error) {
	if err := s.CheckSSHParam(param); err != nil {
		return nil, err
	}

	insClientGroup := sshClientGroup{
		clientMap:      make(map[string]*ssh.Client),
		clientMapMutex: sync.Mutex{},
	}
	channel := make(chan *api.SSHResultRes, len(*param))
	wg := sync.WaitGroup{}
	var err error
	var result []api.SSHResultRes
	// data := make(map[string]string)
	for i := 0; i < len(*param); i++ {
		if err = global.Sem.Acquire(context.Background(), 1); err != nil {
			return nil, fmt.Errorf("获取信号失败，错误为: %v", err)
		}
		wg.Add(1)
		go s.RunSSHCmd(&(*param)[i], channel, &wg, &insClientGroup)
	}
	wg.Wait()
	close(channel)
	for res := range channel {
		result = append(result, *res)
	}
	return &result, err
}

func (s *SSHService) RunSSHCmd(param *api.SSHRunReq, ch chan *api.SSHResultRes, wg *sync.WaitGroup, insClientGroup *sshClientGroup) {
	defer func() {
		if r := recover(); r != nil {
			logger.Log().Error("ssh", "RunSSHCmd执行失败", r)
			result := &api.SSHResultRes{
				HostIp:   param.HostIp,
				Status:   consts.SSHCustomCmdError,
				Response: fmt.Sprintf("触发了recover(): %v", r),
			}
			ch <- result
			wg.Done()
			global.Sem.Release(1)
		}
	}()
	result := &api.SSHResultRes{
		HostIp: param.HostIp,
		Status: 0,
	}
	// client, err := util.SSHNewClient(param.HostIp, param.Username, param.SSHPort, param.Password, param.Key, param.Passphrase)
	client, err := s.getSSHClient(param.HostIp, param.Username, param, insClientGroup)
	if err != nil {
		if exitError, ok := err.(*ssh.ExitError); ok {
			result.Status = exitError.ExitStatus()
		} else {
			result.Status = consts.SSHCustomCmdError
		}
		result.Response = fmt.Sprintf("建立/获取SSH客户端错误: %s", err.Error())
		ch <- result
		wg.Done()
		global.Sem.Release(1)
		return
	}
	defer client.Close()
	session, err := util.SSHNewSession(client)
	if err != nil {
		if exitError, ok := err.(*ssh.ExitError); ok {
			result.Status = exitError.ExitStatus()
		} else {
			result.Status = consts.SSHCustomCmdError
		}
		result.Response = fmt.Sprintf("建立SSH会话错误: %s", err.Error())
		ch <- result
		wg.Done()
		global.Sem.Release(1)
		return
	}
	defer session.Close()
	output, err := session.CombinedOutput(param.Cmd)
	if err != nil {
		if exitError, ok := err.(*ssh.ExitError); ok {
			result.Status = exitError.ExitStatus()
		} else {
			result.Status = consts.SSHCustomCmdError
		}
		result.Response = fmt.Sprintf("Failed to execute command: %s %s", string(output), err.Error())
		ch <- result
		wg.Done()
		global.Sem.Release(1)
		return
	}

	result.Response = string(output)
	ch <- result
	wg.Done()
	global.Sem.Release(1)
}

// 检查是否符合执行条件
func (s *SSHService) CheckSSHParam(param *[]api.SSHRunReq) error {
	for _, p := range *param {
		if p.HostIp == "" || p.Username == "" || p.SSHPort == 0 || p.Cmd == "" || p.Key == nil {
			keyStr := "key存在"
			if p.Key == nil {
				keyStr = "key不存在"
			}
			return fmt.Errorf("执行参数中存在空值: \n HostIp: %s\nUsername: %s\nSSHPort: %d\nCmd: %s\n%s: ", p.HostIp, p.Username, p.SSHPort, p.Cmd, keyStr)
		}
	}
	return nil
}

func (s *SSHService) getSSHClient(hostIp string, username string, param any, insClientGroup *sshClientGroup) (client *ssh.Client, err error) {
	insClientGroup.clientMapMutex.Lock()
	defer insClientGroup.clientMapMutex.Unlock()
	var ok bool
	// 判断对应hostIp的client是否正常存活
	if _, ok = insClientGroup.clientMap[hostIp+"_"+username]; ok {
		if !s.isClientOpen(insClientGroup.clientMap[hostIp+"_"+username]) {
			delete(insClientGroup.clientMap, hostIp+"_"+username)
		}
	}

	// 检查Map中是否已经存在对应hostIp的client
	if client, ok = insClientGroup.clientMap[hostIp+"_"+username]; ok {
		// 如果存在，则直接返回已有的client
		return client, err
	}

	if param, ok := param.(*api.SSHRunReq); ok {
		client, _, _, err = util.SSHNewClient(param.HostIp, param.Username, param.SSHPort, param.Key, param.Passphrase, "")
	}
	if client == nil {
		return nil, errors.New("未能成功获取到ssh.Client")
	}
	insClientGroup.clientMap[hostIp+"_"+username] = client
	return client, err
}

func (s *SSHService) isClientOpen(client *ssh.Client) bool {
	// 发送一个 "keepalive@openssh.com" 请求，这是一个 OpenSSH 定义的全局请求，用于检查连接是否仍然活动
	_, _, err := client.Conn.SendRequest("keepalive@openssh.com", true, nil)
	return err == nil
}
