package middleware

import (
	"FreeOps/internal/consts"
	"FreeOps/internal/service"
	"FreeOps/pkg/api"
	"FreeOps/pkg/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

var casbinService = service.CasbinServiceApp()

// 拦截器
func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 获取请求的URI
		obj := c.Request.URL.RequestURI()
		// 获取请求方法
		act := c.Request.Method

		var sub string

		// 获取用户的角色
		roles, err := util.GetClaimsRole(c)
		if err != nil {
			c.JSON(200, api.Response{
				Code: consts.SERVICE_MODAL_LOGOUT_CODE,
				Msg:  err.Error(),
			})
			c.Abort()
			return
		}

		var isAdmin bool
		// 超级用户判断
		for _, role := range *roles {
			if role.RoleCode == "ADMIN" {
				sub = role.RoleCode
				e := casbinService.Casbin()
				if success, _ := e.Enforce(sub, obj, act); success {
					isAdmin = true
					c.Set("isAdmin", isAdmin)
					c.Next()
					return
				}
			} else {
				sub = strconv.FormatUint(uint64(role.ID), 10)
				e := casbinService.Casbin()
				if success, _ := e.Enforce(sub, obj, act); success {
					c.Set("isAdmin", isAdmin)
					c.Next()
					return
				}
			}
		}
		c.JSON(403, api.Response{
			Code: consts.SERVICE_ERROR_CODE,
			Msg:  "权限不足",
		})
		c.Abort()
	}
}
