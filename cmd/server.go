/*
Copyright © 2020 Ambor <saltbo@foxmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
//go:generate statik -src=../../moreu-front/dist -dest .. -p assets
package cmd

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/saltbo/gopkg/ginutil"
	"github.com/saltbo/gopkg/gormutil"
	"github.com/saltbo/gopkg/jwtutil"
	"github.com/saltbo/gopkg/mailutil"
	"github.com/spf13/cobra"

	_ "github.com/saltbo/moreu/assets"
	"github.com/saltbo/moreu/config"
	"github.com/saltbo/moreu/model"
	"github.com/saltbo/moreu/rest"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		serverRun()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func serverRun() {
	conf, err := config.Parse()
	if err != nil {
		log.Fatalf("conf: %v", err)
	}

	jwtutil.Init(conf.Secret)
	mailutil.Init(conf.Email)
	gormutil.Init(conf.Database, &model.User{}, &model.UserProfile{}, &model.UserInvitation{})
	rest.RBACInit("roles.yml")

	ge := gin.Default()
	ginutil.SetupSwagger(ge)
	ginutil.SetupPing(ge)

	// system api
	apiRouter := ge.Group("/moreu/api")
	ginutil.SetupResource(apiRouter,
		rest.NewTokenResource(conf),
		rest.NewUserResource(conf),
	)

	// system front
	sysRouter := ge.Group("/moreu")
	simpleRouter := ginutil.NewSimpleRouter()
	if conf.MoreuRoot != "" {
		ginutil.SetupStaticAssets(sysRouter, conf.MoreuRoot)
		simpleRouter.StaticIndex("/moreu", conf.MoreuRoot)
	} else {
		ginutil.SetupEmbedAssets(sysRouter, "/css", "/js", "/fonts")
		simpleRouter.StaticFsIndex("/moreu", ginutil.EmbedFS())
	}

	// reverse proxy
	for _, router := range conf.Routers {
		if router.Pattern == "/" {
			simpleRouter.Route("/", rest.LoginAuth, rest.ReverseProxy(router))
			continue
		}

		ge.Any(router.Pattern+"/*action", rest.LoginAuth, rest.ReverseProxy(router))
	}

	// static serve
	for _, static := range conf.Statics {
		assetsRouter := ge.Group(static.Pattern)
		ginutil.SetupStaticAssets(assetsRouter, static.DistDir)
		simpleRouter.StaticIndex(static.Pattern, static.DistDir)
	}

	// server run
	ge.NoRoute(simpleRouter.Handler)
	ginutil.Startup(ge, ":8081")
}
