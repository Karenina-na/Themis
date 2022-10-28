package token

import (
	"Themis/src/config"
	"Themis/src/exception"
	"Themis/src/util"
	"github.com/dgrijalva/jwt-go"
)

// GenerateToken
//
//	@Description: 生成令牌
//	@param account	账号
//	@param password	密码
//	@return s	令牌
//	@return E		错误
func GenerateToken(account string, password string) (s string, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("token-GetToken", util.Strval(r))
		}
	}()
	claims := NewClaims(account, password)
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.Root.TokenSignKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

// ParseToken
//
//	@Description:	解析令牌
//	@param token	令牌
//	@return account	账号
//	@return password	密码
//	@return E	错误
func ParseToken(token string) (account string, password string, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("token-ParseToken", util.Strval(r))
		}
	}()
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Root.TokenSignKey), nil
	})
	if err != nil {
		return "", "", err
	}
	if tokenClaims != nil {
		if c, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return c.Username, c.Password, nil
		}
	}
	return "", "", err
}
