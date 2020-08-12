package rest

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/saltbo/gopkg/httputil"

	"github.com/saltbo/moreu/config"
)

type ReverseProxy struct {
	routers config.Routers
}

func NewReverseProxy(routers config.Routers) *ReverseProxy {
	return &ReverseProxy{
		routers: routers,
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
		rRouter.Use(LoginAuth, RoleAuth).Any("/*action", func(c *gin.Context) {
			upstream.ServeHTTP(c.Writer, c.Request)
		})
	}
}
