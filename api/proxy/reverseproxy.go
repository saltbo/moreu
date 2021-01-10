package proxy

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/saltbo/gopkg/httputil"

	"github.com/saltbo/moreu/config"
)

func ReverseProxy(router config.Router) gin.HandlerFunc {
	u, err := url.Parse(router.Upstream.Address)
	if err != nil {
		log.Fatalf("[upstream] invalid address: %s", err)
	}

	header := http.Header{}
	for k, v := range router.Upstream.Headers {
		header.Set(k, v)
	}

	upstream := httputil.NewReverseProxy(u, header)
	return func(c *gin.Context) {
		upstream.ServeHTTP(c.Writer, c.Request)
	}
}
