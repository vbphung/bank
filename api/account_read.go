package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type readAccountReq struct {
	ID int64 `json:"id" binding:"required,min=1"`
}

func (server *Server) readAccount(ctx *gin.Context) {
	var req readAccountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, failedResponse(err))
		return
	}

	acc, err := server.store.ReadAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, failedResponse(err))
			return
		}

		ctx.JSON(http.StatusBadRequest, failedResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, successResponse(acc))
}
