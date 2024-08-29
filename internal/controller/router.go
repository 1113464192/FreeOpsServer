package controller

import (
	"FreeOps/internal/middleware"
	_ "FreeOps/swagger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRoute() *gin.Engine {
	r := gin.Default()
	// if configs.Conf.System.Mode == "product" {
	// gin.SetMode(gin.ReleaseMode)
	// swagger.SwaggerInfo.Host = "127.0.0.1:9080"
	// }
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Use(middleware.Cors())
	r.GET("ping", Test)
	// ---------登录----------
	authRoute := r.Group("auth")
	{
		authRoute.POST("login", UserLogin)
		authRoute.POST("refreshToken", RefreshToken)
		authRoute.GET("error", CustomError)
		authRoute.GET("constant-routes", GetConstantRoutes) // 获取所有常量路由
	}

	// ---------Git-Webhook相关----------
	//gitWebhookRouter := r.Group("git-gitWebhook")
	//{
	//	gitWebhookRouter.POST("github/:pid/:hid", HandleGithubWebhook)          // 接收github的webhook做处理
	//	gitWebhookRouter.POST("gitlab/:pid/:hid", HandleGitlabWebhook)          // 接收gitlab的webhook做处理
	//	gitWebhookRouter.PATCH("project-update-status", UpdateGitWebhookStatus) // 更改git-webhook记录的状态码
	//}
	// ------------验证相关------------
	r.Use(middleware.JWTAuthMiddleware()).Use(middleware.CasbinHandler()).Use(middleware.UserRecord())
	{
		// -------------接口权限测试--------------
		r.GET("ping2", Test2)

		// ------------用户相关------------
		userRoute := r.Group("users")
		{
			userRoute.POST("", UpdateUser)                          // 新增/修改用户
			userRoute.GET("", GetUsers)                             // 查询用户切片
			userRoute.GET("privilege", GetUserPrivilege)            // 查询用户权限
			userRoute.DELETE("", DeleteUsers)                       // 删除用户
			userRoute.PATCH("password", ChangeUserPassword)         // 修改用户密码
			userRoute.POST("logout", UserLogout)                    // 登出
			userRoute.GET("history-action", GetUserRecordLogs)      // 查询用户所有的历史操作
			userRoute.GET("history-month-exist", GetUserRecordDate) // 查询有多少个月份表可供查询
			userRoute.PUT("ssh-key", UpdateSSHKey)                  // 添加私钥
			userRoute.PUT("bind-roles", BindUserRoles)              // 用户绑定角色
			userRoute.GET("roles", GetUserRoles)                    // 查看用户所有角色
		}

		// ------------角色相关--------------
		roleRoute := r.Group("roles")
		{
			roleRoute.POST("", UpdateRole)                   // 新增/修改角色
			roleRoute.GET("", GetRoles)                      // 获取角色列表
			roleRoute.GET("all-summary", GetAllRolesSummary) // 获取所有角色的简略信息
			roleRoute.DELETE("", DeleteRoles)                // 删除角色
			roleRoute.PUT("bind", BindRoleRelation)          // 角色绑定关系
			roleRoute.GET("menus", GetRoleMenus)             // 获取角色的菜单
			roleRoute.GET("apis", GetRoleApis)               // 获取角色的API
			roleRoute.GET("buttons", GetRoleButtons)         // 获取角色的按钮
			roleRoute.GET("users", GetRoleUsers)             // 获取角色绑定的用户
		}

		// ------------Button相关------------
		buttonRoute := r.Group("buttons")
		{
			buttonRoute.POST("", UpdateButton)    // 新增/修改按钮
			buttonRoute.GET("", GetButtons)       // 获取按钮列表
			buttonRoute.DELETE("", DeleteButtons) // 删除按钮
		}

		// ------------菜单相关--------------
		menuRoute := r.Group("menus")
		{
			menuRoute.POST("", UpdateMenu)                // 新增/修改组
			menuRoute.GET("", GetMenus)                   // 获取菜单信息
			menuRoute.DELETE("", DeleteMenu)              // 删除菜单
			menuRoute.GET("buttons", GetMenuButtons)      // 获取菜单下所有按钮
			menuRoute.GET("all-pages", GetAllPages)       // 获取所有页面
			menuRoute.GET("user-routes", GetUserRoutes)   // 获取用户路由
			menuRoute.GET("tree", GetMenuTree)            // 获取菜单树
			menuRoute.GET("is-route-exist", IsRouteExist) // 判断路由是否存在
		}

		// -------------API相关---------------
		apiRoute := r.Group("apis")
		{
			apiRoute.POST("", UpdateApi)       // 新增/修改API
			apiRoute.GET("", GetApis)          // 获取API列表
			apiRoute.DELETE("", DeleteApi)     // 删除API
			apiRoute.GET("group", GetApiGroup) // 获取存在的API组
			apiRoute.GET("tree", GetApiTree)   // 获取API树
		}
	}
	return r
}
