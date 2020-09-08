package moreu

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/saltbo/gopkg/ginutil"
	"github.com/saltbo/gopkg/jwtutil"
	"github.com/storyicon/grbac"

	"github.com/saltbo/moreu/assets"
	"github.com/saltbo/moreu/model"
	"github.com/saltbo/moreu/rest"
)

type Engine struct {
	r  *gin.Engine
	db *gorm.DB
}

func New(r *gin.Engine, db *gorm.DB) *Engine {
	jwtutil.Init("test123") // todo save me on the fisrt launch.
	db.AutoMigrate(model.Tables()...)
	return &Engine{
		r:  r,
		db: db,
	}
}

// system api
func (e *Engine) SetupAPI(emailAct, invitation bool) {
	apiRouter := e.r.Group("/api/moreu")
	ginutil.SetupResource(apiRouter,
		rest.NewTokenResource(emailAct),
		rest.NewUserResource(emailAct, invitation),
	)
}

// static assets
func (e *Engine) SetupStatic(root string) {
	e.r.Static("/moreu", root)
}

// static assets
func (e *Engine) SetupEmbedStatic() {
	e.r.StaticFS("/moreu", assets.EmbedFS())
}

func (e Engine) NoRoute(c *gin.Context) {
	if strings.HasPrefix(c.Request.URL.Path, "/moreu/") {
		c.FileFromFS("/", assets.EmbedFS())
		c.Abort()
		return
	}
}

func (e *Engine) Auth(roles grbac.Rules) gin.HandlerFunc {
	return rest.LoginAuthWithRoles(roles)
}
