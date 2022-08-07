package service

import (
	"context"
	db "ginbank/db/sqlc"
	"ginbank/token"
	"ginbank/utils"
)

type Service struct {
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewService(config utils.Config, store db.Store, tokenMaker token.Maker) *Service {
	return &Service{config: config, store: store, tokenMaker: tokenMaker}
}

type IService interface {
	IUserService
	IAccountService
	ITokenService
	ITransferService
}

type IUserService interface {
	CreateUser(ctx context.Context, userReq *CreateUserReq) (res *Res)
	LoginUser(ctx context.Context, loginReq *LoginUserReq) (res *Res)
}

type IAccountService interface {
	CreateAccount(ctx context.Context, accountReq *db.CreateAccountParams) (res *Res)
	GetAccount(ctx context.Context, req *GetAccountReq) (res *Res)
	ListAccounts(ctx context.Context, req *ListAccountReq) (res *Res)
}

type ITransferService interface {
	CreateTransfer(ctx context.Context, transferReq *TransferReq) (res *Res)
}

type ITokenService interface {
	RenewAccessToken(ctx context.Context, refreshToken string) (res *Res)
}
