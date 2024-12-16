package global

type Config struct {
	Mysql        Mysql        `json:"mysql"`
	Logger       Logger       `json:"logger"`
	SshConfig    SshConfig    `json:"sshConfig"`
	WebSSH       WebSSH       `json:"webSSH "`
	System       System       `json:"system"`
	Concurrency  Concurrency  `json:"concurrency"`
	GitWebhook   GitWebhook   `json:"git_webhook"`
	SecurityVars SecurityVars `json:"security_vars"`
}

type Mysql struct {
	Conf                 string
	CreateBatchSize      int
	SlowThreshold        int
	LogLevel             int
	SoftDeleteRetainDays int
}

type Logger struct {
	Level string
}

type SshConfig struct {
	SshClientTimeout string
	OpsSSHUsername   string
	OpsKeyPath       string
	OpsKeyPassphrase string
}

// 保留，暂时不用改，有空再从以前的代码中迁移webSSH代码过来
type WebSSH struct {
	ReadBufferSize   int
	WriteBufferSize  int
	HandshakeTimeout string
	SshEcho          uint32
	SshTtyOpISpeed   uint32
	SshTtyOpOSpeed   uint32
	MaxConnNumber    uint32
}

type Concurrency struct {
	Number int64
}

type System struct {
	Mode string
}

// 保留，暂时不用改，有空再从以前的代码中迁移这段CI代码过来
type GitWebhook struct {
	GithubSecret   string
	GitlabSecret   string
	GitCiScriptDir string
	GitCiRepo      string
}

type SecurityVars struct {
	AesKey                     string
	AesIv                      string
	JwtIssuer                  string
	TokenExpireDuration        string
	RefreshTokenExpireDuration string
	TokenKey                   string
	AllowedCIDR                string
}

var Conf = new(Config)
