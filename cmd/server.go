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
	"time"

	"github.com/saltbo/gopkg/ginutil"
	"github.com/spf13/cobra"
	"github.com/storyicon/grbac"

	"github.com/saltbo/goubase/config"
	"github.com/saltbo/goubase/model"
	"github.com/saltbo/goubase/pkg/ormutil"
	"github.com/saltbo/goubase/rest"
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

	ormutil.Init(conf.Database.Driver, conf.Database.DSN)
	ormutil.DB().AutoMigrate(&model.User{}, &model.UserProfile{})

	rs := ginutil.NewServer(":8080")
	rs.SetupGroupRS("/ubase/api", rest.NewUserResource(conf))
	rs.SetupGroupRS("/ubase/api", rest.NewTokenResource(conf))
	rs.SetupStatic("/ubase", conf.Root)
	rs.SetupSwagger()
	rs.SetupPing()

	// upstream routers
	roleLoader := grbac.WithYAML(conf.Roles.Loader, time.Second)
	rbac, err := grbac.New(roleLoader)
	if err != nil {
		log.Fatalln(err)
	}
	rs.SetupEngineRS(rest.NewReverseProxy(conf.Routers, rbac))
	rs.SetupStatic("/", conf.Root)
	rs.SetupIndex(conf.Root)

	// server run
	if err := rs.Run(); err != nil {
		log.Fatalln(err)
	}
}
