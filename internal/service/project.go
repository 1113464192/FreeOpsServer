package service

import (
	"FreeOps/internal/model"
	"FreeOps/pkg/api"
	"errors"
	"fmt"
	"strings"
)

type ProjectService struct{}

var insProject ProjectService

func ProjectServiceApp() *ProjectService {
	return &insProject
}

// 修改/添加项目
func (s *ProjectService) UpdateProject(params *api.UpdateProjectReq) (err error) {

	var (
		project model.Project
		count   int64
	)
	if params.ID != 0 {
		if err = model.DB.Model(&model.Project{}).Where("id = ?", params.ID).Count(&count).Error; count != 1 || err != nil {
			return fmt.Errorf("project ID不存在: %d, 或有错误信息: %v", params.ID, err)
		}

		if err = model.DB.Model(&project).Where("name = ? AND id != ?", params.Name, params.ID).Count(&count).Error; err != nil {
			return fmt.Errorf("查询项目失败: %v", err)
		} else if count > 0 {
			return fmt.Errorf("项目名已被使用: %s", params.Name)
		}

		if err := model.DB.Where("id = ?", params.ID).First(&project).Error; err != nil {
			return fmt.Errorf("项目查询失败: %v", err)
		}
		project.Name = params.Name

		if err = model.DB.Save(&project).Error; err != nil {
			return fmt.Errorf("数据保存失败: %v", err)
		}
		return err
	} else {
		err = model.DB.Model(&project).Where("name = ?", params.Name).Count(&count).Error
		if err != nil {
			return fmt.Errorf("查询项目失败: %v", err)
		} else if count > 0 {
			return fmt.Errorf("项目名(%s)已存在", params.Name)
		}

		project = model.Project{
			Name: params.Name,
		}
		if err = model.DB.Create(&project).Error; err != nil {
			return fmt.Errorf("创建项目失败: %v", err)
		}
		return err
	}
}

func (s *ProjectService) GetProjects(params *api.GetProjectsReq) (*api.GetProjectsRes, error) {
	var projects []model.Project
	var err error
	var count int64

	getDB := model.DB.Model(&model.Project{})
	if params.ID != 0 {
		getDB = getDB.Where("id = ?", params.ID)
	}

	if params.Name != "" {
		sqlName := "%" + strings.ToUpper(params.Name) + "%"
		getDB = getDB.Where("UPPER(name) LIKE ?", sqlName)
	}

	if err = getDB.Count(&count).Error; err != nil {
		return nil, fmt.Errorf("查询项目总数失败: %v", err)

	}
	if params.Page != 0 && params.PageSize != 0 {
		if err = getDB.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize).Find(&projects).Error; err != nil {
			return nil, fmt.Errorf("查询项目失败: %v", err)
		}
	} else {
		if err = getDB.Find(&projects).Error; err != nil {
			return nil, fmt.Errorf("查询项目失败: %v", err)
		}
	}
	var res *[]api.GetProjectReq
	var result api.GetProjectsRes
	res, err = s.GetResults(&projects)
	if err != nil {
		return nil, err
	}
	result = api.GetProjectsRes{
		Records:  *res,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    count,
	}
	return &result, err
}

func (s *ProjectService) DeleteProjects(ids []uint) (err error) {
	var count int64
	if err = model.DB.Model(&model.Host{}).Where("project_id IN (?)", ids).Count(&count).Error; err != nil {
		return fmt.Errorf("查询项目关联服务器失败: %v", err)
	}
	if count > 0 {
		return errors.New("项目下还有服务器存在")
	}

	tx := model.DB.Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()
	if err = tx.Where("project_id IN (?)", ids).Delete(&model.RoleProject{}).Error; err != nil {
		return fmt.Errorf("删除角色项目关系失败 %d: %v", ids, err)
	}
	if err = tx.Where("id IN (?)", ids).Delete(&model.Project{}).Error; err != nil {
		return fmt.Errorf("删除项目失败 %d: %v", ids, err)
	}
	tx.Commit()
	return nil
}

func (s *ProjectService) GetProjectHosts(params api.IdsReq) (hostIds []uint, err error) {
	if err = model.DB.Model(model.Host{}).Where("project_id IN (?)", params.Ids).Pluck("id", &hostIds).Error; err != nil {
		return nil, fmt.Errorf("查询项目服务器IDs失败: %v", err)
	}
	return hostIds, err
}

func (s *ProjectService) GetProjectGames(params api.IdsReq) (gameIds []uint, err error) {
	if err = model.DB.Model(model.Game{}).Where("project_id IN (?)", params.Ids).Pluck("id", &gameIds).Error; err != nil {
		return nil, fmt.Errorf("查询项目游戏服IDs失败: %v", err)
	}
	return gameIds, err
}

func (s *ProjectService) GetResults(projectObj any) (*[]api.GetProjectReq, error) {
	var result []api.GetProjectReq
	var err error
	if projects, ok := projectObj.(*[]model.Project); ok {
		for _, project := range *projects {
			res := api.GetProjectReq{
				ID:   project.ID,
				Name: project.Name,
			}
			result = append(result, res)
		}
		return &result, err
	}
	if project, ok := projectObj.(*model.Project); ok {
		res := api.GetProjectReq{
			ID:   project.ID,
			Name: project.Name,
		}
		result = append(result, res)
		return &result, err
	}
	return &result, errors.New("转换项目结果失败")
}
