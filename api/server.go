package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/herbi-dino/bank/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func CreateServer(svStore *db.Store) *Server {
	sever := &Server{store: svStore}

	router := gin.Default()

	router.POST("/accounts", sever.createAccount)

	sever.router = router

	return sever
}

func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}

func successResponse(result any) gin.H {
	return gin.H{
		"status": "success",
		"result": result,
	}
}

func failedResponse(err error) gin.H {
	return gin.H{
		"status": "failed",
		"error":  err.Error(),
	}
}
