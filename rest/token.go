package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/saltbo/gopkg/ginutil"
	_ "github.com/saltbo/gopkg/httputil"

	"github.com/saltbo/moreu/config"
	"github.com/saltbo/moreu/rest/bind"
	"github.com/saltbo/moreu/service"
)

type TokenResource struct {
	conf *config.Config
}

func NewTokenResource(conf *config.Config) *TokenResource {
	return &TokenResource{
		conf: conf,
	}
}

func (rs *TokenResource) Register(router *gin.RouterGroup) {
	router.POST("/tokens", rs.create)
	router.DELETE("/tokens", rs.delete)
}

// patch godoc
// @Tags Tokens
// @Summary 登录/密码重置
// @Description 用于账户登录和申请密码重置
// @Accept json
// @Produce json
// @Param body body bind.BodyToken true "参数"
// @Success 200 {object} httputil.JSONResponse
// @Failure 400 {object} httputil.JSONResponse
// @Failure 500 {object} httputil.JSONResponse
// @Router /tokens [post]
func (rs *TokenResource) create(c *gin.Context) {
	p := new(bind.BodyToken)
	if err := c.ShouldBindJSON(p); err != nil {
		ginutil.JSONBadRequest(c, err)
		return
	}

	// issue a signIn token into cookies
	if p.Password != "" {
		user, err := service.UserSignIn(p.Email, p.Password)
		if err != nil {
			ginutil.JSONBadRequest(c, err)
			return
		}

		expireSec := 7 * 24 * 3600
		token, err := service.TokenCreate(user.Username, expireSec, user.RolesSplit()...)
		if err != nil {
			ginutil.JSONServerError(c, err)
			return
		}

		tokenCookieSet(c, token, expireSec)
		ginutil.JSON(c)
		return
	}

	// issue a short-term token for password reset
	token, err := service.TokenCreate(p.Email, 300)
	if err != nil {
		ginutil.JSONServerError(c, err)
		return
	}

	resetLink := service.PasswordRestLink(rs.conf.Host, p.Email, token)
	if err := service.PasswordResetNotify(p.Email, resetLink); err != nil {
		ginutil.JSONServerError(c, err)
		return
	}

	ginutil.JSON(c)
}

func (rs *TokenResource) delete(c *gin.Context) {
	tokenCookieClean(c)
	return
}
