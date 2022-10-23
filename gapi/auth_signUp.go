package gapi

import (
	"context"

	db "github.com/vbph/bank/db/sqlc"
	"github.com/vbph/bank/pb"
	"github.com/vbph/bank/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) SignUp(ctx context.Context, req *pb.SignUpReq) (*pb.SignUpRes, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	acc, err := server.store.CreateAccount(ctx, db.CreateAccountParams{
		Email:    req.Email,
		Password: hashedPassword,
		Balance:  int64(0),
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	accessTk, refreshTk, err := server.generateToken(ctx, acc.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.SignUpRes{
		AccessToken:  accessTk,
		RefreshToken: refreshTk,
		Account:      accountResponse(acc),
	}, nil
}
