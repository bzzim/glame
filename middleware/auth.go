package middleware

import (
	"log/slog"
	"strings"

	"github.com/bzzim/glame/pkg/token"
	"github.com/gin-gonic/gin"
)

const AuthKey = "isAuth"

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool(AuthKey) {
			c.Next()
			return
		}
		c.AbortWithStatus(401)
	}
}

func CheckAuth(jwt *token.JWT, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(AuthKey, false)
		header := c.GetHeader("Authorization-Flame")
		if len(header) == 0 {
			log.Debug("empty auth header")
			c.Next()
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			log.Warn("invalid auth header")
			c.Next()
			return
		}

		if jwt.Validate(headerParts[1]) {
			c.Set(AuthKey, true)
		}

		c.Next()
	}
}
