package gapi

import (
	"fmt"
	db "ginbank/db/sqlc"
	"ginbank/pb"
	"ginbank/service"
	"ginbank/token"
	"ginbank/utils"
)

type Server struct {
	pb.UnimplementedGinBankServer
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
	Service    service.IService
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.Service = service.NewService(config, store, tokenMaker)
	return server, nil
}
