package api

type SSHRunReq struct {
	HostIp     string `json:"host_ip"`
	Username   string `json:"username"`
	SSHPort    uint16 `json:"ssh_port"`
	Key        []byte `json:"key"`
	Passphrase []byte `json:"passphrase"`
	Cmd        string `json:"cmd"` // webssh不用填Cmd
}

// 返回更改
type SSHResultRes struct {
	HostIp   string
	Status   int
	Response string
}
