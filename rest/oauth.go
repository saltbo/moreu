package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/saltbo/gopkg/ginutil"

	"github.com/saltbo/authcar/pkg/oauth2"
	"github.com/saltbo/authcar/pkg/rolec"
)

type RoleLoader func(uid string) ([]string, error)

type Oauth struct {
	oauth2  *oauth2.Oauth2
	jwtRole *rolec.JWTRole

	roleLoader RoleLoader
}

func NewOauth(oauth2 *oauth2.Oauth2, jwtRole *rolec.JWTRole) *Oauth {
	return &Oauth{
		oauth2:  oauth2,
		jwtRole: jwtRole,
		roleLoader: func(uid string) ([]string, error) {
			return []string{"admin"}, nil
		},
	}
}

func (o *Oauth) SetupRoleLoader(roleLoader RoleLoader) {
	o.roleLoader = roleLoader
}

func (o *Oauth) Register(router *gin.RouterGroup) {
	router.GET("/authorize", o.authorize)
	router.GET("/signin", o.signIn)
}

func (o *Oauth) authorize(c *gin.Context) {
	redirect := fmt.Sprintf("http://%s/oauth/signin?redirect=%s", c.Request.Host, c.Query("redirect"))
	c.Redirect(http.StatusFound, o.oauth2.Authorize(redirect))
}

func (o *Oauth) signIn(c *gin.Context) {
	redirect := c.Query("redirect")
	// 通过SSO获取用户信息
	userInfo, err := o.oauth2.GetUserInfo(c.Query("code"), redirect)
	if err != nil {
		ginutil.JSONBadRequest(c, err)
		return
	}

	userMail, ok := userInfo["mail"]
	if !ok {
		ginutil.JSONServerError(c, fmt.Errorf("not found username from sso system"))
		return
	}

	// 查询用户权限
	userRole, err := o.roleLoader(userMail)
	if err != nil {
		ginutil.JSONServerError(c, err)
		return
	}

	// 签发Token
	token, err := o.jwtRole.Issue(userMail, userRole)
	if err != nil {
		ginutil.JSONServerError(c, err)
		return
	}

	setCookie := func(name, value string) {
		c.SetCookie(name, value, 7*24*3600, "", "", false, false)
	}

	// 设置登录标识
	setCookie("token", token)
	setCookie("username", userInfo["name"])
	setCookie("role-level", strings.Join(userRole, ","))
	c.Redirect(http.StatusFound, redirect)
}
