package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/vbph/bank/db/sqlc"
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
