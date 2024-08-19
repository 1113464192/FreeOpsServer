package service

import (
	"FreeOps/internal/consts"
	"FreeOps/internal/model"
	"FreeOps/pkg/api"
	"encoding/json"
	"errors"
	"fmt"
)

type MenuService struct{}

var insMenu MenuService

func MenuServiceApp() *MenuService {
	return &insMenu
}

func (s *MenuService) UpdateMenu(params *api.UpdateMenuReq) (err error) {
	var (
		menu  model.Menu
		count int64
		props string
	)

	// 判断父ID是否存在于菜单
	if params.ParentId != 0 {
		if err = model.DB.Model(&model.Menu{}).Where("id = ?", params.ParentId).Count(&count).Error; err != nil || count < 1 {
			return fmt.Errorf("menu ID不存在: %d, 或有错误信息: %v", params.ID, err)
		}
	}

	// 判断ActiveMenu是否存在于路由
	if params.ActiveMenu != "" {
		if err = model.DB.Model(&model.Menu{}).Where("route_name = ?", params.ActiveMenu).Count(&count).Error; err != nil || count < 1 {
			return fmt.Errorf("route_name不存在: %s, 或有错误信息: %v", params.ActiveMenu, err)
		}
	}

	// 判断Props是否符合规范
	switch params.Props.(type) {
	case nil:
		props = ""
	case bool:
		if propsBool := params.Props.(bool); propsBool == true {
			props = consts.MenuModelPropsIsTrue
		} else {
			props = ""
		}
	case string:
		props = params.Props.(string)
	default:
		return fmt.Errorf("props类型错误: %T", params.Props)
	}

	if params.ID != 0 {
		if err = model.DB.Model(&model.Menu{}).Where("id = ?", params.ID).Count(&count).Error; count != 1 || err != nil {
			return fmt.Errorf("menu ID不存在: %d, 或有错误信息: %v", params.ID, err)
		}

		// 判断menu_name是否和现有菜单重复
		if err = model.DB.Model(&menu).Where("menu_name = ? AND id != ?", params.MenuName, params.ID).Count(&count).Error; err != nil || count > 0 {
			return fmt.Errorf("菜单名(%s)已被使用, 或有错误信息: %v", params.MenuName, err)
		}

		if err = model.DB.Where("id = ?", params.ID).First(&menu).Error; err != nil {
			return fmt.Errorf("数据库查询失败: %v", err)
		}

		menu.Status = params.Status
		menu.ParentId = params.ParentId
		menu.MenuType = params.MenuType
		menu.MenuName = params.MenuName
		menu.RouteName = params.RouteName
		menu.RoutePath = params.RoutePath
		menu.Component = params.Component
		menu.Order = params.Order
		menu.I18nKey = params.I18nKey
		menu.IconType = params.IconType
		menu.MultiTab = params.MultiTab
		menu.HideInMenu = params.HideInMenu
		menu.KeepAlive = params.KeepAlive
		menu.ShowRole = params.ShowRole
		if params.Icon == "" {
			menu.Icon = nil
		} else {
			menu.Icon = &params.Icon
		}
		if params.ActiveMenu == "" {
			menu.ActiveMenu = nil
		} else {
			menu.ActiveMenu = &params.ActiveMenu
		}
		menu.IsConstantRoute = params.IsConstantRoute
		menu.FixedIndexInTab = params.FixedIndexInTab
		if props == "" {
			menu.Props = nil
		} else {
			menu.Props = &props
		}
		if params.Query == "" {
			menu.Query = nil
		} else {
			menu.Query = &params.Query
		}
		if params.Href == "" {
			menu.Href = nil
		} else {
			menu.Href = &params.Href
		}
		if err = model.DB.Save(&menu).Error; err != nil {
			return fmt.Errorf("数据保存失败: %v", err)
		}

		return err
	} else {
		err = model.DB.Model(&menu).Where("menu_name = ? OR route_name = ?", params.MenuName, params.RouteName).Count(&count).Error
		// 总数大于0或者有错误就返回
		if err != nil {
			return fmt.Errorf("查询菜单失败: %v", err)
		} else if count > 0 {
			return fmt.Errorf("菜单名(%s)或路由名(%s)已存在", params.MenuName, params.RouteName)
		}

		menu = model.Menu{
			Status:     params.Status,
			ParentId:   params.ParentId,
			MenuType:   params.MenuType,
			MenuName:   params.MenuName,
			RouteName:  params.RouteName,
			RoutePath:  params.RoutePath,
			Component:  params.Component,
			Order:      params.Order,
			I18nKey:    params.I18nKey,
			IconType:   params.IconType,
			MultiTab:   params.MultiTab,
			HideInMenu: params.HideInMenu,
			ShowRole:   params.ShowRole,
			KeepAlive:  params.KeepAlive,
		}
		if params.Icon == "" {
			menu.Icon = nil
		} else {
			menu.Icon = &params.Icon
		}
		if params.ActiveMenu == "" {
			menu.ActiveMenu = nil
		} else {
			menu.ActiveMenu = &params.ActiveMenu
		}
		menu.IsConstantRoute = params.IsConstantRoute
		menu.FixedIndexInTab = params.FixedIndexInTab
		if props == "" {
			menu.Props = nil
		} else {
			menu.Props = &props
		}
		if params.Query == "" {
			menu.Query = nil
		} else {
			menu.Query = &params.Query
		}
		if params.Href == "" {
			menu.Href = nil
		} else {
			menu.Href = &params.Href
		}
		tx := model.DB.Begin()

		if err = tx.Create(&menu).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("创建菜单失败: %v", err)
		}
		// 菜单默认绑定管理员role
		var adminRoleId uint
		if err = tx.Model(&model.Role{}).Where("role_code = ?", consts.RoleModelAdminCode).Select("id").Scan(&adminRoleId).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("查询管理员角色ID失败: %v", err)
		}
		if err = tx.Create(&model.MenuRole{
			MenuId: menu.ID,
			RoleId: adminRoleId,
		}).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("创建菜单角色关联失败: %v", err)
		}
		tx.Commit()
		return err
	}
}

func (s *MenuService) GetMenus(param *api.IdPageReq) (*api.GetMenuRes, error) {
	var err error
	var count int64
	var menus []model.Menu
	if err = model.DB.Model(&model.Menu{}).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("查询菜单总数失败: %v", err)
	}
	if param.Id != 0 {
		if err = model.DB.First(&menus, param.Id).Error; err != nil {
			return nil, fmt.Errorf("查询菜单失败: %v", err)
		}
	} else if param.Page == 0 && param.PageSize == 0 {
		if err = model.DB.Find(&menus).Error; err != nil {
			return nil, fmt.Errorf("查询菜单失败: %v", err)
		}
	} else {
		if err = model.DB.Offset((param.Page - 1) * param.PageSize).Limit(param.PageSize).Find(&menus).Error; err != nil {
			return nil, fmt.Errorf("查询菜单失败: %v", err)
		}
	}
	var res *[]api.MenuRes
	var result api.GetMenuRes
	res, err = s.GetResults(&menus)
	if err != nil {
		return nil, err
	}
	result = api.GetMenuRes{
		MenuRes:  *res,
		Page:     param.Page,
		PageSize: param.PageSize,
		Total:    count,
	}
	return &result, err
}

// 删除菜单目录时查询子菜单中的菜单并删除，如果子菜单中有其它菜单目录，则递归返回执行
func (s MenuService) DeleteMenuDirectory(id uint) (err error) {
	var childIds []uint
	if err = model.DB.Model(&model.Menu{}).Where("parent_id = ?", id).Pluck("id", &childIds).Error; err != nil {
		return fmt.Errorf("查询子菜单失败: %v", err)
	}
	if len(childIds) > 0 {
		// 递归删除子菜单
		if err = s.DeleteMenus(childIds); err != nil {
			return fmt.Errorf("删除子菜单失败: %v", err)
		}
	}
	return err
}

func (s MenuService) DeleteMenus(ids []uint) (err error) {
	tx := model.DB.Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()
	var buttonIds []uint
	if err = tx.Model(&model.Button{}).Where("menu_id IN (?)", ids).Pluck("id", &buttonIds).Error; err != nil {
		return fmt.Errorf("查询按钮失败: %v", err)
	}
	if err = tx.Where("button_id IN (?)", buttonIds).Delete(&model.RoleButton{}).Error; err != nil {
		return fmt.Errorf("删除角色按钮关系失败 %d: %v", ids, err)
	}
	if err = tx.Where("id IN (?)", buttonIds).Delete(&model.Button{}).Error; err != nil {
		return fmt.Errorf("删除菜单按钮失败: %v", err)
	}
	for _, id := range ids {
		// 删除角色菜单关联
		if err = tx.Where("menu_id = ?", id).Delete(&model.MenuRole{}).Error; err != nil {
			return fmt.Errorf("删除角色菜单关联失败: %v", err)
		}

		var menuType int
		if err = model.DB.Model(&model.Menu{}).Where("id = ?", id).Select("menu_type").Scan(&menuType).Error; err != nil {
			return fmt.Errorf("查询菜单类型失败: %v", err)
		}
		if menuType == consts.MenuModelMenuTypeIsDirectory {
			if err = s.DeleteMenuDirectory(id); err != nil {
				return err
			}
			if err = tx.Where("id = ?", id).Delete(&model.Menu{}).Error; err != nil {
				return fmt.Errorf("删除菜单失败: %v", err)
			}
		} else {
			// 删除菜单
			if err = tx.Where("id IN (?)", ids).Delete(&model.Menu{}).Error; err != nil {
				return fmt.Errorf("删除菜单失败: %v", err)
			}
		}
	}
	tx.Commit()
	return err
}

func (s *MenuService) GetMenuButtons(id uint) (res *[]api.GetButtonReq, err error) {
	var count int64

	if err = model.DB.Model(&model.Menu{}).Where("id = ?", id).Count(&count).Error; count != 1 || err != nil {
		return nil, fmt.Errorf("menu ID不存在: %d, 或有错误信息: %v", id, err)

	}

	var buttons []model.Button
	if err = model.DB.Model(&model.Button{}).Where("menu_id = ?", id).Find(&buttons).Error; err != nil {
		return nil, fmt.Errorf("查询按钮失败: %v", err)
	}

	res, err = ButtonServiceApp().GetResults(&buttons)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (s *MenuService) GetAllPages() (routeNames []string, err error) {
	if err = model.DB.Model(&model.Menu{}).Where("menu_type = ?", consts.MenuModelMenuTypeIsMenu).Pluck("route_name", &routeNames).Error; err != nil {
		return routeNames, fmt.Errorf("查询页面失败: %v", err)
	}
	return routeNames, err
}

func (s *MenuService) GetConstantRoutes() (result *[]api.GetRoutesRes, err error) {
	var menus []model.Menu
	if err = model.DB.Model(&model.Menu{}).Where("is_constant_route = ?", consts.MysqlGormBoolTrue).Find(&menus).Error; err != nil {
		return nil, fmt.Errorf("查询路由失败: %v", err)
	}
	result, err = s.GetRoutesRes(&menus)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (s *MenuService) GetUserRoutes(roles *[]model.Role) (result *[]api.GetRoutesRes, err error) {
	var menus []model.Menu
	if err = model.DB.Model(&model.Menu{}).Where("show_role = ?", consts.MysqlGormBoolFalse).Find(&menus).Error; err != nil {
		return nil, fmt.Errorf("查询路由失败: %v", err)
	}
	var limitRoleMenus []model.Menu
	if err = model.DB.Model(&model.Menu{}).Where("show_role = ?", consts.MysqlGormBoolTrue).Find(&limitRoleMenus).Error; err != nil {
		return nil, fmt.Errorf("查询路由失败: %v", err)
	}

	var count int64
	for _, menu := range limitRoleMenus {
		// 看看menu绑定的角色和roles是否有重合，有的话append到menus
		for _, role := range *roles {
			if role.RoleCode == consts.RoleModelAdminCode {
				menus = append(menus, menu)
				break
			}
			if err = model.DB.Model(&model.Menu{}).
				Joins("JOIN menu_role ON menu_role.menu_id = menu.id").
				Where("menu_role.role_id = ? AND menu.id = ?", role.ID, menu.ID).Count(&count).Error; err != nil {
				return nil, fmt.Errorf("查询菜单角色关联失败: %v", err)
			}
			if count > 0 {
				menus = append(menus, menu)
				break
			}
		}
	}

	result, err = s.GetRoutesRes(&menus)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (s *MenuService) GetMenuTree() (*[]api.GetMenuTreeRes, error) {
	var menus []model.Menu
	var err error
	if err = model.DB.Find(&menus).Error; err != nil {
		return nil, fmt.Errorf("查询菜单失败: %v", err)
	}
	menuMap := make(map[uint]*api.GetMenuTreeRes)
	// 将所有菜单项转换为指针，便于修改
	for _, menu := range menus {
		menuMap[menu.ID] = &api.GetMenuTreeRes{
			Id:    menu.ID,
			Label: menu.MenuName,
			Pid:   menu.ParentId,
		}
	}

	// 构建父子关系
	for _, menu := range menus {
		if menu.ParentId != 0 { // 非顶级菜单
			if parent, ok := menuMap[menu.ParentId]; ok {
				if parent.Children == nil {
					parent.Children = &[]api.GetMenuTreeRes{}
				}
				*parent.Children = append(*parent.Children, *menuMap[menu.ID])
			}
		}
	}

	var result []api.GetMenuTreeRes
	// 从map中提取所有顶级菜单
	for _, menu := range menuMap {
		if menu.Pid == 0 {
			result = append(result, *menu)
		}
	}

	return &result, err
}

func (s *MenuService) IsRouteExist(routeName string) (bool, error) {
	var count int64
	if err := model.DB.Model(&model.Menu{}).Where("route_name = ?", routeName).Count(&count).Error; err != nil {
		return false, fmt.Errorf("查询路由失败: %v", err)
	}
	return count > 0, nil
}

func (s *MenuService) GetRoutesRes(menus *[]model.Menu) (*[]api.GetRoutesRes, error) {
	var (
		result []api.GetRoutesRes
		err    error
	)
	menuMap := make(map[uint]*api.GetRoutesRes)
	for _, menu := range *menus {
		res := api.GetRoutesRes{
			Name:      menu.RouteName,
			Path:      menu.RoutePath,
			Component: menu.Component,
			ParentId:  menu.ParentId,
			Meta: api.GetRoutesMetaRes{
				Title:   menu.RouteName,
				I18nKey: menu.I18nKey,
				Order:   menu.Order,
			},
		}
		if menu.Props != nil {
			if *menu.Props == "true" {
				res.Props = true
			} else {
				if err = json.Unmarshal([]byte(*menu.Props), &res.Props); err != nil {
					return nil, fmt.Errorf("解析Props属性失败: %v", err)
				}
			}
		}
		if menu.Icon != nil {
			if menu.IconType == consts.MenuModeIconTypeIsIconify {
				res.Meta.Icon = *menu.Icon
			} else if menu.IconType == consts.MenuModeIconTypeIsLocal {
				res.Meta.LocalIcon = *menu.Icon
			} else {
				return nil, fmt.Errorf("菜单ID %d 的图标类型错误为: %d", menu.ID, menu.IconType)
			}
		}
		if menu.Href != nil {
			res.Meta.Href = *menu.Href
		}
		res.Meta.HideInMenu = menu.HideInMenu
		if menu.ActiveMenu != nil {
			res.Meta.ActiveMenu = *menu.ActiveMenu
		}
		res.Meta.MultiTab = menu.MultiTab
		res.Meta.KeepAlive = menu.KeepAlive
		if menu.ShowRole {
			if err = model.DB.Model(&model.Role{}).
				Joins("JOIN menu_role ON menu_role.role_id = role.id").
				Where("menu_role.menu_id = ?", menu.ID).
				Pluck("role_code", &res.Meta.RoleCodes).Error; err != nil {
				return nil, fmt.Errorf("查询菜单角色代码失败: %v", err)
			}
		}
		res.Meta.IsConstantRoute = menu.IsConstantRoute
		res.Meta.FixedIndexInTab = menu.FixedIndexInTab
		if menu.Query != nil {
			res.Meta.Query = *menu.Query
		}
		menuMap[menu.ID] = &res
	}

	for _, value := range menuMap {
		if value.ParentId != 0 {
			if parent, ok := menuMap[value.ParentId]; ok {
				if parent.Children == nil {
					parent.Children = &[]api.GetRoutesRes{}
				}
				*parent.Children = append(*parent.Children, *value)
			} else {
				return nil, fmt.Errorf("路由: %s 的父ID路由(%d)不在menuMap中, menuMap：\n %v", value.Name, value.ParentId, menuMap)
			}
		}
	}

	for _, value := range menuMap {
		if value.ParentId == 0 {
			result = append(result, *value)
		}
	}

	return &result, err
}

// 返回菜单结果
func (s *MenuService) GetResults(menuObj any) (*[]api.MenuRes, error) {
	var result []api.MenuRes
	var err error
	var count int64
	var buttons []model.Button
	if menus, ok := menuObj.(*[]model.Menu); ok {
		// 创建一个映射，用于存储菜单ID到菜��项的映射
		menuMap := make(map[uint]*api.MenuRes)
		for _, menu := range *menus {
			res := api.MenuRes{
				ID:              menu.ID,
				Status:          menu.Status,
				ParentId:        menu.ParentId,
				MenuType:        menu.MenuType,
				MenuName:        menu.MenuName,
				RouteName:       menu.RouteName,
				RoutePath:       menu.RoutePath,
				Component:       menu.Component,
				Order:           menu.Order,
				I18nKey:         menu.I18nKey,
				IconType:        menu.IconType,
				MultiTab:        menu.MultiTab,
				KeepAlive:       menu.KeepAlive,
				HideInMenu:      menu.HideInMenu,
				IsConstantRoute: menu.IsConstantRoute,
				FixedIndexInTab: menu.FixedIndexInTab,
			}
			if menu.Icon != nil {
				res.Icon = *menu.Icon
			}
			if menu.ActiveMenu != nil {
				res.ActiveMenu = *menu.ActiveMenu
			}

			if menu.Props != nil {
				if *menu.Props == "true" {
					res.Props = true
				} else {
					if err = json.Unmarshal([]byte(*menu.Props), &res.Props); err != nil {
						return nil, fmt.Errorf("解析Props属性失败: %v", err)
					}
				}
			}
			if menu.Href != nil {
				res.Href = *menu.Href
			}
			if menu.ActiveMenu != nil {
				res.ActiveMenu = *menu.ActiveMenu
			}
			if menu.ShowRole {
				if err = model.DB.Model(&model.Role{}).
					Joins("JOIN menu_role ON menu_role.role_id = role.id").
					Where("menu_role.menu_id = ?", menu.ID).
					Pluck("role_code", &res.RoleCodes).Error; err != nil {
					return nil, fmt.Errorf("查询菜单角色代码失败: %v", err)
				}
			}
			if menu.Query != nil {
				res.Query = *menu.Query
			}

			if err = model.DB.Model(&model.Button{}).Where("menu_id = ?", menu.ID).Count(&count).Error; err == nil && count > 0 {
				if err = model.DB.Model(&model.Button{}).Where("menu_id = ?", menu.ID).Select("button_code", "button_desc").Find(&buttons).Error; err != nil {
					return nil, fmt.Errorf("查询按钮失败: %v", err)
				}
				for _, button := range buttons {
					res.Buttons = append(res.Buttons, api.MenuButtonRes{
						Code: button.ButtonCode,
						Desc: button.ButtonDesc,
					})
				}
			}
			// result = append(result, res)
			menuMap[menu.ID] = &res
		}
		// 构建父子关系
		if len(menuMap) == 1 {
			for _, menu := range menuMap {
				result = append(result, *menu)
			}
			return &result, err
		}
		for _, menu := range menuMap {
			if menu.ParentId != 0 {
				if parent, ok := menuMap[menu.ParentId]; ok {
					if parent.Children == nil {
						parent.Children = &[]api.MenuRes{}
					}
					*parent.Children = append(*parent.Children, *menu)
				}
			}
		}

		for _, value := range menuMap {
			if value.ParentId == 0 {
				result = append(result, *value)
			}
		}

		return &result, err
	}
	if menu, ok := menuObj.(*model.Menu); ok {
		res := api.MenuRes{
			ID:              menu.ID,
			Status:          menu.Status,
			ParentId:        menu.ParentId,
			MenuType:        menu.MenuType,
			MenuName:        menu.MenuName,
			RouteName:       menu.RouteName,
			RoutePath:       menu.RoutePath,
			Component:       menu.Component,
			Order:           menu.Order,
			I18nKey:         menu.I18nKey,
			IconType:        menu.IconType,
			MultiTab:        menu.MultiTab,
			KeepAlive:       menu.KeepAlive,
			HideInMenu:      menu.HideInMenu,
			IsConstantRoute: menu.IsConstantRoute,
			FixedIndexInTab: menu.FixedIndexInTab,
		}
		if menu.Icon != nil {
			res.Icon = *menu.Icon
		}
		if menu.ActiveMenu != nil {
			res.ActiveMenu = *menu.ActiveMenu
		}

		if menu.Props != nil {
			if *menu.Props == "true" {
				res.Props = true
			} else {
				if err = json.Unmarshal([]byte(*menu.Props), &res.Props); err != nil {
					return nil, fmt.Errorf("解析Props属性失败: %v", err)
				}
			}
		}
		if menu.Href != nil {
			res.Href = *menu.Href
		}
		if menu.ActiveMenu != nil {
			res.ActiveMenu = *menu.ActiveMenu
		}
		if menu.ShowRole {
			if err = model.DB.Model(&model.Role{}).
				Joins("JOIN menu_role ON menu_role.role_id = role.id").
				Where("menu_role.menu_id = ?", menu.ID).
				Pluck("role_code", &res.RoleCodes).Error; err != nil {
				return nil, fmt.Errorf("查询菜单角色代码失败: %v", err)
			}
		}
		if menu.Query != nil {
			res.Query = *menu.Query
		}

		if err = model.DB.Model(&model.Button{}).Where("menu_id = ?", menu.ID).Count(&count).Error; err == nil && count > 0 {
			if err = model.DB.Model(&model.Button{}).Where("menu_id = ?", menu.ID).Select("button_code", "button_desc").Find(&buttons).Error; err != nil {
				return nil, fmt.Errorf("查询按钮失败: %v", err)
			}
			for _, button := range buttons {
				res.Buttons = append(res.Buttons, api.MenuButtonRes{
					Code: button.ButtonCode,
					Desc: button.ButtonDesc,
				})
			}
		}
		result = append(result, res)
		return &result, err
	}
	return &result, errors.New("转换菜单结果失败")
}
