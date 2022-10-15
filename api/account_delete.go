package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vbph/bank/utils"
)

type deleteAccountReq struct {
	ID int64 `json:"id" binding:"required,min=1"`
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FailedResponse(err))
		return
	}

	acc, err := server.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, utils.FailedResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, utils.FailedResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(acc))
}
