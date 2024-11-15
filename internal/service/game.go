package service

import (
	"FreeOps/internal/consts"
	"FreeOps/internal/model"
	"FreeOps/pkg/api"
	"FreeOps/pkg/util"
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
	// 判断projectId是否和hostId匹配
	if err = model.DB.Model(&model.Host{}).Where("id = ? AND project_id = ?", params.HostId, params.ProjectId).Count(&count).Error; count != 1 || err != nil {
		return fmt.Errorf("host ID不存在: %d, 或不在项目ID: %d 下, 或有错误信息: %v", params.HostId, params.ProjectId, err)
	}
	if params.ActionType == consts.ActionTypeIsCreate {
		// 只需要检查以下，避免不同机器创建了相同的。lb、sPort无需检查，因为必然是脚本创建成功了才会请求本接口。
		if err = model.DB.Model(&model.Game{}).Where("id = ? AND project_id = ? AND type = ?", params.Id, params.ProjectId, params.Type).Count(&count).Error; count > 0 || err != nil {
			return fmt.Errorf("game 已存在,id: %d projectId: %d type: %d , 或有错误信息: %v", params.Id, params.ProjectId, params.Type, err)
		}
		game = model.Game{
			Id:         params.Id,
			Type:       params.Type,
			Status:     params.Status,
			ServerPort: params.ServerPort,
			ProjectId:  params.ProjectId,
			HostId:     params.HostId,
		}
		if params.CrossId != 0 {
			if game.CrossId == nil {
				game.CrossId = new(uint)
			}
			*game.CrossId = params.CrossId
		}
		if params.CommonId != 0 {
			if game.CommonId == nil {
				game.CommonId = new(uint)
			}
			*game.CommonId = params.CommonId
		}
		if params.LbName != "" {
			if game.LbName == nil {
				game.LbName = new(string)
			}
			*game.LbName = params.LbName
		}
		if params.LbListenerPort != 0 {
			if game.LbListenerPort == nil {
				game.LbListenerPort = new(uint)
			}
			*game.LbListenerPort = params.LbListenerPort
		}
		if err = model.DB.Create(&game).Error; err != nil {
			return fmt.Errorf("创建游戏服记录失败: %v", err)
		}
		return err
	} else if params.ActionType == consts.ActionTypeIsUpdate {
		if err = model.DB.Model(&model.Game{}).Where("id = ? AND project_id = ? AND type = ?", params.Id, params.ProjectId, params.Type).Count(&count).Error; count != 1 || err != nil {
			return fmt.Errorf("game 记录数!=1,id: %d projectId: %d type: %d , 或有错误信息: %v", params.Id, params.ProjectId, params.Type, err)
		}
		game.Id = params.Id
		game.Type = params.Type
		game.Status = params.Status
		game.ServerPort = params.ServerPort
		game.ProjectId = params.ProjectId
		game.HostId = params.HostId
		if params.CrossId != 0 {
			if game.CrossId == nil {
				game.CrossId = new(uint)
			}
			*game.CrossId = params.CrossId
		}
		if params.CommonId != 0 {
			if game.CommonId == nil {
				game.CommonId = new(uint)
			}
			*game.CommonId = params.CommonId
		}
		if params.LbName != "" {
			if game.LbName == nil {
				game.LbName = new(string)
			}
			*game.LbName = params.LbName
		}
		if params.LbListenerPort != 0 {
			if game.LbListenerPort == nil {
				game.LbListenerPort = new(uint)
			}
			*game.LbListenerPort = params.LbListenerPort
		}

		if err = model.DB.Save(&game).Error; err != nil {
			return fmt.Errorf("数据保存失败: %v", err)
		}
		return err
	} else {
		return fmt.Errorf("未知操作类型: %d", params.ActionType)
	}
}

func (s *GameService) UpdateGameStatus(params *api.UpdateGameStatusReq) (err error) {
	var game model.Game
	if err = model.DB.Model(&model.Game{}).Where("id = ?", params.Id).First(&game).Error; err != nil {
		return fmt.Errorf("查询游戏服失败: %v", err)
	}
	game.Status = params.Status
	if err = model.DB.Save(&game).Error; err != nil {
		return fmt.Errorf("更新游戏服状态失败: %v", err)
	}
	return nil
}

func (s *GameService) GetGames(params *api.GetGamesReq, bindProjectIds []uint) (*api.GetGamesRes, error) {
	var games []model.Game
	var err error
	var count int64

	getDB := model.DB.Model(&model.Game{})
	if params.Id != 0 {
		getDB = getDB.Where("id = ?", params.Id)
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

	if params.CommonId != 0 {
		getDB = getDB.Where("common_id = ?", params.CommonId)
	}

	if params.HostName != "" {
		sqlHostName := "%" + strings.ToUpper(params.HostName) + "%"
		var hostId []uint
		if err = model.DB.Model(model.Host{}).Where("UPPER(name) LIKE ?", sqlHostName).Pluck("id", &hostId).Error; err != nil {
			return nil, fmt.Errorf("查询服务器ID失败: %v", err)
		}
		getDB = getDB.Where("host_id IN (?)", hostId)
	}

	if params.Ipv4 != "" {
		sqlIpv4 := "%" + strings.ToUpper(params.Ipv4) + "%"
		var hostId []uint
		if err = model.DB.Model(model.Host{}).Where("UPPER(ipv4) LIKE ?", sqlIpv4).Pluck("id", &hostId).Error; err != nil {
			return nil, fmt.Errorf("查询服务器ID失败: %v", err)
		}
		getDB = getDB.Where("host_id IN (?)", hostId)

	}

	if params.ProjectId != 0 {
		if !util.IsUintSliceContain(bindProjectIds, params.ProjectId) {
			return nil, errors.New("用户无权限查看该项目的游戏信息")
		}
		getDB = getDB.Where("project_id = ?", params.ProjectId)
	} else {
		getDB = getDB.Where("project_id IN (?)", bindProjectIds)
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

func (s *GameService) DeleteGames(ids []uint) (err error) {
	var games []model.Game
	if err = model.DB.Model(&model.Game{}).Where("id IN (?)", ids).Find(&games).Error; err != nil {
		return fmt.Errorf("查询游戏服失败: %v", err)
	}
	for _, game := range games {
		// 如果类型不是游服，则判断下面还有没有未合服的游服
		if game.Type != consts.GameModeTypeIsGame {
			var count int64
			switch game.Type {
			case consts.GameModelTypeIsCross:
				if err = model.DB.Model(&model.Game{}).Where("cross_id = ?", game.Id).Count(&count).Error; err != nil {
					return fmt.Errorf("查询跨服下游服失败: %v", err)
				}
				if count > 0 {
					return fmt.Errorf("跨服下还有游服存在: %d", game.Id)
				}
			case consts.GameModelTypeIsCommon:
				if err = model.DB.Model(&model.Game{}).Where("common_id = ?", game.Id).Count(&count).Error; err != nil {
					return fmt.Errorf("查询公共服下游服失败: %v", err)
				}
				if count > 0 {
					return fmt.Errorf("公共服下还有游服存在: %d", game.Id)
				}
			default:
				return fmt.Errorf("未知游服类型: %d", game.Type)
			}
		}
	}
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
	type tmpDataStruct struct {
		Name string
		Ipv4 string
	}

	var result []api.GetGameRes
	var err error
	var tmpData tmpDataStruct
	if games, ok := gameObj.(*[]model.Game); ok {
		for _, game := range *games {
			res := api.GetGameRes{
				Id:         game.Id,
				Type:       game.Type,
				Status:     game.Status,
				ServerPort: game.ServerPort,
				ProjectId:  game.ProjectId,
				HostId:     game.HostId,
			}
			if game.CrossId != nil {
				res.CrossId = *game.CrossId
			}
			if game.CommonId != nil {
				res.CommonId = *game.CommonId
			}
			if game.LbName != nil {
				res.LbName = *game.LbName
			}
			if game.LbListenerPort != nil {
				res.LbListenerPort = *game.LbListenerPort
			}
			if err = model.DB.Model(model.Project{}).Where("id = ?", game.ProjectId).Pluck("name", &res.ProjectName).Error; err != nil {
				return nil, fmt.Errorf("查询项目名称失败: %v", err)
			}
			if err = model.DB.Model(model.Host{}).Where("id = ?", game.HostId).Select("name", "ipv4").Scan(&tmpData).Error; err != nil {
				return nil, fmt.Errorf("查询服务器名称IP失败: %v", err)
			}
			res.Ipv4 = tmpData.Ipv4
			res.HostName = tmpData.Name
			result = append(result, res)
		}
		return &result, err
	}
	if game, ok := gameObj.(*model.Game); ok {
		res := api.GetGameRes{
			Id:         game.Id,
			Type:       game.Type,
			Status:     game.Status,
			ServerPort: game.ServerPort,
			ProjectId:  game.ProjectId,
			HostId:     game.HostId,
		}
		if game.CrossId != nil {
			res.CrossId = *game.CrossId
		}
		if game.CommonId != nil {
			res.CommonId = *game.CommonId
		}
		if game.LbName != nil {
			res.LbName = *game.LbName
		}
		if game.LbListenerPort != nil {
			res.LbListenerPort = *game.LbListenerPort
		}
		if err = model.DB.Model(model.Project{}).Where("id = ?", game.ProjectId).Pluck("name", &res.ProjectName).Error; err != nil {
			return nil, fmt.Errorf("查询项目名称失败: %v", err)
		}
		if err = model.DB.Model(model.Host{}).Where("id = ?", game.HostId).Select("name", "ipv4").Scan(&tmpData).Error; err != nil {
			return nil, fmt.Errorf("查询服务器名称IP失败: %v", err)
		}
		res.Ipv4 = tmpData.Ipv4
		res.HostName = tmpData.Name
		result = append(result, res)
		return &result, err
	}
	return &result, errors.New("转换游戏服结果失败")
}
