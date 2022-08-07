package rest

import (
	db "ginbank/db/sqlc"
	"ginbank/service"
	"ginbank/token"
	"ginbank/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type createAccountReq struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (s *Server) createAccount(ctx *gin.Context) {
	var req createAccountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	res := s.Service.CreateAccount(ctx, &arg)
	if res.Error.Value != nil {
		switch res.Error.Kind {
		case utils.KindErrUniqueViolation:
			ctx.JSON(http.StatusForbidden, errorResponse(res.Error.Value))
			return
		case utils.KindErrForeignKeyViolation:
			ctx.JSON(http.StatusForbidden, errorResponse(res.Error.Value))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, errorResponse(res.Error.Value))
			return
		}
	}

	ctx.JSON(http.StatusCreated, res.Data)
}

type getAccountReq struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getAccount(ctx *gin.Context) {
	var req getAccountReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := service.GetAccountReq{
		AccountID:   req.ID,
		AuthPayload: authPayload,
	}

	res := s.Service.GetAccount(ctx, &arg)
	if res.Error.Value != nil {
		switch res.Error.Kind {
		case utils.KindErrNotFound:
			ctx.JSON(http.StatusNotFound, errorResponse(res.Error.Value))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, errorResponse(res.Error.Value))
			return
		}
	}

	ctx.JSON(http.StatusOK, res.Data)
}

type listAccountReq struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (s *Server) listAccounts(ctx *gin.Context) {
	var req listAccountReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := service.ListAccountReq{
		PageID:      req.PageID,
		PageSize:    req.PageSize,
		AuthPayload: authPayload,
	}

	res := s.Service.ListAccounts(ctx, &arg)
	if res.Error.Value != nil {
		switch res.Error.Kind {
		case utils.KindErrNotFound:
			ctx.JSON(http.StatusNotFound, errorResponse(res.Error.Value))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, errorResponse(res.Error.Value))
			return
		}
	}

	ctx.JSON(http.StatusOK, res.Data)
}
