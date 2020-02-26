package middleware

import (
	"github.com/credo-science/credo-metadata-server/cache"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// before request
		ah := c.GetHeader("Authorization")
		if len(ah) == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Missing authorization header."})
			return
		}

		s := strings.Split(ah, "Bearer ")
		if len(s) != 2 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid authorization header."})
			return
		}

		token := s[1]
		canWrite, err := cache.CanWrite(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Could not validate token. Internal server error."})
			return
		}

		if !canWrite {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token."})
			return
		}

		c.Next()
		// after request
	}
}
