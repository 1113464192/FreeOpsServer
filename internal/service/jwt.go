package service

import (
	"FreeOps/internal/model"
	"errors"

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
	err = model.DB.Create(jwtList).Error
	return
}

func (s *JwtService) IsBlacklist(jwt string) bool {
	err := model.DB.Where("jwt = ?", jwt).First(&model.JwtBlacklist{}).Error
	isNotFound := errors.Is(err, gorm.ErrRecordNotFound)
	return !isNotFound
}
