package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/saltbo/gopkg/jwtutil"
)

func TokenCreate(uid int64, ttl int, roles ...string) (string, error) {
	return jwtutil.Issue(newRoleClaims(uid, ttl, roles))
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

func newRoleClaims(subject int64, ttl int, roles []string) *roleClaims {
	timeNow := time.Now()
	return &roleClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "moreu",
			Audience:  "moreuUsers",
			ExpiresAt: timeNow.Add(time.Duration(ttl) * time.Second).Unix(),
			IssuedAt:  timeNow.Unix(),
			NotBefore: timeNow.Unix(),
			Subject:   strconv.FormatInt(subject, 10),
		},
		Roles: roles,
	}
}
