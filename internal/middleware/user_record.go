package middleware

import (
	"FreeOps/internal/consts"
	"FreeOps/internal/model"
	"FreeOps/internal/service"
	"FreeOps/pkg/api"
	"FreeOps/pkg/logger"
	"FreeOps/pkg/util"
	"bytes"
	"io"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var recordService = service.UserRecordApp()

var respPool = sync.Pool{
	New: func() any {
		return make([]byte, 1024)
	},
}

type responseBodyWriter struct {
	// 嵌入 gin.ResponseWriter，表示它将继承 gin.ResponseWriter 的所有字段和方法
	gin.ResponseWriter
	// 用于存储响应的内容
	body *bytes.Buffer
}

var writerPool = sync.Pool{
	New: func() any {
		return &responseBodyWriter{
			body: &bytes.Buffer{},
		}
	},
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	// 写入到 responseBodyWriter 的 body 字段中的缓冲区
	r.body.Write(b)
	// 同时将内容存储到 responseBodyWriter 的 body 字段中的缓冲区中，以便后续获取响应内容
	return r.ResponseWriter.Write(b)
}

func shouldSkipLogging(method, path string) bool {
	if method == "GET" || method == "OPTIONS" {
		return true
	}

	if method == "POST" {
		for _, p := range consts.SkipLoggingPostPaths {
			if path == p {
				return true
			}
		}
	}

	return false
}

func UserRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		if shouldSkipLogging(c.Request.Method, c.Request.URL.Path) {
			c.Next()
			return
		}

		var body []byte
		user, err := util.GetClaimsUser(c)
		if err != nil {
			c.JSON(200, api.Response{
				Code: consts.SERVICE_MODAL_LOGOUT_CODE,
				Msg:  err.Error(),
			})
			c.Abort()
			return
		}
		body, err = io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Log().Error("userRecord", "记录用户请求body失败", err)
			c.Abort()
			return
		} else {
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		record := &model.UserRecord{
			Ip:       c.ClientIP(),
			Method:   c.Request.Method,
			Path:     c.Request.URL.Path,
			Agent:    c.Request.UserAgent(),
			Body:     string(body),
			UserID:   user.ID,
			Username: user.Username,
		}

		writer := writerPool.Get().(*responseBodyWriter)
		writer.ResponseWriter = c.Writer
		writer.body.Reset()
		c.Writer = writer

		startNow := time.Now().Local()
		c.Next()
		record.Latency = time.Since(startNow)
		record.Status = c.Writer.Status()
		record.Resp = writer.body.String()
		if err = recordService.CreateRecord(record); err != nil {
			logger.Log().Error("userRecord", "创建用户记录失败", err)
		}
		if len(record.Resp) > 1024 {
			newBody := respPool.Get().([]byte)
			copy(newBody, record.Resp)
			record.Resp = string(newBody)
			defer respPool.Put(newBody)
		}

	}
}
