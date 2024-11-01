package util

import (
	"FreeOps/global"
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"net"
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

// func SSHNewClient(config *api.SSHExecReq) (client *ssh.Client, err error) {
func SSHNewClient(hostIp string, username string, sshPort string, password string, priKey []byte, passphrase []byte, sockPath string) (client *ssh.Client, netConn net.Conn, sshAgentPointer *agent.ExtendedAgent, err error) {
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
	// 2. 密码方式 放在key之后,这样密钥失败之后可以使用Password方式
	if password != "" {
		clientConfig.Auth = append(clientConfig.Auth, ssh.Password(password))
	}

	// 3. agent 模式放在最后,意味着websocket连接，需要使用 openssh agent forwarding
	if sockPath != "" {
		if auth, netConn, sshAgentPointer, err = AuthWithAgent(sockPath); err != nil {
			return nil, nil, nil, fmt.Errorf("agent模式生成ssh.AuthMethod失败: %v", err)
		}
		clientConfig.Auth = append(clientConfig.Auth, auth)
	}

	if clientConfig.Auth == nil {
		return nil, nil, nil, errors.New("未能生成clientConfig.Auth")
	}
	client, err = ssh.Dial("tcp", net.JoinHostPort(hostIp, sshPort), clientConfig)
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

func CreateSFTPClient(client *ssh.Client) (sftpClient *sftp.Client, err error) {
	sftpClient, err = sftp.NewClient(client)
	if err != nil {
		return nil, fmt.Errorf("生成sftp.Client失败: %v", err)
	}

	return sftpClient, err
}
