package service

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var defaultJwt = NewToken("2333")

type Token struct {
	secret string
}

func NewToken(secret string) *Token {
	return &Token{secret: secret}
}

func (p *Token) create(claims jwt.Claims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString([]byte(p.secret))
}

func (p *Token) validation(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	return []byte(p.secret), nil
}

func TokenCreate(subject string, ttl int, roles ...string) (string, error) {
	return defaultJwt.create(newRoleClaims(subject, ttl, roles))
}

func TokenVerify(tokenStr string) (*roleClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &roleClaims{}, defaultJwt.validation)
	if err != nil {
		return nil, fmt.Errorf("token valid failed: %s", err)
	}

	return token.Claims.(*roleClaims), nil
}

type roleClaims struct {
	jwt.StandardClaims

	Roles []string `json:"roles"`
}

func newRoleClaims(subject string, ttl int, roles []string) *roleClaims {
	timeNow := time.Now()
	return &roleClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "moreu",
			Audience:  "moreuUsers",
			ExpiresAt: timeNow.Add(time.Duration(ttl) * time.Second).Unix(),
			IssuedAt:  timeNow.Unix(),
			NotBefore: timeNow.Unix(),
			Subject:   subject,
		},
		Roles: roles,
	}
}
