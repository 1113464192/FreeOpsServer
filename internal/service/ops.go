package service

import (
	"FreeOps/internal/consts"
	"FreeOps/internal/model"
	"FreeOps/pkg/api"
	"FreeOps/pkg/util"
	"encoding/json"
	"errors"
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
			return fmt.Errorf("template 已存在, projectId: %d name: %s , 或有错误信息: %v", params.ProjectId, params.Name, err)
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

// 修改/添加 运维操作任务信息
func (s *OpsService) UpdateOpsTask(params api.UpdateOpsTaskReq) (err error) {
	var (
		task  model.OpsTask
		count int64
	)
	// 确保是json构造
	if _, err = json.Marshal(params.TemplateIds); err != nil {
		return fmt.Errorf("templateIds 不是json格式: %v", err)
	}

	if params.ID != 0 {
		if err = model.DB.Model(&model.OpsTask{}).Where("id != ? AND project_id = ? AND name = ?", params.ID, params.ProjectId, params.Name).Count(&count).Error; count > 0 || err != nil {
			return fmt.Errorf("同项目下除自身外仍有相同的名称的 TaskName,id: %d projectId: %d name: %s , 或有错误信息: %v", params.ID, params.ProjectId, params.Name, err)
		}
		if err = model.DB.Model(&model.OpsTask{}).Where("id == ?", params.ID).First(&task).Error; err != nil {
			return fmt.Errorf("查询运维操作任务信息失败: %v", err)
		}
		if params.Auditors != "" {
			if _, err = json.Marshal(params.Auditors); err != nil {
				return fmt.Errorf("auditors 不是json格式: %v", err)
			}
			*task.Auditors = params.Auditors
		}

		task.Name = params.Name
		task.TemplateIds = params.TemplateIds
		task.HostId = params.HostId
		task.IsIntranet = params.IsIntranet
		task.ProjectId = params.ProjectId
		if err = model.DB.Save(&task).Error; err != nil {
			return fmt.Errorf("数据保存失败: %v", err)
		}
		return err
	} else {
		if err = model.DB.Model(&model.OpsTemplate{}).Where("project_id = ? AND name = ?", params.ProjectId, params.Name).Count(&count).Error; count > 0 || err != nil {
			return fmt.Errorf("task 已存在, projectId: %d name: %s , 或有错误信息: %v", params.ProjectId, params.Name, err)
		}

		task = model.OpsTask{
			Name:        params.Name,
			TemplateIds: params.TemplateIds,
			HostId:      params.HostId,
			IsIntranet:  params.IsIntranet,
			ProjectId:   params.ProjectId,
		}
		if params.Auditors != "" {
			if _, err = json.Marshal(params.Auditors); err != nil {
				return fmt.Errorf("auditors 不是json格式: %v", err)
			}
			*task.Auditors = params.Auditors
		}

		if err = model.DB.Create(&task).Error; err != nil {
			return fmt.Errorf("创建运维操作任务信息失败: %v", err)
		}
		return err
	}
}

// 查询运维操作任务信息，不需要content则不传ID
func (s *OpsService) GetOpsTask(params api.GetOpsTaskReq) (*api.GetOpsTasksRes, error) {
	var (
		err   error
		count int64
	)

	if params.ID != 0 {
		var task model.OpsTask
		if err = model.DB.Model(&model.OpsTemplate{}).Where("id = ?", params.ID).First(&task).Error; err != nil {
			return nil, fmt.Errorf("查询运维操作任务信息失败: %v", err)
		}

		res, err := s.getOpsTaskResult(&task)
		if res == nil {
			return nil, errors.New("运维操作任务信息转换结果失败, res为nil")
		}
		if err != nil {
			return nil, fmt.Errorf("运维操作任务信息转换结果失败: err: %v", err)
		}
		result := api.GetOpsTasksRes{
			Records:  *res,
			Page:     1,
			PageSize: 1,
			Total:    1,
		}
		return &result, err
	}

	getDB := model.DB.Model(&model.OpsTask{})
	if params.Name != "" {
		sqlName := "%" + strings.ToUpper(params.Name) + "%"
		getDB = getDB.Where("UPPER(name) LIKE ?", sqlName)
	}

	if params.ProjectId != 0 {
		getDB = getDB.Where("project_id = ?", params.ProjectId)
	}

	if err = getDB.Count(&count).Error; err != nil {
		return nil, fmt.Errorf("查询运维操作任务信息总数失败: %v", err)
	}

	var tasks []model.OpsTask
	if params.Page != 0 && params.PageSize != 0 {
		if err = getDB.Omit("template_ids", "auditors").Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize).Find(&tasks).Error; err != nil {
			return nil, fmt.Errorf("查询运维操作任务信息失败: %v", err)
		}
	} else {
		if err = getDB.Omit("template_ids", "auditors").Find(&tasks).Error; err != nil {
			return nil, fmt.Errorf("查询运维操作任务信息失败: %v", err)
		}
	}
	res, err := s.getOpsTaskResult(&tasks)
	if res == nil {
		return nil, errors.New("运维操作任务信息转换结果失败, res为nil")
	}
	if err != nil {
		return nil, fmt.Errorf("运维操作任务信息转换结果失败: err: %v", err)
	}
	result := api.GetOpsTasksRes{
		Records:  *res,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    count,
	}

	return &result, err
}

func (s *OpsService) getOpsTaskResult(opsObj any) (*[]api.GetOpsTaskRes, error) {
	var result []api.GetOpsTaskRes
	var err error
	// 批量查询不需要temIds和auditors
	if tasks, ok := opsObj.(*[]model.OpsTask); ok {
		for _, task := range *tasks {
			res := api.GetOpsTaskRes{
				ID:        task.ID,
				Name:      task.Name,
				HostId:    task.HostId,
				ProjectId: task.ProjectId,
			}
			result = append(result, res)
		}
		return &result, err
	}
	if task, ok := opsObj.(*model.OpsTask); ok {
		res := api.GetOpsTaskRes{
			ID:        task.ID,
			Name:      task.Name,
			HostId:    task.HostId,
			ProjectId: task.ProjectId,
		}
		if err = json.Unmarshal([]byte(task.TemplateIds), &res.TemplateIds); err != nil {
			return nil, fmt.Errorf("task中的TemplateIds 不符合 []uint 格式: %v", err)
		}
		if task.Auditors != nil {
			if err = json.Unmarshal([]byte(*task.Auditors), &res.Auditors); err != nil {
				return nil, fmt.Errorf("task中的Auditors 不符合 []uint 格式: %v", err)
			}
		}

		result = append(result, res)
	}
	return &result, errors.New("转换运维操作任务信息结果失败")
}
