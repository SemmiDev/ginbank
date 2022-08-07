package rest

import (
	"fmt"
	db "ginbank/db/sqlc"
	"ginbank/service"
	"ginbank/token"
	"ginbank/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
	Service    service.IService
	router     *gin.Engine
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

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server, nil
}

func (s *Server) setupRouter() {
	s.router = gin.Default()

	s.router.POST("/users", s.createUser)
	s.router.POST("/users/login", s.loginUser)
	s.router.POST("/tokens/renew_access", s.renewAccessToken)

	authRoutes := s.router.Group("/").Use(authMiddleware(s.tokenMaker))
	authRoutes.POST("/accounts", s.createAccount)
	authRoutes.GET("/accounts/:id", s.getAccount)
	authRoutes.GET("/accounts", s.listAccounts)

	authRoutes.POST("/transfers", s.createTransfer)
}
