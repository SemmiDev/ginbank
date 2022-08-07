package rest

import (
	"ginbank/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type renewAccessTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (s *Server) renewAccessToken(c *gin.Context) {
	var req renewAccessTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res := s.Service.RenewAccessToken(c, req.RefreshToken)
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
