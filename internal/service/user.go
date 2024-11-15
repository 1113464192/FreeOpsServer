package service

import (
	"FreeOps/global"
	"FreeOps/internal/consts"
	"FreeOps/internal/model"
	"FreeOps/pkg/api"
	"FreeOps/pkg/logger"
	"FreeOps/pkg/util"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"strings"
	"time"
)

type UserService struct{}

var insUser UserService

func UserServiceApp() *UserService {
	return &insUser
}

// 登录
func (s *UserService) Login(u *model.User) (res *api.AuthLoginRes, err error) {
	var user model.User
	if err = model.DB.Where("username = ?", u.Username).First(&user).Error; err != nil {
		return res, fmt.Errorf("获取用户对象失败: %v", err)
	}
	if !util.CheckPassword(user.Password, u.Password) {
		return res, errors.New("密码错误")
	}
	if user.Status != consts.UserModelStatusEnabled {
		return res, errors.New("账号状态非Enable")
	}

	// 设置JWT-token
	token, refreshToken, err := util.GenToken(user)
	if err != nil {
		return res, err
	}

	res = &api.AuthLoginRes{
		Token:        token,
		RefreshToken: refreshToken,
	}
	return res, err
}

func (s *UserService) RefreshToken(refreshToken string) (res *api.AuthLoginRes, code string, err error) {
	// 判断refreshToken是否存在于黑名单
	if JwtServiceApp().IsBlacklist(refreshToken) {
		return res, consts.SERVICE_LOGOUT_CODE, errors.New("refreshToken失效")

	}
	claims, err := util.ParseToken(refreshToken)
	if err != nil {
		return res, consts.SERVICE_LOGOUT_CODE, err
	}

	var user model.User
	if err = model.DB.Model(&model.User{}).Where("id = ?", claims.User.ID).First(&user).Error; err != nil {
		return res, consts.SERVICE_LOGOUT_CODE, fmt.Errorf("根据UID: %d 获取用户失败: %v", claims.User.ID, err)
	}

	jwtBlack := &model.JwtBlacklist{Jwt: refreshToken}
	if err = JwtServiceApp().JwtAddBlacklist(jwtBlack); err != nil {
		return res, consts.SERVICE_MODAL_LOGOUT_CODE, err
	}

	token, refreshToken, err := util.GenToken(user)
	if err != nil {
		return res, consts.SERVICE_LOGOUT_CODE, err
	}
	res = &api.AuthLoginRes{
		Token:        token,
		RefreshToken: refreshToken,
	}
	return res, consts.SERVICE_SUCCESS_CODE, nil
}

// 修改/添加用户
func (s *UserService) UpdateUser(params *api.UpdateUserReq) (string, error) {
	var err error
	// 判断电话和邮箱是否正确
	if params.UserPhone != "" && !util.CheckMobile(params.UserPhone) {
		return "", errors.New("电话格式错误")
	}

	if params.UserEmail != "" && !util.CheckEmail(params.UserEmail) {
		return "", errors.New("邮箱格式错误")
	}

	var user model.User
	var count int64
	if params.ID != 0 {
		if err = model.DB.Model(&model.User{}).Where("id = ?", params.ID).Count(&count).Error; count != 1 || err != nil {
			return "", fmt.Errorf("user ID不存在: %d, 或有错误信息: %v", params.ID, err)
		}

		// 判断username是否和现有用户重复
		err = model.DB.Model(&user).Where("username = ? AND id != ?", params.Username, params.ID).Count(&count).Error
		if err != nil {
			return "", fmt.Errorf("查询用户失败: %v", err)
		} else if count > 0 {
			return "", fmt.Errorf("用户名已被使用: %s", params.Username)
		}

		if err = model.DB.Model(&model.User{}).Where("id = ?", params.ID).Update("updated_at", time.Now()).Error; err != nil {
			return "", fmt.Errorf("修改用户更新时间失败: %v", err)
		}

		if err := model.DB.Where("id = ?", params.ID).First(&user).Error; err != nil {
			return "", fmt.Errorf("用户数据库查询失败: %v", err)
		}
		user.Status = params.Status
		user.Username = params.Username
		user.UserGender = params.UserGender
		user.Nickname = params.Nickname
		user.UserPhone = params.UserPhone
		user.UserEmail = params.UserEmail

		if err = model.DB.Save(&user).Error; err != nil {
			return "", fmt.Errorf("数据保存失败: %v", err)
		}

		return "", err
	} else {
		// 判断username是否和现有用户重复
		err = model.DB.Model(&user).Where("username = ?", params.Username).Count(&count).Error
		if err != nil {
			return "", fmt.Errorf("查询用户失败: %v", err)
		} else if count > 0 {
			return "", fmt.Errorf("用户名(%s)已存在", params.Username)
		}

		user = model.User{
			Status:     params.Status,
			Username:   params.Username,
			UserGender: params.UserGender,
			Nickname:   params.Nickname,
			UserPhone:  params.UserPhone,
			UserEmail:  params.UserEmail,
		}
		// 生成初始化密码
		password := util.RandStringRunes(18)
		user.Password, err = util.GenerateFromPassword(password)
		if err != nil {
			return "", fmt.Errorf("用户密码bcrypt加密失败: %v", err)
		}
		if err = model.DB.Create(&user).Error; err != nil {
			return "", fmt.Errorf("创建用户失败: %v", err)
		}
		return password, err
	}
}

// 获取用户关联的角色
func (s *UserService) GetUserRoles(userId uint) (roles []model.Role, err error) {
	if err = model.DB.Joins("JOIN user_role ON user_role.role_id = role.id").
		Where("user_role.user_id = ?", userId).
		Find(&roles).Error; err != nil {
		return nil, fmt.Errorf("查询用户角色失败: %v", err)
	}
	return roles, nil
}

func (s *UserService) GetUsers(params api.GetUsersReq) (result *api.GetUsersRes, err error) {
	var users []model.User
	var total int64
	// 有ID优先ID
	getDB := model.DB.Model(&model.User{})
	if params.ID != 0 {
		getDB = getDB.Where("id = ?", params.ID)
	}
	if params.Status != 0 {
		getDB = getDB.Where("status = ?", params.Status)
	}
	if params.Username != "" {
		// rolename不为空则模糊查询
		sqlUsername := "%" + strings.ToUpper(params.Username) + "%"
		getDB = getDB.Where("UPPER(username) LIKE ?", sqlUsername)
	}
	if params.UserGender != "" {
		getDB = getDB.Where("user_gender = ?", params.UserGender)
	}
	if params.Nickname != "" {
		sqlNickname := "%" + strings.ToUpper(params.Nickname) + "%"
		getDB = getDB.Where("UPPER(nickname) LIKE ?", sqlNickname)
	}
	if params.UserPhone != "" {
		getDB = getDB.Where("user_phone = ?", params.UserPhone)
	}
	if params.UserEmail != "" {
		getDB = getDB.Where("user_email = ?", params.UserEmail)
	}
	// 获取符合上面叠加条件的总数
	if err = getDB.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("查询用户总数失败: %v", err)
	}
	if err = getDB.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	// 过滤结果
	res, err := s.GetResults(&users)
	if err != nil {
		return nil, err
	}
	result = &api.GetUsersRes{
		Records:  *res,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	}
	return result, err
}

func (s *UserService) GetUserPrivilege(user *model.User, roles *[]model.Role) (res *api.GetUserPrivilegeRes, err error) {
	var roleCodes []string
	var roleIds []uint
	for _, role := range *roles {
		roleCodes = append(roleCodes, role.RoleCode)
		roleIds = append(roleIds, role.ID)
	}

	var buttonCodes []string
	if err = model.DB.Model(&model.Button{}).
		Joins("JOIN role_button ON role_button.button_id = button.id").
		Where("role_button.role_id IN (?)", roleIds).
		Pluck("button_code", &buttonCodes).Error; err != nil {
		return res, fmt.Errorf("查询角色按钮失败: %v", err)
	}

	res = &api.GetUserPrivilegeRes{
		UserId:   user.ID,
		Username: user.Username,
		Roles:    roleCodes,
		Buttons:  buttonCodes,
	}
	return res, err
}

// 删除用户ID
func (s *UserService) DeleteUsers(ids []uint) (err error) {
	tx := model.DB.Begin()
	// 删除用户角色关联
	if err = tx.Where("user_id IN (?)", ids).Delete(&model.UserRole{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除用户角色关联失败 %d: %v", ids, err)
	}

	if err = tx.Where("id IN (?)", ids).Delete(&model.User{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除用户失败 %d: %v", ids, err)
	}
	tx.Commit()
	return nil
}

// 修改指定用户密码
func (s *UserService) ChangeUserPassword(params api.ChangeUserPasswordReq) (err error) {
	var (
		user     model.User
		password string
	)
	if err = model.DB.Where("id = ?", params.ID).First(&user).Error; err != nil {
		return fmt.Errorf("查询用户 %d 失败: %v", params.ID, err)
	}

	if password, err = util.GenerateFromPassword(params.NewPassword); err != nil {
		return fmt.Errorf("密码加密失败: %v", err)
	}

	if err = model.DB.Model(&user).Update("password", password).Error; err != nil {
		return fmt.Errorf("更新用户 %d 密码失败: %v", params.ID, err)
	}
	return err
}

// 通过文件更新私钥
func (s *UserService) UpdateKeyFileContext(file *multipart.FileHeader, passphrase string, id uint) error {
	fileP, err := file.Open()
	if err != nil {
		return err
	}
	defer fileP.Close()

	fileBytes, err := io.ReadAll(fileP)
	if err != nil {
		return err
	}

	// AES加密并写入prikey
	var data []byte
	data, err = util.EncryptAESCBC(fileBytes, []byte(global.Conf.SecurityVars.AesKey), []byte(global.Conf.SecurityVars.AesIv))
	if err != nil {
		return fmt.Errorf("用户私钥加密失败: %v", err)
	}
	if err = model.DB.Model(&model.User{}).Where("id = ?", id).Update("pri_key", data).Error; err != nil {
		return fmt.Errorf("私钥写入数据库失败: %v", err)
	}

	// AES加密并写入passphrase
	data, err = util.EncryptAESCBC([]byte(passphrase), []byte(global.Conf.SecurityVars.AesKey), []byte(global.Conf.SecurityVars.AesIv))
	if err != nil {
		return fmt.Errorf("用户passphrase加密失败: %v", err)
	}
	if err = model.DB.Model(&model.User{}).Where("id = ?", id).Update("passphrase", data).Error; err != nil {
		return errors.New("通行证密码写入数据库失败")
	}
	return nil
}

// 通过字符串更新私钥内容
func (s *UserService) UpdateSSHKey(key []byte, passphrase string, id uint) (err error) {
	// AES加密并写入prikey
	var data []byte
	data, err = util.EncryptAESCBC(key, []byte(global.Conf.SecurityVars.AesKey), []byte(global.Conf.SecurityVars.AesIv))
	if err != nil {
		return fmt.Errorf("用户私钥加密失败: %v", err)
	}
	if err = model.DB.Model(&model.User{}).Where("id = ?", id).Update("pri_key", data).Error; err != nil {
		return fmt.Errorf("私钥字符串写入数据库失败: %v", err)
	}

	// AES加密并写入passphrase
	data, err = util.EncryptAESCBC([]byte(passphrase), []byte(global.Conf.SecurityVars.AesKey), []byte(global.Conf.SecurityVars.AesIv))
	if err != nil {
		return fmt.Errorf("用户passphrase加密失败: %v", err)
	}
	if err = model.DB.Model(&model.User{}).Where("id = ?", id).Update("passphrase", data).Error; err != nil {
		return fmt.Errorf("通行证密码写入数据库失败: %v", err)
	}
	return nil
}

func (s *UserService) BindUserRoles(uid uint, roleIds []uint) (err error) {
	// 先传的id是否都存在
	var count int64
	if err = model.DB.Model(&model.User{}).Where("id = ?", uid).Count(&count).Error; count != 1 || err != nil {
		return fmt.Errorf("user 不存在ID: %d, 如果查询用户失败: %v", uid, err)
	}

	if err = model.DB.Model(&model.Role{}).Where("id IN (?)", roleIds).Count(&count).Error; count != int64(len(roleIds)) || err != nil {
		notExistIds, err2 := util.FindNotExistIDs(consts.MysqlTableNameRole, roleIds)
		if err2 != nil {
			return fmt.Errorf("查询角色失败: %v", err2)
		}
		return fmt.Errorf("role 不存在ID: %d, 如果查询角色失败: %v", notExistIds, err)
	}

	tx := model.DB.Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	// 修改uid对应用户的updateAt
	if err = tx.Model(&model.User{}).Where("id = ?", uid).Update("updated_at", time.Now()).Error; err != nil {
		return fmt.Errorf("修改用户更新时间失败: %v", err)
	}

	// 要先清空用户的当前关联
	if err = tx.Where("user_id = ?", uid).Delete(&model.UserRole{}).Error; err != nil {
		return fmt.Errorf("清空用户角色关联失败: %v", err)
	}

	var userRoles []model.UserRole
	for _, id := range roleIds {
		userRole := model.UserRole{
			UserId: uid,
			RoleId: id,
		}
		userRoles = append(userRoles, userRole)
	}
	if err = tx.Create(&userRoles).Error; err != nil {
		return fmt.Errorf("绑定用户角色失败: %v", err)
	}

	tx.Commit()
	return nil
}

// 获取用户对应的项目ID
func (s *UserService) GetUserProjectIDs(c *gin.Context) (bindProjectIds []uint, err error) {
	var (
		roleIds []uint
	)
	// 获取角色对应的项目ID
	if roleIds, err = RoleServiceApp().GetSelfRoleIDs(c); err != nil {
		logger.Log().Error("role", "获取用户的角色IDs失败", err)
		return nil, fmt.Errorf("获取用户的角色IDs失败: %v", err)
	}
	if bindProjectIds, err = RoleServiceApp().GetRoleProjects(roleIds); err != nil {
		logger.Log().Error("role", "获取角色对应的项目ID失败", err)
		return nil, fmt.Errorf("获取角色对应的项目ID失败: %v", err)
	}
	return bindProjectIds, err
}

// 获取用户关联的角色
func (s *UserService) GetUserProjectOptions(c *gin.Context) (*[]api.GetUserProjectOptionsRes, error) {
	var (
		res            []api.GetUserProjectOptionsRes
		projects       []model.Project
		bindProjectIds []uint
		err            error
	)
	if bindProjectIds, err = s.GetUserProjectIDs(c); err != nil {
		return nil, fmt.Errorf("获取用户项目ID失败: %v", err)
	}
	if err = model.DB.Model(&model.Project{}).Where("id IN (?)", bindProjectIds).Select("id", "name").Find(&projects).Error; err != nil {
		return nil, fmt.Errorf("查询项目失败: %v", err)
	}
	for _, project := range projects {
		res = append(res, api.GetUserProjectOptionsRes{
			Label: project.Name,
			Value: project.ID,
		})
	}
	return &res, nil
}

// 返回用户结果
func (s *UserService) GetResults(userObj any) (*[]api.UserRes, error) {
	var result []api.UserRes
	var err error
	if users, ok := userObj.(*[]model.User); ok {
		for _, user := range *users {
			res := api.UserRes{
				ID:         user.ID,
				Status:     user.Status,
				Username:   user.Username,
				UserGender: user.UserGender,
				Nickname:   user.Nickname,
				UserPhone:  user.UserPhone,
				UserEmail:  user.UserEmail,
			}
			result = append(result, res)
		}
		return &result, err
	}
	if user, ok := userObj.(*model.User); ok {
		res := api.UserRes{
			ID:         user.ID,
			Status:     user.Status,
			Username:   user.Username,
			UserGender: user.UserGender,
			Nickname:   user.Nickname,
			UserPhone:  user.UserPhone,
			UserEmail:  user.UserEmail,
		}
		result = append(result, res)
		return &result, err
	}
	return &result, errors.New("转换用户结果失败")
}
