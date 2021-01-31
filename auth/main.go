package auth

import (
	"github.com/gin-gonic/gin"
	"meme/helpers"
	"net/http"
	"strings"
	"time"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessTokenBearer := c.Request.Header.Get("Authorization")

		if accessTokenBearer == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": "NO_ACCESS",
				"message": "Restricted access",
			})
			return
		}

		accessToken := strings.Split(accessTokenBearer, "Bearer ")
		accessTokenClaims, err := helpers.ParseToken(accessToken[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code": "BAD_REQUEST",
				"message": "Invalid token",
			})
			return
		}

		expireTime, _ := time.Parse(time.RFC3339, accessTokenClaims["expire"].(string))

		if isTokenExpired := (expireTime.Unix() - time.Now().Unix()) < 0; isTokenExpired {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": "SESSION_EXPIRED",
			})
			return
		}

		c.Next()
	}
}