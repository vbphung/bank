package utils

import "github.com/gin-gonic/gin"

func SuccessResponse(result any) gin.H {
	return gin.H{
		"status": "success",
		"result": result,
	}
}

func FailedResponse(err error) gin.H {
	return gin.H{
		"status": "failed",
		"error":  err.Error(),
	}
}
