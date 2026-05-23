package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"chat-system-pro/config"

	"github.com/gin-gonic/gin"
)

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	// 每分钟允许的请求数
	MaxRequests int
	// 限流窗口（分钟）
	Window time.Duration
}

// 限流存储（使用Redis）
type RateLimiter struct {
	config *RateLimitConfig
}

// 创建新的限流器
func NewRateLimiter(maxRequests int, windowMinutes int) *RateLimiter {
	return &RateLimiter{
		config: &RateLimitConfig{
			MaxRequests: maxRequests,
			Window:      time.Duration(windowMinutes) * time.Minute,
		},
	}
}

// 限流中间件
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端标识（优先使用X-Forwarded-For，然后是X-Real-IP，最后是IP）
		clientIP := c.ClientIP()
		
		// 构建Redis key
		key := fmt.Sprintf("ratelimit:%s:%s", c.FullPath(), clientIP)
		
		ctx := context.Background()
		
		// 获取当前请求计数
		count, err := config.RDB.Get(ctx, key).Int()
		if err != nil {
			count = 0
		}
		
		// 检查是否超过限制
		if count >= rl.config.MaxRequests {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code": 429,
				"msg":  "请求过于频繁，请稍后再试",
				"retry_after": int(rl.config.Window.Seconds()),
			})
			c.Abort()
			return
		}
		
		// 增加计数
		count++
		pipe := config.RDB.Pipeline()
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, rl.config.Window)
		_, err = pipe.Exec(ctx)
		
		if err != nil {
			// Redis错误时记录日志但不阻止请求
			fmt.Printf("Rate limiter Redis error: %v\n", err)
		}
		
		// 设置限流头信息
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.config.MaxRequests))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", rl.config.MaxRequests-count))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(rl.config.Window).Unix()))
		
		c.Next()
	}
}

// 登录限流中间件（更严格的限制）
func LoginRateLimiter() gin.HandlerFunc {
	limiter := NewRateLimiter(5, 1) // 每分钟最多5次登录尝试
	return limiter.Middleware()
}

// 注册限流中间件
func RegisterRateLimiter() gin.HandlerFunc {
	limiter := NewRateLimiter(3, 1) // 每分钟最多3次注册尝试
	return limiter.Middleware()
}

// API通用限流中间件
func APIRateLimiter() gin.HandlerFunc {
	limiter := NewRateLimiter(100, 1) // 每分钟最多100次请求
	return limiter.Middleware()
}

// 动态限流中间件（基于用户ID）
func DynamicRateLimiter(maxRequests int) gin.HandlerFunc {
	limiter := NewRateLimiter(maxRequests, 1)
	return limiter.Middleware()
}

// 暴力破解检测器
type BruteForceDetector struct {
	maxAttempts int
	window      time.Duration
	lockout     time.Duration
}

func NewBruteForceDetector(maxAttempts int, windowMinutes, lockoutMinutes int) *BruteForceDetector {
	return &BruteForceDetector{
		maxAttempts: maxAttempts,
		window:      time.Duration(windowMinutes) * time.Minute,
		lockout:    time.Duration(lockoutMinutes) * time.Minute,
	}
}

// 检测暴力破解
func (bf *BruteForceDetector) Detect(c *gin.Context, identifier string) bool {
	ctx := context.Background()
	key := fmt.Sprintf("bruteforce:%s", identifier)
	
	// 检查是否已被锁定
	locked, err := config.RDB.Exists(ctx, key+":locked").Result()
	if err == nil && locked > 0 {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"code": 429,
			"msg":  "账户已被锁定，请15分钟后再试",
		})
		c.Abort()
		return true
	}
	
	// 获取失败尝试次数
	attempts, _ := config.RDB.Get(ctx, key).Int()
	
	if attempts >= bf.maxAttempts {
		// 锁定账户
		pipe := config.RDB.Pipeline()
		pipe.Set(ctx, key+":locked", "1", bf.lockout)
		pipe.Del(ctx, key)
		pipe.Exec(ctx)
		
		c.JSON(http.StatusTooManyRequests, gin.H{
			"code": 429,
			"msg":  "登录尝试次数过多，账户已锁定15分钟",
		})
		c.Abort()
		return true
	}
	
	return false
}

// 记录失败尝试
func (bf *BruteForceDetector) RecordFailure(c *gin.Context, identifier string) {
	ctx := context.Background()
	key := fmt.Sprintf("bruteforce:%s", identifier)
	
	pipe := config.RDB.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, bf.window)
	pipe.Exec(ctx)
}

// 清除失败记录（登录成功后调用）
func (bf *BruteForceDetector) ClearFailures(c *gin.Context, identifier string) {
	ctx := context.Background()
	key := fmt.Sprintf("bruteforce:%s", identifier)
	config.RDB.Del(ctx, key, key+":locked")
}