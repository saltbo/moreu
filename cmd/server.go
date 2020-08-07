/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"
	"time"

	"github.com/saltbo/gopkg/ginutil"
	"github.com/spf13/cobra"
	"github.com/storyicon/grbac"

	"github.com/saltbo/authcar/config"
	"github.com/saltbo/authcar/pkg/rolec"
	"github.com/saltbo/authcar/rest"
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func serverRun() {
	conf, err := config.Parse()
	if err != nil {
		log.Fatalf("conf: %v", err)
	}

	roleLoader := grbac.WithYAML(conf.Roles.Loader, time.Second)
	jwtRole, err := rolec.NewJWTRole(conf.Secret, roleLoader)
	if err != nil {
		log.Fatalln(err)
	}

	// define and start server
	rs := ginutil.NewServer(":8080")
	oauthRouter := rest.NewOauth(conf.Oauth2, jwtRole)
	oauthRouter.SetupRoleLoader(roleLoaderInit())
	rs.SetupRS("/oauth", oauthRouter)
	for _, router := range conf.Routers {
		rs.SetupRS(router.Pattern, rest.NewReverseProxy(router.Upstream, jwtRole))
	}

	// serve the static files
	if conf.Root != "" {
		rs.SetupStatic("/", conf.Root)
		rs.SetupIndex(conf.Root)
	}

	// server run
	if err := rs.Run(); err != nil {
		log.Fatalln(err)
	}
}

func roleLoaderInit() rest.RoleLoader {
	return func(uid string) ([]string, error) {
		// 查询用户角色
		return nil, nil
	}
}
