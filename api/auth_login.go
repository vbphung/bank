package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vbph/bank/utils"
)

type loginReq struct {
	Email    string `uri:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginRes struct {
	AccessToken  tokenRes `json:"access_token"`
	RefreshToken tokenRes `json:"refresh_token"`
}

func (server *Server) login(ctx *gin.Context) {
	var req loginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FailedResponse(err))
		return
	}

	acc, err := server.store.ReadAccount(ctx, req.Email)
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

	res, err := server.generateToken(ctx, acc.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FailedResponse(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		utils.SuccessResponse(res),
	)
}
