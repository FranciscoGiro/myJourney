package middlewares

import (
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("No authorization header set"))
			return
		}

		authParts := strings.Split(authHeader, " ")
		if authParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("Invalid access token type"))
			return
		}

		token := authParts[1]

		payload, err := auth.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("Invalid token"))
			return
		}

		c.Set("user", (*payload).User)
		c.Next()
	}
}