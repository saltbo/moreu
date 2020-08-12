package rest

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/saltbo/gopkg/ginutil"
	"github.com/storyicon/grbac"

	"github.com/saltbo/moreu/service"
)

var defaultRBAC *grbac.Controller

func RBACInit(name string) {
	roleLoader := grbac.WithYAML(name, time.Second)
	rbac, err := grbac.New(roleLoader)
	if err != nil {
		log.Fatalln(err)
	}

	defaultRBAC = rbac
}

func LoginAuth(c *gin.Context) {
	token, err := c.Cookie("token")
	if errors.Is(err, http.ErrNoCookie) {
		ginutil.JSONUnauthorized(c, fmt.Errorf("none token!"))
		return
	} else if err != nil {
		ginutil.JSONUnauthorized(c, err)
		return
	}

	rc, err := service.TokenVerify(token)
	if err != nil {
		ginutil.JSONForbidden(c, err)
		return
	}

	c.Set("roles", rc.Roles)
	c.Request.Header.Set("X-Auth-Sub", rc.Subject)
}

func RoleAuth(c *gin.Context) {
	state, err := defaultRBAC.IsRequestGranted(c.Request, c.GetStringSlice("roles"))
	if err != nil {
		ginutil.JSONForbidden(c, err)
		return
	}

	if !state.IsGranted() {
		ginutil.JSONForbidden(c, fmt.Errorf("您没有权限进行此操作，请联系管理员"))
		return
	}
}
