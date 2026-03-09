package middleware

import (
	"net/http"
	"strings"

	"attendance-system/backend/utils"
	"github.com/gin-gonic/gin"
)

const authUserKey = "auth_user"

type AuthUser struct {
	UserID   uint
	Role     string
	CenterID *uint
}

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			if token := c.Query("token"); token != "" {
				header = "Bearer " + token
			}
		}
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "missing authorization header"})
			return
		}
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid authorization header"})
			return
		}
		claims, err := utils.ParseToken(secret, parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
			return
		}
		c.Set(authUserKey, AuthUser{
			UserID:   claims.UserID,
			Role:     claims.Role,
			CenterID: claims.CenterID,
		})
		c.Next()
	}
}

func MustGetAuthUser(c *gin.Context) AuthUser {
	value, _ := c.Get(authUserKey)
	return value.(AuthUser)
}
