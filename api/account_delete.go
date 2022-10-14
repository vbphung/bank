package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type deleteAccountReq struct {
	ID int64 `json:"id" binding:"required,min=1"`
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, failedResponse(err))
		return
	}

	acc, err := server.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, failedResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, failedResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, successResponse(acc))
}
