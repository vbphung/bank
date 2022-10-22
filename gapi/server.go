package gapi

import (
	db "github.com/vbph/bank/db/sqlc"
	"github.com/vbph/bank/pb"
	"github.com/vbph/bank/token"
	"github.com/vbph/bank/utils"
)

type Server struct {
	pb.UnimplementedBankServer
	store      *db.Store
	tokenMaker token.Maker
	config     utils.Config
}

func CreateServer(config utils.Config, dbStore *db.Store) (*Server, error) {
	tokenMaker, err := token.CreatePasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		store:      dbStore,
		tokenMaker: tokenMaker,
		config:     config,
	}

	return server, nil
}
