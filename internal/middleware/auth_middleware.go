package middleware

import (
	"database/sql"
	"net/http"
	"strings"

	"hospital/internal/models"
	"hospital/internal/repository"
	"hospital/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("supersecret")

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		tokenStr := parts[1]

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing user in token"})
			return
		}

		user, err := repository.GetUserByID(userID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to load user"})
			return
		}

		c.Set("currentUser", user)
		c.Set("userID", user.ID)
		c.Set("role", user.Role)

		c.Next()
	}
}

func CurrentUser(c *gin.Context) *models.User {
	value, exists := c.Get("currentUser")
	if !exists {
		return nil
	}

	user, ok := value.(*models.User)
	if !ok {
		return nil
	}

	return user
}

func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := CurrentUser(c)
		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user context not found"})
			return
		}

		if !service.HasPermission(user.Role, permission) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			return
		}

		c.Next()
	}
}

func RequireClinicScope(paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := CurrentUser(c)
		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user context not found"})
			return
		}

		if user.Role == service.RoleSystemAdmin {
			c.Next()
			return
		}

		targetClinicID := c.Param(paramName)
		if targetClinicID == "" {
			c.Next()
			return
		}

		if !user.ClinicID.Valid || user.ClinicID.String != targetClinicID {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "clinic scope violation"})
			return
		}

		c.Next()
	}
}
