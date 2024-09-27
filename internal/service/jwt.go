package service

import (
	"FreeOps/internal/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type JwtService struct {
}

var (
	insJwt = JwtService{}
)

func JwtServiceApp() *JwtService {
	return &insJwt
}

func (s *JwtService) JwtAddBlacklist(jwtList *model.JwtBlacklist) (err error) {
	if err = model.DB.Create(jwtList).Error; err != nil {
		return fmt.Errorf("添加jwt到黑名单失败: %v", err)
	}
	return err
}

func (s *JwtService) IsBlacklist(jwt string) bool {
	err := model.DB.Where("jwt = ?", jwt).First(&model.JwtBlacklist{}).Error
	isNotFound := errors.Is(err, gorm.ErrRecordNotFound)
	return !isNotFound
}
