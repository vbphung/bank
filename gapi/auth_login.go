package gapi

import (
	"context"
	"database/sql"

	"github.com/vbph/bank/pb"
	"github.com/vbph/bank/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRes, error) {
	acc, err := server.store.ReadAccount(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	err = utils.VerifyPassword(req.Password, acc.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.LoginRes{}, nil
}
