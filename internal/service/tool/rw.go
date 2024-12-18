package tool

import (
	"FreeOps/internal/consts"
	"FreeOps/pkg/api"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

func flushCombOutput(w *WebsshBufferWriter, wsConn *websocket.Conn) error {
	w.Mu.Lock()
	defer w.Mu.Unlock()
	if w.Buffer.Len() != 0 {
		err := wsConn.WriteMessage(websocket.TextMessage, w.Buffer.Bytes())
		if err != nil {
			return err
		}
		w.Buffer.Reset()
	}
	return nil
}

func (s *SSHConnect) wsQuit(ch chan struct{}) {
	s.Once.Do(func() {
		close(ch)
	})
}

// 向websocket发送服务器返回的信息
func (s *SSHConnect) WsSend(wsConn *websocket.Conn, quitCh chan struct{}) {
	defer s.wsQuit(quitCh)
	pingTick := time.NewTicker(consts.WebSSHPingPeriod)
	defer pingTick.Stop()

	tick := time.NewTicker(consts.WebSSHReadMessageTickerDuration)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			//write combine output bytes into websocket response
			if err := flushCombOutput(s.CombineOutput, wsConn); err != nil {
				if e := Tool().WebSSHSendErr(wsConn, "发送服务器返回信息到websocket失败: "+err.Error()); e != nil {
					s.Logger.Error("tool", "发送错误信息至websocket失败", err)
				}
				s.Logger.Error("tool", "发送服务器返回信息到websocket失败", err)
				return
			}
			// 发送ping至websocket
		case <-pingTick.C:
			if err := wsConn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				if e := Tool().WebSSHSendErr(wsConn, "发送Ping至websocket失败: "+err.Error()); e != nil {
					s.Logger.Error("tool", "发送错误信息至websocket失败", err)
				}
				s.Logger.Error("tool", "发送Ping至websocket失败", err)
				return
			}

		case <-quitCh:
			return
		}
	}
}

func (s *SSHConnect) WsRec(wsConn *websocket.Conn, quitCh chan struct{}) {
	//tells other go routine quit
	defer s.wsQuit(quitCh)
	// 处理pong消息
	wsConn.SetReadDeadline(time.Now().Add(consts.WebSSHPongWait))
	wsConn.SetPongHandler(func(appData string) error {
		// 重置读取存活时间
		wsConn.SetReadDeadline(time.Now().Add(consts.WebSSHPongWait))
		return nil
	})

	for {
		select {
		case <-quitCh:
			return
		default:
			// read websocket msg
			_, wsData, err := wsConn.ReadMessage()
			if err != nil {
				if e := Tool().WebSSHSendErr(wsConn, "接收websocket发送的信息失败: "+err.Error()); e != nil {
					s.Logger.Error("tool", "发送错误信息至websocket失败", err)
				}
				s.Logger.Error("tool", "接收websocket发送的信息失败", err)
				return
			}

			// 每次传输一个或多个char
			if len(wsData) > 0 {
				// resize 或者 粘贴
				resize := api.WindowSize{}
				err = json.Unmarshal(wsData, &resize)
				if err != nil {
					goto SEND
				}
				if resize.Height > 0 && resize.Width > 0 {
					if err = s.Session.WindowChange(resize.Height, resize.Width); err != nil {
						if e := Tool().WebSSHSendErr(wsConn, "变更WindowSize失败: "+err.Error()); e != nil {
							s.Logger.Error("tool", "发送错误信息至websocket失败", err)
						}
						s.Logger.Error("tool", "变更WindowSize失败", err)
						break
					}
				} else {
					goto SEND
				}
				break
			}
			// 服务器的返回发送给websocket
		SEND:
			if _, err = s.StdinPipe.Write(wsData); err != nil {
				if e := Tool().WebSSHSendErr(wsConn, "发送服务器信息到前端失败: "+err.Error()); e != nil {
					s.Logger.Error("tool", "发送错误信息至websocket失败", err)
				}
				s.Logger.Error("tool", "发送服务器信息到前端失败", err)
				break
			}
		}
	}
}

// 等待 shell 会话的结束，并在会话结束后执行一些清理操作
func (s *SSHConnect) SessionWait(quitChan chan struct{}) {
	// 等待远程命令的执行完成。它会阻塞当前的 Go 协程，直到远程命令执行完毕或会话被关闭。
	if err := s.Session.Wait(); err != nil {
		s.Logger.Error("tool", "Session.Wait报错", err)
		s.wsQuit(quitChan)
	}
}
