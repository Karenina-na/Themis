package token

import (
	"Themis/src/config"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// Claims
// @Description: 令牌声明
type Claims struct {
	Username           string `json:"username"`
	Password           string `json:"password"`
	jwt.StandardClaims        // jwt中标准格式,主要是设置token的过期时间
}

// NewClaims
//
//	@Description: 创建令牌声明
//	@return *Claims	令牌声明
func NewClaims(username string, password string) *Claims {
	expireTime := time.Now().Add(time.Duration(config.Root.TokenExpireTime) * time.Second)
	issuer := "Themis"
	return &Claims{
		Username: username,
		Password: password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    issuer,
		},
	}
}
