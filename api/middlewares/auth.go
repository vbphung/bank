package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vbph/bank/token"
	"github.com/vbph/bank/utils"
)

func Auth(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("auth")
		if len(authHeader) == 0 {
			unauthorizedJsonAbort("auth token is not provided", ctx)
			return
		}

		authFields := strings.Fields(authHeader)
		if len(authFields) <= 1 {
			unauthorizedJsonAbort("invalid auth header format", ctx)
			return
		}

		if strings.ToLower(authFields[0]) != "bearer" {
			unauthorizedJsonAbort("unsupported auth type", ctx)
			return
		}

		accessTk := authFields[1]

		payload, err := tokenMaker.VerifyToken(accessTk)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.FailedResponse(err))
			return
		}

		ctx.Set("auth_payload", payload)

		ctx.Next()
	}
}

func unauthorizedJsonAbort(errMessage string, ctx *gin.Context) {
	ctx.AbortWithStatusJSON(
		http.StatusUnauthorized,
		utils.FailedResponse(errors.New(errMessage)),
	)
}
