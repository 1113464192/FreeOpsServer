package service

import (
	"FreeOps/internal/consts"
	"FreeOps/internal/model"
	"FreeOps/pkg/api"
	"fmt"
)

type HomeService struct{}

var insHome HomeService

func HomeServiceApp() *HomeService {
	return &insHome
}

func (s *HomeService) GetHomeInfo() (result api.GetHomeInfoRes, err error) {
	if err = model.DB.Model(&model.Project{}).Count(&result.ProjectCount).Error; err != nil {
		return result, fmt.Errorf("获取项目数量失败: %w", err)
	}
	if err = model.DB.Model(&model.Host{}).Count(&result.HostCount).Error; err != nil {
		return result, fmt.Errorf("获取主机数量失败: %w", err)
	}
	if err = model.DB.Model(&model.OpsTaskLog{}).Where("status = ?", consts.OpsTaskStatusIsWaiting).Count(&result.TaskNeedApproveCount).Error; err != nil {
		return result, fmt.Errorf("获取待审批任务数量失败: %w", err)
	}
	if err = model.DB.Model(&model.OpsTaskLog{}).Where("status = ?", consts.OpsTaskStatusIsRunning).Count(&result.TaskRunningCount).Error; err != nil {
		return result, fmt.Errorf("获取运行中任务数量失败: %w", err)
	}
	return result, err
}
