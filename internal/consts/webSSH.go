package consts

import "time"

// webSSH
const (
	WebSSHLinuxTerminal = "linux"
	// 终端窗口
	WebSSHXTerminal                   = "xterm"
	WebSSHPingPeriod                  = 20 * time.Second
	WebSSHPongWait                    = WebSSHPingPeriod * 2
	WebSSHWriteWait                   = 10 * time.Second
	WebSSHReadMessageTickerDuration   = time.Millisecond * time.Duration(40)
	WebSSHSockPath                    = `/tmp/agent.%d`
	WebSSHIdKeyPath                   = `/tmp/%d_key`
	WebSSHMaxRecordLength             = 2048
	WebSSHGenerateLocalSSHAgentSocket = "cd %s/server/shellScript && ./ssh_agent.sh"
)
