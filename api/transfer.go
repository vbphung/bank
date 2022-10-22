package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/vbph/bank/db/sqlc"
	"github.com/vbph/bank/token"
	"github.com/vbph/bank/utils"
)

type transferReq struct {
	FromAcc string `json:"from_account" binding:"required,email"`
	ToAcc   string `json:"to_account" binding:"required,email"`
	Amount  int64  `json:"amount" binding:"required,gt=0"`
}

func (server *Server) transfer(ctx *gin.Context) {
	var req transferReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FailedResponse(err))
		return
	}

	fromId, toId, err := server.verifyTransfer(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FailedResponse(err))
	}

	trfRes, err := server.store.Transfer(ctx, db.TransferParams{
		FromAcc: fromId,
		ToAcc:   toId,
		Amount:  req.Amount,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FailedResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(trfRes))
}

func (server *Server) verifyTransfer(ctx *gin.Context, req transferReq) (fromId int64, toId int64, err error) {
	authPayload := ctx.MustGet("auth_payload").(*token.Payload)
	if authPayload.Email != req.FromAcc {
		err = errors.New("cannot transfer money from other's account")
		return
	}

	if req.FromAcc == req.ToAcc {
		err = errors.New("cannot transfer to yourself")
		return
	}

	fromAcc, err := server.store.ReadAccount(ctx, req.FromAcc)
	if err != nil {
		err = errors.New("cannot find from account")
		return
	}

	if fromAcc.Balance < req.Amount {
		err = errors.New("insufficient money to transfer")
		return
	}

	toAcc, err := server.store.ReadAccount(ctx, req.ToAcc)
	if err != nil {
		err = errors.New("cannot find to account")
		return
	}

	return fromAcc.ID, toAcc.ID, nil
}
