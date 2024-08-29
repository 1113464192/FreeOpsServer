package service

import (
	"FreeOps/internal/model"
	"FreeOps/pkg/api"
	"fmt"
)

type ApiService struct{}

var insApi ApiService

func ApiServiceApp() *ApiService {
	return &insApi
}

// UpdateApi
func (s *ApiService) UpdateApi(params *api.UpdateApiReq) (err error) {
	var apiModel model.Api
	var total int64
	if params.ID != 0 {
		var count int64
		if err = model.DB.Model(&model.Api{}).Where("id = ?", params.ID).Count(&count).Error; count != 1 || err != nil {
			return fmt.Errorf("api ID不存在: %d, 或有错误信息: %v", params.ID, err)
		}

		err = model.DB.Model(&apiModel).Where("path = ? AND method = ? AND id != ?", params.Path, params.Method, params.ID).Count(&total).Error
		if err != nil {
			return fmt.Errorf("查询api失败: %v", err)
		} else if total > 0 {
			return fmt.Errorf("api路径+方法已被使用: %s+%s", params.Path, params.Method)
		}

		if err := model.DB.Where("id = ?", params.ID).First(&apiModel).Error; err != nil {
			return fmt.Errorf("数据库查询失败: %v", err)
		}
		apiModel.Path = params.Path
		apiModel.Method = params.Method
		apiModel.ApiGroup = params.ApiGroup
		apiModel.Description = params.Description

		if err = CasbinServiceApp().UpdateCasbinApi(apiModel.Path, params.Path, apiModel.Method, params.Method); err != nil {
			return fmt.Errorf("casbin api随动更新失败: %v", err)
		}

		if err = model.DB.Save(&apiModel).Error; err != nil {
			return fmt.Errorf("数据保存失败: %v", err)
		}

		return err
	} else {
		err = model.DB.Model(&apiModel).Where("path = ? AND method = ?", params.Path, params.Method).Count(&total).Error
		// 总数大于0或者有错误就返回
		if err != nil {
			return fmt.Errorf("查询api失败: %v", err)
		} else if total > 0 {
			return fmt.Errorf("api路径+方法已被使用: %s+%s", params.Path, params.Method)
		}

		apiModel = model.Api{
			Path:        params.Path,
			Method:      params.Method,
			ApiGroup:    params.ApiGroup,
			Description: params.Description,
		}

		if err = model.DB.Create(&apiModel).Error; err != nil {
			return fmt.Errorf("创建api失败: %v", err)
		}
		return err
	}
}

func (s *ApiService) GetApis(params *api.GetApiReq) (*api.GetApiRes, error) {
	var apis []model.Api
	var err error
	var count int64
	// 调用Where方法时，它并不会直接修改原始的DB对象，而是返回一个新的*gorm.DB实例，这个新的实例包含了新的查询条件。所以，当你连续调用Where方法时，每次都会返回一个新的*gorm.DB实例，这个新的实例包含了所有之前的查询条件
	getDB := model.DB
	if params.ID != 0 {
		getDB = getDB.Where("id = ?", params.ID)
	}
	if params.ApiGroup != "" {
		getDB = getDB.Where("api_group = ?", params.ApiGroup)
	}

	// 获取符合上面叠加条件的总数
	if err = getDB.Model(&model.Api{}).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("查询角色总数失败: %v", err)

	}
	if err = getDB.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize).Find(&apis).Error; err != nil {
		return nil, fmt.Errorf("查询角色失败: %v", err)
	}
	result := api.GetApiRes{
		Records:  apis,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    count,
	}
	return &result, err
}

func (s *ApiService) GetApiGroup() (*[]string, error) {
	var apiGroups []string
	if err := model.DB.Model(&model.Api{}).Select("DISTINCT api_group").Pluck("api_group", &apiGroups).Error; err != nil {
		return nil, fmt.Errorf("查询api组失败: %v", err)
	}
	return &apiGroups, nil

}

func (s *ApiService) DeleteApi(params api.IdsReq) (err error) {
	var apis []model.Api
	if err = model.DB.Where("id IN (?)", params.Ids).Find(&apis).Error; err != nil {
		return fmt.Errorf("查询api失败: %v", err)
	}
	// 先判断是否有关联角色
	var count int64
	for _, apiModel := range apis {
		count = 0
		if err = model.DB.Model(&model.Api{}).
			Joins("JOIN casbin_rule ON casbin_rule.v1 = api.path").
			Where("casbin_rule.v1 = ? AND casbin_rule.v2 = ?", apiModel.Path, apiModel.Method).
			Count(&count).Error; err != nil {
			return fmt.Errorf("查询关联角色失败: %v", err)
		}
		if count > 0 {
			if err = CasbinServiceApp().DeleteApiPolicy(apiModel.Path, apiModel.Method); err != nil {
				return fmt.Errorf("删除casbin api失败: %v", err)
			}
		}
	}
	if err = model.DB.Where("id IN (?)", params.Ids).Delete(&model.Api{}).Error; err != nil {
		return fmt.Errorf("删除api失败: %v", err)
	}
	return err
}

func (s *ApiService) GetApiTree() (*[]api.GetApiTreeRes, error) {
	var apis []model.Api
	var err error
	if err = model.DB.Find(&apis).Error; err != nil {
		return nil, fmt.Errorf("查询Api失败: %v", err)
	}
	var apiGroup *[]string
	if apiGroup, err = s.GetApiGroup(); err != nil {
		return nil, fmt.Errorf("查询Api组失败: %v", err)
	}

	apiMap := make(map[int]api.GetApiTreeRes)
	for _, value := range apis {
		apiMap[int(value.ID)] = api.GetApiTreeRes{
			Id:    int(value.ID),
			Label: fmt.Sprintf("%s————————%s", value.Path, value.Method),
			Group: value.ApiGroup,
		}
	}
	// 创建父目录
	vid := -1
	for _, value := range *apiGroup {
		apiMap[vid] = api.GetApiTreeRes{
			Id:       vid,
			Label:    value,
			Group:    value,
			Children: &[]api.GetApiTreeRes{},
		}
		// 获取children
		for _, apiValue := range apis {
			if apiValue.ApiGroup == value {
				*apiMap[vid].Children = append(*apiMap[vid].Children, apiMap[int(apiValue.ID)])
			}
		}
		vid--
	}

	var result []api.GetApiTreeRes
	// 从map中提取所有顶级菜单
	for _, value := range apiMap {
		if value.Id < 0 {
			result = append(result, value)
		}
	}
	return &result, nil
}
