package service

import (
	"context"
	"database/sql"
	"errors"
	db "ginbank/db/sqlc"
	"ginbank/token"
	"ginbank/utils"
	"github.com/lib/pq"
)

func (s *Service) CreateAccount(ctx context.Context, accountReq *db.CreateAccountParams) (res *Res) {
	res = &Res{Error: &utils.WrapError{}, Data: &db.Account{}}

	account, err := s.store.CreateAccount(ctx, *accountReq)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				res = &Res{Error: &utils.WrapError{Value: err, Kind: utils.KindErrUniqueViolation}}
				return
			}
		}
		res = &Res{Error: &utils.WrapError{Value: err, Kind: utils.KindErrInternalServerError}}
		return
	}

	res.Data = &account
	return
}

type GetAccountReq struct {
	AccountID   int64
	AuthPayload *token.Payload
}

func (s *Service) GetAccount(ctx context.Context, req *GetAccountReq) (res *Res) {
	res = &Res{Error: &utils.WrapError{}, Data: &db.Account{}}

	account, err := s.store.GetAccount(ctx, req.AccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			res = &Res{Error: &utils.WrapError{Value: err, Kind: utils.KindErrNotFound}}
			return
		}
		res = &Res{Error: &utils.WrapError{Value: err, Kind: utils.KindErrInternalServerError}}
		return
	}

	if account.Owner != req.AuthPayload.Username {
		res = &Res{Error: &utils.WrapError{Value: errors.New("account doesn't belong to the authenticated user"), Kind: utils.KindErrUnauthorized}}
		return
	}

	res.Data = &account
	return
}

type ListAccountReq struct {
	PageID   int32
	PageSize int32

	AuthPayload *token.Payload
}

func (s *Service) ListAccounts(ctx context.Context, req *ListAccountReq) (res *Res) {
	res = &Res{Error: &utils.WrapError{}, Data: &[]db.Account{}}

	arg := db.ListAccountsParams{
		Owner:  req.AuthPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := s.store.ListAccounts(ctx, arg)
	if err != nil {
		res = &Res{Error: &utils.WrapError{Value: err, Kind: utils.KindErrInternalServerError}}
		return
	}

	res.Data = &accounts
	return
}
