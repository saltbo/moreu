package rest

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/saltbo/gopkg/ginutil"
	"github.com/saltbo/gopkg/httputil"
	"github.com/storyicon/grbac"

	"github.com/saltbo/moreu/config"
	"github.com/saltbo/moreu/service"
)

type ReverseProxy struct {
	routers config.Routers
	rbac    *grbac.Controller
}

func NewReverseProxy(routers config.Routers, rbac *grbac.Controller) *ReverseProxy {
	return &ReverseProxy{
		routers: routers,
		rbac:    rbac,
	}
}

func (rp *ReverseProxy) Register(router *gin.RouterGroup) {
	for _, r := range rp.routers {
		u, err := url.Parse(r.Upstream.Address)
		if err != nil {
			log.Fatalf("[upstream] invalid address: %s", err)
		}

		header := http.Header{}
		for k, v := range r.Upstream.Headers {
			header.Set(k, v)
		}

		upstream := httputil.NewReverseProxy(u, header)
		rRouter := router.Group(r.Pattern)
		rRouter.Use(rp.Auth).Any("/*action", func(c *gin.Context) {
			upstream.ServeHTTP(c.Writer, c.Request)
		})
	}
}

func (rp *ReverseProxy) Auth(c *gin.Context) {
	token, err := c.Cookie("token")
	if errors.Is(err, http.ErrNoCookie) {
		ginutil.JSONUnauthorized(c, fmt.Errorf("none token!"))
		return
	} else if err != nil {
		ginutil.JSONUnauthorized(c, err)
		return
	}

	rc, err := service.TokenVerify(token)
	if err != nil {
		ginutil.JSONForbidden(c, err)
		return
	}

	state, err := rp.rbac.IsRequestGranted(c.Request, rc.Roles)
	if err != nil {
		ginutil.JSONForbidden(c, err)
		return
	}

	if !state.IsGranted() {
		ginutil.JSONForbidden(c, fmt.Errorf("您没有权限进行此操作，请联系管理员"))
		return
	}

	c.Request.Header.Set("X-Auth-Sub", rc.Subject)
	c.Next()
}
