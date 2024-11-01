package api

type SSHRunReq struct {
	HostIp     string `json:"host_ip"`
	Username   string `json:"username"`
	SSHPort    string `json:"ssh_port"`
	Key        []byte `json:"key"`
	Passphrase []byte `json:"passphrase"`
	Cmd        string `json:"cmd"` // webssh不用填Cmd
}

type SFTPRunReq struct {
	HostIp      string `json:"host_ip"`
	Username    string `json:"username"`
	SSHPort     string `json:"ssh_port"`
	Key         []byte `json:"key"`
	Passphrase  []byte `json:"passphrase"`
	Path        string `json:"path"`
	FileContent string `json:"file_content"`
}

// 返回更改
type SSHResultRes struct {
	HostIp   string
	Status   int
	Response string
}
