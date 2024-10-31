package service

import (
	"FreeOps/internal/consts"
	"FreeOps/internal/model"
	"FreeOps/pkg/api"
	"FreeOps/pkg/util"
	"fmt"
	"strings"
)

type OpsService struct{}

var insOps OpsService

func OpsServiceApp() *OpsService {
	return &insOps
}

// 修改/添加运维操作模板
func (s *OpsService) UpdateOpsTemplate(params *api.UpdateOpsTemplateReq) (err error) {
	var (
		template model.OpsTemplate
		count    int64
	)
	if params.ID != 0 {
		if err = model.DB.Model(&model.OpsTemplate{}).Where("id != ? AND project_id = ? AND name = ?", params.ID, params.ProjectId, params.Name).Count(&count).Error; count > 0 || err != nil {
			return fmt.Errorf("同项目下除自身外仍有相同的 template name,id: %d projectId: %d name: %s , 或有错误信息: %v", params.ID, params.ProjectId, params.Name, err)
		}
		if err = model.DB.Model(&model.OpsTemplate{}).Where("id == ?", params.ID).First(&template).Error; err != nil {
			return fmt.Errorf("查询运维操作模板失败: %v", err)
		}
		template.Name = params.Name
		template.Content = params.Content
		template.ProjectId = params.ProjectId
		if err = model.DB.Save(&template).Error; err != nil {
			return fmt.Errorf("数据保存失败: %v", err)
		}
		return err
	} else {
		if err = model.DB.Model(&model.OpsTemplate{}).Where("project_id = ? AND name = ?", params.ProjectId, params.Name).Count(&count).Error; count > 0 || err != nil {
			return fmt.Errorf("template 已存在,id: %d projectId: %d name: %s , 或有错误信息: %v", params.ID, params.ProjectId, params.Name, err)
		}
		template = model.OpsTemplate{
			Name:      params.Name,
			Content:   params.Content,
			ProjectId: params.ProjectId,
		}

		if err = model.DB.Create(&template).Error; err != nil {
			return fmt.Errorf("创建运维操作模板失败: %v", err)
		}
		return err
	}
}

// 查询运维模板，不需要content则不传ID
func (s *OpsService) GetOpsTemplate(params *api.GetOpsTemplatesReq) (*api.GetOpsTemplatesRes, error) {
	var (
		err   error
		count int64
	)

	if params.ID != 0 {
		var template model.OpsTemplate
		if err = model.DB.Model(&model.OpsTemplate{}).Where("id = ?", params.ID).First(&template).Error; err != nil {
			return nil, fmt.Errorf("查询运维操作模板失败: %v", err)
		}
		result := api.GetOpsTemplatesRes{
			Records: []api.GetOpsTemplateRes{
				{
					ID:        template.ID,
					Name:      template.Name,
					Content:   template.Content,
					ProjectId: template.ProjectId,
				},
			},
			Page:     1,
			PageSize: 1,
			Total:    1,
		}
		return &result, err
	}

	getDB := model.DB.Model(&model.OpsTemplate{})
	if params.Name != "" {
		sqlName := "%" + strings.ToUpper(params.Name) + "%"
		getDB = getDB.Where("UPPER(name) LIKE ?", sqlName)
	}

	if params.ProjectId != 0 {
		getDB = getDB.Where("project_id = ?", params.ProjectId)
	}

	if err = getDB.Count(&count).Error; err != nil {
		return nil, fmt.Errorf("查询运维操作模板总数失败: %v", err)
	}

	var templates []model.OpsTemplate
	if params.Page != 0 && params.PageSize != 0 {
		if err = getDB.Omit("content").Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize).Find(&templates).Error; err != nil {
			return nil, fmt.Errorf("查询运维操作模板失败: %v", err)
		}
	} else {
		if err = getDB.Omit("content").Find(&templates).Error; err != nil {
			return nil, fmt.Errorf("查询运维操作模板失败: %v", err)
		}
	}
	var res []api.GetOpsTemplateRes
	var result api.GetOpsTemplatesRes
	for _, tem := range templates {
		v := api.GetOpsTemplateRes{
			ID:        tem.ID,
			UpdatedAt: tem.UpdatedAt.Format("2006-01-02 15:04:05"),
			Name:      tem.Name,
			ProjectId: tem.ProjectId,
		}
		res = append(res, v)
	}

	result = api.GetOpsTemplatesRes{
		Records:  res,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    count,
	}

	return &result, err
}

// 修改/添加运维操作的参数模板
func (s *OpsService) UpdateOpsParamsTemplate(params model.OpsParam) (err error) {
	var (
		template model.OpsParam
		count    int64
	)
	if params.ID != 0 {
		if err = model.DB.Model(&model.OpsParam{}).Where("id != ? AND keyword = ? AND variable = ?", params.ID, params.Keyword, params.Variable).Count(&count).Error; count > 0 || err != nil {
			return fmt.Errorf("有相同的 TemplateParam,id: %d Param: %s variable: %s , 或有错误信息: %v", params.ID, params.Keyword, params.Variable, err)
		}
		if err = model.DB.Model(&model.OpsParam{}).Where("id == ?", params.ID).First(&template).Error; err != nil {
			return fmt.Errorf("查询运维操作的参数模板失败: %v", err)
		}
		template.Keyword = params.Keyword
		template.Variable = params.Variable
		if err = model.DB.Save(&template).Error; err != nil {
			return fmt.Errorf("数据保存失败: %v", err)
		}
		return err
	} else {
		if err = model.DB.Model(&model.OpsParam{}).Where("keyword = ? AND variable = ?", params.Keyword, params.Variable).Count(&count).Error; count > 0 || err != nil {
			return fmt.Errorf("TemplateParam 已存在,id: %d Param: %s variable: %s , 或有错误信息: %v", params.ID, params.Keyword, params.Variable, err)
		}
		template = model.OpsParam{
			Keyword:  params.Keyword,
			Variable: params.Variable,
		}

		if err = model.DB.Create(&template).Error; err != nil {
			return fmt.Errorf("创建运维操作的参数模板失败: %v", err)
		}
		return err
	}
}

// 查询运维操作的参数模板
func (s *OpsService) GetOpsParamsTemplate(params api.GetOpsParamsTemplatesReq) (*api.GetOpsParamsTemplatesRes, error) {
	var (
		err   error
		count int64
	)

	getDB := model.DB.Model(&model.OpsParam{})
	if params.ID != 0 {
		getDB = getDB.Where("id = ?", params.ID)
	}
	if params.Keyword != "" {
		sqlName := "%" + strings.ToUpper(params.Keyword) + "%"
		getDB = getDB.Where("UPPER(keyword) LIKE ?", sqlName)
	}

	if err = getDB.Count(&count).Error; err != nil {
		return nil, fmt.Errorf("查询运维操作的参数模板总数失败: %v", err)
	}

	var templates []model.OpsParam
	if params.Page != 0 && params.PageSize != 0 {
		if err = getDB.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize).Find(&templates).Error; err != nil {
			return nil, fmt.Errorf("查询运维操作的参数模板失败: %v", err)
		}
	} else {
		if err = getDB.Find(&templates).Error; err != nil {
			return nil, fmt.Errorf("查询运维操作的参数模板失败: %v", err)
		}
	}

	result := api.GetOpsParamsTemplatesRes{
		Records:  templates,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    count,
	}

	return &result, err
}

func (s *OpsService) BindTemplateParams(TemplateID uint, ParamIDs []uint) (err error) {
	// 先传的id是否都存在
	var count int64
	if err = model.DB.Model(&model.OpsTemplate{}).Where("id = ?", TemplateID).Count(&count).Error; count != 1 || err != nil {
		return fmt.Errorf("template 不存在ID: %d, 如果查询template失败: %v", TemplateID, err)
	}

	if err = model.DB.Model(&model.OpsParam{}).Where("id IN (?)", ParamIDs).Count(&count).Error; count != int64(len(ParamIDs)) || err != nil {
		notExistIds, err2 := util.FindNotExistIDs(consts.MysqlTableNameOpsParam, ParamIDs)
		if err2 != nil {
			return fmt.Errorf("查询OpsParam失败: %v", err2)
		}
		return fmt.Errorf("opsParam 不存在ID: %d, 如果查询OpsParam失败: %v", notExistIds, err)
	}

	tx := model.DB.Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	// 要先清空用户的当前关联
	if err = tx.Where("template_id = ?", TemplateID).Delete(&model.OpsTemplateParam{}).Error; err != nil {
		return fmt.Errorf("清空运维操作模板与运维操作的参数模板关联失败: %v", err)
	}

	var templateParams []model.OpsTemplateParam
	for _, id := range ParamIDs {
		templateParam := model.OpsTemplateParam{
			TemplateId: TemplateID,
			ParamId:    id,
		}
		templateParams = append(templateParams, templateParam)
	}
	if err = tx.Create(&templateParams).Error; err != nil {
		return fmt.Errorf("绑定运维操作模板与运维操作的参数模板失败: %v", err)
	}

	tx.Commit()
	return nil
}

// 查询运维操作模板对应的参数模板
func (s *OpsService) GetTemplateParams(TemplateID uint) (res []model.OpsParam, err error) {
	if err = model.DB.Joins("JOIN ops_template_param ON ops_template_param.param_id = ops_param.id").
		Where("ops_template_param.template_id = ?", TemplateID).
		Find(&res).Error; err != nil {
		return nil, fmt.Errorf("查询运维操作模板对应的参数模板失败: %v", err)
	}
	return res, nil
}
