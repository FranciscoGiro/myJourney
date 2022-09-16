package middlewares

import (
	"fmt"
	"net/http"
	"errors"
	"github.com/gin-gonic/gin"

	"github.com/FranciscoGiro/myJourney/backend/src/auth"
)


func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		token, err := ctx.Cookie("Authorization")
		if err != nil {
			fmt.Println("Unable to retrieve auth cookie. Error:", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("Unable to read auth cookie"))
			return
		}


		payload, err := auth.ValidateToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("Invalid token"))
			return
		}

		ctx.Set("user", (*payload).User)
		ctx.Next()
	}
}