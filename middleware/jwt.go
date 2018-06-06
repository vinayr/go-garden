package middleware

import (
	"log"
	"strings"
	"time"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/vinayr/go-garden/services"
	"golang.org/x/crypto/bcrypt"
)

// Auth Middleware
func Auth(s *services.UserService, jwtSecret string) *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm: "garden.io",
		Key:   []byte(jwtSecret),
		// Timeout: time.Hour,
		Timeout:    time.Second * 15,
		MaxRefresh: time.Hour,

		Authenticator: authentication(s),
		PayloadFunc:   payload(s),

		Authorizator:  authorization,
		Unauthorized:  unauthorized,
		TokenLookup:   "header:Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	}

}

// authentication (called during signin)
func authentication(s *services.UserService) func(string, string, *gin.Context) (string, bool) {
	return func(username string, password string, c *gin.Context) (string, bool) {
		username = strings.ToLower(username)
		user, err := s.FindUserByUsername(username)
		if err != nil {
			return username, false
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if err != nil {
			return username, false
		}

		return user.Username, true
	}
}

// payload (called during signin)
func payload(s *services.UserService) func(string) map[string]interface{} {
	return func(username string) map[string]interface{} {
		username = strings.ToLower(username)
		user, _ := s.FindUserByUsername(username)
		return map[string]interface{}{
			"is_admin": user.IsAdmin,
		}
	}
}

// authorize authenticated user
func authorization(username string, c *gin.Context) bool {
	return true
}

// unauthorized ...
func unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

// JWTGetCurrentUser extracts username from JWT
func JWTGetCurrentUser(c *gin.Context) string {
	claims := jwt.ExtractClaims(c)
	log.Print("USERNAME", claims["id"])
	log.Print("IS_ADMIN", claims["is_admin"])
	return claims["id"].(string)
}

// Admin middleware
func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		if !claims["is_admin"].(bool) {
			c.AbortWithStatus(401)
			return
		}
		c.Next()
	}
}
