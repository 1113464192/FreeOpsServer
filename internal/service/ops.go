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
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"strings"
	"sync"
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
	// 确保是[]uint格式
	if _, ok := interface{}(params.TemplateIds).([]uint); !ok {
		return fmt.Errorf("templateIds 不是[]uint格式")
	}

	if params.ID != 0 {
		if err = model.DB.Model(&model.OpsTask{}).Where("id != ? AND project_id = ? AND name = ?", params.ID, params.ProjectId, params.Name).Count(&count).Error; count > 0 || err != nil {
			return fmt.Errorf("同项目下除自身外仍有相同的名称的 TaskName,id: %d projectId: %d name: %s , 或有错误信息: %v", params.ID, params.ProjectId, params.Name, err)
		}
		if err = model.DB.Model(&model.OpsTask{}).Where("id == ?", params.ID).First(&task).Error; err != nil {
			return fmt.Errorf("查询运维操作任务信息失败: %v", err)
		}
		if params.Auditors != "" {
			// 确保是[]uint格式
			if _, ok := interface{}(params.Auditors).([]uint); !ok {
				return fmt.Errorf("auditors 不是[]uint格式")
			}
			*task.Auditors = params.Auditors
		}

		task.Name = params.Name
		task.CheckTemplateId = params.CheckTemplateId
		task.TemplateIds = params.TemplateIds
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
			TemplateIds:     params.TemplateIds,
			HostId:          params.HostId,
			IsIntranet:      params.IsIntranet,
			IsConcurrent:    params.IsConcurrent,
			ProjectId:       params.ProjectId,
		}
		if params.Auditors != "" {
			// 确保是[]uint格式
			if _, ok := interface{}(params.Auditors).([]uint); !ok {
				return fmt.Errorf("auditors 不是[]uint格式")
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
		getDB = getDB.Where("project_id = ?", params.ProjectId)
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

func (s *OpsService) getOpsTaskResult(opsObj any) (*[]api.GetOpsTaskRes, error) {
	var (
		result []api.GetOpsTaskRes
		err    error
	)
	// 批量查询不需要temIds和auditors
	if tasks, ok := opsObj.(*[]model.OpsTask); ok {
		for _, task := range *tasks {
			res := api.GetOpsTaskRes{
				ID:           task.ID,
				Name:         task.Name,
				HostId:       task.HostId,
				IsIntranet:   task.IsIntranet,
				IsConcurrent: task.IsConcurrent,
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
	}
	return &result, errors.New("转换运维操作任务信息结果失败")
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

// 提交运维操作任务
func (s *OpsService) SubmitOpsTask(params api.SubmitOpsTaskReq, submitter uint) (err error) {
	var (
		task           model.OpsTask
		taskLog        model.OpsTaskLog
		stepStatusList []api.OpsTaskLogtepStatus
		commands       []string
	)
	// 确保是[]uint格式
	if _, ok := interface{}(params.TemplateIds).([]uint); !ok {
		return fmt.Errorf("templateIds 不是[]uint格式")
	}

	if err = model.DB.Model(&model.OpsTask{}).Where("id = ?", params.TaskId).First(&task).Error; err != nil {
		return fmt.Errorf("查询task失败，taskID: %d, err: %v", params.TaskId, err)
	}

	taskLog = model.OpsTaskLog{
		Name:            task.Name,
		Status:          consts.OpsTaskStatusIsWaiting,
		Auditors:        "",
		PendingAuditors: "",
		RejectAuditor:   0,
		TaskId:          task.ID,
		ProjectId:       task.ProjectId,
		Submitter:       submitter,
	}
	if params.Auditors != nil {
		taskLog.Auditors = util.UintSliceToString(params.Auditors)
		taskLog.PendingAuditors = util.UintSliceToString(params.Auditors)
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
	return err
}

func (s *OpsService) updateStepData(stepStatusList *[]api.OpsTaskLogtepStatus, command string, status int, res string, sshResStatus int) {
	for i := range *stepStatusList {
		stepStatus := &(*stepStatusList)[i]
		if stepStatus.Command == command {
			if status != consts.OpsTaskStatusIsWaiting {
				stepStatus.Status = status
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
	taskLog.Status = consts.OpsTaskStatusIsSuccess
	if err = json.Unmarshal([]byte(taskLog.StepStatus), &stepStatusList); err != nil {
		logger.Log().Error("ops", "RunOpsTaskSequential查询任务日志中的StepStatus失败", err)
		taskLog.Status = consts.OpsTaskStatusIsFailed
		return
	}
	for _, sshReq := range *sshReqs {
		result, err := SSH().RunSSHCmdAsync(sshReqs)
		var isBreak bool
		if err != nil {
			logger.Log().Error("ops", "RunOpsTaskSequential RunSSH失败", err)
			s.updateStepData(&stepStatusList, sshReq.Cmd, consts.OpsTaskStatusIsFailed, err.Error(), consts.SSHCustomCmdError)
			taskLog.Status = consts.OpsTaskStatusIsFailed
			isBreak = true
		}
		if (*result)[0].Status != 0 {
			s.updateStepData(&stepStatusList, sshReq.Cmd, consts.OpsTaskStatusIsFailed, (*result)[0].Response, (*result)[0].Status)
			taskLog.Status = consts.OpsTaskStatusIsFailed
			isBreak = true
		} else {
			s.updateStepData(&stepStatusList, sshReq.Cmd, consts.OpsTaskStatusIsSuccess, (*result)[0].Response, (*result)[0].Status)
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
			result, err := SSH().RunSSHCmdAsync(&[]api.SSHRunReq{sshReq})
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				logger.Log().Error("ops", "RunOpsTaskConcurrent RunSSH失败", err)
				s.updateStepData(&stepStatusList, sshReq.Cmd, consts.OpsTaskStatusIsFailed, err.Error(), consts.SSHCustomCmdError)
				taskFailed = true
				return
			}
			if (*result)[0].Status != 0 {
				s.updateStepData(&stepStatusList, sshReq.Cmd, consts.OpsTaskStatusIsFailed, (*result)[0].Response, (*result)[0].Status)
				taskFailed = true
				return
			}
			s.updateStepData(&stepStatusList, sshReq.Cmd, consts.OpsTaskStatusIsSuccess, (*result)[0].Response, (*result)[0].Status)

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
	if err := model.DB.Model(&model.OpsTaskLog{}).Where("id = ?", taskLog.ID).Update("status", taskLog.Status).Error; err != nil {
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

	if err = tx.Model(&model.OpsTaskLog{}).Where("id = ?", params.TaskId).First(&taskLog).Set("gorm:query_option", "FOR UPDATE").Error; err != nil {
		return fmt.Errorf("查询taskLog并加锁失败，taskLogID: %d, err: %v", params.TaskId, err)
	}

	if taskLog.Status != consts.OpsTaskStatusIsWaiting {
		return fmt.Errorf("taskLog %d 不是等待审批状态", params.TaskId)
	}

	if err = json.Unmarshal([]byte(taskLog.PendingAuditors), &pendingAuditors); err != nil {
		return fmt.Errorf("taskLog %d PendingAuditors转换失败: %v", params.TaskId, err)
	}

	if !util.IsUintSliceContain(pendingAuditors, uid) {
		return fmt.Errorf("taskLog %d PendingAuditors中不包含uid: %d", params.TaskId, uid)
	}
	util.DeleteUintSliceByPtr(&pendingAuditors, uid)
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

	go func() {}()
	if taskLog.Status == consts.OpsTaskStatusIsRunning {
		if err = s.runOpsTaskCommands(&taskLog); err != nil {
			return fmt.Errorf("最后一人审批完毕，但是启动执行task命令失败: %v", err)
		}
	}

	return err
}

func (s *OpsService) GetTaskPendingApprovers() (*[]api.GetTaskPendingApproversRes, error) {
	var (
		taskLogs             []model.OpsTaskLog
		result               []api.GetTaskPendingApproversRes
		pendingApprovers     []uint
		pendingApproverNames []string
		SubmitterName        string
		err                  error
	)
	if err = model.DB.Model(&model.OpsTaskLog{}).Where("status = ?", consts.OpsTaskStatusIsWaiting).Find(&taskLogs).Error; err != nil {
		return nil, fmt.Errorf("查询等待审批的任务失败: %v", err)
	}
	for _, taskLog := range taskLogs {
		if err = json.Unmarshal([]byte(taskLog.PendingAuditors), &pendingApprovers); err != nil {
			return nil, fmt.Errorf("taskLog %d pendingApprovers转换失败: %v", taskLog.TaskId, err)
		}
		if err = model.DB.Model(&model.User{}).Where("id IN (?)", pendingApprovers).Pluck("username", &pendingApproverNames).Error; err != nil {
			return nil, fmt.Errorf("查询等待审批的任务的审批人失败: %v", err)
		}
		if err = model.DB.Model(&model.User{}).Where("id = ?", taskLog.Submitter).Pluck("username", &SubmitterName).Error; err != nil {
			return nil, fmt.Errorf("查询等待审批的任务的提交人失败: %v", err)
		}
		res := api.GetTaskPendingApproversRes{
			TaskName:     taskLog.Name,
			PendingUsers: pendingApproverNames,
			Submitter:    SubmitterName,
		}
		result = append(result, res)
	}
	return &result, err
}

func (s *OpsService) getOpsTaskLogResult(opsObj any) (*[]api.GetOpsTaskLogRes, error) {
	var (
		result []api.GetOpsTaskLogRes
		err    error
	)
	// 批量查询不需要temIds和auditors
	if taskLogs, ok := opsObj.(*[]model.OpsTaskLog); ok {
		for _, taskLog := range *taskLogs {
			res := api.GetOpsTaskLogRes{
				ID:            taskLog.ID,
				Name:          taskLog.Name,
				Status:        taskLog.Status,
				RejectAuditor: taskLog.RejectAuditor,
				ProjectId:     taskLog.ProjectId,
				Submitter:     taskLog.Submitter,
			}
			if err = json.Unmarshal([]byte(taskLog.Auditors), &res.Auditors); err != nil {
				return nil, fmt.Errorf("taskLog中的Auditors 不符合 json 格式: %v", err)
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
			RejectAuditor: taskLog.RejectAuditor,
			ProjectId:     taskLog.ProjectId,
			Submitter:     taskLog.Submitter,
		}
		if err = json.Unmarshal([]byte(taskLog.Commands), &res.Commands); err != nil {
			return nil, fmt.Errorf("taskLog中的Commands 不符合 json 格式: %v", err)
		}
		if err = json.Unmarshal([]byte(taskLog.StepStatus), &res.StepStatus); err != nil {
			return nil, fmt.Errorf("taskLog中的StepStatus 不符合 json 格式: %v", err)
		}
		if err = json.Unmarshal([]byte(taskLog.Auditors), &res.Auditors); err != nil {
			return nil, fmt.Errorf("taskLog中的Auditors 不符合 json 格式: %v", err)
		}
		if err = json.Unmarshal([]byte(taskLog.PendingAuditors), &res.PendingAuditors); err != nil {
			return nil, fmt.Errorf("taskLog中的PendingAuditors 不符合 json 格式: %v", err)
		}
		result = append(result, res)
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

		res, err = s.getOpsTaskLogResult(taskLog)
		if res == nil {
			return nil, errors.New("运维操作任务信息转换结果失败, res为nil")
		}
		if err != nil {
			return nil, fmt.Errorf("运维操作任务信息转换结果失败: err: %v", err)
		}

		result = api.GetOpsTaskLogsRes{
			Records:  *res,
			Page:     1,
			PageSize: 1,
			Total:    1,
		}
		return &result, err
	}

	getDB := model.DB.Model(&model.OpsTaskLog{})
	if params.Name != "" {
		sqlName := "%" + strings.ToUpper(params.Name) + "%"
		getDB = getDB.Where("UPPER(name) LIKE ?", sqlName)
	}
	if params.Status != 0 {
		getDB = getDB.Where("status = ?", params.Status)
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
		if err = getDB.Omit("commands", "step_status", "pending_auditors").Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize).Find(&taskLogs).Error; err != nil {
			return nil, fmt.Errorf("查询运维操作任务日志失败: %v", err)
		}
	} else {
		if err = getDB.Omit("commands", "step_status", "pending_auditors").Find(&taskLogs).Error; err != nil {
			return nil, fmt.Errorf("查询运维操作任务日志失败: %v", err)
		}
	}
	res, err = s.getOpsTaskLogResult(&taskLogs)
	if res == nil {
		return nil, errors.New("运维操作任务日志转换结果失败, res为nil")
	}
	if err != nil {
		return nil, fmt.Errorf("运维操作任务日志转换结果失败: err: %v", err)
	}
	result = api.GetOpsTaskLogsRes{
		Records:  *res,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    count,
	}

	return &result, err
}

func (s *OpsService) GetOpsTaskRunningWS(wsConn *websocket.Conn, c *gin.Context, bindProjectIds []uint) (err error) {
	var (
		taskLogs []model.OpsTaskLog
	)
	if err = model.DB.Model(&model.OpsTaskLog{}).Where("status = ? AND project_id IN (?)", consts.OpsTaskStatusIsRunning, bindProjectIds).Find(&taskLogs).Error; err != nil {
		return fmt.Errorf("查询运维操作任务日志失败: %v", err)
	}
	// 不应该无限循环，应该每次有数据变化再传吧？
	for {
		for _, taskLog := range taskLogs {

		}
	}

	for _, taskLog := range taskLogs {
		if !util.IsUintSliceContain(bindProjectIds, taskLog.ProjectId) {
			continue
		}
		res, err := s.getOpsTaskLogResult(&taskLog)
		if res == nil {
			logger.Log().Error("ops", "GetOpsTaskRunningWS运维操作任务日志转换结果失败, res为nil", err)
			continue
		}
		if err != nil {
			logger.Log().Error("ops", "GetOpsTaskRunningWS运维操作任务日志转换结果失败", err)
			continue
		}
		err = wsConn.WriteJSON(res)
		if err != nil {
			logger.Log().Error("ops", "GetOpsTaskRunningWS发送websocket消息失败", err)
			break
		}
	}
}
