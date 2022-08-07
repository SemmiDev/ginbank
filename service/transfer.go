package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	db "ginbank/db/sqlc"
	"ginbank/token"
	"ginbank/utils"
)

type TransferReq struct {
	FromAccountID int64
	ToAccountID   int64
	Amount        int64
	Currency      string

	AuthPayload *token.Payload
}

func (s *Service) CreateTransfer(ctx context.Context, transferReq *TransferReq) (res *Res) {
	res = &Res{Error: &utils.WrapError{}, Data: &db.TransferTxResult{}}

	fromAccount, wrapErr := s.validAccount(ctx, transferReq.FromAccountID, transferReq.Currency)
	if wrapErr != nil {
		res.Error = wrapErr
		return
	}

	if fromAccount.Owner != transferReq.AuthPayload.Username {
		err := errors.New("from account doesn't belong to the authenticated user")
		res.Error = &utils.WrapError{Value: err, Kind: utils.KindErrUnauthorized}
		return
	}

	if fromAccount.Balance < transferReq.Amount {
		err := errors.New("insufficient balance")
		res.Error = &utils.WrapError{Value: err, Kind: utils.KindErrInsufficientBalance}
		return
	}

	_, wrapErr = s.validAccount(ctx, transferReq.ToAccountID, transferReq.Currency)
	if wrapErr != nil {
		res.Error = wrapErr
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: transferReq.FromAccountID,
		ToAccountID:   transferReq.ToAccountID,
		Amount:        transferReq.Amount,
	}

	result, err := s.store.TransferTx(ctx, arg)
	if err != nil {
		res.Error = &utils.WrapError{Value: err, Kind: utils.KindErrInternalServerError}
		return
	}
	res.Data = result
	return
}

func (s *Service) validAccount(ctx context.Context, accountID int64, currency string) (*db.Account, *utils.WrapError) {
	account, err := s.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &utils.WrapError{Value: err, Kind: utils.KindErrNotFound}
		}
		return nil, &utils.WrapError{Value: err, Kind: utils.KindErrInternalServerError}
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		return nil, &utils.WrapError{Value: err, Kind: utils.KindErrBadRequest}
	}

	return &account, nil
}
