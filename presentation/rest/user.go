package rest

import (
	"ginbank/service"
	"ginbank/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type createUserReq struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func (s *Server) createUser(c *gin.Context) {
	var req createUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res := s.Service.CreateUser(c, &service.CreateUserReq{
		Username: req.Username,
		Password: req.Password,
		FullName: req.FullName,
		Email:    req.Email,
	})

	if res.Error.Value != nil {
		switch res.Error.Kind {
		case utils.KindErrUniqueViolation:
			c.JSON(http.StatusForbidden, errorResponse(res.Error.Value))
			return
		default:
			c.JSON(http.StatusInternalServerError, errorResponse(res.Error.Value))
			return
		}
	}

	c.JSON(http.StatusCreated, res.Data)
}

type loginUserReq struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

func (s *Server) loginUser(c *gin.Context) {
	var req loginUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res := s.Service.LoginUser(c, &service.LoginUserReq{
		Req: service.Req{
			ClientIP:  c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
		},
		Username: req.Username,
		Password: req.Password,
	})

	if res.Error.Value != nil {
		switch res.Error.Kind {
		case utils.KindErrNotFound:
			c.JSON(http.StatusNotFound, errorResponse(res.Error.Value))
			return
		case utils.KindErrUnauthorized:
			c.JSON(http.StatusUnauthorized, errorResponse(res.Error.Value))
			return
		default:
			c.JSON(http.StatusInternalServerError, errorResponse(res.Error.Value))
			return
		}
	}

	c.JSON(http.StatusOK, res.Data)
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
