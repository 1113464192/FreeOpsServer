package service

import (
	"FreeOps/global"
	"FreeOps/internal/consts"
	"FreeOps/internal/model"
	"FreeOps/pkg/api"
	"FreeOps/pkg/logger"
	"FreeOps/pkg/util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"strings"
	"sync"
	"time"
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
		if err = model.DB.Model(&model.OpsTemplate{}).Where("id = ?", params.ID).First(&template).Error; err != nil {
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

// 通用函数，用于查询项目名称
func getProjectNames[T any](items *[]T, getProjectId func(T) uint) (map[uint]string, error) {
	var (
		projectIds []uint
		err        error
	)
	for _, item := range *items {
		projectIds = append(projectIds, getProjectId(item))
	}

	var projects []model.Project
	if err = model.DB.Model(&model.Project{}).Select("id", "name").Where("id IN (?)", projectIds).Find(&projects).Error; err != nil {
		return nil, fmt.Errorf("查询项目名称失败: %v", err)
	}

	projectNameMap := make(map[uint]string)
	for _, project := range projects {
		projectNameMap[project.ID] = project.Name
	}
	return projectNameMap, nil
}

// 获取OpsTask的项目名称
func (s *OpsService) getOpsTasksProjectName(tasks *[]model.OpsTask) (map[uint]string, error) {
	return getProjectNames(tasks, func(task model.OpsTask) uint {
		return task.ProjectId
	})
}

// 获取OpsTemplate的项目名称
func (s *OpsService) getOpsTemplatesProjectName(templates *[]model.OpsTemplate) (map[uint]string, error) {
	return getProjectNames(templates, func(template model.OpsTemplate) uint {
		return template.ProjectId
	})
}

// 获取OpsTaskLog的项目名称
func (s *OpsService) getOpsTaskLogsProjectName(logs *[]model.OpsTaskLog) (map[uint]string, error) {
	return getProjectNames(logs, func(log model.OpsTaskLog) uint {
		return log.ProjectId
	})
}

// 查询运维模板，不需要content则不传ID
func (s *OpsService) GetOpsTemplate(params *api.GetOpsTemplatesReq, bindProjectIds []uint) (*api.GetOpsTemplatesRes, error) {
	var (
		err   error
		count int64
	)

	if params.ID != 0 {
		var (
			template    model.OpsTemplate
			projectName string
		)
		if err = model.DB.Model(&model.OpsTemplate{}).Where("id = ?", params.ID).First(&template).Error; err != nil {
			return nil, fmt.Errorf("查询运维操作模板失败: %v", err)
		}
		if err = model.DB.Model(&model.Project{}).Where("id = ?", template.ProjectId).Pluck("name", &projectName).Error; err != nil {
			return nil, fmt.Errorf("查询项目名称失败: %v", err)
		}
		result := api.GetOpsTemplatesRes{
			Records: []api.GetOpsTemplateRes{
				{
					ID:          template.ID,
					Name:        template.Name,
					Content:     template.Content,
					ProjectName: projectName,
					ProjectId:   template.ProjectId,
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
		if !util.IsUintSliceContain(bindProjectIds, params.ProjectId) {
			return nil, errors.New("用户无权限查看该项目的运维操作模板")
		}
		getDB = getDB.Where("project_id = ?", params.ProjectId)
	} else {
		getDB = getDB.Where("project_id IN (?)", bindProjectIds)
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
	var (
		res            []api.GetOpsTemplateRes
		result         api.GetOpsTemplatesRes
		projectNameMap map[uint]string
	)
	if projectNameMap, err = s.getOpsTemplatesProjectName(&templates); err != nil {
		return nil, err
	}
	for _, tem := range templates {
		v := api.GetOpsTemplateRes{
			ID:          tem.ID,
			UpdatedAt:   tem.UpdatedAt.Format("2006-01-02 15:04:05"),
			Name:        tem.Name,
			ProjectName: projectNameMap[tem.ProjectId],
			ProjectId:   tem.ProjectId,
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

// 删除运维操作模板
func (s *OpsService) DeleteOpsTemplate(ids []uint) (err error) {
	tx := model.DB.Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()
	if err = tx.Where("template_id IN (?)", ids).Delete(&model.OpsTemplateParam{}).Error; err != nil {
		return fmt.Errorf("删除操作模板与参数模板的关联失败 %d: %v", ids, err)
	}
	if err = tx.Where("id IN (?)", ids).Delete(&model.OpsTemplate{}).Error; err != nil {
		return fmt.Errorf("删除模板失败 %d: %v", ids, err)
	}
	tx.Commit()
	return nil
}

// 修改/添加运维操作的参数模板
func (s *OpsService) UpdateOpsParamsTemplate(params model.OpsParam) (err error) {
	var (
		template model.OpsParam
		count    int64
	)
	if params.ID != 0 {
		if err = model.DB.Model(&model.OpsParam{}).Where("id != ? AND keyword = ?", params.ID, params.Keyword).Count(&count).Error; count > 0 || err != nil {
			return fmt.Errorf("有相同的 TemplateParam,id: %d Param: %s, 或有错误信息: %v", params.ID, params.Keyword, err)
		}
		if err = model.DB.Model(&model.OpsParam{}).Where("id = ?", params.ID).First(&template).Error; err != nil {
			return fmt.Errorf("查询运维操作的参数模板失败: %v", err)
		}
		template.Keyword = params.Keyword
		template.Variable = params.Variable
		if err = model.DB.Save(&template).Error; err != nil {
			return fmt.Errorf("数据保存失败: %v", err)
		}
		return err
	} else {
		if err = model.DB.Model(&model.OpsParam{}).Where("keyword = ?", params.Keyword).Count(&count).Error; count > 0 || err != nil {
			return fmt.Errorf("TemplateParam 已存在,id: %d Param: %s, 或有错误信息: %v", params.ID, params.Keyword, err)
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

// 删除运维操作的参数模板
func (s *OpsService) DeleteOpsParamsTemplate(ids []uint) (err error) {
	var count int64
	if err = model.DB.Model(&model.OpsTemplateParam{}).Where("param_id IN (?)", ids).Count(&count).Error; err != nil {
		return fmt.Errorf("查询 运维操作的参数模板 与 运维操作模板 的关联关系失败: %v", err)
	}
	if count > 0 {
		return errors.New("要删除的param有关联的 运维操作模板 关系存在")
	}
	if err = model.DB.Where("id IN (?)", ids).Delete(&model.OpsParam{}).Error; err != nil {
		return fmt.Errorf("删除 运维操作的参数模板 失败 %d: %v", ids, err)
	}
	return nil
}

func (s *OpsService) BindTemplateParams(TemplateID uint, ParamIDs []uint) (err error) {
	// 先传的id是否都存在
	var count int64
	if err = model.DB.Model(&model.OpsTemplate{}).Where("id = ?", TemplateID).Count(&count).Error; count != 1 || err != nil {
		return fmt.Errorf("template 不存在ID: %d, 如果查询template失败: %v", TemplateID, err)
	}

	if len(ParamIDs) > 0 {
		if err = model.DB.Model(&model.OpsParam{}).Where("id IN (?)", ParamIDs).Count(&count).Error; count != int64(len(ParamIDs)) || err != nil {
			notExistIds, err2 := util.FindNotExistIDs(consts.MysqlTableNameOpsParam, ParamIDs)
			if err2 != nil {
				return fmt.Errorf("查询OpsParam失败: %v", err2)
			}
			return fmt.Errorf("opsParam 不存在ID: %d, 如果查询OpsParam失败: %v", notExistIds, err)
		}
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

	if len(ParamIDs) > 0 {
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

// 判断创建/修改 task是否符合规范
func (s *OpsService) validateTask(params api.UpdateOpsTaskReq) (auditors string, templateIds string, err error) {
	var count int64
	if len(params.Auditors) > 0 {
		var auditorBytes []byte
		if auditorBytes, err = json.Marshal(params.Auditors); err != nil {
			return "", "", fmt.Errorf("auditors 转换为 string 失败: %v", err)
		}
		auditors = string(auditorBytes)
	} else {
		auditors = ""
	}
	if err = model.DB.Model(&model.OpsTemplate{}).Where("id IN (?)", params.TemplateIds).Count(&count).Error; count != int64(len(params.TemplateIds)) || err != nil {
		notExistIds, err2 := util.FindNotExistIDs(consts.MysqlTableNameOpsTemplate, params.TemplateIds)
		if err2 != nil {
			return "", "", fmt.Errorf("查询OpsTemplate失败: %v", err2)
		}
		return "", "", fmt.Errorf("OpsTemplate 不存在ID: %d, 如果查询OpsTemplate失败: %v", notExistIds, err)
	}
	var templateIdsBytes []byte
	if templateIdsBytes, err = json.Marshal(params.TemplateIds); err != nil {
		return "", "", fmt.Errorf("TemplateIds 转换为 string 失败: %v", err)
	}
	templateIds = string(templateIdsBytes)
	return auditors, templateIds, err
}

// 修改/添加 运维操作任务信息
func (s *OpsService) UpdateOpsTask(params api.UpdateOpsTaskReq) (err error) {
	var (
		task  model.OpsTask
		host  model.Host
		count int64
	)
	if err = model.DB.Where("id = ?", params.HostId).Select("vip").First(&host).Error; err != nil {
		return fmt.Errorf("查询运维管理机信息失败: %v", err)
	}
	if params.IsIntranet && host.Vip == "" {
		return errors.New("运维管理机没有内网IP")
	}
	if params.ID != 0 {
		if err = model.DB.Model(&model.OpsTask{}).Where("id != ? AND project_id = ? AND name = ?", params.ID, params.ProjectId, params.Name).Count(&count).Error; count > 0 || err != nil {
			return fmt.Errorf("同项目下除自身外仍有相同的名称的 TaskName,id: %d projectId: %d name: %s , 或有错误信息: %v", params.ID, params.ProjectId, params.Name, err)
		}
		if err = model.DB.Model(&model.OpsTask{}).Where("id = ?", params.ID).First(&task).Error; err != nil {
			return fmt.Errorf("查询运维操作任务信息失败: %v", err)
		}
		if task.Auditors == nil {
			task.Auditors = new(string)
		}
		if *task.Auditors, task.TemplateIds, err = s.validateTask(params); err != nil {
			return err
		}
		if *task.Auditors == "" {
			task.Auditors = nil
		}
		task.Name = params.Name
		task.CheckTemplateId = params.CheckTemplateId
		task.HostId = params.HostId
		task.IsIntranet = params.IsIntranet
		task.IsConcurrent = params.IsConcurrent
		task.ProjectId = params.ProjectId
		if err = model.DB.Save(&task).Error; err != nil {
			return fmt.Errorf("数据保存失败: %v", err)
		}
		return err
	} else {
		if err = model.DB.Model(&model.OpsTask{}).Where("project_id = ? AND name = ?", params.ProjectId, params.Name).Count(&count).Error; count > 0 || err != nil {
			return fmt.Errorf("task 已存在, projectId: %d name: %s , 或有错误信息: %v", params.ProjectId, params.Name, err)
		}
		task = model.OpsTask{
			Name:            params.Name,
			CheckTemplateId: params.CheckTemplateId,
			HostId:          params.HostId,
			IsIntranet:      params.IsIntranet,
			IsConcurrent:    params.IsConcurrent,
			ProjectId:       params.ProjectId,
		}
		if task.Auditors == nil {
			task.Auditors = new(string)
		}
		if *task.Auditors, task.TemplateIds, err = s.validateTask(params); err != nil {
			return err
		}
		if *task.Auditors == "" {
			task.Auditors = nil
		}
		if err = model.DB.Create(&task).Error; err != nil {
			return fmt.Errorf("创建运维操作任务信息失败: %v", err)
		}
		return err
	}
}

// 查询运维操作任务信息，不需要content则不传ID
func (s *OpsService) GetOpsTask(params api.GetOpsTaskReq, bindProjectIds []uint) (*api.GetOpsTasksRes, error) {
	var (
		err   error
		count int64
	)

	if params.ID != 0 {
		var task model.OpsTask
		if err = model.DB.Model(&model.OpsTask{}).Where("id = ?", params.ID).First(&task).Error; err != nil {
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
		if !util.IsUintSliceContain(bindProjectIds, params.ProjectId) {
			return nil, errors.New("用户无权限查看该项目的运维操作任务")
		}
		getDB = getDB.Where("project_id = ?", params.ProjectId)
	} else {
		getDB = getDB.Where("project_id IN (?)", bindProjectIds)
	}

	if err = getDB.Count(&count).Error; err != nil {
		return nil, fmt.Errorf("查询运维操作任务信息总数失败: %v", err)
	}

	var tasks []model.OpsTask
	if params.Page != 0 && params.PageSize != 0 {
		if err = getDB.Omit("template_ids", "auditors", "check_template_id", "commands").Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize).Find(&tasks).Error; err != nil {
			return nil, fmt.Errorf("查询运维操作任务信息失败: %v", err)
		}
	} else {
		if err = getDB.Omit("template_ids", "auditors", "check_template_id", "commands").Find(&tasks).Error; err != nil {
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

// 通用函数，用于查询项目名称
func getHostNames[T any](items *[]T, getHostId func(T) uint) (map[uint]string, error) {
	var (
		hostIds []uint
		err     error
	)
	for _, item := range *items {
		hostIds = append(hostIds, getHostId(item))
	}

	var hosts []model.Host
	if err = model.DB.Model(&model.Host{}).Select("id", "name").Where("id IN (?)", hostIds).Find(&hosts).Error; err != nil {
		return nil, fmt.Errorf("查询服务器名称失败: %v", err)
	}

	hostNameMap := make(map[uint]string)
	for _, host := range hosts {
		hostNameMap[host.ID] = host.Name
	}
	return hostNameMap, nil
}

// 获取OpsTask的项目名称
func (s *OpsService) getOpsTaskHostNames(tasks *[]model.OpsTask) (map[uint]string, error) {
	return getHostNames(tasks, func(task model.OpsTask) uint {
		return task.HostId
	})
}

func (s *OpsService) getOpsTaskResult(opsObj any) (*[]api.GetOpsTaskRes, error) {
	var (
		result []api.GetOpsTaskRes
		err    error
	)
	// 批量查询不需要temIds和auditors
	if tasks, ok := opsObj.(*[]model.OpsTask); ok {
		var projectNameMap map[uint]string
		if projectNameMap, err = s.getOpsTasksProjectName(tasks); err != nil {
			return nil, err
		}
		var hostNameMap map[uint]string
		if hostNameMap, err = s.getOpsTaskHostNames(tasks); err != nil {
			return nil, err
		}
		for _, task := range *tasks {
			res := api.GetOpsTaskRes{
				ID:           task.ID,
				Name:         task.Name,
				HostId:       task.HostId,
				HostName:     hostNameMap[task.HostId],
				IsIntranet:   task.IsIntranet,
				IsConcurrent: task.IsConcurrent,
				ProjectName:  projectNameMap[task.ProjectId],
				ProjectId:    task.ProjectId,
			}
			result = append(result, res)
		}
		return &result, err
	}
	if task, ok := opsObj.(*model.OpsTask); ok {
		var (
			count       int64
			uintIds     []uint
			notExistIds []uint
			err2        error
		)
		res := api.GetOpsTaskRes{
			ID:              task.ID,
			Name:            task.Name,
			CheckTemplateId: task.CheckTemplateId,
			HostId:          task.HostId,
			IsIntranet:      task.IsIntranet,
			IsConcurrent:    task.IsConcurrent,
			ProjectId:       task.ProjectId,
		}
		if err = model.DB.Model(&model.Project{}).Where("id = ?", task.ProjectId).Pluck("name", &res.ProjectName).Error; err != nil {
			return nil, fmt.Errorf("查询项目名称失败: %v", err)
		}
		if err = model.DB.Model(&model.Host{}).Where("id = ?", task.HostId).Pluck("name", &res.HostName).Error; err != nil {
			return nil, fmt.Errorf("查询服务器名称失败: %v", err)
		}
		// 判断模板ID和审批人ID还是否在template表和user表中
		if uintIds, err = util.StringToUintSlice(task.TemplateIds); err != nil {
			return nil, fmt.Errorf("task中的TemplateIds 格式不合规: %v", err)
		}

		if err = model.DB.Model(&model.OpsTemplate{}).Where("id IN (?)", uintIds).Count(&count).Error; count != int64(len(uintIds)) || err != nil {
			if notExistIds, err2 = util.FindNotExistIDs(consts.MysqlTableNameOpsTemplate, uintIds); err2 != nil {
				return nil, fmt.Errorf("查询 %s 失败: %v", consts.MysqlTableNameOpsTemplate, err2)
			}
			return nil, fmt.Errorf("%s 不存在ID: %d, 如果查询失败: %v", consts.MysqlTableNameOpsTemplate, notExistIds, err)
		}
		if err = json.Unmarshal([]byte(task.TemplateIds), &res.TemplateIds); err != nil {
			return nil, fmt.Errorf("task中的TemplateIds 不符合 json 格式: %v", err)
		}
		if task.Auditors != nil {
			if uintIds, err = util.StringToUintSlice(*task.Auditors); err != nil {
				return nil, fmt.Errorf("task中的TemplateIds 格式不合规: %v", err)
			}

			if err = model.DB.Model(&model.User{}).Where("id IN (?)", uintIds).Count(&count).Error; count != int64(len(uintIds)) || err != nil {
				if notExistIds, err2 = util.FindNotExistIDs(consts.MysqlTableNameUser, uintIds); err2 != nil {
					return nil, fmt.Errorf("查询 %s 失败: %v", consts.MysqlTableNameUser, err2)
				}
				return nil, fmt.Errorf("%s 不存在ID: %d, 如果查询失败: %v", consts.MysqlTableNameUser, notExistIds, err)
			}
			if err = json.Unmarshal([]byte(*task.Auditors), &res.Auditors); err != nil {
				return nil, fmt.Errorf("task中的Auditors 不符合 json 格式: %v", err)
			}
		}

		result = append(result, res)
		return &result, err
	}
	return &result, errors.New("转换运维操作任务信息结果失败")
}

// 删除运维操作模板
func (s *OpsService) DeleteOpsTask(ids []uint) (err error) {
	if err = model.DB.Where("id IN (?)", ids).Delete(&model.OpsTask{}).Error; err != nil {
		return fmt.Errorf("删除运维操作任务信息失败 %d: %v", ids, err)
	}
	return nil
}

// OpsTemplate按照自己提供的id顺序排序
func sortTemplatesByIds(templates []model.OpsTemplate, temIds []uint) []model.OpsTemplate {
	templateMap := make(map[uint]model.OpsTemplate)
	for _, template := range templates {
		templateMap[template.ID] = template
	}

	sortedTemplates := make([]model.OpsTemplate, 0, len(temIds))
	for _, id := range temIds {
		if template, exists := templateMap[id]; exists {
			sortedTemplates = append(sortedTemplates, template)
		}
	}
	return sortedTemplates
}

// 提取运营文案中的参数分给关联的模板
func (s *OpsService) extractOpsTaskParams(temIds []uint, content string) (commands []string, err error) {
	var (
		params          []model.OpsParam
		templates       []model.OpsTemplate
		templateContent string
		paramMap        = make(map[string]string)
	)
	if err = model.DB.Model(&model.OpsTemplate{}).Where("id IN (?)", temIds).Find(&templates).Error; err != nil {
		return nil, fmt.Errorf("查询运维操作模板失败: %v", err)
	}

	// 按照temIds的顺序对templates进行排序
	templates = sortTemplatesByIds(templates, temIds)

	// 解析content，提取参数值
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		// 兼容中文的：
		if len(parts) != 2 {
			parts = strings.SplitN(line, "：", 2)
		}
		if len(parts) == 2 {
			paramMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	// JOIN获取每一个模板对应的参数
	for _, tem := range templates {
		if params, err = s.GetTemplateParams(tem.ID); err != nil {
			return nil, fmt.Errorf("查询运维操作模板 %d 对应的参数模板失败: %v", tem.ID, err)
		}
		templateContent = tem.Content
		for _, param := range params {
			// 从文案中找到对应的keyword，取出: 后面的值一直到那一行换行位置，并去除值左右的空格
			if value, ok := paramMap[param.Keyword]; ok {
				templateContent = strings.ReplaceAll(templateContent, "${"+param.Variable+"}", value)
			}
		}
		commands = append(commands, templateContent)
	}
	return commands, err
}

// 生成sshCmd
func (s *OpsService) generateSSHCmd(host *model.Host, commands []string, isIntranet bool) (*[]api.SSHRunReq, error) {
	var (
		err     error
		sshReqs []api.SSHRunReq
	)
	for _, command := range commands {
		sshReq := api.SSHRunReq{
			Username:   global.Conf.SshConfig.OpsSSHUsername,
			SSHPort:    host.SSHPort,
			Key:        global.OpsSSHKey,
			Passphrase: nil,
			Cmd:        command,
		}
		if isIntranet {
			sshReq.HostIp = host.Vip
		} else {
			sshReq.HostIp = host.Ipv4
		}
		if global.Conf.SshConfig.OpsKeyPassphrase != "" {
			sshReq.Passphrase = []byte(global.Conf.SshConfig.OpsKeyPassphrase)
		}
		sshReqs = append(sshReqs, sshReq)
	}
	if len(sshReqs) != len(commands) {
		return nil, fmt.Errorf("生成的命令和数量不匹配")
	}
	return &sshReqs, err
}

// 执行单个运维任务
func (s *OpsService) RunOpsTaskCheckScript(params api.RunOpsTaskCheckScriptReq) (result *[]api.SSHResultRes, err error) {
	var (
		commands []string
		task     model.OpsTask
		host     model.Host
		sshReqs  *[]api.SSHRunReq
	)
	if err = model.DB.Model(&model.OpsTask{}).Where("id = ?", params.TaskId).First(&task).Error; err != nil {
		return nil, fmt.Errorf("查询task失败, taskID: %d, err: %v", params.TaskId, err)
	}
	if task.CheckTemplateId == 0 {
		return nil, fmt.Errorf("task %d 的CheckTemplateId为0", task.ID)
	}

	if commands, err = s.extractOpsTaskParams([]uint{task.CheckTemplateId}, params.ExecContext); err != nil {
		return nil, err
	}
	if err = model.DB.Model(&model.Host{}).Where("id = ?", task.HostId).First(&host).Error; err != nil {
		return nil, fmt.Errorf("查询host失败, hostID: %d, err: %v", task.HostId, err)
	}
	if sshReqs, err = s.generateSSHCmd(&host, commands, task.IsIntranet); err != nil {
		return nil, fmt.Errorf("生成SSH命令失败: %v", err)
	}
	if result, err = SSH().RunSSHCmdAsync(sshReqs); err != nil {
		return nil, fmt.Errorf("执行SSH命令失败: %v", err)
	}
	return result, err
}

// 查看执行运维任务的命令
func (s *OpsService) GetOpsTaskTmpCommands(params api.GetOpsTaskTmpCommandsReq) (commands []string, err error) {
	var count int64
	if err = model.DB.Model(&model.OpsTemplate{}).Where("id IN (?)", params.TemplateIds).Count(&count).Error; count != int64(len(params.TemplateIds)) || err != nil {
		notExistIds, err2 := util.FindNotExistIDs(consts.MysqlTableNameOpsTemplate, params.TemplateIds)
		if err2 != nil {
			return nil, fmt.Errorf("查询OpsTemplate失败: %v", err2)
		}
		return nil, fmt.Errorf("OpsTemplate 不存在ID: %d, 如果查询OpsTemplate失败: %v", notExistIds, err)
	}
	// 提取运营文案中的参数分给关联的模板
	if commands, err = s.extractOpsTaskParams(params.TemplateIds, params.ExecContext); err != nil {
		return nil, fmt.Errorf("提取运营文案中的参数分给关联的模板失败: %v", err)
	}
	if len(params.TemplateIds) != len(commands) {
		return nil, fmt.Errorf("模板ID和生成的命令数量不匹配")
	}
	return commands, err
}

// 提交运维操作任务
func (s *OpsService) SubmitOpsTask(params api.SubmitOpsTaskReq, submitter uint) (err error) {
	var (
		task           model.OpsTask
		host           model.Host
		taskLog        model.OpsTaskLog
		stepStatusList []api.OpsTaskLogtepStatus
		commands       []string
	)
	if len(params.TemplateIds) < 1 {
		return errors.New("模板ID不能为空")
	}
	if err = model.DB.Model(&model.OpsTask{}).Where("id = ?", params.TaskId).First(&task).Error; err != nil {
		return fmt.Errorf("查询task失败，taskID: %d, err: %v", params.TaskId, err)
	}

	taskLog = model.OpsTaskLog{
		Name:          task.Name,
		RejectAuditor: 0,
		TaskId:        task.ID,
		ProjectId:     task.ProjectId,
		Submitter:     submitter,
		StartTime:     nil,
		ExecContext:   params.ExecContext,
		CheckResponse: params.CheckResponse,
	}
	if err = model.DB.Where("id = ?", task.HostId).Select("ipv4", "vip", "ipv6").First(&host).Error; err != nil {
		return fmt.Errorf("查询host失败，hostID: %d, err: %v", task.HostId, err)
	}
	if task.IsIntranet {
		if host.Vip == "" {
			return errors.New("host没有内网IP")
		}
		taskLog.HostIp = host.Vip
	} else if host.Ipv4 != "" {
		taskLog.HostIp = host.Ipv4
	} else if host.Ipv6 != nil && *host.Ipv6 != "" {
		taskLog.HostIp = *host.Ipv6
	} else {
		return errors.New("host的参数错误")
	}
	if params.Auditors != nil && len(params.Auditors) > 0 {
		if taskLog.Auditors, err = util.UintSliceToString(params.Auditors); err != nil {
			return fmt.Errorf("auditors转换为string失败: %v", err)
		}
		if taskLog.PendingAuditors, err = util.UintSliceToString(params.Auditors); err != nil {
			return fmt.Errorf("pendingAuditors转换为string失败: %v", err)
		}
	} else {
		taskLog.Auditors = "[]"
		taskLog.PendingAuditors = "[]"
	}
	if len(params.Auditors) == 0 {
		taskLog.Status = consts.OpsTaskStatusIsRunning
	} else {
		taskLog.Status = consts.OpsTaskStatusIsWaiting
	}
	// 提取运营文案中的参数分给关联的模板
	if commands, err = s.extractOpsTaskParams(params.TemplateIds, params.ExecContext); err != nil {
		return fmt.Errorf("提取运营文案中的参数分给关联的模板失败: %v", err)
	}
	if len(params.TemplateIds) != len(commands) {
		return fmt.Errorf("模板ID和生成的命令数量不匹配")
	}
	commandsByte, err := json.Marshal(commands)
	if err != nil {
		return fmt.Errorf("commands转换json失败: %v", err)
	}
	taskLog.Commands = string(commandsByte)

	// 做成命令和状态的对应关系
	for _, command := range commands {
		stepStatus := api.OpsTaskLogtepStatus{
			Command: command,
			Status:  consts.OpsTaskStatusIsWaiting,
		}
		stepStatusList = append(stepStatusList, stepStatus)
	}
	if len(stepStatusList) != len(commands) {
		return fmt.Errorf("命令和状态的对应关系不匹配")
	}
	stepStatusByte, err := json.Marshal(stepStatusList)
	if err != nil {
		return fmt.Errorf("stepStatus转换json失败: %v", err)
	}
	taskLog.StepStatus = string(stepStatusByte)
	if err = model.DB.Create(&taskLog).Error; err != nil {
		return fmt.Errorf("创建运维操作任务日志失败: %v", err)
	}
	if taskLog.Status == consts.OpsTaskStatusIsRunning {
		if err = s.runOpsTaskCommands(&taskLog); err != nil {
			return fmt.Errorf("无需审批直接执行，但是启动执行task命令失败: %v", err)
		}
	}
	return err
}

func (s *OpsService) updateStepData(stepStatusList *[]api.OpsTaskLogtepStatus, command string, startTime string,
	endTime string, status int, res string, sshResStatus int) {
	for i := range *stepStatusList {
		stepStatus := &(*stepStatusList)[i]
		if stepStatus.Command == command {
			if status != consts.OpsTaskStatusIsWaiting {
				stepStatus.Status = status
			}
			if startTime != "" {
				stepStatus.StartTime = startTime
			}
			if endTime != "" {
				stepStatus.EndTime = endTime
			}
			if res != "" {
				stepStatus.Response = res
			}
			stepStatus.SSHResponseStatus = sshResStatus
			break
		}
	}
}

// 执行串行的运维操作任务
func (s *OpsService) RunOpsTaskSequential(sshReqs *[]api.SSHRunReq, taskLog *model.OpsTaskLog) {
	var (
		err            error
		stepStatusList []api.OpsTaskLogtepStatus
		stepStatusByte []byte
	)
	if err = json.Unmarshal([]byte(taskLog.StepStatus), &stepStatusList); err != nil {
		logger.Log().Error("ops", "RunOpsTaskSequential查询任务日志中的StepStatus失败", err)
		taskLog.Status = consts.OpsTaskStatusIsFailed
		return
	}
	for _, sshReq := range *sshReqs {
		cmdReq := []api.SSHRunReq{sshReq}
		s.updateStepData(&stepStatusList, sshReq.Cmd, time.Now().Format("2006-01-02 15:04:05"), "", consts.OpsTaskStatusIsWaiting, "", 0)
		result, err := SSH().RunSSHCmdAsync(&cmdReq)
		s.updateStepData(&stepStatusList, sshReq.Cmd, "", time.Now().Format("2006-01-02 15:04:05"), consts.OpsTaskStatusIsRunning, "", 0)
		var isBreak bool
		if err != nil {
			logger.Log().Error("ops", "RunOpsTaskSequential RunSSH失败", err)
			s.updateStepData(&stepStatusList, sshReq.Cmd, "", "", consts.OpsTaskStatusIsFailed, err.Error(), consts.SSHCustomCmdError)
			taskLog.Status = consts.OpsTaskStatusIsFailed
			isBreak = true
		}
		if (*result)[0].Status != 0 {
			s.updateStepData(&stepStatusList, sshReq.Cmd, "", "", consts.OpsTaskStatusIsFailed, (*result)[0].Response, (*result)[0].Status)
			taskLog.Status = consts.OpsTaskStatusIsFailed
			isBreak = true
		} else {
			s.updateStepData(&stepStatusList, sshReq.Cmd, "", "", consts.OpsTaskStatusIsSuccess, (*result)[0].Response, (*result)[0].Status)
		}
		if stepStatusByte, err = json.Marshal(stepStatusList); err != nil {
			logger.Log().Error("ops", "RunOpsTaskSequential转换stepStatusByte失败", err)
		}
		taskLog.StepStatus = string(stepStatusByte)
		if err = model.DB.Save(taskLog).Error; err != nil {
			logger.Log().Error("ops", "RunOpsTaskSequential保存taskLog失败", err)
			isBreak = true
		}
		if isBreak {
			break
		}
	}
	if taskLog.Status != consts.OpsTaskStatusIsFailed {
		endTime := time.Now()
		if err := model.DB.Model(&model.OpsTaskLog{}).Where("id = ?", taskLog.ID).Updates(model.OpsTaskLog{
			Status:  taskLog.Status,
			EndTime: &endTime,
		}).Error; err != nil {
			logger.Log().Error("ops", "RunOpsTaskConcurrent更新taskLog失败", err)
		}
	}
}

// 执行并行的运维操作任务
func (s *OpsService) RunOpsTaskConcurrent(sshReqs *[]api.SSHRunReq, taskLog *model.OpsTaskLog) {
	var (
		wg             sync.WaitGroup
		mu             sync.Mutex
		stepStatusList []api.OpsTaskLogtepStatus
		taskFailed     bool
	)
	if err := json.Unmarshal([]byte(taskLog.StepStatus), &stepStatusList); err != nil {
		logger.Log().Error("ops", "RunOpsTaskConcurrent查询任务日志中的StepStatus失败", err)
		taskLog.Status = consts.OpsTaskStatusIsFailed
		return
	}

	for _, sshReq := range *sshReqs {
		wg.Add(1)
		go func(sshReq api.SSHRunReq) {
			defer wg.Done()
			s.updateStepData(&stepStatusList, sshReq.Cmd, time.Now().Format("2006-01-02 15:04:05"), "", consts.OpsTaskStatusIsRunning, "", 0)
			result, err := SSH().RunSSHCmdAsync(&[]api.SSHRunReq{sshReq})
			mu.Lock()
			defer mu.Unlock()
			s.updateStepData(&stepStatusList, sshReq.Cmd, "", time.Now().Format("2006-01-02 15:04:05"), consts.OpsTaskStatusIsRunning, "", 0)
			if err != nil {
				logger.Log().Error("ops", "RunOpsTaskConcurrent RunSSH失败", err)
				s.updateStepData(&stepStatusList, sshReq.Cmd, "", "", consts.OpsTaskStatusIsFailed, err.Error(), consts.SSHCustomCmdError)
				taskFailed = true
				return
			}
			if (*result)[0].Status != 0 {
				s.updateStepData(&stepStatusList, sshReq.Cmd, "", "", consts.OpsTaskStatusIsFailed, (*result)[0].Response, (*result)[0].Status)
				taskFailed = true
				return
			}
			s.updateStepData(&stepStatusList, sshReq.Cmd, "", "", consts.OpsTaskStatusIsSuccess, (*result)[0].Response, (*result)[0].Status)

			var stepStatusByte []byte
			if stepStatusByte, err = json.Marshal(stepStatusList); err != nil {
				logger.Log().Error("ops", "RunOpsTaskConcurrent转换stepStatusByte失败", err)
			}
			taskLog.StepStatus = string(stepStatusByte)
			if err = model.DB.Save(taskLog).Error; err != nil {
				logger.Log().Error("ops", "RunOpsTaskConcurrent保存taskLog失败", err)
			}

		}(sshReq)
	}
	wg.Wait()

	if taskFailed {
		taskLog.Status = consts.OpsTaskStatusIsFailed
	} else {
		taskLog.Status = consts.OpsTaskStatusIsSuccess
	}
	endTime := time.Now()
	if err := model.DB.Model(&model.OpsTaskLog{}).Where("id = ?", taskLog.ID).Updates(model.OpsTaskLog{
		Status:  taskLog.Status,
		EndTime: &endTime,
	}).Error; err != nil {
		logger.Log().Error("ops", "RunOpsTaskConcurrent更新taskLog失败", err)
	}
}

// 执行task命令
func (s *OpsService) runOpsTaskCommands(taskLog *model.OpsTaskLog) (err error) {
	var (
		commands []string
		task     model.OpsTask
		host     model.Host
		sshReqs  *[]api.SSHRunReq
	)

	if err = model.DB.Model(&model.OpsTask{}).Where("id = ?", taskLog.TaskId).First(&task).Error; err != nil {
		return fmt.Errorf("查询task失败, taskID: %d, err: %v", taskLog.TaskId, err)
	}

	if err = json.Unmarshal([]byte(taskLog.Commands), &commands); err != nil {
		return fmt.Errorf("taskLog中的Commands不符合 json 格式: %v", err)
	}
	if err = model.DB.Model(&model.Host{}).Where("id = ?", task.HostId).First(&host).Error; err != nil {
		return fmt.Errorf("查询host失败, hostID: %d, err: %v", task.HostId, err)
	}

	if sshReqs, err = s.generateSSHCmd(&host, commands, task.IsIntranet); err != nil {
		return fmt.Errorf("生成SSH命令失败: %v", err)
	}

	if err = model.DB.Model(&model.OpsTaskLog{}).Where("id = ?", taskLog.ID).Update("start_time", time.Now()).Error; err != nil {
		return fmt.Errorf("更新taskLog的startTime失败: %v", err)
	}

	// 按照属性决定并发执行还是串行执行
	go func() {
		if task.IsConcurrent {
			s.RunOpsTaskConcurrent(sshReqs, taskLog)
		} else {
			s.RunOpsTaskSequential(sshReqs, taskLog)
		}
	}()
	return err
}

func (s *OpsService) ApproveOpsTask(params api.ApproveOpsTaskReq, uid uint) (err error) {
	var (
		taskLog         model.OpsTaskLog
		pendingAuditors []uint
	)
	tx := model.DB.Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()
	// 对当前操作的taskLog加行锁
	if err = tx.Model(&model.OpsTaskLog{}).Where("id = ?", params.ID).First(&taskLog).Set("gorm:query_option", "FOR UPDATE").Error; err != nil {
		return fmt.Errorf("查询taskLog并加锁失败，taskLogID: %d, err: %v", params.ID, err)
	}

	if taskLog.Status != consts.OpsTaskStatusIsWaiting {
		return fmt.Errorf("taskLog %d 不是等待审批状态", params.ID)
	}

	if err = json.Unmarshal([]byte(taskLog.PendingAuditors), &pendingAuditors); err != nil {
		return fmt.Errorf("taskLog %d PendingAuditors转换失败: %v", params.ID, err)
	}

	if !util.IsUintSliceContain(pendingAuditors, uid) {
		return fmt.Errorf("taskLog %d PendingAuditors中不包含uid: %d", params.ID, uid)
	}
	util.DeleteUintSliceByPtr(&pendingAuditors, uid)
	if taskLog.PendingAuditors, err = util.UintSliceToString(pendingAuditors); err != nil {
		return fmt.Errorf("pendingAuditors转换为string失败: %v", err)
	}
	if !params.IsAllow {
		taskLog.RejectAuditor = uid
		taskLog.Status = consts.OpsTaskStatusIsRejected
	}

	if len(pendingAuditors) == 0 && taskLog.Status == consts.OpsTaskStatusIsWaiting && taskLog.RejectAuditor == 0 {
		taskLog.Status = consts.OpsTaskStatusIsRunning
	}

	if err = tx.Save(&taskLog).Error; err != nil {
		return fmt.Errorf("更新运维操作任务日志失败: %v", err)
	}
	tx.Commit()
	if taskLog.Status == consts.OpsTaskStatusIsRunning {
		if err = s.runOpsTaskCommands(&taskLog); err != nil {
			return fmt.Errorf("最后一人审批完毕，但是启动执行task命令失败: %v", err)
		}
	}

	return err
}

func (s *OpsService) GetOpsTaskNeedApprove(wsConn *websocket.Conn, uid uint) error {
	if wsConn == nil {
		return fmt.Errorf("WebSocket connection is nil")
	}
	var (
		count int64
		err   error
	)
	// 处理pong消息
	wsConn.SetReadDeadline(time.Now().Add(consts.WebSocketPongWait))
	wsConn.SetPongHandler(func(appData string) error {
		// 重置读取存活时间
		wsConn.SetReadDeadline(time.Now().Add(consts.WebSocketPongWait))
		return nil
	})
	// 实时同步任务状态
	ticker := time.NewTicker(consts.WebSocketPingWait)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if wsConn == nil {
				return fmt.Errorf("WebSocket connection is nil")
			}
			if err = wsConn.WriteMessage(websocket.PingMessage, nil); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					// 前端主动断开连接，返回nil
					return nil
				}
				return fmt.Errorf("发送Ping失败: %v", err)
			}
		default:
			if err = model.DB.Model(&model.OpsTaskLog{}).
				Where("status = ?", consts.OpsTaskStatusIsWaiting).
				Where("JSON_CONTAINS(pending_auditors, ?)", uid).
				Count(&count).Error; err != nil {
				logger.Log().Error("ops", "查询用户是否有需要审批的任务失败", err)
				return fmt.Errorf("查询用户是否有需要审批的任务失败: %v", err)
			}
			if count > 0 {
				message := []byte(fmt.Sprintf("%t", count > 0))
				if wsConn == nil {
					return fmt.Errorf("WebSocket connection is nil")
				}
				if err = wsConn.WriteMessage(websocket.TextMessage, message); err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
						// 前端主动断开连接，返回nil
						return nil
					}
					return fmt.Errorf("发送任务状态失败: %v", err)
				}
			}
			time.Sleep(consts.WebSocketPingWait)
		}
	}
}

func (s *OpsService) GetUserTaskPending(uid uint, bindProjectIds []uint, params api.GetUserTaskPendingReq) (*api.GetOpsTaskLogsRes, error) {
	var (
		taskLogs []model.OpsTaskLog
		count    int64
		res      *[]api.GetOpsTaskLogRes
		err      error
	)
	getDB := model.DB.Model(&model.OpsTaskLog{}).Where("status = ?", consts.OpsTaskStatusIsWaiting).
		Where("JSON_CONTAINS(pending_auditors, ?)", uid).Find(&taskLogs)
	if params.TaskName != "" {
		sqlName := "%" + strings.ToUpper(params.TaskName) + "%"
		getDB = getDB.Where("UPPER(name) LIKE ?", sqlName)
	}
	if params.ProjectId != 0 {
		if !util.IsUintSliceContain(bindProjectIds, params.ProjectId) {
			return nil, errors.New("用户无权限查看该项目的运维操作任务日志")
		}
		getDB = getDB.Where("project_id = ?", params.ProjectId)
	} else {
		getDB = getDB.Where("project_id IN (?)", bindProjectIds)
	}

	if err = getDB.Count(&count).Error; err != nil {
		return nil, fmt.Errorf("查询用户待审批任务失败: %v", err)
	}
	if params.Page != 0 && params.PageSize != 0 {
		if err = getDB.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize).Find(&taskLogs).Error; err != nil {
			return nil, fmt.Errorf("查询用户待审批任务失败: %v", err)
		}
	} else {
		if err = getDB.Find(&taskLogs).Error; err != nil {
			return nil, fmt.Errorf("查询用户待审批任务失败: %v", err)
		}
	}
	if res, err = s.getOpsTaskLogResult(&taskLogs, true); err != nil {
		return nil, fmt.Errorf("运维操作任务日志转换结果失败: err: %v", err)
	} else if res == nil {
		return nil, errors.New("运维操作任务日志转换结果失败, res为nil")
	}
	result := api.GetOpsTaskLogsRes{
		Records:  *res,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    count,
	}

	return &result, err
}

func (s *OpsService) getOpsTaskLogUserNames(auditors []uint, rejectAuditor uint, submitter uint, pendingAuditors []uint) (
	auditorNames []string, rejectAuditorName string, submitterName string, pendingAuditorNames []string, err error) {
	// 创建一个集合来存储所有唯一的用户ID
	userIdsMap := make(map[uint]struct{})

	// 添加审核员ID
	for _, auditor := range auditors {
		if auditor != 0 {
			userIdsMap[auditor] = struct{}{}
		}
	}

	// 添加提交者ID
	if submitter != 0 {
		userIdsMap[submitter] = struct{}{}
	}

	// 将集合中的用户ID转换为切片
	userIds := make([]uint, 0, len(userIdsMap))
	for userId := range userIdsMap {
		userIds = append(userIds, userId)
	}

	// 查询用户信息
	var users []model.User
	if err = model.DB.Model(&model.User{}).Where("id IN (?)", userIds).Select("id", "nickname").Find(&users).Error; err != nil {
		return nil, "", "", nil, fmt.Errorf("查询用户信息失败: %v", err)
	}

	// 构建用户ID到用户名的映射
	userNameMap := make(map[uint]string)
	for _, user := range users {
		userNameMap[user.ID] = user.Nickname
	}

	// 根据传入的用户ID列表，获取对应的用户名
	for _, auditor := range auditors {
		auditorNames = append(auditorNames, userNameMap[auditor])
	}

	if rejectAuditor != 0 {
		rejectAuditorName = userNameMap[rejectAuditor]
	}

	submitterName = userNameMap[submitter]

	if pendingAuditors != nil && len(pendingAuditors) > 0 {
		for _, pendingAuditor := range pendingAuditors {
			pendingAuditorNames = append(pendingAuditorNames, userNameMap[pendingAuditor])
		}
	}

	return auditorNames, rejectAuditorName, submitterName, pendingAuditorNames, err
}

func (s *OpsService) getOpsTaskLogResult(opsObj any, isDetail bool) (*[]api.GetOpsTaskLogRes, error) {
	var (
		result []api.GetOpsTaskLogRes
		err    error
	)
	// 批量查询不需要temIds和auditors
	if taskLogs, ok := opsObj.(*[]model.OpsTaskLog); ok {
		var projectNameMap map[uint]string
		if projectNameMap, err = s.getOpsTaskLogsProjectName(taskLogs); err != nil {
			return nil, err
		}
		for _, taskLog := range *taskLogs {
			res := api.GetOpsTaskLogRes{
				ID:            taskLog.ID,
				Name:          taskLog.Name,
				Status:        taskLog.Status,
				HostIp:        taskLog.HostIp,
				RejectAuditor: taskLog.RejectAuditor,
				ProjectName:   projectNameMap[taskLog.ProjectId],
				ProjectId:     taskLog.ProjectId,
				Submitter:     taskLog.Submitter,
			}
			if taskLog.StartTime != nil {
				res.StartTime = taskLog.StartTime.Format("2006-01-02 15:04:05")
			}
			if taskLog.EndTime != nil {
				res.EndTime = taskLog.EndTime.Format("2006-01-02 15:04:05")
			}
			if err = json.Unmarshal([]byte(taskLog.Auditors), &res.Auditors); err != nil {
				return nil, fmt.Errorf("taskLog中的Auditors 不符合 json 格式: %v", err)
			}
			if isDetail {
				if err = json.Unmarshal([]byte(taskLog.Commands), &res.Commands); err != nil {
					return nil, fmt.Errorf("taskLog中的Commands 不符合 json 格式: %v", err)
				}
				if err = json.Unmarshal([]byte(taskLog.StepStatus), &res.StepStatus); err != nil {
					return nil, fmt.Errorf("taskLog中的StepStatus 不符合 json 格式: %v", err)
				}
				if err = json.Unmarshal([]byte(taskLog.PendingAuditors), &res.PendingAuditors); err != nil {
					return nil, fmt.Errorf("taskLog中的PendingAuditors 不符合 json 格式: %v", err)
				}
				if res.AuditorNames, res.RejectAuditorName, res.SubmitterName, res.PendingAuditorNames, err = s.
					getOpsTaskLogUserNames(res.Auditors, res.RejectAuditor, res.Submitter, res.PendingAuditors); err != nil {
					return nil, fmt.Errorf("查询用户信息失败: %v", err)
				}
				res.ExecContext = taskLog.ExecContext
				res.CheckResponse = taskLog.CheckResponse
			} else {
				if res.AuditorNames, res.RejectAuditorName, res.SubmitterName, _, err = s.
					getOpsTaskLogUserNames(res.Auditors, res.RejectAuditor, res.Submitter, res.PendingAuditors); err != nil {
					return nil, fmt.Errorf("查询用户信息失败: %v", err)
				}
			}
			result = append(result, res)
		}
		return &result, err
	}
	if taskLog, ok := opsObj.(*model.OpsTaskLog); ok {
		res := api.GetOpsTaskLogRes{
			ID:            taskLog.ID,
			Name:          taskLog.Name,
			Status:        taskLog.Status,
			HostIp:        taskLog.HostIp,
			RejectAuditor: taskLog.RejectAuditor,
			ProjectId:     taskLog.ProjectId,
			Submitter:     taskLog.Submitter,
		}
		if taskLog.StartTime != nil {
			res.StartTime = taskLog.StartTime.Format("2006-01-02 15:04:05")
		}
		if taskLog.EndTime != nil {
			res.EndTime = taskLog.EndTime.Format("2006-01-02 15:04:05")
		}
		if err = model.DB.Model(&model.Project{}).Where("id = ?", taskLog.ProjectId).Pluck("name", &res.ProjectName).Error; err != nil {
			return nil, fmt.Errorf("查询项目名称失败: %v", err)
		}
		if err = json.Unmarshal([]byte(taskLog.Auditors), &res.Auditors); err != nil {
			return nil, fmt.Errorf("taskLog中的Auditors 不符合 json 格式: %v", err)
		}
		if isDetail {
			if err = json.Unmarshal([]byte(taskLog.Commands), &res.Commands); err != nil {
				return nil, fmt.Errorf("taskLog中的Commands 不符合 json 格式: %v", err)
			}
			if err = json.Unmarshal([]byte(taskLog.StepStatus), &res.StepStatus); err != nil {
				return nil, fmt.Errorf("taskLog中的StepStatus 不符合 json 格式: %v", err)
			}
			if err = json.Unmarshal([]byte(taskLog.PendingAuditors), &res.PendingAuditors); err != nil {
				return nil, fmt.Errorf("taskLog中的PendingAuditors 不符合 json 格式: %v", err)
			}
			res.ExecContext = taskLog.ExecContext
			res.CheckResponse = taskLog.CheckResponse
		}
		result = append(result, res)
		return &result, err
	}
	return &result, errors.New("转换运维操作任务日志的结果失败")
}

func (s *OpsService) GetOpsTaskLog(params api.GetOpsTaskLogReq, bindProjectIds []uint) (*api.GetOpsTaskLogsRes, error) {
	var (
		err    error
		result api.GetOpsTaskLogsRes
		count  int64
		res    *[]api.GetOpsTaskLogRes
	)

	if params.ID != 0 {
		var taskLog model.OpsTaskLog
		if err = model.DB.Model(&model.OpsTaskLog{}).Where("id = ?", params.ID).First(&taskLog).Error; err != nil {
			return nil, fmt.Errorf("查询运维操作任务日志失败: %v", err)
		}
		res, err = s.getOpsTaskLogResult(&taskLog, true)
		if err != nil {
			return nil, fmt.Errorf("运维操作任务信息转换结果失败: err: %v", err)
		}
		if res == nil {
			return nil, errors.New("运维操作任务信息转换结果失败, res为nil")
		}

		result = api.GetOpsTaskLogsRes{
			Records:  *res,
			Page:     1,
			PageSize: 1,
			Total:    1,
		}
		return &result, err
	}

	getDB := model.DB.Model(&model.OpsTaskLog{}).Order("id DESC")
	if params.Name != "" {
		sqlName := "%" + strings.ToUpper(params.Name) + "%"
		getDB = getDB.Where("UPPER(name) LIKE ?", sqlName)
	}
	if params.Status != 0 {
		getDB = getDB.Where("status = ?", params.Status)
	} else {
		getDB = getDB.Where("status != ?", consts.OpsTaskStatusIsWaiting)
	}
	if params.ProjectId != 0 {
		if !util.IsUintSliceContain(bindProjectIds, params.ProjectId) {
			return nil, errors.New("用户无权限查看该项目的运维操作任务日志")
		}
		getDB = getDB.Where("project_id = ?", params.ProjectId)
	} else {
		getDB = getDB.Where("project_id IN (?)", bindProjectIds)
	}
	if params.Username != "" {
		sqlUsername := "%" + strings.ToUpper(params.Username) + "%"
		getDB = getDB.Joins("JOIN user ON user.id = ops_task_log.submitter").Where("UPPER(user.username) LIKE ?", sqlUsername)
	}
	if err = getDB.Count(&count).Error; err != nil {
		return nil, fmt.Errorf("查询运维操作任务日志总数失败: %v", err)
	}
	var taskLogs []model.OpsTaskLog
	if params.Page != 0 && params.PageSize != 0 {
		if err = getDB.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize).Find(&taskLogs).Error; err != nil {
			return nil, fmt.Errorf("查询运维操作任务日志失败: %v", err)
		}
	} else {
		if err = getDB.Find(&taskLogs).Error; err != nil {
			return nil, fmt.Errorf("查询运维操作任务日志失败: %v", err)
		}
	}
	res, err = s.getOpsTaskLogResult(&taskLogs, true)
	if err != nil {
		return nil, fmt.Errorf("运维操作任务日志转换结果失败: err: %v", err)
	}
	if res == nil {
		return nil, errors.New("运维操作任务日志转换结果失败, res为nil")
	}
	result = api.GetOpsTaskLogsRes{
		Records:  *res,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    count,
	}

	return &result, err
}

func (s *OpsService) GetOpsTaskRunningWS(wsConn *websocket.Conn, bindProjectIds []uint) (err error) {
	var (
		taskLogs []model.OpsTaskLog
		ticker   = time.NewTicker(3 * time.Second)
	)
	defer ticker.Stop()

	pingTicker := time.NewTicker(consts.WebSocketPingWait)
	defer pingTicker.Stop()

	// 处理pong消息
	wsConn.SetReadDeadline(time.Now().Add(consts.WebSocketPongWait))
	wsConn.SetPongHandler(func(appData string) error {
		// 重置读取存活时间
		wsConn.SetReadDeadline(time.Now().Add(consts.WebSocketPongWait))
		return nil
	})

	for {
		select {
		case <-ticker.C:
			// 获取当前时间和6秒前的时间
			now := time.Now()
			tenSecondsAgo := now.Add(-6 * time.Second)

			if err = model.DB.Model(&model.OpsTaskLog{}).
				Where("((status = ? OR status = ? OR status = ?) AND project_id IN (?)) AND (status = ? OR updated_at >= ?)",
					consts.OpsTaskStatusIsRunning, consts.OpsTaskStatusIsSuccess, consts.OpsTaskStatusIsFailed, bindProjectIds, consts.OpsTaskStatusIsRunning, tenSecondsAgo).
				Find(&taskLogs).Error; err != nil {
				return fmt.Errorf("查询运维操作任务日志失败: %v", err)
			}
			result, err := s.getOpsTaskLogResult(&taskLogs, true)
			if result == nil {
				return errors.New("运维操作任务日志转换结果失败, result为nil")
			}
			if err != nil {
				return fmt.Errorf("运维操作任务日志转换结果失败: err: %v", err)
			}
			for _, res := range *result {
				if err = wsConn.WriteJSON(res); err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
						// 前端主动断开连接，返回nil
						return nil
					}
					return fmt.Errorf("发送任务状态失败: %v", err)
				}
			}
		case <-pingTicker.C:
			if err = wsConn.WriteMessage(websocket.PingMessage, nil); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					// 前端主动断开连接，返回nil
					return nil
				}
				return fmt.Errorf("发送Ping失败: %v", err)
			}
		}

	}
}
