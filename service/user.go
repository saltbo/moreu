package service

import (
	"fmt"
	"strings"

	"github.com/saltbo/gopkg/cryptoutil"
	"github.com/saltbo/gopkg/gormutil"
	"github.com/saltbo/gopkg/randutil"

	"github.com/saltbo/moreu/model"
)

func UserEmailExist(email string) (*model.User, bool) {
	user := new(model.User)
	if !gormutil.DB().Where("email = ?", email).First(user).RecordNotFound() {
		return user, true
	}

	return nil, false
}

func UserCreate(email, password string, roles ...string) (*model.User, error) {
	_, exist := UserEmailExist(email)
	if exist {
		return nil, fmt.Errorf("email already exist")
	}

	user := &model.User{
		Email:    email,
		Username: fmt.Sprintf("mu%s", randutil.RandString(18)),
		Password: cryptoutil.Md5Hex(password),
		Roles:    strings.Join(roles, ","),
	}
	if err := gormutil.DB().Create(user).Error; err != nil {
		return nil, err
	}

	userProfile := &model.UserProfile{
		UserId:   user.ID,
		Nickname: email[:strings.Index(email, "@")],
	}
	if err := gormutil.DB().Create(userProfile).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func UserGet(username string) (*model.User, error) {
	user := new(model.User)
	if gormutil.DB().Where("username = ?", username).First(user).RecordNotFound() {
		return nil, fmt.Errorf("user not exist")
	}

	return user, nil
}

func UserSignIn(email, password string) (*model.User, error) {
	user, exist := UserEmailExist(email)
	if !exist {
		return nil, fmt.Errorf("user not exist")
	}

	if user.Password != cryptoutil.Md5Hex(password) {
		return nil, fmt.Errorf("invalid password")
	}

	if !user.Activated {
		return nil, fmt.Errorf("account is not activated")
	}

	return user, nil
}

func UserActivate(username string) error {
	user, err := UserGet(username)
	if err != nil {
		return err
	}

	if err := gormutil.DB().Model(user).Update("activated", true).Error; err != nil {
		return err
	}

	return nil
}

// ResetPassword update the new password
func UserPasswordReset(username, newPwd string) error {
	user, err := UserGet(username)
	if err != nil {
		return err
	}

	if err := gormutil.DB().Model(user).Update("password", cryptoutil.Md5Hex(newPwd)).Error; err != nil {
		return err
	}
	// record the old password

	return nil
}
