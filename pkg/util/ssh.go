package util

import (
	"FreeOps/global"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"net"
	"strconv"
	"time"
)

func AuthWithPrivateKeyBytes(key []byte, passphrase []byte) (ssh.AuthMethod, error) {
	var signer ssh.Signer
	var err error
	if passphrase == nil {
		signer, err = ssh.ParsePrivateKey(key)
	} else {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(key, passphrase)
	}
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}

func AuthWithAgent(sockPath string) (ssh.AuthMethod, net.Conn, *agent.ExtendedAgent, error) {
	socks, err := net.Dial("unix", sockPath)
	if err != nil {
		return nil, nil, nil, err
	}
	// 1. 返回Signers函数的结果
	sshAgent := agent.NewClient(socks)
	return ssh.PublicKeysCallback(sshAgent.Signers), socks, &sshAgent, nil
}

// 不添加密码连接，安全性太低
func SSHNewClient(hostIp string, username string, sshPort uint16, priKey []byte, passphrase []byte, sockPath string) (client *ssh.Client, netConn net.Conn, sshAgentPointer *agent.ExtendedAgent, err error) {
	// 将sshPort变成string
	sshPortStr := strconv.FormatUint(uint64(sshPort), 10)
	duration, err := time.ParseDuration(global.Conf.SshConfig.SshClientTimeout)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("超时时间获取失败: %v", err)
	}

	clientConfig := &ssh.ClientConfig{
		User:            username,
		Timeout:         duration,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 忽略public key的安全验证
	}

	var auth ssh.AuthMethod
	if priKey != nil {
		if auth, err = AuthWithPrivateKeyBytes(priKey, passphrase); err == nil {
			clientConfig.Auth = append(clientConfig.Auth, auth)
		}
	}
	// 2. agent 模式放在key之后,意味着websocket连接，需要使用 openssh agent forwarding
	if sockPath != "" {
		if auth, netConn, sshAgentPointer, err = AuthWithAgent(sockPath); err != nil {
			return nil, nil, nil, fmt.Errorf("agent模式生成ssh.AuthMethod失败: %v", err)
		}
		clientConfig.Auth = append(clientConfig.Auth, auth)
	}

	if clientConfig.Auth == nil {
		return nil, nil, nil, errors.New("未能生成clientConfig.Auth")
	}
	client, err = ssh.Dial("tcp", net.JoinHostPort(hostIp, sshPortStr), clientConfig)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("生成ssh.Client失败: %v", err)
	}
	return client, netConn, sshAgentPointer, err
}

func SSHNewSession(client *ssh.Client) (session *ssh.Session, err error) {
	session, err = client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("生成ssh.Session失败: %v", err)
	}
	return session, err
}
