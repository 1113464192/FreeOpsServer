package consts

import "time"

const (
	WebSocketHandshakeTimeout = 30 * time.Second
	WebSocketReadBufferSize   = 4096
	WebSocketWriteBufferSize  = 1024
	WebSocketPingWait         = 1 * time.Minute // 修改为 1 分钟
	WebSocketPongWait         = WebSocketPingWait * 2
)
