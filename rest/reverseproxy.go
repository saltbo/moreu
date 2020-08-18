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
	router config.Router
}

func NewReverseProxy(router config.Router) *ReverseProxy {
	return &ReverseProxy{
		router: router,
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
	rRouter.Use(LoginAuth).Any("/*action", func(c *gin.Context) {
		upstream.ServeHTTP(c.Writer, c.Request)
	})
}
