package middleware

import (
	"fmt"
	"github.com/bzzim/glame/pkg/token"
	"github.com/gin-gonic/gin"
	"strings"
)

func Auth(jwt *token.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization-Flame")
		if len(header) == 0 {
			fmt.Println("parseAuthHeader: invalid auth header")
			c.AbortWithStatus(401)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			fmt.Println("parseAuthHeader: invalid auth header")
			c.AbortWithStatus(401)
			return
		}

		if jwt.Validate(headerParts[1]) {
			c.Next()
			return
		}

		c.AbortWithStatus(401)
	}
}

func Check(jwt *token.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("isAuth", false)
		header := c.GetHeader("Authorization-Flame")
		if len(header) == 0 {
			fmt.Println("parseAuthHeader: invalid auth header")
			c.Next()
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			fmt.Println("parseAuthHeader: invalid auth header")
			c.Next()
			return
		}

		if jwt.Validate(headerParts[1]) {
			c.Set("isAuth", true)
		}

		c.Next()
	}
}
