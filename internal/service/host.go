package service

import (
	"FreeOps/internal/model"
	"FreeOps/pkg/api"
	"errors"
	"fmt"
	"strings"
)

type HostService struct{}

var insHost HostService

func HostServiceApp() *HostService {
	return &insHost
}

// 修改/添加服务器
func (s *HostService) UpdateHost(params *api.UpdateHostReq) (err error) {
	var (
		host  model.Host
		count int64
	)
	if params.ID != 0 {
		if err = model.DB.Model(&model.Host{}).Where("id = ?", params.ID).Count(&count).Error; count != 1 || err != nil {
			return fmt.Errorf("host ID不存在: %d, 或有错误信息: %v", params.ID, err)
		}

		if err = model.DB.Model(&host).Where("name = ? AND id != ?", params.Name, params.ID).Count(&count).Error; err != nil {
			return fmt.Errorf("查询服务器失败: %v", err)
		} else if count > 0 {
			return fmt.Errorf("服务器名已被使用: %s", params.Name)
		}

		if err := model.DB.Where("id = ?", params.ID).First(&host).Error; err != nil {
			return fmt.Errorf("服务器查询失败: %v", err)
		}

		host.Name = params.Name
		host.Ipv4 = params.Ipv4
		if params.Ipv6 != "" {
			*host.Ipv6 = params.Ipv6
		}
		host.Vip = params.Vip
		host.Zone = params.Zone
		host.Cloud = params.Cloud
		host.System = params.System
		host.Cores = params.Cores
		host.DataDisk = params.DataDisk
		host.Mem = params.Mem
		host.ProjectId = params.ProjectId

		if err = model.DB.Save(&host).Error; err != nil {
			return fmt.Errorf("数据保存失败: %v", err)
		}
		return err
	} else {
		err = model.DB.Model(&host).Where("name = ?", params.Name).Count(&count).Error
		if err != nil {
			return fmt.Errorf("查询服务器失败: %v", err)
		} else if count > 0 {
			return fmt.Errorf("服务器名(%s)已存在", params.Name)
		}
		host = model.Host{
			Name:      params.Name,
			Ipv4:      params.Ipv4,
			Vip:       params.Vip,
			Zone:      params.Zone,
			Cloud:     params.Cloud,
			System:    params.System,
			Cores:     params.Cores,
			DataDisk:  params.DataDisk,
			Mem:       params.Mem,
			ProjectId: params.ProjectId,
		}
		if params.Ipv6 != "" {
			*host.Ipv6 = params.Ipv6
		}

		if err = model.DB.Create(&host).Error; err != nil {
			return fmt.Errorf("创建服务器失败: %v", err)
		}
		return err
	}
}

func (s *HostService) GetHosts(params *api.GetHostsReq) (*api.GetHostsRes, error) {
	var hosts []model.Host
	var err error
	var count int64

	getDB := model.DB.Model(&model.Host{})
	if params.ID != 0 {
		getDB = getDB.Where("id = ?", params.ID)
	}

	if params.Name != "" {
		sqlName := "%" + strings.ToUpper(params.Name) + "%"
		getDB = getDB.Where("UPPER(name) LIKE ?", sqlName)
	}

	if params.Ipv4 != "" {
		sqlIpv4 := "%" + params.Ipv4 + "%"
		getDB = getDB.Where("ipv4 LIKE ?", sqlIpv4)
	}

	if params.Ipv6 != "" {
		sqlIpv6 := "%" + params.Ipv6 + "%"
		getDB = getDB.Where("ipv6 LIKE ?", sqlIpv6)
	}

	if params.Vip != "" {
		sqlVip := "%" + params.Vip + "%"
		getDB = getDB.Where("vip LIKE ?", sqlVip)
	}

	if params.Zone != "" {
		sqlZone := "%" + strings.ToUpper(params.Zone) + "%"
		getDB = getDB.Where("UPPER(zone) LIKE ?", sqlZone)
	}

	if params.Cloud != "" {
		sqlCloud := "%" + strings.ToUpper(params.Cloud) + "%"
		getDB = getDB.Where("UPPER(cloud) LIKE ?", sqlCloud)
	}

	if params.System != "" {
		sqlSystem := "%" + strings.ToUpper(params.System) + "%"
		getDB = getDB.Where("UPPER(system) LIKE ?", sqlSystem)
	}

	if params.ProjectName != "" {
		sqlProjectName := "%" + strings.ToUpper(params.ProjectName) + "%"
		var projectId uint
		if err = model.DB.Model(model.Project{}).Where("UPPER(name) LIKE ?", sqlProjectName).Pluck("id", &projectId).Error; err != nil {
			return nil, fmt.Errorf("查询项目ID失败: %v", err)
		}
		getDB = getDB.Where("project_id = ?", projectId)
	}

	if err = getDB.Count(&count).Error; err != nil {
		return nil, fmt.Errorf("查询服务器总数失败: %v", err)

	}
	if params.Page != 0 && params.PageSize != 0 {
		if err = getDB.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize).Find(&hosts).Error; err != nil {
			return nil, fmt.Errorf("查询服务器失败: %v", err)
		}
	} else {
		if err = getDB.Find(&hosts).Error; err != nil {
			return nil, fmt.Errorf("查询服务器失败: %v", err)
		}
	}
	var res *[]api.GetHostRes
	var result api.GetHostsRes
	res, err = s.GetResults(&hosts)
	if err != nil {
		return nil, err
	}
	result = api.GetHostsRes{
		Records:  *res,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    count,
	}
	return &result, err
}

func (s *HostService) DeleteHosts(ids []uint) (err error) {
	tx := model.DB.Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()
	if err = tx.Where("id IN (?)", ids).Delete(&model.Host{}).Error; err != nil {
		return fmt.Errorf("删除服务器失败 %d: %v", ids, err)
	}
	tx.Commit()
	return nil
}

func (s *HostService) GetResults(hostObj any) (*[]api.GetHostRes, error) {
	var result []api.GetHostRes
	var err error
	if hosts, ok := hostObj.(*[]model.Host); ok {
		for _, host := range *hosts {
			res := api.GetHostRes{
				ID:       host.ID,
				Name:     host.Name,
				Ipv4:     host.Ipv4,
				Vip:      host.Vip,
				Zone:     host.Zone,
				Cloud:    host.Cloud,
				System:   host.System,
				Cores:    host.Cores,
				DataDisk: host.DataDisk,
				Mem:      host.Mem,
			}
			if host.Ipv6 != nil {
				res.Ipv6 = *host.Ipv6
			}
			if err = model.DB.Model(model.Project{}).Where("id = ?", host.ProjectId).Pluck("name", &res.ProjectName).Error; err != nil {
				return nil, fmt.Errorf("查询项目名称失败: %v", err)
			}
			result = append(result, res)
		}
		return &result, err
	}
	if host, ok := hostObj.(*model.Host); ok {
		res := api.GetHostRes{
			ID:       host.ID,
			Name:     host.Name,
			Ipv4:     host.Ipv4,
			Vip:      host.Vip,
			Zone:     host.Zone,
			Cloud:    host.Cloud,
			System:   host.System,
			Cores:    host.Cores,
			DataDisk: host.DataDisk,
			Mem:      host.Mem,
		}
		if host.Ipv6 != nil {
			res.Ipv6 = *host.Ipv6
		}
		if err = model.DB.Model(model.Project{}).Where("id = ?", host.ProjectId).Pluck("name", &res.ProjectName).Error; err != nil {
			return nil, fmt.Errorf("查询项目名称失败: %v", err)
		}
		result = append(result, res)
		return &result, err
	}
	return &result, errors.New("转换服务器结果失败")
}