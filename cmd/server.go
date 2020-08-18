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

	"github.com/saltbo/gopkg/ginutil"
	"github.com/saltbo/gopkg/gormutil"
	"github.com/saltbo/gopkg/jwtutil"
	"github.com/saltbo/gopkg/mailutil"
	"github.com/spf13/cobra"

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
	gormutil.Init(conf.Database, &model.User{}, &model.UserProfile{})

	rs := ginutil.NewServer(":8081")
	rs.SetupGroupRS("/moreu/api",
		rest.NewUserResource(conf),
		rest.NewTokenResource(conf),
	)
	rs.SetupStatic("/moreu", conf.Moreu)
	rs.SetupIndex("/moreu", ginutil.NewIndex(conf.Moreu))
	rs.SetupSwagger()
	rs.SetupPing()

	rest.RBACInit("roles.yml")
	for _, router := range conf.Routers {
		rs.SetupRS(rest.NewReverseProxy(router))
	}

	for _, static := range conf.Statics {
		rs.SetupStatic(static.Pattern, static.DistDir)
		rs.SetupIndex(static.Pattern, ginutil.NewIndex(static.DistDir, rest.StaticAuth))
	}

	// server run
	if err := rs.Run(); err != nil {
		log.Fatalln(err)
	}
}
