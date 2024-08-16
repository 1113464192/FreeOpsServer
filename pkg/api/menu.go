package api

type UpdateMenuReq struct {
	ID              uint   `json:"id" form:"id"`                                          // 修改才需要传，没有传算新增
	Status          uint   `json:"status" form:"status" binding:"required,oneof=1 2"`     // 1(enabled),2(disabled)
	ParentId        uint   `json:"parentId" form:"parentId"`                              // 一级菜单填0
	MenuType        uint   `json:"menuType" form:"menuType" binding:"required,oneof=1 2"` // 1(directory),2(menu)
	MenuName        string `json:"menuName" form:"menuName" binding:"required"`
	RouteName       string `json:"routeName" form:"routeName" binding:"required"`
	RoutePath       string `json:"routePath" form:"routePath" binding:"required"`
	Component       string `json:"component" form:"component"`
	Order           uint   `json:"order" form:"order" binding:"required"`
	I18nKey         string `json:"i18nKey" form:"i18nKey"`
	Icon            string `json:"icon" form:"icon"`
	IconType        uint   `json:"iconType" form:"iconType" binding:"oneof=1 2"` // 1(iconify),2(local)
	MultiTab        bool   `json:"multiTab" form:"multiTab"`                     // 是多标签页则传1，否则不传
	HideInMenu      bool   `json:"hideInMenu" form:"hideInMenu"`                 // 隐藏标签页则传1，否则不传
	KeepAlive       bool   `json:"keepAlive" form:"keepAlive"`
	ShowRole        bool   `json:"showRole" form:"showRole"`
	ActiveMenu      string `json:"activeMenu" form:"activeMenu"`                             // 没有则不传
	IsConstantRoute bool   `json:"constant" form:"constant"`                                 // 是常量(访问该路由时将不会进行登录验证和权限验证)路由则传1，否则不传
	FixedIndexInTab uint8  `json:"fixedIndexInTab" form:"fixedIndexInTab" binding:"max=127"` // 固定在页签中的序号
	Props           any    `json:"props" form:"props"`                                       // 路由属性作跳转路由则直接传true，否则传完整json
	Query           string `json:"query" form:"query"`                                       // 查询条件，json传送
	Href            string `json:"href" form:"href"`
}

type MenuButtonRes struct {
	Code string `json:"code" form:"code"`
	Desc string `json:"desc" form:"desc"`
}

type GetMenusReq struct {
	ID       uint   `form:"id" json:"id"`
	Status   uint   `form:"status" json:"status"  binding:"required,oneof=1 2"`
	RoleName string `form:"roleName" json:"roleName"`
	RoleCode string `form:"roleCode" json:"roleCode"`
	PageInfo
}

type MenuRes struct {
	ID              uint            `json:"id" form:"id"` // 修改才需要传，没有传算新增
	Status          uint            `json:"status" form:"status" binding:"required"`
	ParentId        uint            `json:"parentId" form:"parentId" binding:"required"`
	MenuType        uint            `json:"menuType" form:"menuType" binding:"required"`
	MenuName        string          `json:"menuName" form:"menuName" binding:"required"`
	RouteName       string          `json:"routeName" form:"routeName" binding:"required"`
	RoutePath       string          `json:"routePath" form:"routePath" binding:"required"`
	Component       string          `json:"component" form:"component"`
	Order           uint            `json:"order" form:"order" binding:"required"`
	I18nKey         string          `json:"i18nKey" form:"i18nKey"`
	Icon            string          `json:"icon" form:"icon"`
	IconType        uint            `json:"iconType" form:"iconType"`
	MultiTab        bool            `json:"multiTab,omitempty" form:"multiTab"`
	HideInMenu      bool            `json:"hideInMenu,omitempty" form:"hideInMenu"`
	KeepAlive       bool            `json:"keepAlive,omitempty" form:"keepAlive"`
	ActiveMenu      string          `json:"activeMenu,omitempty" form:"activeMenu"`
	IsConstantRoute bool            `json:"constant,omitempty" form:"constant"`               // 是常量(访问该路由时将不会进行登录验证和权限验证)路由则传1，否则不传
	FixedIndexInTab uint8           `json:"fixedIndexInTab,omitempty" form:"fixedIndexInTab"` // 固定在页签中的序号
	Props           any             `json:"props,omitempty" form:"props"`                     // 路由属性作跳转路由则直接传true，否则传完整json
	Query           string          `json:"query,omitempty" form:"query"`                     // 查询条件，json传送
	Href            string          `json:"href,omitempty" form:"href"`
	RoleCodes       []string        `json:"roles,omitempty" form:"roles"` // 角色Codes
	Buttons         []MenuButtonRes `json:"buttons,omitempty"`
	Children        []MenuRes       `json:"children,omitempty"`
}

type GetMenuRes struct {
	MenuRes  []MenuRes `json:"records" form:"records"`
	Page     int       `json:"current" form:"current"`
	PageSize int       `json:"size" form:"size"`
	Total    int64     `json:"total"`
}

type GetConstantRoutesMetaRes struct {
	Title           string   `json:"title" form:"title"`
	I18nKey         string   `json:"i18nKey" form:"i18nKey"`
	Order           uint     `json:"order" form:"order"`
	Icon            string   `json:"icon,omitempty" form:"icon"`
	LocalIcon       string   `json:"localIcon,omitempty" form:"localIcon"`
	Href            string   `json:"href,omitempty" form:"href"`
	HideInMenu      bool     `json:"hideInMenu,omitempty" form:"hideInMenu"`
	ActiveMenu      string   `json:"activeMenu,omitempty" form:"activeMenu"`
	MultiTab        bool     `json:"multiTab,omitempty" form:"multiTab"`
	KeepAlive       bool     `json:"keepAlive,omitempty" form:"keepAlive"`
	RoleCodes       []string `json:"roles,omitempty" form:"roles"`                     // 角色Codes
	IsConstantRoute bool     `json:"constant,omitempty" form:"constant"`               // 是常量(访问该路由时将不会进行登录验证和权限验证)路由则传1，否则不传
	FixedIndexInTab uint8    `json:"fixedIndexInTab,omitempty" form:"fixedIndexInTab"` // 固定在页签中的序号
	Query           string   `json:"query,omitempty" form:"query"`                     // 查询条件，json传送
}

type GetConstantRoutesRes struct {
	Name      string                   `json:"name" form:"name"`
	Path      string                   `json:"path" form:"path"`
	Component string                   `json:"component" form:"component"`
	Props     any                      `json:"props,omitempty" form:"props"` // 路由属性作跳转路由则直接传true，否则传完整json
	ParentId  uint                     `json:"-"`
	Meta      GetConstantRoutesMetaRes `json:"meta,omitempty" form:"meta"`
	Children  []GetConstantRoutesRes   `json:"children,omitempty" form:"children"`
}

type GetMenuTreeRes struct {
	Id       uint             `json:"id" form:"id"`
	Label    string           `json:"label" form:"label"`
	Pid      uint             `json:"pid" form:"pid"`
	Children []GetMenuTreeRes `json:"children,omitempty" form:"children"`
}
