package util

import (
	"FreeOps/internal/consts"
	"FreeOps/internal/model"
	"FreeOps/pkg/api"
	"FreeOps/pkg/logger"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"time"
)

func IsDir(path string) bool {
	dirStat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return dirStat.IsDir()
}

// interface切片是否包含指定值
func IsSliceContain(slice interface{}, value interface{}) bool {
	sliceValue := reflect.ValueOf(slice)
	if sliceValue.Kind() != reflect.Slice {
		return false
	}
	for i := 0; i < sliceValue.Len(); i++ {
		item := sliceValue.Index(i).Interface()
		if reflect.DeepEqual(item, value) {
			return true
		}
	}

	return false
}

// uint切片是否包含指定的值
func IsUintSliceContain(slice []uint, value uint) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// 字符串切片去重
func UniqueStrings(stringSlice []string) (result []string) {
	unique := make(map[string]bool)
	for _, v := range stringSlice {
		if !unique[v] {
			unique[v] = true
			result = append(result, v)
		}
	}
	return result
}

// uint切片去重
func UniqueUint(uintSlice []uint) (result []uint) {
	unique := make(map[uint]bool)
	for _, v := range uintSlice {
		if !unique[v] {
			unique[v] = true
			result = append(result, v)
		}
	}
	return result
}

// uint切片转string
func UintSliceToString(uintSlice []uint) (result string, err error) {
	bytes, err := json.Marshal(uintSlice)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// string切片转uint
func StringToUintSlice(str string) (result []uint, err error) {
	err = json.Unmarshal([]byte(str), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 绑定参数错误的返回
func BindErrorResponse(err error) api.Response {
	var errMessage string
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		errMessage = "JSON类型不匹配"
	} else {
		errMessage = "参数错误"
	}

	return api.Response{
		Code: consts.SERVICE_ERROR_CODE,
		Msg:  errMessage,
		Data: map[string]string{
			"error": err.Error(),
		},
	}
}

// 后端服务错误的返回
func ServerErrorResponse(msg string, err error) api.Response {
	return api.Response{
		Code: consts.SERVICE_ERROR_CODE,
		Msg:  msg,
		Data: map[string]string{
			"error": err.Error(),
		},
	}
}

// 匹配手机号
func CheckMobile(phone string) bool {
	reg := `^1(3\d{2}|4[14-9]\d|5([0-35689]\d|7[1-79])|66\d|7[2-35-8]\d|8\d{2}|9[13589]\d)\d{7}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(phone)
}

// 匹配电子邮箱
func CheckEmail(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// RandStringRunes 返回随机字符串
func RandStringRunes(n int) string {
	var letterRunes = []byte("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	// go1.20开始不再需要显式的调用seed
	// rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// 删除uint切片中的指定元素
func DeleteUintSlice(s []uint, i uint) []uint {
	j := 0
	for _, v := range s {
		if v != i {
			s[j] = v
			j++
		}
	}
	return s[:j]
}

// 通过指针删除uint切片中的指定元素
func DeleteUintSliceByPtr(s *[]uint, id uint) {
	j := 0
	for _, v := range *s {
		if v != id {
			(*s)[j] = v
			j++
		}
	}
	*s = (*s)[:j]
}

func IsSelf(c *gin.Context, userId uint) (bool, error) {
	user, err := GetClaimsUser(c)
	if err != nil {
		return false, err
	}
	return user.ID == userId, nil
}

// websocket不可用该函数
func IsSelfAdmin(c *gin.Context) (bool, error) {
	getAdmin, isExist := c.Get("isAdmin")
	if !isExist {
		return false, errors.New("未获取到isAdmin的中间件预设变量")
	}
	isAdmin, ok := getAdmin.(bool)
	if !ok {
		return false, errors.New("中间件预设变量isAdmin不合法")
	}
	return isAdmin, nil
}

func GetClaimsUser(c *gin.Context) (*model.User, error) {
	getUser, isExist := c.Get("user")
	if !isExist {
		return nil, errors.New("未获取到Token中claims的User对象")
	}
	user, ok := getUser.(*model.User)
	if !ok {
		return nil, errors.New("Token中claims的User对象不合法")
	}
	return user, nil
}

func GetClaimsRole(c *gin.Context) (*[]model.Role, error) {
	getRoles, isExist := c.Get("roles")
	if !isExist {
		return nil, errors.New("未获取到Token中claims的Roles对象")
	}
	roles, ok := getRoles.(*[]model.Role)
	if !ok {
		return nil, errors.New("未获取到Token中claims的Roles对象")
	}
	return roles, nil
}

// 查询并返回存在于指定表中的ID切片
func ExtractExistIDs(tableName string, ids []uint) ([]uint, error) {
	var existIDs []uint
	result := model.DB.Table(tableName).Where("id IN (?)", ids).Pluck("id", &existIDs)
	if result.Error != nil {
		return nil, result.Error
	}
	return existIDs, nil
}

// 返回两个uint切片的差集（第一个切片中存在而第二个切片中不存在的元素）
func FindDifference(original, exist []uint) []uint {
	existMap := make(map[uint]bool)
	for _, id := range exist {
		existMap[id] = true
	}

	var difference []uint
	for _, id := range original {
		if _, exists := existMap[id]; !exists {
			difference = append(difference, id)
		}
	}
	return difference
}

// 查找指定ID切片中不存在的ID，使用ExtractExistIDs和FindDifference重构
func FindNotExistIDs(tableName string, ids []uint) ([]uint, error) {
	existIDs, err := ExtractExistIDs(tableName, ids)
	if err != nil {
		return nil, err
	}

	return FindDifference(ids, existIDs), nil
}

func validateToken(token string) (*model.User, *[]model.Role, bool) {
	// 判断是否在黑名单
	err := model.DB.Where("jwt = ?", token).First(&model.JwtBlacklist{}).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, false
	}
	// 判断token是否已过期
	var claims *CustomClaims
	claims, err = ParseToken(token)
	if err != nil {
		return nil, nil, false
	}

	// 判断用户是否被禁用
	if claims.User.Status != consts.UserModelStatusEnabled {
		return nil, nil, false
	}
	// 判断用户的updateAt是否在token签发之后
	var updateAt time.Time
	if err = model.DB.Model(&model.User{}).Where("id = ?", claims.User.ID).Select("updated_at").Scan(&updateAt).Error; err != nil || updateAt != claims.User.UpdatedAt {
		return nil, nil, false
	}
	return &claims.User, &claims.Roles, true
}

// UpgraderWebsocket
func UpgraderWebSocket(c *gin.Context, isAuth bool) (*websocket.Conn, *model.User, *[]model.Role, error) {
	var (
		conn *websocket.Conn
		err  error
	)
	var upgrader = websocket.Upgrader{
		HandshakeTimeout: consts.WebSocketHandshakeTimeout,
		ReadBufferSize:   consts.WebSocketReadBufferSize,
		WriteBufferSize:  consts.WebSocketWriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
			w.WriteHeader(status)
			w.Write([]byte("WebSocket upgrade failed: " + reason.Error()))
		},
	}

	if conn, err = upgrader.Upgrade(c.Writer, c.Request, nil); err != nil {
		return nil, nil, nil, fmt.Errorf("websocket连接失败: %v", err)
	}
	logger.Log().Info("common", "websocket连接成功")
	if isAuth {
		_, message, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			return nil, nil, nil, fmt.Errorf("读取认证消息失败: %v", err)
		}
		var authMsg struct {
			Type  string `json:"type"`
			Token string `json:"token"`
		}
		if err = json.Unmarshal(message, &authMsg); err != nil || authMsg.Type != "auth" {
			conn.Close()
			return nil, nil, nil, fmt.Errorf("无效的认证消息")
		}
		user, roles, ok := validateToken(authMsg.Token)
		if !ok {
			conn.WriteMessage(websocket.TextMessage, []byte(`{"type": "auth_result", "success": false}`))
			conn.Close()
			return nil, nil, nil, fmt.Errorf("token认证失败")
		}

		conn.WriteMessage(websocket.TextMessage, []byte(`{"type": "auth_result", "success": true}`))
		logger.Log().Info("common", "websocket auth认证通过")
		return conn, user, roles, err
	}
	return conn, nil, nil, err
}
