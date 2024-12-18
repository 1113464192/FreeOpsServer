package middleware

import (
	"FreeOps/internal/consts"
	"FreeOps/internal/model"
	"FreeOps/internal/service"
	"FreeOps/pkg/api"
	"FreeOps/pkg/logger"
	"FreeOps/pkg/util"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	jwtService = service.JwtServiceApp()
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(200,
				api.Response{
					Code: consts.SERVICE_MODAL_LOGOUT_CODE,
					Msg:  "用户访问令牌缺失",
				},
			)
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(200,
				api.Response{
					Code: consts.SERVICE_MODAL_LOGOUT_CODE,
					Msg:  "用户访问令牌格式有误",
				},
			)
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		// 判断是否在黑名单
		if jwtService.IsBlacklist(parts[1]) {
			c.JSON(200,
				api.Response{
					Code: consts.SERVICE_MODAL_LOGOUT_CODE,
					Msg:  "您的账户Token失效",
				},
			)
			c.Abort()
			return
		}
		// 判断token是否已过期
		var claims *util.CustomClaims
		var err error
		claims, err = util.ParseToken(parts[1])
		if err != nil {
			var code string
			if err.Error() == "invalid Issuer" {
				code = consts.SERVICE_MODAL_LOGOUT_CODE
			} else {
				code = consts.SERVICE_REFRESH_TOKEN_CODES
				jwt := &model.JwtBlacklist{Jwt: parts[1]}
				if err := service.JwtServiceApp().JwtAddBlacklist(jwt); err != nil {
					logger.Log().Error("jwt", "jwt没有拉入黑名单", err)
					c.JSON(200, util.ServerErrorResponse("用户登出失败,jwt没有成功拉入黑名单", err))
					c.Abort()
					return
				}
			}
			logger.Log().Info("jwt", "invalid", err)
			c.JSON(200,
				api.Response{
					Code: code,
					Msg:  err.Error(),
				},
			)
			c.Abort()
			return
		}

		// 判断用户是否被禁用
		if claims.User.Status != consts.UserModelStatusEnabled {
			c.JSON(200,
				api.Response{
					Code: consts.SERVICE_MODAL_LOGOUT_CODE,
					Msg:  "用户已被禁用",
				},
			)
			c.Abort()
			return
		}

		// 判断用户的updateAt是否在token签发之后
		var updateAt time.Time
		if err = model.DB.Model(&model.User{}).Where("id = ?", claims.User.ID).Select("updated_at").Scan(&updateAt).Error; err != nil || updateAt != claims.User.UpdatedAt {
			c.JSON(200,
				api.Response{
					Code: consts.SERVICE_REFRESH_TOKEN_CODES,
					Msg:  "用户信息或关联角色有变更",
				},
			)
			c.Abort()
			return
		}

		// 将当前请求的claims信息保存到请求的上下文c上
		c.Set("user", &claims.User)
		c.Set("roles", &claims.Roles)
		c.Next() // 后续的处理函数可以用过c.Get("user")来获取当前请求的用户信息
	}
}
