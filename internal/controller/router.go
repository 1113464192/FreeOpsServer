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
	// 所有请求加上api做前缀
	apiRoute := r.Group("api")
	r.GET("ping", Test)
	// ---------登录----------
	authRoute := apiRoute.Group("auth")
	{
		authRoute.POST("login", UserLogin)
		authRoute.POST("refreshToken", RefreshToken)
		authRoute.GET("error", CustomError)
		authRoute.GET("constant-routes", GetConstantRoutes) // 获取所有常量路由
	}
	// --------不走头信息权限验证的功能接口----------
	//websocket
	nonAuthOpsRoute := apiRoute.Group("ops")
	{
		nonAuthOpsRoute.GET("task-need-approve", GetOpsTaskNeedApprove) // 查询用户是否有待审批的任务
		nonAuthOpsRoute.GET("task-running-ws", GetOpsTaskRunningWS)     // 实时查看运行中的任务
	}

	// ---------Git-Webhook相关----------
	//gitWebhookRouter := r.Group("git-gitWebhook")
	//{
	//	gitWebhookRouter.POST("github/:pid/:hid", HandleGithubWebhook)          // 接收github的webhook做处理
	//	gitWebhookRouter.POST("gitlab/:pid/:hid", HandleGitlabWebhook)          // 接收gitlab的webhook做处理
	//	gitWebhookRouter.PATCH("project-update-status", UpdateGitWebhookStatus) // 更改git-webhook记录的状态码
	//}
	// ------------验证相关------------
	apiRoute.Use(middleware.JWTAuthMiddleware()).Use(middleware.CasbinHandler()).Use(middleware.UserRecord())
	{
		// -------------接口权限测试--------------
		apiRoute.GET("ping2", Test2)
		// ------------home相关------------
		homeRoute := apiRoute.Group("home")
		{
			homeRoute.GET("info", GetHomeInfo) // 获取首页信息
		}

		// ------------用户相关------------
		userRoute := apiRoute.Group("users")
		{
			userRoute.POST("", UpdateUser)                          // 新增/修改用户
			userRoute.GET("", GetUsers)                             // 查询用户切片
			userRoute.GET("privilege", GetUserPrivilege)            // 查询用户权限
			userRoute.DELETE("", DeleteUsers)                       // 删除用户
			userRoute.PATCH("password", ChangeUserPassword)         // 修改用户密码
			userRoute.POST("logout", UserLogout)                    // 登出
			userRoute.GET("history-action", GetUserRecordLogs)      // 查询用户所有的历史操作
			userRoute.GET("history-month-exist", GetUserRecordDate) // 查询有多少个月份表可供查询
			userRoute.PUT("bind-roles", BindUserRoles)              // 用户绑定角色
			userRoute.GET("roles", GetUserRoles)                    // 查看用户所有角色
			userRoute.GET("project-options", GetUserProjectOptions) // 查看用户所有项目选项
		}

		// ------------角色相关--------------
		roleRoute := apiRoute.Group("roles")
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
			roleRoute.GET("projects", GetRoleProjects)       // 获取角色绑定的项目
		}

		// ------------Button相关------------
		buttonRoute := apiRoute.Group("buttons")
		{
			buttonRoute.POST("", UpdateButton)             // 新增/修改按钮
			buttonRoute.GET("", GetButtons)                // 获取按钮列表
			buttonRoute.DELETE("menus", DeleteMenuButtons) // 删除按钮
		}

		// ------------菜单相关--------------
		menuRoute := apiRoute.Group("menus")
		{
			menuRoute.POST("", UpdateMenu)                // 新增/修改菜单
			menuRoute.GET("", GetMenus)                   // 获取菜单信息
			menuRoute.DELETE("", DeleteMenu)              // 删除菜单
			menuRoute.GET("buttons", GetMenuButtons)      // 获取菜单下所有按钮
			menuRoute.GET("all-pages", GetAllPages)       // 获取所有页面
			menuRoute.GET("user-routes", GetUserRoutes)   // 获取用户路由
			menuRoute.GET("tree", GetMenuTree)            // 获取菜单树
			menuRoute.GET("is-route-exist", IsRouteExist) // 判断路由是否存在
		}

		// -------------API相关---------------
		apisRoute := apiRoute.Group("apis")
		{
			apisRoute.POST("", UpdateApi)       // 新增/修改API
			apisRoute.GET("", GetApis)          // 获取API列表
			apisRoute.DELETE("", DeleteApi)     // 删除API
			apisRoute.GET("group", GetApiGroup) // 获取存在的API组
			apisRoute.GET("tree", GetApiTree)   // 获取API树
		}

		// ------------项目相关--------------
		projectRoute := apiRoute.Group("projects")
		{
			projectRoute.POST("", UpdateProject)                    // 新增/修改项目
			projectRoute.GET("", GetProjects)                       // 查询项目
			projectRoute.GET("all-summary", GetProjectList)         // 获取项目列表
			projectRoute.DELETE("", DeleteProjects)                 // 删除项目
			projectRoute.GET("hosts", GetProjectHosts)              // 查询项目关联的服务器
			projectRoute.GET("games", GetProjectGames)              // 查询项目关联的游戏
			projectRoute.GET("assets-total", GetProjectAssetsTotal) // 查询项目各资产总数
		}
		// ----------服务器相关---------------
		hostRoute := apiRoute.Group("hosts")
		{
			hostRoute.POST("", UpdateHost)              // 新增/修改服务器
			hostRoute.GET("", GetHosts)                 // 查询服务器
			hostRoute.DELETE("", DeleteHosts)           // 删除服务器
			hostRoute.GET("summary", GetHostList)       // 获取服务器列表
			hostRoute.GET("game-info", GetHostGameInfo) // 获取服务器的游戏信息
		}
		// ---------游戏服务相关------------
		gameRoute := apiRoute.Group("games")
		{
			gameRoute.POST("", UpdateGame)              // 新增/修改游戏
			gameRoute.GET("", GetGames)                 // 查询游戏
			gameRoute.DELETE("", DeleteGames)           // 删除游戏
			gameRoute.PATCH("status", UpdateGameStatus) // 更新游戏状态
		}
		// --------运维操作服务相关----------
		// 便于运维自定义脚本路径参数模板、设定运营参数自动填入模板变量、设定任务涵盖模板执行顺序、设定每个任务的检查脚本路径
		opsRoute := apiRoute.Group("ops")
		{
			opsRoute.POST("template", UpdateOpsTemplate)               // 创建/修改 模板
			opsRoute.GET("template", GetOpsTemplate)                   // 查看模板
			opsRoute.DELETE("template", DeleteOpsTemplate)             // 删除模板
			opsRoute.POST("param-template", UpdateOpsParamsTemplate)   // 创建/修改 获取参数模板 (从运营文案信息获取参数的正则模板)
			opsRoute.GET("param-template", GetOpsParamsTemplate)       // 查看参数
			opsRoute.DELETE("param-template", DeleteOpsParamsTemplate) // 删除参数
			opsRoute.PUT("bind-template-params", BindTemplateParams)   // 绑定模板参数
			opsRoute.GET("template-params", GetTemplateParams)         // 查看模板关联的参数
			// 还需要拼接执行模板顺序的任务接口、执行任务接口、获取任务运行状态接口、任务日志接口
			opsRoute.POST("task", UpdateOpsTask)                          // 创建/修改 任务(拼接执行模板顺序的任务)
			opsRoute.DELETE("task", DeleteOpsTask)                        // 删除任务
			opsRoute.GET("task", GetOpsTask)                              // 查看任务
			opsRoute.POST("commands", GetOpsTaskTmpCommands)              // 查看根据参数会生成的命令
			opsRoute.POST("submit-task", SubmitOpsTask)                   // 提交任务
			opsRoute.PUT("approve-task", ApproveOpsTask)                  // 用户审批任务
			opsRoute.GET("task-pending", GetUserTaskPending)              // 查询用户待审批的任务
			opsRoute.POST("run-task-check-script", RunOpsTaskCheckScript) // 执行并等待运营检查脚本返回结果
			opsRoute.GET("task-log", GetOpsTaskLog)                       // 查看任务日志
		}
		// ---------云平台相关------------
		// 云平台一切操作运维脚本(因为脚本变动频繁，且便于运维随时配合自动化修改,平台只需要注意传参的参数即可)，云平台相关运维脚本路径和参数写死在service层
		// 不建议平台写死代码，否则改动过于频繁，不能及时配合各项目运维自动化脚本实时改动
		cloudRoute := apiRoute.Group("clouds")
		{
			// 创建
			CloudCreateRoute := cloudRoute.Group("create")
			{
				CloudCreateRoute.POST("project", CreateCloudProject) // 创建云项目
				// 一切配置如: vpc、dataSize、period等，都通过运维的ini文件，便于随时变更与批量购买
				CloudCreateRoute.POST("host", CreateCloudHost) // 创建云服务器
			}
			// 更新
			CloudUpdateRoute := cloudRoute.Group("update")
			{
				CloudUpdateRoute.POST("project", UpdateCloudProject) // 更新云项目
			}
			// 查询
			CloudQueryRoute := cloudRoute.Group("query")
			{
				CloudQueryRoute.GET("project", GetCloudProjectId) // 查询云项目ID
			}
		}
	}
	return r
}
