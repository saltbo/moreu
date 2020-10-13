package rest

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/saltbo/gopkg/ginutil"
	_ "github.com/saltbo/gopkg/httputil"

	"github.com/saltbo/moreu/rest/bind"
	"github.com/saltbo/moreu/service"
)

type TokenResource struct {
	emailAct bool
}

func NewTokenResource(emailAct bool) *TokenResource {
	return &TokenResource{
		emailAct: emailAct,
	}
}

func (rs *TokenResource) Register(router *gin.RouterGroup) {
	router.POST("/tokens", rs.create)
	router.DELETE("/tokens", rs.delete)
}

// create godoc
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
		} else if rs.emailAct && !user.Activated() {
			ginutil.JSONBadRequest(c, fmt.Errorf("account is not activated"))
			return
		}

		expireSec := 7 * 24 * 3600
		token, err := service.TokenCreate(user.Ux, expireSec, user.RolesSplit()...)
		if err != nil {
			ginutil.JSONServerError(c, err)
			return
		}

		tokenCookieSet(c, token, expireSec)
		ginutil.Cookie(c, cookieRoleKey, user.Roles, expireSec)
		ginutil.JSON(c)
		return
	}

	user, ok := service.UserEmailExist(p.Email)
	if !ok {
		ginutil.JSONBadRequest(c, fmt.Errorf("email not exist"))
		return
	}

	// issue a short-term token for password reset
	token, err := service.TokenCreate(user.Ux, 300)
	if err != nil {
		ginutil.JSONServerError(c, err)
		return
	}

	resetLink := service.PasswordRestLink(ginutil.GetOrigin(c), p.Email, token)
	if err := service.PasswordResetNotify(p.Email, resetLink); err != nil {
		ginutil.JSONServerError(c, err)
		return
	}

	ginutil.JSON(c)
}

// delete godoc
// @Tags Tokens
// @Summary 退出登录
// @Description 用户状态登出
// @Accept json
// @Produce json
// @Success 200 {object} httputil.JSONResponse
// @Failure 400 {object} httputil.JSONResponse
// @Failure 500 {object} httputil.JSONResponse
// @Router /tokens [delete]
func (rs *TokenResource) delete(c *gin.Context) {
	ginutil.Cookie(c, cookieTokenKey, "", 1)
	ginutil.Cookie(c, cookieRoleKey, "", 1)
	return
}
