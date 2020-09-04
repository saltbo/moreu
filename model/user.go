package model

import (
	"strings"
	"time"
)

const (
	RoleAdmin  = "admin"
	RoleMember = "member"
	RoleGuest  = "guest"
)

type User struct {
	Id        int64      `json:"id"`
	Ux        string     `json:"ux" gorm:"size:32;unique_index;not null"` // Global unique user ID
	Email     string     `json:"email" gorm:"size:32;unique_index;not null"`
	Username  string     `json:"username" gorm:"size:20;unique_index;not null"`
	Password  string     `json:"-" gorm:"size:32;not null"`
	Activated bool       `json:"activated" gorm:"not null"`
	Roles     string     `json:"roles" gorm:"size:64;not null"`
	Ticket    string     `json:"ticket" gorm:"size:6;unique_index;not null"`
	Deleted   *time.Time `json:"-" gorm:"column:deleted_at"`
	Created   time.Time  `json:"created" gorm:"column:created_at;not null"`
	Updated   time.Time  `json:"updated" gorm:"column:updated_at;not null"`
}

func (u *User) RolesSplit() []string {
	return strings.Split(u.Roles, ",")
}
