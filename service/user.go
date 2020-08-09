package service

import (
	"fmt"
	"strings"

	"github.com/saltbo/gopkg/cryptoutil"

	"github.com/saltbo/goubase/model"
	"github.com/saltbo/goubase/pkg/ormutil"
)

type User struct {
}

func UserExist(email string) (*model.User, bool) {
	user := new(model.User)
	if !ormutil.DB().Where("email = ?", email).First(user).RecordNotFound() {
		return user, true
	}

	return nil, false
}

func UserCreate(email string) (*model.User, error) {
	user, exist := UserExist(email)
	if exist {
		return nil, fmt.Errorf("user already exist")
	}

	user.Email = email
	user.Password = cryptoutil.Md5Hex(email)
	if err := ormutil.DB().Create(user).Error; err != nil {
		return nil, err
	}

	userProfile := &model.UserProfile{
		UserId:   user.ID,
		Nickname: email[:strings.Index(email, "@")],
		//Roles:    ROLE_MEMBER,
	}
	if err := ormutil.DB().Create(userProfile).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func UserSignIn(email, password string) (*model.User, error) {
	user, exist := UserExist(email)
	if !exist {
		return nil, fmt.Errorf("user not exist")
	}

	if !user.Enabled {
		return nil, fmt.Errorf("account is not activated")
	}

	if user.Password != password {
		return nil, fmt.Errorf("invalid password")
	}

	return user, nil
}

func UserActivate(email string) error {
	user, exist := UserExist(email)
	if !exist {
		return fmt.Errorf("user not exist")
	}

	if err := ormutil.DB().Model(user).Update("enabled", true).Error; err != nil {
		return err
	}

	return nil
}

// ResetPassword update the new password
func UserPasswordReset(email, newPwd string) error {
	user, exist := UserExist(email)
	if !exist {
		return fmt.Errorf("user not exist")
	}

	if err := ormutil.DB().Model(user).Update("password", cryptoutil.Md5Hex(newPwd)).Error; err != nil {
		return err
	}
	// record the old password

	return nil
}
