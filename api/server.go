package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/vbph/bank/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func CreateServer(svStore *db.Store) *Server {
	server := &Server{store: svStore}
	server.initRouter()

	return server
}

func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}

func (server *Server) initRouter() {
	router := gin.Default()

	router.POST("/account", server.createAccount)
	router.GET("/account/:id", server.readAccount)
	router.DELETE("/account", server.deleteAccount)

	router.POST("/transfer", server.transfer)

	server.router = router
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
