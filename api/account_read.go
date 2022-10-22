package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vbph/bank/utils"
)

type readAccountReq struct {
	Account string `uri:"account" binding:"required,email"`
}

func (server *Server) readAccount(ctx *gin.Context) {
	var req readAccountReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FailedResponse(err))
		return
	}

	acc, err := server.store.ReadAccount(ctx, req.Account)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, utils.FailedResponse(err))
			return
		}

		ctx.JSON(http.StatusBadRequest, utils.FailedResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(acc))
}
