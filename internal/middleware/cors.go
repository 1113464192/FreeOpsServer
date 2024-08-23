package middleware

import (
	"FreeOps/global"
	"FreeOps/internal/consts"
	"FreeOps/pkg/api"
	"FreeOps/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许放通的CIDR
		_, cidrNet, err := net.ParseCIDR(global.Conf.SecurityVars.AllowedCIDR)
		if err != nil {
			logger.Log().Error("net", "生成cidrNet错误: %s", err.Error())
			c.JSON(500, api.Response{
				Code: consts.SERVICE_ERROR_CODE,
				Msg:  fmt.Sprintf("生成cidrNet错误: %s", err.Error()),
			})
			c.Abort()
			return
		}
		ip := net.ParseIP(c.ClientIP())
		if ip == nil || !cidrNet.Contains(ip) {
			c.JSON(403, api.Response{
				Code: consts.SERVICE_ERROR_CODE,
				Msg:  fmt.Sprintf("非法请求,IP:%s  不在白名单内", c.ClientIP()),
			})
			c.Abort()
			return
		}

		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")

		// 如果请求的头信息中包含了不被服务器支持的头信息，浏览器会拦截请求，并阻止JavaScript代码对返回结果的访问
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Headers", "Origin,Content-Length,Content-Type,Cookie,Authorization,token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE,PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, New-Token, New-Expires-At")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(204)
		}
		// 处理请求
		c.Next()
	}
}
