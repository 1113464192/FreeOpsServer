package service

import (
	"FreeOps/internal/consts"
	"FreeOps/internal/model"
	"FreeOps/pkg/api"
	"FreeOps/pkg/util"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

type RoleService struct{}

var insRole RoleService

func RoleServiceApp() *RoleService {
	return &insRole
}

func (s *RoleService) UpdateRole(params *api.UpdateRoleReq) (err error) {
	var (
		role  model.Role
		count int64
	)
	if params.ID != 0 {
		// 修改
		tx := model.DB.Begin()
		defer func() {
			if r := recover(); r != nil || err != nil {
				tx.Rollback()
			}
		}()
		// 更改当前角色对应的所有用户的updateAt，先获取根据第三方表获取对应的所有用户，再更新用户的updateAt
		mBool, err := s.HasBoundUsers(params.ID)
		if err != nil {
			return err
		}
		if mBool {
			var userIds []uint
			if err = tx.Model(&model.User{}).
				Joins("JOIN user_role ON user_role.user_id = user.id").
				Where("user_role.role_id = ?", params.ID).
				Pluck("id", &userIds).Error; err != nil {
				return fmt.Errorf("查询用户失败: %v", err)
			}
			if err = tx.Model(&model.User{}).Where("id IN (?)", userIds).Update("updated_at", time.Now()).Error; err != nil {
				return fmt.Errorf("修改用户更新时间失败: %v", err)
			}
		}

		if err = model.DB.Model(&model.Role{}).Where("id = ?", params.ID).Count(&count).Error; count != 1 || err != nil {
			return fmt.Errorf("role ID不存在: %d, 或有错误信息: %v", params.ID, err)
		}

		// 判断role_name是否和现有角色重复
		err = tx.Model(&model.Role{}).Where("role_name = ? AND id != ? OR role_code = ? AND id != ?", params.RoleName, params.ID, params.RoleCode, params.ID).Count(&count).Error
		if err != nil {
			return fmt.Errorf("查询角色失败: %v", err)
		} else if count > 0 {
			return fmt.Errorf("角色代码(%s)或角色名(%s)已被使用", params.RoleCode, params.RoleName)
		}

		if err := tx.Where("id = ?", params.ID).First(&role).Error; err != nil {
			return fmt.Errorf("数据库查询失败: %v", err)
		}
		role.RoleName = params.RoleName
		role.RoleCode = params.RoleCode
		role.RoleDesc = params.RoleDesc

		if err = tx.Save(&role).Error; err != nil {
			return fmt.Errorf("数据保存失败: %v", err)
		}
		tx.Commit()
		return err
	} else {
		err = model.DB.Model(&model.Role{}).Where("role_name = ? OR role_code = ?", params.RoleName, params.RoleCode).Count(&count).Error
		// 总数大于0或者有错误就返回
		if err != nil {
			return fmt.Errorf("查询角色失败: %v", err)
		} else if count > 0 {
			return fmt.Errorf("角色代码(%s)或角色名(%s)已存在", params.RoleCode, params.RoleName)
		}

		role = model.Role{
			RoleName: params.RoleName,
			RoleCode: params.RoleCode,
			RoleDesc: params.RoleDesc,
		}

		if err = model.DB.Create(&role).Error; err != nil {
			return fmt.Errorf("创建角色失败: %v", err)
		}
		return err
	}
}

func (s *RoleService) GetRoles(params *api.GetRolesReq) (*api.GetRolesRes, error) {
	var roles []model.Role
	var err error
	var count int64
	// 调用Where方法时，它并不会直接修改原始的DB对象，而是返回一个新的*gorm.DB实例，这个新的实例包含了新的查询条件。所以，当你连续调用Where方法时，每次都会返回一个新的*gorm.DB实例，这个新的实例包含了所有之前的查询条件
	getDB := model.DB
	if params.ID != 0 {
		getDB = getDB.Where("id = ?", params.ID)
	}
	if params.RoleName != "" {
		// rolename不为空则模糊查询
		sqlRoleName := "%" + strings.ToUpper(params.RoleName) + "%"
		getDB = getDB.Where("UPPER(role_name) LIKE ?", sqlRoleName)
	}
	if params.RoleCode != "" {
		getDB = getDB.Where("role_code = ?", params.RoleCode)
	}
	// 获取符合上面叠加条件的总数
	if err = getDB.Model(&model.Role{}).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("查询角色总数失败: %v", err)

	}
	if err = getDB.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize).Find(&roles).Error; err != nil {
		return nil, fmt.Errorf("查询角色失败: %v", err)
	}
	var res *[]api.UpdateRoleReq
	var result api.GetRolesRes
	res, err = s.GetResults(&roles)
	if err != nil {
		return nil, err
	}
	result = api.GetRolesRes{
		Records:  *res,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    count,
	}
	return &result, err
}

func (s *RoleService) GetAllRolesSummary() (*[]api.GetAllRolesSummaryRes, error) {
	var roles []model.Role
	var err error
	if err = model.DB.Find(&roles).Error; err != nil {
		return nil, fmt.Errorf("查询角色失败: %v", err)
	}
	var res api.GetAllRolesSummaryRes
	var result []api.GetAllRolesSummaryRes
	for _, role := range roles {
		res = api.GetAllRolesSummaryRes{
			ID:       role.ID,
			RoleName: role.RoleName,
			RoleCode: role.RoleCode,
		}
		result = append(result, res)
	}
	return &result, err
}

func (s *RoleService) bindRoleMenus(roleId uint, menuIds []uint) (err error) {
	var count int64
	if err = model.DB.Model(&model.Menu{}).Where("id IN (?)", menuIds).Count(&count).Error; count != int64(len(menuIds)) || err != nil {
		notExistIds, err2 := util.FindNotExistIDs(consts.MysqlTableNameMenu, menuIds)
		if err2 != nil {
			return fmt.Errorf("查询菜单失败: %v", err2)
		}
		return fmt.Errorf("menu 不存在ID: %d, 如果查询菜单失败: %v", notExistIds, err)
	}

	tx := model.DB.Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()
	// 要先清空角色的当前关联
	if err = tx.Where("role_id = ?", roleId).Delete(&model.MenuRole{}).Error; err != nil {
		return fmt.Errorf("清空角色菜单关联失败: %v", err)
	}

	var menus []model.Menu
	if err = model.DB.Where("id in (?)", menuIds).Select("id").Find(&menus).Error; err != nil {
		return fmt.Errorf("查询菜单失败: %v", err)
	}
	var roleMenus []model.MenuRole
	for _, menu := range menus {
		roleMenu := model.MenuRole{
			MenuId: menu.ID,
			RoleId: roleId,
		}
		roleMenus = append(roleMenus, roleMenu)
	}

	if err = tx.Create(&roleMenus).Error; err != nil {
		return fmt.Errorf("绑定角色菜单失败: %v", err)
	}

	tx.Commit()
	return nil
}

func (s *RoleService) bindRoleButtons(roleId uint, buttonIds []uint) (err error) {
	var count int64
	if err = model.DB.Model(&model.Button{}).Where("id IN (?)", buttonIds).Count(&count).Error; count != int64(len(buttonIds)) || err != nil {
		notExistIds, err2 := util.FindNotExistIDs(consts.MysqlTableNameButton, buttonIds)
		if err2 != nil {
			return fmt.Errorf("查询按钮失败: %v", err2)
		}
		return fmt.Errorf("按钮 不存在ID: %d, 如果查询按钮失败: %v", notExistIds, err)
	}

	tx := model.DB.Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()
	// 要先清空角色的当前关联
	if err = tx.Where("role_id = ?", roleId).Delete(&model.RoleButton{}).Error; err != nil {
		return fmt.Errorf("清空角色按钮关联失败: %v", err)
	}

	var buttons []model.Button
	if err = model.DB.Where("id in (?)", buttonIds).Select("id").Find(&buttons).Error; err != nil {
		return fmt.Errorf("查询按钮失败: %v", err)
	}
	var roleButtons []model.RoleButton
	for _, button := range buttons {
		roleButton := model.RoleButton{
			ButtonId: button.ID,
			RoleId:   roleId,
		}
		roleButtons = append(roleButtons, roleButton)
	}

	if err = tx.Create(&roleButtons).Error; err != nil {
		return fmt.Errorf("绑定角色按钮失败: %v", err)
	}

	tx.Commit()
	return nil
}

func (s *RoleService) bindRoleProjects(roleId uint, projectIds []uint) (err error) {
	var count int64
	if err = model.DB.Model(&model.Project{}).Where("id IN (?)", projectIds).Count(&count).Error; count != int64(len(projectIds)) || err != nil {
		notExistIds, err2 := util.FindNotExistIDs(consts.MysqlTableNameProject, projectIds)
		if err2 != nil {
			return fmt.Errorf("查询项目失败: %v", err2)
		}
		return fmt.Errorf("项目 不存在ID: %d, 如果查询项目失败: %v", notExistIds, err)
	}

	tx := model.DB.Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()
	// 要先清空角色的当前关联
	if err = tx.Where("role_id = ?", roleId).Delete(&model.RoleProject{}).Error; err != nil {
		return fmt.Errorf("清空角色项目关联失败: %v", err)
	}

	var projects []model.Project
	if err = model.DB.Where("id in (?)", projectIds).Select("id").Find(&projects).Error; err != nil {
		return fmt.Errorf("查询项目失败: %v", err)
	}
	var roleProjects []model.RoleProject
	for _, project := range projects {
		roleProject := model.RoleProject{
			ProjectId: project.ID,
			RoleId:    roleId,
		}
		roleProjects = append(roleProjects, roleProject)
	}

	if err = tx.Create(&roleProjects).Error; err != nil {
		return fmt.Errorf("绑定角色项目失败: %v", err)
	}

	tx.Commit()
	return nil
}

func (s *RoleService) BindRoleRelation(param api.BindRoleRelationReq) (err error) {
	var count int64
	if err = model.DB.Model(&model.Role{}).Where("id = ?", param.RoleId).Count(&count).Error; count < 1 || err != nil {
		return fmt.Errorf("角色ID不存在, 如果查询角色失败: %v", err)
	}

	switch param.AssociationType {
	case consts.RoleAssociationTypeApi:
		roleIdStr := strconv.Itoa(int(param.RoleId))
		if err = CasbinServiceApp().UpdateCasbin(roleIdStr, param.ObjectIds); err != nil {
			return err
		}
	case consts.RoleAssociationTypeMenu:
		if err = s.bindRoleMenus(param.RoleId, param.ObjectIds); err != nil {
			return err
		}
	case consts.RoleAssociationTypeButton:
		if err = s.bindRoleButtons(param.RoleId, param.ObjectIds); err != nil {
			return err
		}
	case consts.RoleAssociationTypeProject:
		if err = s.bindRoleProjects(param.RoleId, param.ObjectIds); err != nil {
			return err
		}
	default:
		return fmt.Errorf("关联类型错误: %d", param.AssociationType)
	}
	return err
}

func (s *RoleService) GetRoleUsers(params api.IdPageReq) ([]uint, error) {
	var users []model.User
	var err error

	if err = model.DB.Model(&model.User{}).
		Joins("JOIN user_role ON user_role.user_id = user.id").
		Where("user_role.role_id = ?", params.Id).
		Select("DISTINCT id").
		Find(&users).Error; err != nil {
		return nil, fmt.Errorf("查询角色用户失败: %v", err)
	}

	var res []uint
	for _, user := range users {
		res = append(res, user.ID)
	}
	return res, err
}

func (s *RoleService) GetSelfRoleIDs(c *gin.Context) (roleIds []uint, err error) {
	var roles *[]model.Role
	// 获取角色对应的项目ID
	if roles, err = util.GetClaimsRole(c); err != nil {
		return nil, err
	}
	// 取出所有roles的id
	for _, role := range *roles {
		roleIds = append(roleIds, role.ID)
	}
	return roleIds, err
}

func (s *RoleService) GetRoleProjects(ids []uint) ([]uint, error) {
	var projects []model.Project
	var err error

	if err = model.DB.Model(&model.Project{}).
		Joins("JOIN role_project ON role_project.project_id = project.id").
		Where("role_project.role_id IN ?", ids).
		Select("DISTINCT id").
		Find(&projects).Error; err != nil {
		return nil, fmt.Errorf("查询角色项目失败: %v", err)
	}

	var res []uint
	for _, project := range projects {
		res = append(res, project.ID)
	}
	return res, err
}

func (s *RoleService) GetRoleMenus(params api.IdsReq) ([]uint, error) {
	var menus []model.Menu
	var err error

	if err = model.DB.Model(&model.Menu{}).
		Joins("JOIN menu_role ON menu_role.menu_id = menu.id").
		Where("menu_role.role_id IN (?)", params.Ids).
		Select("DISTINCT id").
		Find(&menus).Error; err != nil {
		return nil, fmt.Errorf("查询菜单角色失败: %v", err)
	}

	var res []uint
	for _, menu := range menus {
		res = append(res, menu.ID)
	}
	return res, err
}

func (s *RoleService) GetRoleButtons(params api.IdsReq) ([]uint, error) {
	var buttons []model.Button
	var err error

	if err = model.DB.Model(&model.Button{}).
		Joins("JOIN role_button ON role_button.button_id = button.id").
		Where("role_button.role_id IN (?)", params.Ids).
		Select("DISTINCT id").
		Find(&buttons).Error; err != nil {
		return nil, fmt.Errorf("查询角色按钮失败: %v", err)
	}

	var res []uint
	for _, button := range buttons {
		res = append(res, button.ID)
	}
	return res, err
}

// 删除角色ID
func (s *RoleService) DeleteRoles(ids []uint) (err error) {
	tx := model.DB.Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()
	// 删除用户角色关联
	if err = tx.Where("role_id IN (?)", ids).Delete(&model.UserRole{}).Error; err != nil {
		return fmt.Errorf("删除用户角色关联失败 %d: %v", ids, err)
	}

	// 清空角色按钮关联
	if err = tx.Where("role_id IN (?)", ids).Delete(&model.RoleButton{}).Error; err != nil {
		return fmt.Errorf("删除角色按钮关系失败 %d: %v", ids, err)
	}

	// 删除角色菜单关联
	if err = tx.Where("role_id IN (?)", ids).Delete(&model.MenuRole{}).Error; err != nil {
		return fmt.Errorf("删除角色菜单关联失败: %v", err)
	}

	if err = tx.Where("id IN (?)", ids).Delete(&model.Role{}).Error; err != nil {
		return fmt.Errorf("删除角色失败 %d: %v", ids, err)
	}

	for _, id := range ids {
		idStr := strconv.Itoa(int(id))
		if !CasbinServiceApp().ClearCasbin(0, idStr) {
			return fmt.Errorf("role %d: 删除Casbin失败", id)
		}
	}
	tx.Commit()
	return nil
}

// 角色是否有绑定用户
func (s *RoleService) HasBoundUsers(roleId uint) (bool, error) {
	var count int64
	err := model.DB.Model(&model.User{}).
		Joins("JOIN user_role ON user_role.user_id = user.id").
		Where("user_role.role_id = ?", roleId).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// 返回角色结果
func (s *RoleService) GetResults(roleObj any) (*[]api.UpdateRoleReq, error) {
	var result []api.UpdateRoleReq
	var err error
	if roles, ok := roleObj.(*[]model.Role); ok {
		for _, role := range *roles {
			res := api.UpdateRoleReq{
				ID:       role.ID,
				RoleName: role.RoleName,
				RoleCode: role.RoleCode,
				RoleDesc: role.RoleDesc,
			}
			result = append(result, res)
		}
		return &result, err
	}
	if role, ok := roleObj.(*model.Role); ok {
		res := api.UpdateRoleReq{
			ID:       role.ID,
			RoleName: role.RoleName,
			RoleCode: role.RoleCode,
			RoleDesc: role.RoleDesc,
		}
		result = append(result, res)
		return &result, err
	}
	return &result, errors.New("转换角色结果失败")
}
