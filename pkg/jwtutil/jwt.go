package jwtutil

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

var defaultJWTUtil *JWTUtil

func Init(secret string) {
	defaultJWTUtil = NewJWTUtil(secret)
}

func Issue(claims jwt.Claims) (string, error) {
	return defaultJWTUtil.issue(claims)
}

func Verify(token string, claims jwt.Claims) (*jwt.Token, error) {
	return defaultJWTUtil.parse(token, claims)
}

type JWTUtil struct {
	secret string
}

func NewJWTUtil(secret string) *JWTUtil {
	return &JWTUtil{secret: secret}
}

func (p *JWTUtil) issue(claims jwt.Claims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString([]byte(p.secret))
}

func (p *JWTUtil) parse(token string, claims jwt.Claims) (*jwt.Token, error) {
	validation := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(p.secret), nil
	}

	return jwt.ParseWithClaims(token, claims, validation)
}
