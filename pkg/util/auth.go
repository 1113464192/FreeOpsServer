package util

import (
	"FreeOps/global"
	"FreeOps/internal/model"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// CustomSecret 用于加盐的字符串
var CustomSecret = []byte(global.Conf.SecurityVars.TokenKey)

type CustomClaims struct {
	// 可根据需要自行添加字段
	User                 model.User
	Roles                []model.Role
	jwt.RegisteredClaims // 内嵌标准的声明
}

// GenToken 生成JWT
func GenToken(user model.User) (string, string, error) {
	tokenDuration, err := time.ParseDuration(global.Conf.SecurityVars.TokenExpireDuration)
	refreshTokenDuration, err := time.ParseDuration(global.Conf.SecurityVars.RefreshTokenExpireDuration)

	if err != nil {
		return "", "", fmt.Errorf("生成Token过期时间失败: %v", err)
	}
	var roles []model.Role
	if err := model.DB.Joins("JOIN user_role on user_role.role_id = role.id").
		Where("user_role.user_id = ?", user.ID).
		Find(&roles).Error; err != nil {
		return "", "", fmt.Errorf("查询用户角色失败: %v", err)
	}
	// 创建一个我们自己的声明
	claims := CustomClaims{
		User:  user,
		Roles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenDuration)),
			Issuer:    global.Conf.SecurityVars.JwtIssuer, // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(CustomSecret)
	if err != nil {
		return "", "", fmt.Errorf("生成Token失败: %v", err)
	}

	refreshClaims := CustomClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenDuration)),
			Issuer:    global.Conf.SecurityVars.JwtIssuer, // 签发人
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	refreshTokenString, err := refreshToken.SignedString(CustomSecret)
	if err != nil {
		return "", "", fmt.Errorf("生成RefreshToken失败: %v", err)
	}
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return tokenString, refreshTokenString, nil
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析token
	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法，接受三个参数：要解析的JWT字符串、一个空的CustomClaims实例（用于存储解析后的payload）、一个函数（用于返回签名密钥）。
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i any, err error) {
		// 直接使用标准的Claim则可以直接使用Parse方法
		return CustomSecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 对token对象中的Claim进行类型断言
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid { // 校验token
		if claims.Issuer != global.Conf.SecurityVars.JwtIssuer {
			return nil, errors.New("invalid Issuer")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// StringCost 字符串加密难度
var StringCost = 12

//@function: GenerateFromPassword
//@description: 字符串加密
//@param: str string
//@return: string

func GenerateFromPassword(str string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), StringCost)
	return string(bytes), err
}

// CheckAdminPassword 校验密码
func CheckPassword(userPwd, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userPwd), []byte(password))
	return err == nil
}

func Pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func Unpad(data []byte) []byte {
	padding := int(data[len(data)-1])
	return data[:len(data)-padding]
}

func EncryptAESCBC(plaintext, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plaintext = Pad(plaintext, block.BlockSize())
	ciphertext := make([]byte, len(plaintext))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	return ciphertext, nil
}

func DecryptAESCBC(ciphertext, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, len(ciphertext))

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)
	plaintext = Unpad(plaintext)

	return plaintext, nil
}
