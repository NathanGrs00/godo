// Package middleware holds Gin middleware used by the API.
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// userIDKey is where RequireAuth stores the logged-in user's ID.
const userIDKey = "userID"

// TokenVerifier checks a login token and returns its user. auth.TokenManager implements it.
type TokenVerifier interface {
	Verify(token string) (primitive.ObjectID, error)
}

// RequireAuth rejects requests without a valid "Authorization: Bearer <token>"
// header and stores the user's ID for handlers (see UserID).
func RequireAuth(tokens TokenVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		token, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		userID, err := tokens.Verify(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
		c.Set(userIDKey, userID)
		c.Next()
	}
}

// UserID returns the user ID stored by RequireAuth.
// It must only be used on routes behind RequireAuth.
func UserID(c *gin.Context) primitive.ObjectID {
	return c.MustGet(userIDKey).(primitive.ObjectID)
}
