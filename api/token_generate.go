package api

import (
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/vbph/bank/db/sqlc"
)

type tokenRes struct {
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (server *Server) generateToken(ctx *gin.Context, email string) (loginRes, error) {
	payload, accessTk, err := server.tokenMaker.CreateToken(email, server.config.AccessTokenExpiredTime)
	if err != nil {
		return loginRes{}, err
	}

	accessTkRes := tokenRes{
		Token:     accessTk,
		ExpiredAt: payload.ExpiredAt,
	}

	payload, refreshTk, err := server.tokenMaker.CreateToken(email, server.config.RefreshTokenExpiredTime)
	if err != nil {
		return loginRes{}, err
	}

	if _, err = server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           payload.ID,
		Email:        payload.Email,
		RefreshToken: refreshTk,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		ExpiredAt:    payload.ExpiredAt,
	}); err != nil {
		return loginRes{}, err
	}

	refreshTkRes := tokenRes{
		Token:     refreshTk,
		ExpiredAt: payload.ExpiredAt,
	}

	return loginRes{
		AccessToken:  accessTkRes,
		RefreshToken: refreshTkRes,
	}, nil
}
