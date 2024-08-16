package model

import "gorm.io/gorm"

type Menu struct {
	gorm.Model
	Status          uint    `json:"status" gorm:"default:1;comment:状态: 1(enabled),2(disabled);type:tinyint"`
	ParentId        uint    `json:"parentId" gorm:"comment:父菜单;index"`
	MenuType        uint    `json:"menuType" gorm:"comment:菜单类型: 1(directory),2(menu);type:tinyint"`
	MenuName        string  `json:"menuName" gorm:"type:varchar(30);unique;comment:菜单名称"`
	RouteName       string  `json:"routeName" gorm:"type:varchar(30);uniqueIndex;comment:路由名称"`
	RoutePath       string  `json:"routePath" gorm:"type:varchar(50);comment:路由地址"`
	Component       string  `json:"component" gorm:"type:varchar(100);comment:组件"`
	Order           uint    `json:"order" gorm:"type:tinyint;comment:排序标记"`
	I18nKey         string  `json:"i18nKey" gorm:"type:varchar(100);comment:国际化key"`
	Icon            string  `json:"icon" gorm:"type:varchar(100);comment:图标"`
	IconType        uint    `json:"iconType" gorm:"type:tinyint;comment:图标类型: 1(iconify),2(local),"`
	MultiTab        bool    `json:"multiTab" gorm:"comment:是否多标签页: 0(no),1(yes)"`
	HideInMenu      bool    `json:"hideInMenu" gorm:"comment:是否隐藏菜单: 0(no),1(yes)"`
	KeepAlive       bool    `json:"keepAlive" gorm:"comment:是否缓存: 0(no),1(yes)"`
	ShowRole        bool    `json:"showRole" gorm:"comment:是否根据角色显示: 0(no),1(yes)"`
	ActiveMenu      *string `json:"activeMenu" gorm:"type:varchar(50);comment:激活菜单(指定当进入某个路由时，哪个菜单项应该被激活)"`
	IsConstantRoute bool    `json:"constant" gorm:"comment:是否是常量(访问该路由时将不会进行登录验证和权限验证)路由: 0(no),1(yes)"`
	FixedIndexInTab uint8   `json:"fixedIndexInTab" gorm:"type:tinyint;comment:固定在页签中的序号"`
	Props           *string `json:"props" gorm:"type:text;comment:路由属性作跳转路由则直接传true，否则传完整json"`
	Query           *string `json:"query" gorm:"type:text;comment:查询条件，json传送"`
	Href            *string `json:"href" gorm:"type:text;comment:外链地址"`
}
