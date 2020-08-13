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
	router    config.Router
	protected bool
}

func NewReverseProxy(router config.Router, protected bool) *ReverseProxy {
	return &ReverseProxy{
		router:    router,
		protected: protected,
	}
}

func (rp *ReverseProxy) Register(router *gin.RouterGroup) {
	u, err := url.Parse(rp.router.Upstream.Address)
	if err != nil {
		log.Fatalf("[upstream] invalid address: %s", err)
	}

	header := http.Header{}
	for k, v := range rp.router.Upstream.Headers {
		header.Set(k, v)
	}

	upstream := httputil.NewReverseProxy(u, header)
	rRouter := router.Group(rp.router.Pattern)
	if rp.protected {
		rRouter.Use(APIAuth, RoleAuth)
	}

	rRouter.Any("/*action", func(c *gin.Context) {
		upstream.ServeHTTP(c.Writer, c.Request)
	})
}
