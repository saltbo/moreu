package service

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/saltbo/gopkg/gormutil"
	"github.com/saltbo/gopkg/strutil"

	"github.com/saltbo/moreu/model"
)

var emailExp = regexp.MustCompile(`^[A-Za-z0-9]+([_\.][A-Za-z0-9]+)*@([A-Za-z0-9\-]+\.)+[A-Za-z]{2,6}$`)

func UserEmailExist(email string) (*model.User, bool) {
	return userExist("email", email)
}

func UsernameExist(username string) (*model.User, bool) {
	return userExist("username", username)
}

func UserTicketExist(ticket string) (*model.User, bool) {
	return userExist("ticket", ticket)
}

func userExist(k, v string) (*model.User, bool) {
	user := new(model.User)
	if !gormutil.DB().Where(k+"=?", v).First(user).RecordNotFound() {
		return user, true
	}

	return nil, false
}

type UserSignUpService struct {
	ux     string
	roles  string
	ticket string
}

func NewUserSignUpService() *UserSignUpService {
	return &UserSignUpService{
		ux: strutil.RandomText(32),
	}
}

func (s *UserSignUpService) SetTicket(ticket string) {
	s.ticket = ticket
}

func (s *UserSignUpService) SetRoles(roles ...string) {
	s.roles = strings.Join(roles, ",")
}

func (s *UserSignUpService) Signup(email, password string) error {
	var parentUser *model.User
	if s.ticket != "" {
		pu, exist := UserTicketExist(s.ticket)
		if !exist {
			return fmt.Errorf("ticket not exist")
		}

		parentUser = pu
	}

	_, exist := UserEmailExist(email)
	if exist {
		return fmt.Errorf("email already exist")
	}

	return gormutil.DB().Transaction(func(tx *gorm.DB) error {
		// 创建基本信息
		user := &model.User{
			Ux:       s.ux,
			Email:    email,
			Username: fmt.Sprintf("mu%s", strutil.RandomText(18)),
			Password: strutil.Md5Hex(password),
			Roles:    s.roles,
			Ticket:   strutil.RandomText(6),
		}
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// 创建个人信息
		userProfile := &model.UserProfile{
			Ux:       user.Ux,
			Nickname: email[:strings.Index(email, "@")],
		}
		if err := tx.Create(userProfile).Error; err != nil {
			return err
		}

		if parentUser == nil {
			return nil
		}

		// 记录邀请来源
		userInvitation := &model.UserInvitation{
			Ux:    parentUser.Ux,
			SubUx: user.Ux,
		}
		return tx.Create(userInvitation).Error
	})
}

func (s *UserSignUpService) Ux() string {
	return s.ux
}

func (s *UserSignUpService) Roles() string {
	return s.roles
}

func UserGet(ux string) (*model.User, error) {
	user := new(model.User)
	if gormutil.DB().Where("ux=?", ux).First(user).RecordNotFound() {
		return nil, fmt.Errorf("user not exist")
	}

	return user, nil
}

func UserSignIn(usernameOrEmail, password string) (*model.User, error) {
	userFinder := UsernameExist
	if emailExp.MatchString(usernameOrEmail) {
		userFinder = UserEmailExist
	}

	user, exist := userFinder(usernameOrEmail)
	if !exist {
		return nil, fmt.Errorf("user not exist")
	}

	if user.Password != strutil.Md5Hex(password) {
		return nil, fmt.Errorf("invalid password")
	}

	if !user.Activated {
		return nil, fmt.Errorf("account is not activated")
	}

	return user, nil
}

func UserActivate(ux string) error {
	user, err := UserGet(ux)
	if err != nil {
		return err
	}

	if err := gormutil.DB().Model(user).Update("activated", true).Error; err != nil {
		return err
	}

	return nil
}

// ResetPassword update the new password
func UserPasswordReset(ux, newPwd string) error {
	user, err := UserGet(ux)
	if err != nil {
		return err
	}

	if err := gormutil.DB().Model(user).Update("password", strutil.Md5Hex(newPwd)).Error; err != nil {
		return err
	}
	// record the old password

	return nil
}
