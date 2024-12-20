package service

import (
	"FreeOps/global"
	"FreeOps/internal/consts"
	"FreeOps/internal/model"
	"FreeOps/pkg/util"
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	casbinUtil "github.com/casbin/casbin/v2/util"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"strconv"
	"strings"
	"sync"
)

type CasbinService struct {
}

var (
	InsCasbin = CasbinService{}
)

func CasbinServiceApp() *CasbinService {
	return &InsCasbin
}

// @function: UpdateCasbin
// @description: 更新casbin权限
// @param: api.UpdateCasbinReq
// @return: error
func (s *CasbinService) UpdateCasbin(roleId string, apiIds []uint) (err error) {
	var count int64
	if err = model.DB.Model(&model.Api{}).Where("id IN (?)", apiIds).Count(&count).Error; count != int64(len(apiIds)) || err != nil {
		notExistIds, err2 := util.FindNotExistIDs(consts.MysqlTableNameApi, apiIds)
		if err2 != nil {
			return fmt.Errorf("查询API失败: %v", err2)
		}
		return fmt.Errorf("API 不存在ID: %d, 如果查询API失败: %v", notExistIds, err)
	}

	s.ClearCasbin(0, roleId)
	var apis []model.Api
	if err = model.DB.Where("id IN (?)", apiIds).Find(&apis).Error; err != nil {
		return err
	}
	var rules [][]string
	for _, v := range apis {
		rules = append(rules, []string{roleId, v.Path, v.Method})
	}
	e := s.Casbin()
	if success, _ := e.AddPolicies(rules); !success {
		return errors.New("存在相同api,添加失败,请联系管理员")
	}
	return nil
}

// @function: UpdateCasbinApi
// @description: API更新随动
// @param: oldPath string, newPath string, oldMethod string, newMethod string
// @return: error
func (s *CasbinService) UpdateCasbinApi(oldPath string, newPath string, oldMethod string, newMethod string) error {
	err := model.DB.Model(&model.CasbinRule{}).Where("v1 = ? AND v2 = ?", oldPath, oldMethod).Updates(map[string]any{
		"v1": newPath,
		"v2": newMethod,
	}).Error
	return err
}

// @function: GetPolicyPathByGroupIds
// @description: 获取权限列表
// @param: roleIds []uint
// @return: res []uint, err error
func (s *CasbinService) GetPolicyPathByGroupIds(roleIds []uint) (res []uint, err error) {
	e := s.Casbin()
	var policies [][]string

	for _, roleId := range roleIds {
		// 0是从0开始，全部检索
		list, err := e.GetFilteredPolicy(0, strconv.FormatUint(uint64(roleId), 10))
		if err != nil {
			return nil, err
		}
		policies = append(policies, list...)
	}

	apiIdMap := make(map[uint]struct{})
	// 去重并取出
	for _, policy := range policies {
		var api model.Api
		if err = model.DB.Where("path = ? AND method = ?", policy[1], policy[2]).Select("id").First(&api).Error; err != nil {
			return nil, err
		}
		apiIdMap[api.ID] = struct{}{}
	}

	// 去重
	for apiId := range apiIdMap {
		res = append(res, apiId)
	}

	return res, nil
}

// @function: ClearCasbin
// @description: 清除匹配的权限
// @param: v int, p ...string
// @return: bool
func (s *CasbinService) ClearCasbin(v int, p ...string) bool {
	e := s.Casbin()
	success, _ := e.RemoveFilteredPolicy(v, p...)
	return success

}

// DeleteApiPolicy 删除特定API路径和方法的策略。
// @param path string 要删除的API路径。
// @param method string 要删除的HTTP方法。
// @return error 如果删除失败返回错误，成功则返回nil。
func (s *CasbinService) DeleteApiPolicy(path, method string) error {
	e := s.Casbin() // 获取Casbin执行器实例。
	// 尝试移除匹配特定路径和方法的策略。
	// 第一个参数(0)表示我们不通过角色ID（或任何其他在第一个位置的标准）进行过滤。
	// "path" 和 "method" 是我们在策略中匹配的值。
	success, err := e.RemoveFilteredPolicy(0, "", path, method)
	if err != nil {
		return err // 如果操作失败返回错误。
	}
	if !success {
		return errors.New("未找到匹配的策略或删除失败")
	}
	return nil // 返回nil表示成功。
}

// @function: Casbin
// @description: 持久化到数据库  引入自定义规则
// @return: *casbin.Enforcer
var (
	syncedEnforcer *casbin.SyncedEnforcer
	once           sync.Once
)

func (s *CasbinService) Casbin() *casbin.SyncedEnforcer {
	once.Do(func() {
		a, _ := gormadapter.NewAdapterByDB(model.DB)
		syncedEnforcer, _ = casbin.NewSyncedEnforcer(global.RootPath+"/configs/casbin.conf", a)
		syncedEnforcer.AddFunction("ParamsMatch", s.ParamsMatchFunc)
	})
	_ = syncedEnforcer.LoadPolicy()
	return syncedEnforcer
}

// @function: ParamsMatch
// @description: 自定义规则函数
// @param: fullNameKey1 string, key2 string
// @return: bool
func (s *CasbinService) ParamsMatch(fullNameKey1 string, key2 string) bool {
	key1 := strings.Split(fullNameKey1, "?")[0]
	// 剥离路径后再使用casbin的keyMatch2
	return casbinUtil.KeyMatch2(key1, key2)
}

// @function: ParamsMatchFunc
// @description: 自定义规则函数
// @param: args ...any
// @return: any, error
func (s *CasbinService) ParamsMatchFunc(args ...any) (any, error) {
	name1 := args[0].(string)
	name2 := args[1].(string)

	return s.ParamsMatch(name1, name2), nil
}
