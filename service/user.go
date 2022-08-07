package service

import (
	"context"
	"database/sql"
	db "ginbank/db/sqlc"
	"ginbank/utils"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"strings"
	"time"
)

type CreateUserReq struct {
	Username string
	Password string
	FullName string
	Email    string
}

func (s *Service) CreateUser(ctx context.Context, userReq *CreateUserReq) (res *Res) {
	res = &Res{Error: &utils.WrapError{}, Data: &db.User{}}

	hashedPassword, err := utils.HashPassword(userReq.Password)
	if err != nil {
		res.Error = &utils.WrapError{Value: err, Kind: utils.KindErrHashedPassword}
		return
	}

	arg := db.CreateUserParams{
		Username:       strings.ToLower(userReq.Username),
		HashedPassword: hashedPassword,
		FullName:       userReq.FullName,
		Email:          userReq.Email,
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				res.Error = &utils.WrapError{Value: err, Kind: utils.KindErrUniqueViolation}
				return
			}
		}
		res.Error = &utils.WrapError{Value: err, Kind: utils.KindErrInternalServerError}
		return
	}

	res.Data = &user
	return
}

type LoginUserReq struct {
	Req
	Username, Password string
}

type LoginUserRes struct {
	SessionID             uuid.UUID `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User                  *UserRes  `json:"user"`
}

func (s *Service) LoginUser(ctx context.Context, loginReq *LoginUserReq) (res *Res) {
	res = &Res{Error: &utils.WrapError{}, Data: &LoginUserRes{}}

	user, err := s.store.GetUser(ctx, loginReq.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			res.Error = &utils.WrapError{Value: err, Kind: utils.KindErrNotFound}
			return
		}
		res.Error = &utils.WrapError{Value: err, Kind: utils.KindErrInternalServerError}
		return
	}

	err = utils.CheckPassword(loginReq.Password, user.HashedPassword)
	if err != nil {
		res.Error = &utils.WrapError{Value: err, Kind: utils.KindErrUnauthorized}
		return
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		user.Username,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		res.Error = &utils.WrapError{Value: err, Kind: utils.KindErrInternalServerError}
		return
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
		user.Username,
		s.config.RefreshTokenDuration,
	)
	if err != nil {
		res.Error = &utils.WrapError{Value: err, Kind: utils.KindErrInternalServerError}
		return
	}

	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    loginReq.UserAgent,
		ClientIp:     loginReq.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		res.Error = &utils.WrapError{Value: err, Kind: utils.KindErrInternalServerError}
		return
	}

	res.Data = &LoginUserRes{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  NewUserRes(&user),
	}
	return
}
