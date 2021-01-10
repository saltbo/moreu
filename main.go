/*
Copyright © 2020 Ambor <saltbo@foxmail.com>

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
package main

import (
	"github.com/saltbo/moreu/cmd"
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

// @BasePath /api

func main() {
	cmd.Execute()
}
