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
	"github.com/saltbo/gopkg/jwtutil"
	"github.com/spf13/cobra"

	"github.com/saltbo/moreu/api/proxy"
	"github.com/saltbo/moreu/config"
	server2 "github.com/saltbo/moreu/internel/app/middleware"
)

// proxyCmd represents the proxy command
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("proxy called")
		proxyRun()
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// proxyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// proxyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func proxyRun() {
	ge := gin.Default()
	ginutil.SetupPing(ge)
	jwtutil.Init("test123") // todo save me on the fisrt launch.

	// reverse proxy
	conf := config.Parse()
	for _, router := range conf.Routers {
		if router.Pattern == "/" {
			ge.NoRoute(server2.LoginAuth(), proxy.ReverseProxy(router))
			continue
		}

		ge.Any(router.Pattern+"/*action", server2.LoginAuth(), proxy.ReverseProxy(router))
	}

	// static serve
	for _, static := range conf.Statics {
		assetsRouter := ge.Group(static.Pattern)
		ginutil.SetupStaticAssets(assetsRouter, static.DistDir)
		//simpleRouter.StaticIndex(static.Pattern, static.DistDir)
		//mu.SetupEmbedStatic()
	}

	ginutil.Startup(ge, ":8082")
}
