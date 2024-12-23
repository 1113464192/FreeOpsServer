package consts

const (
	SERVICE_SUCCESS_CODE        = "0000"          // 后端请求成功的 code
	SERVICE_ERROR_CODE          = "BACKEND_ERROR" // 后端正常返回错误的 code
	SERVICE_LOGOUT_CODE         = "8888"          // 后端请求失败并需要用户退出登录的 code
	SERVICE_MODAL_LOGOUT_CODE   = "7778"          // 后端请求失败并需要用户退出登录的 code（通过弹窗形式提醒)
	SERVICE_REFRESH_TOKEN_CODES = "9998"          // 让前端刷新Token
)
