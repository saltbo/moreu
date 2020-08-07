package rolec

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/storyicon/grbac"
)

type RoleClaims struct {
	jwt.StandardClaims

	Roles []string `json:"roles"`
}

type JWTRole struct {
	rbac *grbac.Controller

	secret string
}

func NewJWTRole(secret string, roleLoader grbac.ControllerOption) (*JWTRole, error) {
	rbac, err := grbac.New(roleLoader)
	if err != nil {
		return nil, err
	}

	return &JWTRole{
		rbac:   rbac,
		secret: secret,
	}, nil
}

func (p *JWTRole) Issue(username string, roles []string) (string, error) {
	timeNow := time.Now()
	claims := &RoleClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "GWProtal",
			Audience:  "GWProtalAPI",
			ExpiresAt: timeNow.Add(7 * 24 * 3600 * time.Second).Unix(),
			IssuedAt:  timeNow.Unix(),
			NotBefore: timeNow.Unix(),
			Subject:   username,
		},
		Roles: roles,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString([]byte(p.secret))
}

func (p *JWTRole) Verify(tokenStr string, req *http.Request) error {
	token, err := jwt.ParseWithClaims(tokenStr, &RoleClaims{}, p.validation)
	if err != nil {
		return fmt.Errorf("token valid failed: %s", err)
	}

	rc := token.Claims.(*RoleClaims)
	state, err := p.rbac.IsRequestGranted(req, rc.Roles)
	if err != nil {
		return err
	}

	if !state.IsGranted() {
		return fmt.Errorf("您没有权限进行此操作，请联系管理员.")
	}

	req.Header.Set("X-Auth-Sub", rc.Subject)
	return nil
}

func (p *JWTRole) validation(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	return []byte(p.secret), nil
}
