package rest

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/saltbo/gopkg/ginutil"
	_ "github.com/saltbo/gopkg/httputil"

	"github.com/saltbo/moreu/config"
	"github.com/saltbo/moreu/model"
	"github.com/saltbo/moreu/pkg/ormutil"
	"github.com/saltbo/moreu/rest/bind"
	"github.com/saltbo/moreu/service"
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

	router.GET("/users", LoginAuth, RoleAuth, rs.findAll) // todo 注入默认规则
	router.GET("/users/:username", rs.find)
	//router.PUT("/users/:username", LoginAuth, rs.update)
}

// findAll godoc
// @Tags Users
// @Summary 用户查询
// @Description 获取一个用户信息
// @Accept json
// @Produce json
// @Param query query bind.QueryUser true "参数"
// @Success 200 {object} httputil.JSONResponse{data=gin.H{list=[]model.UserProfile},total=int64}
// @Failure 400 {object} httputil.JSONResponse
// @Failure 500 {object} httputil.JSONResponse
// @Router /users [get]
func (rs *UserResource) findAll(c *gin.Context) {
	p := new(bind.QueryUser)
	if err := c.BindQuery(p); err != nil {
		ginutil.JSONBadRequest(c, err)
		return
	}

	list := make([]model.UserProfile, 0)
	if err := ormutil.DB().Offset(p.Offset).Limit(p.Limit).Find(&list).Error; err != nil {
		ginutil.JSONBadRequest(c, err)
		return
	}

	var total int64
	ormutil.DB().Model(model.UserProfile{}).Count(&total)

	ginutil.JSONList(c, list, total)
}

// find godoc
// @Tags Users
// @Summary 用户查询
// @Description 获取一个用户信息
// @Accept json
// @Produce json
// @Param username path string true "用户名"
// @Success 200 {object} httputil.JSONResponse{data=model.UserProfile}
// @Failure 400 {object} httputil.JSONResponse
// @Failure 500 {object} httputil.JSONResponse
// @Router /users/{username} [get]
func (rs *UserResource) find(c *gin.Context) {
	user := &model.User{Username: c.Param("username")}
	if err := ormutil.DB().First(user).Error; err != nil {
		ginutil.JSONBadRequest(c, err)
		return
	}

	userProfile := &model.UserProfile{UserId: user.ID}
	if err := ormutil.DB().First(userProfile).Error; err != nil {
		ginutil.JSONBadRequest(c, err)
		return
	}

	ginutil.JSONData(c, userProfile)
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

	user, err := service.UserCreate(p.Email, p.Password)
	if err != nil {
		ginutil.JSONBadRequest(c, err)
		return
	}

	token, err := service.TokenCreate(p.Email, 6*3600, user.RolesSplit()...)
	if err != nil {
		ginutil.JSONServerError(c, err)
		return
	}

	activateLink := service.ActivateLink(rs.conf.Host, p.Email, token)
	if err := service.SignupNotify(p.Email, activateLink); err != nil {
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
	if p.Activated {
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
