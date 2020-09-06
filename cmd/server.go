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
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/saltbo/gopkg/ginutil"
	"github.com/saltbo/gopkg/gormutil"
	"github.com/saltbo/gopkg/jwtutil"
	"github.com/saltbo/gopkg/mailutil"
	"github.com/saltbo/gopkg/strutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/saltbo/moreu/assets"
	"github.com/saltbo/moreu/config"
	"github.com/saltbo/moreu/model"
	"github.com/saltbo/moreu/rest"
	"github.com/saltbo/moreu/service"
)

var conf = &config.Config{}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: Run,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().BoolVar(&conf.Debug, "debug", false, "specify the driver of database")
	serverCmd.Flags().StringVar(&conf.Secret, "secret", strutil.RandomText(6), "specify the driver of database")
	serverCmd.Flags().BoolVar(&conf.Invitation, "invitation", false, "specify the driver of database")
	serverCmd.Flags().StringVar(&conf.Database.Driver, "db-driver", "mysql", "specify the driver of database")
	serverCmd.Flags().StringVar(&conf.Database.DSN, "db-dsn", "", "specify the dsn of database")
	serverCmd.Flags().StringVar(&conf.Email.Host, "email-host", "", "specify the host of email")
	serverCmd.Flags().StringVar(&conf.Email.Sender, "email-sender", "", "specify the sender of email")
	serverCmd.Flags().StringVar(&conf.Email.Username, "email-username", "", "specify the username of email")
	serverCmd.Flags().StringVar(&conf.Email.Password, "email-password", "", "specify the password of email")
	serverCmd.Flags().StringVar(&conf.GRbacFile, "grbac-config", "roles.yml", "specify the filepath of grbac roles file")

	serverCmd.Flags().String("proxy-config", "routers.yml", "specify the path of routers.yml")
	viper.BindPFlag("proxy-config", serverCmd.Flags().Lookup("proxy-config"))
}

func Run(cmd *cobra.Command, args []string) {
	if cmd.Flags().NFlag() == 0 {
		conf = config.Parse()
	} else {
		v := viper.New()
		v.SetConfigFile(viper.GetString("proxy-config"))
		v.ReadInConfig()
		if err := v.Unmarshal(&conf); err == nil {
			fmt.Println("Using proxy-config file:", v.ConfigFileUsed())
		}
	}

	rest.RBACInit(conf.GRbacFile)
	jwtutil.Init(conf.Secret)
	if conf.EmailAct() {
		mailutil.Init(conf.Email)
	}

	gormutil.Init(conf.Database, conf.Debug)
	gormutil.SetupPrefix("mu_")
	gormutil.AutoMigrate(model.Tables())
	service.AdministratorInit() // create the user administrator

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
		ginutil.SetupEmbedAssets(sysRouter, assets.EmbedFS(), "/css", "/js", "/fonts")
		simpleRouter.StaticFsIndex("/moreu", assets.EmbedFS())
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
