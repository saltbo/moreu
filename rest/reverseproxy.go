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

	"github.com/saltbo/authcar/config"
	"github.com/saltbo/authcar/pkg/rolec"
)

type ReverseProxy struct {
	upstream config.Upstream
	jwtRole  *rolec.JWTRole
}

func NewReverseProxy(upstream config.Upstream, jwtRole *rolec.JWTRole) *ReverseProxy {
	return &ReverseProxy{
		upstream: upstream,
		jwtRole:  jwtRole,
	}
}

func (rp *ReverseProxy) Register(router *gin.RouterGroup) {
	u, err := url.Parse(rp.upstream.Address)
	if err != nil {
		log.Fatalf("[upstream] invalid address: %s", err)
	}

	header := http.Header{}
	for k, v := range rp.upstream.Headers {
		header.Set(k, v)
	}

	upstream := httputil.NewReverseProxy(u, header)
	router.Any("/*action", rp.createReverseProxy(upstream))
}

func (rp *ReverseProxy) createReverseProxy(upstream *httputil.ReverseProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if errors.Is(err, http.ErrNoCookie) {
			ginutil.JSONUnauthorized(c, fmt.Errorf("none token!"))
			return
		} else if err != nil {
			ginutil.JSONUnauthorized(c, err)
			return
		}

		if err := rp.jwtRole.Verify(token, c.Request); err != nil {
			ginutil.JSONForbidden(c, err)
			return
		}

		upstream.ServeHTTP(c.Writer, c.Request)
	}
}
