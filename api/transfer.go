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
	FromID int64 `json:"from_id" binding:"required,min=1"`
	ToID   int64 `json:"to_id" binding:"required,min=1"`
	Amount int64 `json:"amount" binding:"required,gt=0"`
}

func (server *Server) transfer(ctx *gin.Context) {
	var req transferReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FailedResponse(err))
		return
	}

	if err := server.verifyTransfer(ctx, req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FailedResponse(err))
	}

	trfRes, err := server.store.Transfer(ctx, db.TransferParams{
		FromAcc: req.FromID,
		ToAcc:   req.ToID,
		Amount:  req.Amount,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FailedResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(trfRes))
}

func (server *Server) verifyTransfer(ctx *gin.Context, req transferReq) error {
	authPayload := ctx.MustGet("auth_payload").(*token.Payload)
	if authPayload.AccountID != req.FromID {
		return errors.New("cannot transfer money from other's account")
	}

	if req.FromID == req.ToID {
		return errors.New("cannot transfer to yourself")
	}

	fromAcc, err := server.store.ReadAccount(ctx, req.FromID)
	if err != nil {
		return err
	}

	if fromAcc.Balance < req.Amount {
		return errors.New("insufficient money to transfer")
	}

	if _, err = server.store.ReadAccount(ctx, req.ToID); err != nil {
		return err
	}

	return nil
}
