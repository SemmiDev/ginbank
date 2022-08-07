package service

import (
	db "ginbank/db/sqlc"
	"ginbank/utils"
	"time"
)

type Req struct {
	ClientIP, UserAgent string
}

type Res struct {
	Error *utils.WrapError
	Data  any
}

type UserRes struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func NewUserRes(user *db.User) *UserRes {
	return &UserRes{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}
