package service

import (
	"FreeOps/internal/model"
	"FreeOps/pkg/api"
	"errors"
	"fmt"
	"strings"
)

type GameService struct{}

var insGame GameService

func GameServiceApp() *GameService {
	return &insGame
}

// 修改/添加游戏服
func (s *GameService) UpdateGame(params *api.UpdateGameReq) (err error) {
	var (
		game  model.Game
		count int64
	)
	if params.ID != 0 {
		if err = model.DB.Model(&model.Game{}).Where("id = ?", params.ID).Count(&count).Error; count != 1 || err != nil {
			return fmt.Errorf("game ID不存在: %d, 或有错误信息: %v", params.ID, err)
		}

		if err = model.DB.Model(&game).Where("(name = ? AND id != ?) OR (lb_listener_port = ? AND lb_name = ? AND id != ?)",
			params.Name, params.ID, params.LbListenerPort, params.LbName, params.ID).Count(&count).Error; err != nil {
			return fmt.Errorf("查询游戏服失败: %v", err)
		} else if count > 0 {
			return fmt.Errorf("游戏服名已被使用: %s", params.Name)
		}

		if err := model.DB.Where("id = ?", params.ID).First(&game).Error; err != nil {
			return fmt.Errorf("游戏服查询失败: %v", err)
		}
		game.Name = params.Name
		game.Type = params.Type
		game.Status = params.Status
		game.LbName = params.LbName
		game.LbListenerPort = params.LbListenerPort
		game.ServerPort = params.ServerPort
		game.ProjectId = params.ProjectId
		game.HostId = params.HostId
		if params.CrossId != 0 {
			*game.CrossId = params.CrossId
		}
		if params.GlobalId != 0 {
			*game.GlobalId = params.GlobalId
		}

		if err = model.DB.Save(&game).Error; err != nil {
			return fmt.Errorf("数据保存失败: %v", err)
		}
		return err
	} else {
		err = model.DB.Model(&game).Where("name = ? OR (lb_listener_port = ? AND lb_name = ?)", params.Name, params.LbListenerPort, params.LbName).Count(&count).Error
		if err != nil {
			return fmt.Errorf("查询游戏服失败: %v", err)
		} else if count > 0 {
			return fmt.Errorf("游戏服名(%s)已存在", params.Name)
		}
		game = model.Game{
			Name:           params.Name,
			Type:           params.Type,
			Status:         params.Status,
			LbName:         params.LbName,
			LbListenerPort: params.LbListenerPort,
			ServerPort:     params.ServerPort,
			ProjectId:      params.ProjectId,
			HostId:         params.HostId,
		}
		if params.CrossId != 0 {
			*game.CrossId = params.CrossId
		}
		if params.GlobalId != 0 {
			*game.GlobalId = params.GlobalId
		}
		if err = model.DB.Create(&game).Error; err != nil {
			return fmt.Errorf("创建游戏服失败: %v", err)
		}
		return err
	}
}

func (s *GameService) GetGames(params *api.GetGamesReq) (*api.GetGamesRes, error) {
	var games []model.Game
	var err error
	var count int64

	getDB := model.DB.Model(&model.Game{})
	if params.ID != 0 {
		getDB = getDB.Where("id = ?", params.ID)
	}

	if params.Name != "" {
		sqlName := "%" + strings.ToUpper(params.Name) + "%"
		getDB = getDB.Where("UPPER(name) LIKE ?", sqlName)
	}

	if params.Type != 0 {
		getDB = getDB.Where("type = ?", params.Type)
	}

	if params.Status != 0 {
		getDB = getDB.Where("status = ?", params.Status)
	}

	if params.CrossId != 0 {
		getDB = getDB.Where("cross_id = ?", params.CrossId)
	}

	if params.GlobalId != 0 {
		getDB = getDB.Where("global_id = ?", params.GlobalId)
	}

	if params.HostName != "" {
		sqlHostName := "%" + strings.ToUpper(params.HostName) + "%"
		var hostId uint
		if err = model.DB.Model(model.Host{}).Where("UPPER(name) LIKE ?", sqlHostName).Pluck("id", &hostId).Error; err != nil {
			return nil, fmt.Errorf("查询服务器ID失败: %v", err)
		}
		getDB = getDB.Where("host_id = ?", hostId)

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
		return nil, fmt.Errorf("查询游戏服总数失败: %v", err)

	}
	if params.Page != 0 && params.PageSize != 0 {
		if err = getDB.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize).Find(&games).Error; err != nil {
			return nil, fmt.Errorf("查询游戏服失败: %v", err)
		}
	} else {
		if err = getDB.Find(&games).Error; err != nil {
			return nil, fmt.Errorf("查询游戏服失败: %v", err)
		}
	}
	var res *[]api.GetGameRes
	var result api.GetGamesRes
	res, err = s.GetResults(&games)
	if err != nil {
		return nil, err
	}
	result = api.GetGamesRes{
		Records:  *res,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    count,
	}
	return &result, err
}

// 记得看看跨服有没有游服还存在
func (s *GameService) DeleteGames(ids []uint) (err error) {
	tx := model.DB.Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()
	if err = tx.Where("id IN (?)", ids).Delete(&model.Game{}).Error; err != nil {
		return fmt.Errorf("删除游戏服失败 %d: %v", ids, err)
	}
	tx.Commit()
	return nil
}

func (s *GameService) GetResults(gameObj any) (*[]api.GetGameRes, error) {
	var result []api.GetGameRes
	var err error
	if games, ok := gameObj.(*[]model.Game); ok {
		for _, game := range *games {
			res := api.GetGameRes{
				ID:             game.ID,
				Name:           game.Name,
				Type:           game.Type,
				Status:         game.Status,
				LbName:         game.LbName,
				LbListenerPort: game.LbListenerPort,
				ServerPort:     game.ServerPort,
			}
			if game.CrossId != nil {
				res.CrossId = *game.CrossId
			}
			if game.GlobalId != nil {
				res.GlobalId = *game.GlobalId
			}
			if err = model.DB.Model(model.Project{}).Where("id = ?", game.ProjectId).Pluck("name", &res.ProjectName).Error; err != nil {
				return nil, fmt.Errorf("查询项目名称失败: %v", err)
			}
			if err = model.DB.Model(model.Host{}).Where("id = ?", game.HostId).Pluck("name", &res.HostName).Error; err != nil {
				return nil, fmt.Errorf("查询服务器名称失败: %v", err)
			}
			result = append(result, res)
		}
		return &result, err
	}
	if game, ok := gameObj.(*model.Game); ok {
		res := api.GetGameRes{
			ID:             game.ID,
			Name:           game.Name,
			Type:           game.Type,
			Status:         game.Status,
			LbName:         game.LbName,
			LbListenerPort: game.LbListenerPort,
			ServerPort:     game.ServerPort,
		}
		if game.CrossId != nil {
			res.CrossId = *game.CrossId
		}
		if game.GlobalId != nil {
			res.GlobalId = *game.GlobalId
		}
		if err = model.DB.Model(model.Project{}).Where("id = ?", game.ProjectId).Pluck("name", &res.ProjectName).Error; err != nil {
			return nil, fmt.Errorf("查询项目名称失败: %v", err)
		}
		if err = model.DB.Model(model.Host{}).Where("id = ?", game.HostId).Pluck("name", &res.HostName).Error; err != nil {
			return nil, fmt.Errorf("查询服务器名称失败: %v", err)
		}
		result = append(result, res)
		return &result, err
	}
	return &result, errors.New("转换游戏服结果失败")
}
