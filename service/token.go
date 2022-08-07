package service

import (
	"context"
	"database/sql"
	"fmt"
	"ginbank/utils"
	"time"
)

type RenewAccessTokenRes struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (s *Service) RenewAccessToken(ctx context.Context, refreshToken string) (res *Res) {
	res = &Res{Error: &utils.WrapError{}, Data: &RenewAccessTokenRes{}}

	refreshPayload, err := s.tokenMaker.VerifyToken(refreshToken)
	if err != nil {
		res.Error = &utils.WrapError{Value: err, Kind: utils.KindErrUnauthorized}
		return
	}

	session, err := s.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			res.Error = &utils.WrapError{Value: err, Kind: utils.KindErrNotFound}
			return
		}
		res.Error = &utils.WrapError{Value: err, Kind: utils.KindErrInternalServerError}
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		res.Error = &utils.WrapError{Value: err, Kind: utils.KindErrUnauthorized}
		return
	}

	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("incorrect session user")
		res.Error = &utils.WrapError{Value: err, Kind: utils.KindErrUnauthorized}
		return
	}

	if session.RefreshToken != refreshToken {
		err := fmt.Errorf("mismatched session token")
		res.Error = &utils.WrapError{Value: err, Kind: utils.KindErrUnauthorized}
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("expired session")
		res.Error = &utils.WrapError{Value: err, Kind: utils.KindErrUnauthorized}
		return
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		refreshPayload.Username,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		res.Error = &utils.WrapError{Value: err, Kind: utils.KindErrInternalServerError}
		return
	}

	res.Data = &RenewAccessTokenRes{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	return
}
