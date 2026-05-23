package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"chat-system-pro/config"

	"github.com/gin-gonic/gin"
)

// CORSConfig CORS配置
type CORSConfig struct {
	// 允许的来源（域名列表）
	AllowOrigins []string
	// 允许的请求方法
	AllowMethods []string
	// 允许的请求头
	AllowHeaders []string
	// 暴露的响应头
	ExposeHeaders []string
	// 是否允许携带凭证
	AllowCredentials bool
	// 预检请求缓存时间
	MaxAge int
}

// 默认CORS配置
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:8080",
			"https://your-domain.com", // 替换为你的域名
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"PATCH",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
			"X-CSRF-Token",
			"X-Forwarded-For",
			"X-Real-IP",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
			"X-RateLimit-Limit",
			"X-RateLimit-Remaining",
			"X-RateLimit-Reset",
		},
		AllowCredentials: true,
		MaxAge:          86400, // 24小时
	}
}

// CORS中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求来源
		origin := c.GetHeader("Origin")
		
		// 检查origin是否在允许列表中
		if !isOriginAllowed(origin) {
			// 如果origin不在列表中，检查是否为开发环境
			if isDevelopmentMode() {
				// 开发环境允许所有origin
				origin = "*"
			} else {
				// 生产环境拒绝不在列表中的origin
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"code": 403,
					"msg":  "不允许的请求来源",
				})
				return
			}
		}
		
		// 设置CORS响应头
		if origin != "*" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
		}
		
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With, X-CSRF-Token")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type, X-RateLimit-Limit, X-RateLimit-Remaining")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")
		
		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}

// 检查origin是否允许
func isOriginAllowed(origin string) bool {
	if origin == "" {
		return true
	}
	
	// 从配置获取允许的域名列表
	allowedOrigins := getAllowedOrigins()
	
	for _, allowed := range allowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
		// 支持通配符子域名
		if strings.HasPrefix(allowed, "https://") || strings.HasPrefix(allowed, "http://") {
			if strings.HasPrefix(origin, allowed) {
				return true
			}
		}
	}
	
	return false
}

// 获取允许的域名列表
func getAllowedOrigins() []string {
	// 从配置中读取，如果没有配置则使用默认值
	if config.Cfg != nil && len(config.Cfg.CORS.AllowOrigins) > 0 {
		return config.Cfg.CORS.AllowOrigins
	}
	
	// 默认允许的域名
	return []string{
		"http://localhost:3000",
		"http://localhost:8080",
		"http://localhost",
		"http://127.0.0.1:3000",
		"http://127.0.0.1:8080",
	}
}

// 检查是否为开发模式
func isDevelopmentMode() bool {
	if config.Cfg != nil {
		return config.Cfg.Server.Mode == "debug"
	}
	return true
}

// 安全的CORS配置（生产环境）
func SecureCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowOrigins: []string{
			"https://your-domain.com",      // 主域名
			"https://www.your-domain.com",  // www子域名
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:          86400,
	}
}

// 开发环境CORS配置
func DevCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
		ExposeHeaders: []string{"*"},
		AllowCredentials: false,
		MaxAge:          86400,
	}
}

// 动态CORS配置中间件
func DynamicCORSMiddleware(config *CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		
		if len(config.AllowOrigins) > 0 && config.AllowOrigins[0] == "*" {
			c.Header("Access-Control-Allow-Origin", "*")
		} else if isOriginAllowed(origin) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
		}
		
		c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ", "))
		c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ", "))
		c.Header("Access-Control-Expose-Headers", strings.Join(config.ExposeHeaders, ", "))
		
		if config.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		
		c.Header("Access-Control-Max-Age", fmt.Sprintf("%d", config.MaxAge))
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}