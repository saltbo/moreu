package rest

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/saltbo/gopkg/ginutil"
	_ "github.com/saltbo/gopkg/httputil"

	"github.com/saltbo/goubase/config"
	"github.com/saltbo/goubase/rest/bind"
	"github.com/saltbo/goubase/service"
)

type UserResource struct {
	conf *config.Config
}

func NewUserResource(conf *config.Config) *UserResource {
	return &UserResource{
		conf: conf,
	}
}

func (rs *UserResource) Register(router *gin.RouterGroup) {
	router.POST("/users", rs.create)        // 账户注册
	router.PATCH("/users/:email", rs.patch) // 账户激活、密码重置

	//router.GET("/users", rs.findAll)  // 管理员权限
	//router.GET("/users/:uid", rs.find) // 管理员权限
}

func (rs *UserResource) findAll(c *gin.Context) {
	p := new(bind.QueryUser)
	if err := c.BindQuery(p); err != nil {
		ginutil.JSONBadRequest(c, err)
		return
	}

	//list := make([]model.User, 0)
	//total, err := dao.DB.Limit(p.Limit, p.Offset).FindAndCount(&list)
	//if err != nil {
	//	ginutil.JSONBadRequest(c, err)
	//	return
	//}

	//ginutil.JSONList(c, list, total)
}

func (rs *UserResource) find(c *gin.Context) {
	//userId := c.Param("uid")
	//
	//user := new(model.User)
	//if _, err := dao.DB.Id(userId).Get(user); err != nil {
	//	ginutil.JSONBadRequest(c, err)
	//	return
	//}
	//
	//user.Password = ""
	//ginutil.JSONData(c, user)
}

// create godoc
// @Tags Users
// @Summary 用户注册
// @Description 注册一个用户
// @Accept json
// @Produce json
// @Param body body bind.BodyUser true "参数"
// @Success 200 {object} httputil.JSONResponse{data=model.User}
// @Failure 400 {object} httputil.JSONResponse
// @Failure 500 {object} httputil.JSONResponse
// @Router /users [post]
func (rs *UserResource) create(c *gin.Context) {
	p := new(bind.BodyUser)
	if err := c.ShouldBindJSON(p); err != nil {
		ginutil.JSONBadRequest(c, err)
		return
	}

	user, err := service.UserCreate(p.Email)
	if err != nil {
		ginutil.JSONBadRequest(c, err)
		return
	}

	token, err := service.TokenCreate(p.Email, 6*3600, user.RolesSplit()...)
	if err != nil {
		ginutil.JSONServerError(c, err)
		return
	}

	ginutil.JSONData(c, token)
	return

	activeLink := fmt.Sprintf("%s/ubase/tokens?email=%s&at=%s", rs.conf.SiteOrigin, p.Email, token)
	if err := service.SignupNotify(p.Email, activeLink); err != nil {
		ginutil.JSONServerError(c, err)
		return
	}

	ginutil.JSON(c)
}

// patch godoc
// @Tags Users
// @Summary 更新一项用户信息
// @Description 用于账户激活和密码重置
// @Accept json
// @Produce json
// @Param email path string true "邮箱"
// @Param body body bind.BodyUserPatch true "参数"
// @Success 200 {object} httputil.JSONResponse
// @Failure 400 {object} httputil.JSONResponse
// @Failure 500 {object} httputil.JSONResponse
// @Router /users/{email} [patch]
func (rs *UserResource) patch(c *gin.Context) {
	p := new(bind.BodyUserPatch)
	if err := c.ShouldBindJSON(p); err != nil {
		ginutil.JSONBadRequest(c, err)
		return
	}

	// valid token
	email := c.Param("email")
	rc, err := service.TokenVerify(p.Token)
	if err != nil {
		ginutil.JSONBadRequest(c, err)
		return
	} else if rc.Subject != email {
		ginutil.JSONBadRequest(c, fmt.Errorf("token not match for the email"))
		return
	}

	// account activate
	if p.Enabled {
		if err := service.UserActivate(email); err != nil {
			ginutil.JSONServerError(c, err)
			return
		}
	}

	// password reset
	if p.Password != "" {
		if err := service.UserPasswordReset(email, p.Password); err != nil {
			ginutil.JSONServerError(c, err)
			return
		}
	}

	ginutil.JSON(c)
}
