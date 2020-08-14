package model

import (
	"strings"
	"time"
)

const (
	RoleAdmin  = "admin"
	RoleMember = "member"
)

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email" gorm:"size:32;unique_index;not null"`
	Username  string    `json:"username" gorm:"size:20;unique_index;not null"`
	Password  string    `json:"-" gorm:"size:32;not null"`
	Roles     string    `json:"roles" gorm:"size:64;not null"`
	Activated bool      `json:"activated" gorm:"not null"`
	Deleted   time.Time `json:"deleted" gorm:"column:deleted_at;not null"`
	Created   time.Time `json:"created" gorm:"column:created_at;not null"`
	Updated   time.Time `json:"updated" gorm:"column:updated_at;not null"`
}

func (User) TableName() string {
	return "user"
}

func (u *User) RolesSplit() []string {
	return strings.Split(u.Roles, ",")
}
