package rest

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/saltbo/gopkg/ginutil"
	"github.com/storyicon/grbac"

	"github.com/saltbo/moreu/client"
	"github.com/saltbo/moreu/service"
)

var defaultRBAC *grbac.Controller

func RBACInit(name string) {
	roleLoader := grbac.WithYAML(name, time.Second)
	rbac, err := grbac.New(roleLoader)
	if err != nil {
		log.Fatalln(err)
	}

	defaultRBAC = rbac
}

func LoginAuth(c *gin.Context) {
	if err := loginAuth(c); err != nil {
		ginutil.JSONUnauthorized(c, err)
		return
	}
}

func StaticAuth(c *gin.Context) {
	if err := loginAuth(c); err != nil {
		c.Redirect(http.StatusFound, service.SignInLink(c.Request.URL.RequestURI()))
		c.Abort()
		return
	}
}

func loginAuth(c *gin.Context) error {
	token, err := tokenCookieGet(c)
	if errors.Is(err, http.ErrNoCookie) {
		return fmt.Errorf("none token")
	}

	rc, err := service.TokenVerify(token)
	if err != nil {
		return err
	}

	userIdSet(c, rc.Subject)
	userRolesSet(c, rc.Roles)
	c.Request.Header.Set(moreu.HeaderUserIdKey, rc.Subject)
	return nil
}

func RoleAuth(c *gin.Context) {
	state, err := defaultRBAC.IsRequestGranted(c.Request, userRolesGet(c))
	if err != nil {
		ginutil.JSONForbidden(c, err)
		return
	}

	if !state.IsGranted() {
		ginutil.JSONForbidden(c, fmt.Errorf("您没有权限进行此操作，请联系管理员"))
		return
	}
}

// auth k-v
const (
	cookieTokenKey = "moreu-token"

	ctxUserIdKey    = "user_id"
	ctxUserRolesKey = "user_roles"
)

func userIdSet(c *gin.Context, userId string) {
	uid, _ := strconv.ParseInt(userId, 10, 64)
	c.Set(ctxUserIdKey, uid)
}

func userIdGet(c *gin.Context) int64 {
	return c.GetInt64(ctxUserIdKey)
}

func userRolesSet(c *gin.Context, roles []string) {
	c.Set(ctxUserRolesKey, roles)
}

func userRolesGet(c *gin.Context) []string {
	return c.GetStringSlice(ctxUserRolesKey)
}

func tokenCookieSet(c *gin.Context, token string, expireSec int) {
	c.SetCookie(cookieTokenKey, token, expireSec, "/", "", false, true)
}

func tokenCookieGet(c *gin.Context) (string, error) {
	return c.Cookie(cookieTokenKey)
}

func tokenCookieClean(c *gin.Context) {
	ginutil.Cookie(c, cookieTokenKey, "", 1)
}
