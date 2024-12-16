package api

type WebSSHConnReq struct {
	Hid        uint `json:"hid" form:"hid" binding:"required"` // 服务器id
	WindowSize      // 屏幕大小
}

type WindowSize struct {
	Height int `json:"height" form:"height" binding:"required"` // 单位为字符
	Weight int `json:"weight" form:"weight" binding:"required"` // 单位为字符
}