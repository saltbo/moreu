package rest

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/saltbo/gopkg/ginutil"
	"github.com/saltbo/gopkg/gormutil"
	_ "github.com/saltbo/gopkg/httputil"

	"github.com/saltbo/moreu/config"
	"github.com/saltbo/moreu/model"
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

	router.Use(LoginAuth)
	router.GET("/users", rs.findAll)        // 查询用户列表，需管理员权限
	router.GET("/users/:username", rs.find) // 查询某一个用户的公开信息

	router.GET("/user", rs.profile) // 获取已登录用户的所有信息
	router.PUT("/user", rs.update)  // 更新已登录用户个人信息
}

// findAll godoc
// @Tags Users
// @Summary 用户列表
// @Description 获取用户列表信息
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
	if err := gormutil.DB().Offset(p.Offset).Limit(p.Limit).Find(&list).Error; err != nil {
		ginutil.JSONBadRequest(c, err)
		return
	}

	var total int64
	gormutil.DB().Model(model.UserProfile{}).Count(&total)

	ginutil.JSONList(c, list, total)
}

// find godoc
// @Tags Users
// @Summary 用户查询
// @Description 获取一个用户的公开信息
// @Accept json
// @Produce json
// @Param username path string true "用户名"
// @Success 200 {object} httputil.JSONResponse{data=model.UserProfile}
// @Failure 400 {object} httputil.JSONResponse
// @Failure 500 {object} httputil.JSONResponse
// @Router /users/{username} [get]
func (rs *UserResource) find(c *gin.Context) {
	user := new(model.User)
	if gormutil.DB().Where("username=?", c.Param("username")).First(user).RecordNotFound() {
		ginutil.JSONServerError(c, fmt.Errorf("user not exist"))
		return
	}

	userProfile := &model.UserProfile{Ux: user.Ux}
	if err := gormutil.DB().First(userProfile).Error; err != nil {
		ginutil.JSONBadRequest(c, err)
		return
	}

	ginutil.JSONData(c, userProfile)
}

// profile godoc
// @Tags Users
// @Summary 当前登录用户信息
// @Description 获取已登录用户的详细信息
// @Accept json
// @Produce json
// @Param username path string true "用户名"
// @Success 200 {object} httputil.JSONResponse{data=model.UserProfile}
// @Failure 400 {object} httputil.JSONResponse
// @Failure 500 {object} httputil.JSONResponse
// @Router /user [get]
func (rs *UserResource) profile(c *gin.Context) {
	user, err := service.UserGet(uxGet(c))
	if err != nil {
		ginutil.JSONServerError(c, err)
		return
	}

	userProfile := new(model.UserProfile)
	if err := gormutil.DB().First(userProfile, "ux=?", user.Ux).Error; err != nil {
		ginutil.JSONBadRequest(c, err)
		return
	}

	ginutil.JSONData(c, userProfile)
}

// update godoc
// @Tags Users
// @Summary 修改个人信息
// @Description 更新用户的个人信息
// @Accept json
// @Produce json
// @Param username path string true "用户名"
// @Success 200 {object} httputil.JSONResponse
// @Failure 400 {object} httputil.JSONResponse
// @Failure 500 {object} httputil.JSONResponse
// @Router /users/{username} [put]
func (rs *UserResource) update(c *gin.Context) {
	p := new(bind.BodyUserProfile)
	if err := c.ShouldBindJSON(p); err != nil {
		ginutil.JSONBadRequest(c, err)
		return
	}

	user, err := service.UserGet(uxGet(c))
	if err != nil {
		ginutil.JSONServerError(c, err)
		return
	}

	userProfile := new(model.UserProfile)
	if err := gormutil.DB().Where("ux=?", user.Ux).First(userProfile).Error; err != nil {
		ginutil.JSONServerError(c, err)
		return
	}

	userProfile.Avatar = p.Avatar
	userProfile.Nickname = p.Nickname
	userProfile.Bio = p.Bio
	userProfile.URL = p.URL
	userProfile.Company = p.Company
	userProfile.Location = p.Location
	if err := gormutil.DB().Save(userProfile).Error; err != nil {
		ginutil.JSONServerError(c, err)
		return
	}

	ginutil.JSON(c)
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

	if rs.conf.Invitation && p.Ticket == "" {
		ginutil.JSONBadRequest(c, fmt.Errorf("ticket required"))
		return
	}

	opt := service.NewUserCreateOption()
	opt.Roles = model.RoleMember
	opt.Ticket = p.Ticket
	if rs.conf.EmailAct() {
		opt.Origin = ginutil.GetOrigin(c)
	}

	if err := service.UserSignup(p.Email, p.Password, opt); err != nil {
		ginutil.JSONBadRequest(c, err)
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
	rc, err := service.TokenVerify(p.Token)
	if err != nil {
		ginutil.JSONBadRequest(c, err)
		return
	}

	// account activate
	if p.Activated {
		if err := service.UserActivate(rc.Subject); err != nil {
			ginutil.JSONServerError(c, err)
			return
		}
	}

	// password reset
	if p.Password != "" {
		if err := service.UserPasswordReset(rc.Subject, p.Password); err != nil {
			ginutil.JSONServerError(c, err)
			return
		}
	}

	ginutil.JSON(c)
}
