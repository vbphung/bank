package gapi

import (
	db "github.com/vbph/bank/db/sqlc"
	"github.com/vbph/bank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func accountResponse(acc db.Account) *pb.Account {
	return &pb.Account{
		Email:             acc.Email,
		PasswordChangedAt: timestamppb.New(acc.PasswordChangedAt),
		CreatedAt:         timestamppb.New(acc.CreatedAt),
	}
}
