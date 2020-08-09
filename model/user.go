package model

import (
	"strings"
	"time"
)

type User struct {
	ID       int64     `json:"id"`
	Email    string    `json:"email" gorm:"type:varchar(100);unique_index;not null"`
	Password string    `json:"password" gorm:"size:32;not null"`
	Roles    string    `json:"roles" gorm:"size:64;not null"`
	Enabled  bool      `json:"enabled" gorm:"not null"`
	Deleted  time.Time `json:"deleted" gorm:"column:deleted_at;not null"`
	Created  time.Time `json:"created" gorm:"column:created_at;not null"`
	Updated  time.Time `json:"updated" gorm:"column:updated_at;not null"`
}

func (User) TableName() string {
	return "user"
}

func (u *User) RolesSplit() []string {
	return strings.Split(u.Roles, ",")
}
