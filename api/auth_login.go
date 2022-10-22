package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/vbph/bank/db/sqlc"
	"github.com/vbph/bank/utils"
)

type loginReq struct {
	ID       int64  `uri:"id" binding:"required,min=1"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginRes struct {
	AccessToken  tokenRes        `json:"access_token"`
	RefreshToken refreshTokenRes `json:"refresh_token"`
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

	payload, accessTk, err := server.tokenMaker.CreateToken(acc.Email, server.config.AccessTokenExpiredTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FailedResponse(err))
		return
	}

	accessTkRes := tokenRes{
		Token:     accessTk,
		ExpiredAt: payload.ExpiredAt,
	}

	payload, refreshTk, err := server.tokenMaker.CreateToken(acc.Email, server.config.RefreshTokenExpiredTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FailedResponse(err))
		return
	}

	if _, err = server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           payload.ID,
		Email:        payload.Email,
		RefreshToken: refreshTk,
		ExpiredAt:    payload.ExpiredAt,
	}); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FailedResponse(err))
		return
	}

	refreshTkRes := refreshTokenRes{
		tokenRes: tokenRes{
			Token:     refreshTk,
			ExpiredAt: payload.ExpiredAt,
		},
		ID: payload.ID,
	}

	ctx.JSON(
		http.StatusOK,
		utils.SuccessResponse(
			loginRes{
				AccessToken:  accessTkRes,
				RefreshToken: refreshTkRes,
			},
		),
	)
}
