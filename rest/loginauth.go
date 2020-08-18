package rest

import (
	"errors"
	"log"
	"net/http"
	"strconv"
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
		token, _ = service.TokenCreate(0, 30, model.RoleAnonymous) // 未登录状态颁发一个匿名Token
	}

	rc, err := service.TokenVerify(token)
	if err != nil {
		ginutil.JSONUnauthorized(c, err)
		return
	}

	userIdSet(c, rc.Subject)
	client.InjectUserId(c.Request, rc.Subject)
	state, err := defaultRBAC.IsRequestGranted(c.Request, rc.Roles)
	if err != nil {
		ginutil.JSONServerError(c, err)
		return
	}

	if !state.IsGranted() {
		ginutil.JSONForbidden(c, err)
		return
	}
}

func StaticAuth(c *gin.Context) {
	token, err := tokenCookieGet(c)
	if errors.Is(err, http.ErrNoCookie) {
		ginutil.FoundRedirect(c, service.Link2SignIn(c.Request.URL.RequestURI()))
		return
	}

	rc, err := service.TokenVerify(token)
	if err != nil {
		ginutil.FoundRedirect(c, service.Link2SignIn(c.Request.URL.RequestURI()))
		return
	}

	state, err := defaultRBAC.IsRequestGranted(c.Request, rc.Roles)
	if err != nil {
		ginutil.FoundRedirect(c, service.Link2ServerError(err))
		return
	}

	if !state.IsGranted() {
		ginutil.FoundRedirect(c, service.Link2Forbidden())
		return
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
