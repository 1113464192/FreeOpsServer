package util

import (
	"FreeOps/internal/consts"
	"FreeOps/internal/model"
	"FreeOps/pkg/api"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"math/rand"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
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
func UintSliceToString(uintSlice []uint) (result string) {
	// 如[]uint{1, 2, 3} -> "1, 2, 3"
	for i, v := range uintSlice {
		if i == 0 {
			result += string(v)
		} else {
			result += ", " + string(v)
		}
	}
	return result
}

// string切片转uint
func StringToUintSlice(str string) (result []uint, err error) {
	// 如"1, 2, 3" -> []uint{1, 2, 3}
	strSlice := strings.Split(str, ",")
	for _, s := range strSlice {
		s = strings.TrimSpace(s)
		num, convErr := strconv.ParseUint(s, 10, 32)
		if convErr != nil {
			return nil, convErr
		}
		result = append(result, uint(num))
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
