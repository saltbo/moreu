package rest

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/saltbo/gopkg/ginutil"
	"github.com/storyicon/grbac"

	"github.com/saltbo/moreu/client"
	"github.com/saltbo/moreu/model"
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
	token, err := tokenCookieGet(c)
	if errors.Is(err, http.ErrNoCookie) {
		token, _ = service.TokenCreate(0, 30, model.RoleGuest) // 未登录状态颁发一个匿名Token
	}

	rc, err := service.TokenVerify(token)
	if err != nil {
		tokenError(c, err)
		return
	}

	state, err := defaultRBAC.IsRequestGranted(c.Request, rc.Roles)
	if err != nil {
		grantedError(c, err)
		return
	}

	if !state.IsGranted() {
		notGrantedError(c)
		return
	}

	userIdSet(c, rc.Subject)
	client.InjectUserId(c.Request, rc.Subject)
}

func tokenError(c *gin.Context, err error) {
	accept := c.Request.Header.Get("Accept")
	if strings.Contains(accept, gin.MIMEJSON) {
		ginutil.JSONUnauthorized(c, err)
	} else {
		ginutil.FoundRedirect(c, service.Link2SignIn(c.Request.URL.RequestURI()))
	}
}

func grantedError(c *gin.Context, err error) {
	accept := c.Request.Header.Get("Accept")
	if strings.Contains(accept, gin.MIMEJSON) {
		ginutil.JSONServerError(c, err)
	} else {
		ginutil.FoundRedirect(c, service.Link2ServerError(err))
	}
}

func notGrantedError(c *gin.Context) {
	accept := c.Request.Header.Get("Accept")
	if strings.Contains(accept, gin.MIMEJSON) {
		ginutil.JSONForbidden(c, fmt.Errorf("access deny"))
	} else {
		ginutil.FoundRedirect(c, service.Link2SignIn(c.Request.URL.RequestURI()))
	}
}

// auth k-v
const (
	ctxUserIdKey = "user_id"

	cookieTokenKey = "moreu-token"
)

func userIdSet(c *gin.Context, userId string) {
	uid, _ := strconv.ParseInt(userId, 10, 64)
	c.Set(ctxUserIdKey, uid)
}

func userIdGet(c *gin.Context) int64 {
	return c.GetInt64(ctxUserIdKey)
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
