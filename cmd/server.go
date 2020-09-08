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
	"github.com/gin-gonic/gin"
	"github.com/saltbo/gopkg/ginutil"
	"github.com/saltbo/gopkg/gormutil"
	"github.com/spf13/cobra"

	"github.com/saltbo/moreu/config"
	"github.com/saltbo/moreu/moreu"
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
}

func Run(cmd *cobra.Command, args []string) {
	conf = config.Parse()
	gormutil.Init(conf.Database, conf.Debug)

	ge := gin.Default()
	ginutil.SetupPing(ge)
	ginutil.SetupSwagger(ge)

	mu := moreu.New(ge, gormutil.DB())
	if conf.EmailAct() {
		mu.SetupMail(conf.Email)
	}
	mu.SetupAPI(conf.EmailAct(), conf.Invitation)
	if conf.MoreuRoot != "" {
		mu.SetupStatic(conf.MoreuRoot)
	} else {
		mu.SetupEmbedStatic()
	}

	ge.NoRoute(mu.NoRoute)
	ginutil.Startup(ge, ":8081")
}
