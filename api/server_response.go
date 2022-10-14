package api

import "github.com/gin-gonic/gin"

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
