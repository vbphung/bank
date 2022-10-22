package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/vbph/bank/db/sqlc"
)

type accessTokenRes struct {
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}

type refreshTokenRes struct {
	accessTokenRes
	ID uuid.UUID `json:"id"`
}

type tokenRes struct {
	AccessToken  accessTokenRes  `json:"access_token"`
	RefreshToken refreshTokenRes `json:"refresh_token"`
}

func (server *Server) generateToken(ctx *gin.Context, email string) (tokenRes, error) {
	payload, accessTk, err := server.tokenMaker.CreateToken(email, server.config.AccessTokenExpiredTime)
	if err != nil {
		return tokenRes{}, err
	}

	accessTkRes := accessTokenRes{
		Token:     accessTk,
		ExpiredAt: payload.ExpiredAt,
	}

	payload, refreshTk, err := server.tokenMaker.CreateToken(email, server.config.RefreshTokenExpiredTime)
	if err != nil {
		return tokenRes{}, err
	}

	if _, err = server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           payload.ID,
		Email:        payload.Email,
		RefreshToken: refreshTk,
		ExpiredAt:    payload.ExpiredAt,
	}); err != nil {
		return tokenRes{}, err
	}

	refreshTkRes := refreshTokenRes{
		accessTokenRes: accessTokenRes{
			Token:     refreshTk,
			ExpiredAt: payload.ExpiredAt,
		},
		ID: payload.ID,
	}

	return tokenRes{
		AccessToken:  accessTkRes,
		RefreshToken: refreshTkRes,
	}, nil
}
