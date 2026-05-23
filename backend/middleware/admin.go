package middleware

import (
	"net/http"

	"chat-system-pro/config"
	"chat-system-pro/models"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware 管理员权限验证中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Unauthorized",
			})
			c.Abort()
			return
		}

		// 查询用户信息
		var user models.User
		if err := config.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "User not found",
			})
			c.Abort()
			return
		}

		// 检查用户角色 (1: 管理员, 2: 超级管理员)
		if user.Role < 1 {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "Permission denied",
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user", &user)
		c.Next()
	}
}

// SuperAdminMiddleware 超级管理员权限验证中间件
func SuperAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Unauthorized",
			})
			c.Abort()
			return
		}

		// 查询用户信息
		var user models.User
		if err := config.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "User not found",
			})
			c.Abort()
			return
		}

		// 检查用户角色 (2: 超级管理员)
		if user.Role < 2 {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "Super admin permission required",
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user", &user)
		c.Next()
	}
}
