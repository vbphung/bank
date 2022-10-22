package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vbph/bank/utils"
)

type refreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (server *Server) refreshToken(ctx *gin.Context) {
	var req refreshTokenReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FailedResponse(err))
		return
	}

	payload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.FailedResponse(err))
		return
	}

	session, err := server.store.ReadSession(ctx, payload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, utils.FailedResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, utils.FailedResponse(err))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		ctx.JSON(http.StatusUnauthorized, utils.FailedResponse(
			errors.New("mismatched refresh token"),
		))

		return
	}

	if session.ExpiredAt.Before(time.Now()) {
		ctx.JSON(http.StatusUnauthorized, utils.FailedResponse(
			errors.New("expired refresh token"),
		))

		return
	}

	if session.Email != payload.Email {
		ctx.JSON(http.StatusUnauthorized, utils.FailedResponse(
			errors.New("incorrect account"),
		))

		return
	}

	payload, accessTk, err := server.tokenMaker.CreateToken(payload.Email, server.config.AccessTokenExpiredTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FailedResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		tokenRes{
			Token:     accessTk,
			ExpiredAt: payload.ExpiredAt,
		},
	))
}
