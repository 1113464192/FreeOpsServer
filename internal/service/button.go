package service

import (
	"FreeOps/internal/consts"
	"FreeOps/internal/model"
	"FreeOps/pkg/api"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type ButtonService struct{}

var insButton ButtonService

func ButtonServiceApp() *ButtonService {
	return &insButton
}

func (s *ButtonService) UpdateButtons(params *api.UpdateButtonsReq) (err error) {
	var (
		button      model.Button
		buttons     []model.Button
		buttonIds   []uint
		buttonCodes []string
		count       int64
		newParams   api.UpdateButtonsReq
	)
	// 判断是否需要操作
	for _, param := range params.Buttons {
		if err = model.DB.Model(&model.Button{}).Where("button_code = ? AND menu_id = ? AND button_desc = ?", param.ButtonCode, param.MenuId, param.ButtonDesc).Count(&count).Error; err != nil {
			return fmt.Errorf("查询按钮是否存在失败: %v", err)
		}
		if count < 1 {
			// 从params删除已存在的param
			newParams.Buttons = append(newParams.Buttons, param)
		}
	}

	if len(newParams.Buttons) == 0 {
		return nil
	}

	tx := model.DB.Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	for _, param := range newParams.Buttons {
		buttonCodes = append(buttonCodes, param.ButtonCode)
	}
	if err = tx.Model(&model.Button{}).Where("button_code IN (?)", buttonCodes).Pluck("id", &buttonIds).Error; err != nil {
		return fmt.Errorf("查询按钮ID失败: %v", err)
	}

	if err = s.deleteButtons(buttonIds, tx); err != nil {
		return fmt.Errorf("删除按钮失败: %v", err)
	}

	for _, param := range newParams.Buttons {
		if err = tx.Model(&model.Menu{}).Where("id = ?", param.MenuId).Count(&count).Error; count != 1 || err != nil {
			return fmt.Errorf("菜单ID不存在: %d, 或有错误信息: %v", param.MenuId, err)
		}

		button = model.Button{
			ButtonCode: param.ButtonCode,
			ButtonDesc: param.ButtonDesc,
			MenuId:     param.MenuId,
		}

		buttons = append(buttons, button)
	}

	if err = tx.Create(&buttons).Error; err != nil {
		return fmt.Errorf("创建按钮失败: %v", err)
	}
	tx.Commit()
	// 绑定按钮和管理员关系
	var adminRoleId uint
	if err = model.DB.Model(&model.Role{}).Where("role_code = ?", consts.RoleModelAdminCode).Select("id").Scan(&adminRoleId).Error; err != nil {
		return fmt.Errorf("查询管理员角色ID失败: %v", err)
	}
	for _, button = range buttons {
		if err = model.DB.Create(&model.RoleButton{
			ButtonId: button.ID,
			RoleId:   adminRoleId,
		}).Error; err != nil {
			return fmt.Errorf("创建按钮角色关联失败: %v", err)
		}
	}
	return err
}

func (s *ButtonService) GetButtons(params *api.GetButtonsReq) (*api.GetButtonsRes, error) {
	var buttons []model.Button
	var err error
	var count int64

	getDB := model.DB
	if params.ID != 0 {
		getDB = getDB.Where("id = ?", params.ID)
	}

	if params.ButtonCode != "" {
		getDB = getDB.Where("button_code = ?", params.ButtonCode)
	}
	if params.MenuId != 0 {
		getDB = getDB.Where("menu_id = ?", params.MenuId)
	}
	// 获取符合上面叠加条件的总数
	if err = getDB.Model(&model.Button{}).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("查询按钮总数失败: %v", err)

	}
	if params.Page != 0 && params.PageSize != 0 {
		if err = getDB.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize).Find(&buttons).Error; err != nil {
			return nil, fmt.Errorf("查询按钮失败: %v", err)
		}
	} else {
		if err = getDB.Find(&buttons).Error; err != nil {
			return nil, fmt.Errorf("查询按钮失败: %v", err)
		}
	}
	var res *[]api.GetButtonReq
	var result api.GetButtonsRes
	res, err = s.GetResults(&buttons)
	if err != nil {
		return nil, err
	}
	result = api.GetButtonsRes{
		Records:  *res,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    count,
	}
	return &result, err
}

// DeleteButtons
// @param ids: 菜单IDs
func (s *ButtonService) DeleteMenuButtons(ids []uint, cModel *gorm.DB) (err error) {
	var buttonIds []uint
	if err = cModel.Model(model.Button{}).Where("menu_id IN (?)", ids).Pluck("id", &buttonIds).Error; err != nil {
		return fmt.Errorf("查询按钮ID失败 %d: %v", ids, err)
	}
	// 创建按钮时执行，不删除角色按钮关系
	if err = cModel.Where("button_id IN (?)", buttonIds).Delete(&model.RoleButton{}).Error; err != nil {
		return fmt.Errorf("删除角色按钮关系失败 %d: %v", ids, err)
	}
	if err = cModel.Where("id IN (?)", buttonIds).Delete(&model.Button{}).Error; err != nil {
		return fmt.Errorf("删除按钮失败 %d: %v", ids, err)
	}
	return nil
}

// deleteButtons
func (s *ButtonService) deleteButtons(ids []uint, cModel *gorm.DB) (err error) {
	// 创建按钮时执行，不删除角色按钮关系
	if err = cModel.Where("button_id IN (?)", ids).Delete(&model.RoleButton{}).Error; err != nil {
		return fmt.Errorf("删除角色按钮关系失败 %d: %v", ids, err)
	}
	if err = cModel.Where("id IN (?)", ids).Delete(&model.Button{}).Error; err != nil {
		return fmt.Errorf("删除按钮失败 %d: %v", ids, err)
	}
	return nil
}

// 返回按钮结果
func (s *ButtonService) GetResults(buttonObj any) (*[]api.GetButtonReq, error) {
	var result []api.GetButtonReq
	var err error
	if buttons, ok := buttonObj.(*[]model.Button); ok {
		for _, button := range *buttons {
			res := api.GetButtonReq{
				ID:         button.ID,
				ButtonCode: button.ButtonCode,
				ButtonDesc: button.ButtonDesc,
			}
			result = append(result, res)
		}
		return &result, err
	}
	if button, ok := buttonObj.(*model.Button); ok {
		res := api.GetButtonReq{
			ID:         button.ID,
			ButtonCode: button.ButtonCode,
			ButtonDesc: button.ButtonDesc,
		}
		result = append(result, res)
		return &result, err
	}
	return &result, errors.New("转换按钮结果失败")
}
