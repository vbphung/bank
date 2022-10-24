package gapi

import (
	"context"

	db "github.com/vbph/bank/db/sqlc"
	"github.com/vbph/bank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) generateToken(ctx context.Context, email string) (accessTkRes *pb.TokenRes, refreshTkRes *pb.TokenRes, err error) {
	payload, accessTk, err := server.tokenMaker.CreateToken(email, server.config.AccessTokenExpiredTime)
	if err != nil {
		return
	}

	accessTkRes = &pb.TokenRes{
		Token:     accessTk,
		ExpiredAt: timestamppb.New(payload.ExpiredAt),
	}

	payload, refreshTk, err := server.tokenMaker.CreateToken(email, server.config.RefreshTokenExpiredTime)
	if err != nil {
		return
	}

	md := server.extractMetadata(ctx)

	if _, err = server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           payload.ID,
		Email:        payload.Email,
		RefreshToken: refreshTk,
		UserAgent:    md.UserAgent,
		ClientIp:     md.ClientIp,
		ExpiredAt:    payload.ExpiredAt,
	}); err != nil {
		return
	}

	refreshTkRes = &pb.TokenRes{
		Token:     refreshTk,
		ExpiredAt: timestamppb.New(payload.ExpiredAt),
	}

	return
}
