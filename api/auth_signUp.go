package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/vbph/bank/db/sqlc"
	"github.com/vbph/bank/utils"
)

type signUpReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type signUpRes struct {
	loginRes
	Account accountRes `json:"account"`
}

func (server *Server) signUp(ctx *gin.Context) {
	var req signUpReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FailedResponse(err))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FailedResponse(err))
		return
	}

	acc, err := server.store.CreateAccount(ctx, db.CreateAccountParams{
		Email:    req.Email,
		Password: hashedPassword,
		Balance:  int64(0),
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FailedResponse(err))
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
			signUpRes{
				Account: accountResponse(acc),
				loginRes: loginRes{
					AccessToken:  accessTkRes,
					RefreshToken: refreshTkRes,
				},
			},
		),
	)
}
