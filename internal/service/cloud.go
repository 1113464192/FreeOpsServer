package service

import (
	"FreeOps/internal/model"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
)

type CloudService struct{}

var insCloud CloudService

func CloudServiceApp() *CloudService {
	return &insCloud
}

func (s *CloudService) CreateCloudProject(name string, cloudPlatform string) error {
	var (
		err   error
		count int64
	)
	if err = model.DB.Model(&model.Project{}).Where("name = ? AND cloud_platform = ?", name, cloudPlatform).Count(&count).Error; err != nil {
		return fmt.Errorf("查询项目失败: %v", err)
	}
	if count < 1 {
		return errors.New("项目名在库中不存在")
	}

	cmd := exec.Command("python3", "/data/scripts/CloudManage.py", "-a", "createProject", "-p", cloudPlatform, "-n", name)
	if _, err = cmd.Output(); err != nil {
		return fmt.Errorf("创建云项目失败: %v", err)
	}
	return err
}

func (s *CloudService) CreateCloudHost(projectId uint, cloudPlatform string, hostCount uint64) (err error) {
	var (
		count int64
	)
	if err = model.DB.Model(&model.Project{}).Where("id = ? AND cloud_platform = ?", projectId, cloudPlatform).Count(&count).Error; err != nil {
		return fmt.Errorf("查询项目失败: %v", err)
	}
	if count < 1 {
		return errors.New("项目ID在库中不存在")
	}

	cid, err := s.GetCloudProjectId("", cloudPlatform, projectId)
	if err != nil {
		return fmt.Errorf("获取云项目ID失败: %v", err)
	}

	cmd := exec.Command("python3", "/data/scripts/CloudManage.py", "-a", "createHost", "-p", cloudPlatform, "-i", strconv.FormatUint(cid, 10), "-c", strconv.FormatUint(hostCount, 10))
	if _, err = cmd.Output(); err != nil {
		return fmt.Errorf("购买服务器失败: %v", err)
	}
	return err
}

func (s *CloudService) UpdateCloudProject(name string, cloudPlatform string) error {
	var (
		err   error
		count int64
	)
	if err = model.DB.Model(&model.Project{}).Where("name = ? AND cloud_platform = ?", name, cloudPlatform).Count(&count).Error; err != nil {
		return fmt.Errorf("查询项目失败: %v", err)
	}
	if count < 1 {
		return errors.New("项目名在库中不存在")
	}

	cid, err := s.GetCloudProjectId(name, cloudPlatform, 0)
	if err != nil {
		return fmt.Errorf("获取云项目ID失败: %v", err)
	}

	cmd := exec.Command("python3", "/data/scripts/CloudManage.py", "-a", "updateProject", "-p", cloudPlatform, "-i", strconv.FormatUint(cid, 10), "-n", name)
	if _, err = cmd.Output(); err != nil {
		return fmt.Errorf("更新云项目失败: %v", err)
	}
	return err
}

func (s *CloudService) GetCloudProjectId(name string, cloudPlatform string, projectId uint) (uint64, error) {
	var (
		err    error
		count  int64
		output []byte
	)
	if projectId == 0 {
		if err = model.DB.Model(&model.Project{}).Where("name = ? AND cloud_platform = ?", name, cloudPlatform).Count(&count).Error; err != nil {
			return 0, fmt.Errorf("查询项目失败: %v", err)
		}
		if count < 1 {
			return 0, errors.New("项目名在库中不存在")
		}
	} else {
		if err = model.DB.Model(&model.Project{}).Where("id = ?", projectId).Select("name").Find(&name).Error; err != nil {
			return 0, fmt.Errorf("查询项目失败: %v", err)
		}
	}

	cmd := exec.Command("python3", "/data/scripts/CloudManage.py", "-a", "createProject", "-p", cloudPlatform, "-n", name)
	if output, err = cmd.Output(); err != nil {
		return 0, fmt.Errorf("查询云项目失败: %v", err)
	}

	projectID, err := strconv.ParseUint(string(output), 10, 32)
	if err != nil {
		return 0, fmt.Errorf("转换项目ID失败: %v", err)
	}
	return projectID, nil
}
