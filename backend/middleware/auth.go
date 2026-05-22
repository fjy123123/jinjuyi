package middleware

import (
	"net/http"
	"strings"

	"chat-system-pro/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未登录"})
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "token格式错误"})
			c.Abort()
			return
		}
		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Cfg.JWT.Secret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "token无效"})
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "token无效"})
			c.Abort()
			return
		}
		userID := uint(claims["user_id"].(float64))
		c.Set("user_id", userID)
		c.Next()
	}
}
