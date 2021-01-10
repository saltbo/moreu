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
package cmd

import (
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/saltbo/gopkg/ginutil"
	"github.com/saltbo/gopkg/gormutil"
	"github.com/saltbo/gopkg/jwtutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/saltbo/moreu/api/server"
	"github.com/saltbo/moreu/assets"
	"github.com/saltbo/moreu/config"
	"github.com/saltbo/moreu/model"
)

// serverCmd represents the middleware command
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
	ge := gin.Default()
	ginutil.SetupPing(ge)
	ginutil.SetupSwagger(ge)
	jwtutil.Init("test123") // todo save me on the fisrt launch.

	conf := config.Parse()
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println(in, in.Name, in.String())
		gormutil.Init(conf.Database, true)
		gormutil.AutoMigrate(model.Tables())
	})

	apiRouter := ge.Group("/api")
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
	ginutil.Startup(ge, ":8081")
}
