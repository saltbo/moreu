package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/saltbo/gopkg/ginutil"
	"github.com/spf13/viper"

	"github.com/saltbo/moreu/api/server"
	"github.com/saltbo/moreu/assets"
	"github.com/saltbo/moreu/config"
	_ "github.com/saltbo/moreu/docs"
)

// @title Moreu API
// @version 1.0.0
// @description This is a moreu server.

// @contact.name More Support
// @contact.url https://saltbo.cn
// @contact.email saltbo@foxmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api/moreu
func SetupServerRoutes(ge *gin.Engine) {
	apiRouter := ge.Group("/api/moreu")
	apiRouter.Use(func(c *gin.Context) {
		if viper.ConfigFileUsed() == "" {
			ginutil.JSONError(c, http.StatusInternalServerError, fmt.Errorf("system is not initialized"))
			return
		}
	})

	conf := config.Parse()
	ginutil.SetupResource(apiRouter,
		server.NewConfigResource(),
		server.NewTokenResource(conf.EmailAct()),
		server.NewUserResource(conf.EmailAct(), conf.Invitation),
	)

	if conf.MoreuRoot != "" {
		ge.Static("/moreu", conf.MoreuRoot)
	} else {
		ge.StaticFS("/moreu", assets.EmbedFS())
	}

	ge.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/moreu/") {
			c.FileFromFS("/", assets.EmbedFS())
			c.Abort()
			return
		}
	})
}
