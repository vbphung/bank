package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/vbph/bank/db/sqlc"
	"github.com/vbph/bank/utils"
)

type createAccountReq struct {
	FullName string `json:"full_name" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, failedResponse(err))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, failedResponse(err))
		return
	}

	acc, err := server.store.CreateAccount(ctx, db.CreateAccountParams{
		FullName: req.FullName,
		Password: hashedPassword,
		Balance:  int64(0),
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, failedResponse(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		successResponse(accountResponse(acc)),
	)
}
