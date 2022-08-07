package rest

import (
	"ginbank/service"
	"ginbank/token"
	"ginbank/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type transferReq struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (s *Server) createTransfer(c *gin.Context) {
	var req transferReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := service.TransferReq{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		AuthPayload:   authPayload,
	}

	res := s.Service.CreateTransfer(c, &arg)
	if res.Error.Value != nil {
		switch res.Error.Kind {
		case utils.KindErrNotFound:
			c.JSON(http.StatusNotFound, errorResponse(res.Error.Value))
			return
		case utils.KindErrUnauthorized:
			c.JSON(http.StatusUnauthorized, errorResponse(res.Error.Value))
			return
		case utils.KindErrBadRequest:
			c.JSON(http.StatusUnauthorized, errorResponse(res.Error.Value))
			return
		default:
			c.JSON(http.StatusInternalServerError, errorResponse(res.Error.Value))
			return
		}
	}

	c.JSON(http.StatusOK, res.Data)
}
