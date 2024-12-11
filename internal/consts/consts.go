package consts

const (
	// 自定义SSH命令错误返回值
	SSHCustomCmdError = 9999
)

// 不记录到用户操作的post请求路径
var SkipLoggingPostPaths = []string{
	"/ops/commands",
	"/auth/login",
	"/auth/refreshToken",
	"/ops/run-task-check-script",
	"/users/logout",
}
