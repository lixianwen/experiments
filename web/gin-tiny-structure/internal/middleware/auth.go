package middleware

import (
	"net/http"
	"strings"

	"gdemo/jwthelper"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token is empty"})
			return
		}

		token, err := jwthelper.ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		if c.Request.URL.Path == "/refresh" {
			// assertion
			claims := *(token.Claims.(*jwt.MapClaims))
			c.Set("exp", claims["exp"].(float64))
		}

		c.Next()
	}
}
