package service

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/saltbo/moreu/pkg/jwtutil"
)

func TokenCreate(username string, ttl int, roles ...string) (string, error) {
	_, exist := UsernameExist(username)
	if !exist {
		return "", fmt.Errorf("user not exist")
	}

	return jwtutil.Issue(newRoleClaims(username, ttl, roles))
}

func TokenVerify(tokenStr string) (*roleClaims, error) {
	token, err := jwtutil.Verify(tokenStr, &roleClaims{})
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
