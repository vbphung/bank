package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vbph/bank/utils"
)

type loginReq struct {
	ID       int64  `uri:"id" binding:"required,min=1"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginRes struct {
	AccessToken string `json:"access_token"`
}

func (server *Server) login(ctx *gin.Context) {
	var req loginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FailedResponse(err))
		return
	}

	acc, err := server.store.ReadAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, utils.FailedResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, utils.FailedResponse(err))
		return
	}

	err = utils.VerifyPassword(req.Password, acc.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.FailedResponse(err))
		return
	}

	accessTk, err := server.tokenMaker.CreateToken(acc.ID, server.config.AccessTokenExpiredTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FailedResponse(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		utils.SuccessResponse(
			loginRes{
				AccessToken: accessTk,
			},
		),
	)
}
