package middleware

import (
	"context"
	"net/http"
	"strings"

	"chat-system-pro/config"
	"chat-system-pro/models"

	"github.com/gin-gonic/gin"
)

// InviteCodeCheckMiddleware 邀请码验证中间件
func InviteCodeCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.Cfg.Security.InviteCodeEnabled {
			inviteCode := c.PostForm("invite_code")
			if inviteCode == "" {
				inviteCode = c.Query("invite_code")
			}
			
			if inviteCode == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"code": 4000,
					"msg":  "Invitation code is required",
				})
				c.Abort()
				return
			}
			
			var code models.InviteCode
			if err := config.DB.Where("code = ? AND status = 0 AND (max_count = 0 OR used_count < max_count)", inviteCode).First(&code).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code": 4001,
					"msg":  "Invalid invitation code",
				})
				c.Abort()
				return
			}
			
			c.Set("invite_code", &code)
		}
		
		c.Next()
	}
}

// CaptchaMiddleware 验证码验证
func CaptchaMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.Cfg.Security.CaptchaEnabled {
			c.Next()
			return
		}
		
		captchaID := c.PostForm("captcha_id")
		captchaValue := c.PostForm("captcha_value")
		
		if captchaID == "" || captchaValue == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 4002,
				"msg":  "Captcha is required",
			})
			c.Abort()
			return
		}
		
		// 从Redis获取验证码
		ctx := context.Background()
		trueValue, err := config.RDB.Get(ctx, "captcha:"+captchaID).Result()
		if err != nil || !strings.EqualFold(trueValue, captchaValue) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 4003,
				"msg":  "Invalid captcha",
			})
			c.Abort()
			return
		}
		
		// 验证成功，删除验证码（防止重复使用）
		config.RDB.Del(ctx, "captcha:"+captchaID)
		
		c.Next()
	}
}

// SecurityHeadersMiddleware 安全头部中间件
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 安全相关的 HTTP 头部
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'")
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		
		c.Next()
	}
}
