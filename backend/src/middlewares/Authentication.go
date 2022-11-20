package middlewares

import (
	"fmt"
	"net/http"
	"errors"
	"strings"
	"github.com/gin-gonic/gin"

	"github.com/FranciscoGiro/myJourney/backend/src/auth"
)


func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			fmt.Println("No Authorization header found.")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("No authorization header set."))
			return
		}

		authParts := strings.Split(authHeader, " ")
		if authParts[0] != "Bearer" {
			fmt.Println("Access token is not of type Bearer $token.")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("Invalid access token type."))
			return
		}

		token := authParts[1]

		payload, err := auth.ValidateToken(token)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("Invalid token."))
			return
		}

		c.Set("user", (*payload).UserID)
		c.Next()
	}
}