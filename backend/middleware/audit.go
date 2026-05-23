package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// AuditLog 审计日志结构
type AuditLog struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	UserID       uint      `gorm:"index" json:"user_id"`
	Username     string    `gorm:"size:50" json:"username"`
	Action       string    `gorm:"size:100;index" json:"action"`
	Resource     string    `gorm:"size:100" json:"resource"`
	ResourceID   string    `gorm:"size:100" index" json:"resource_id"`
	Method       string    `gorm:"size:10" json:"method"`
	Path         string    `gorm:"size:500" json:"path"`
	Query        string    `gorm:"type:text" json:"query"`
	RequestBody  string    `gorm:"type:text" json:"request_body"`
	ResponseCode int       `gorm:"index" json:"response_code"`
	ResponseMsg  string    `gorm:"size:200" json:"response_msg"`
	IP           string    `gorm:"size:50" index" json:"ip"`
	UserAgent    string    `gorm:"size:500" json:"user_agent"`
	Duration     int64     `json:"duration"`
	Status       string    `gorm:"size:20;index" json:"status"` // success, failure, warning
	CreatedAt    time.Time `json:"created_at"`
}

// 敏感字段正则表达式
var sensitivePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(password|passwd|pwd)`),
	regexp.MustCompile(`(?i)(secret|token|key|api_key|apikey)`),
	regexp.MustCompile(`(?i)(authorization|bearer)`),
	regexp.MustCompile(`(?i)(credit_card|card_number|cvv)`),
	regexp.MustCompile(`(?i)(ssn|social_security)`),
}

// 敏感字段列表
var sensitiveFields = []string{
	"password",
	"passwd",
	"pwd",
	"secret",
	"token",
	"key",
	"api_key",
	"apikey",
	"authorization",
	"old_password",
	"new_password",
	"confirm_password",
}

// AuditLogger 审计日志记录器
type AuditLogger struct {
	// 需要记录的操作类型
	actionPatterns map[string]string
}

// 创建审计日志记录器
func NewAuditLogger() *AuditLogger {
	return &AuditLogger{
		actionPatterns: map[string]string{
			"POST":   "create",
			"PUT":    "update",
			"PATCH":  "update",
			"DELETE": "delete",
			"GET":    "view",
		},
	}
}

// 审计中间件
func (al *AuditLogger) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		startTime := time.Now()
		
		// 包装ResponseWriter以捕获响应状态码
		wrapper := &responseWriterWrapper{
			ResponseWriter: c.Writer,
			statusCode:    200,
		}
		c.Writer = wrapper
		
		// 读取请求体
		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			requestBody = al.sanitizeData(string(bodyBytes))
		}
		
		// 处理请求
		c.Next()
		
		// 计算耗时
		duration := time.Since(startTime).Milliseconds()
		
		// 判断是否需要记录
		if !al.shouldLog(c) {
			return
		}
		
		// 获取用户ID
		userID, _ := c.Get("user_id")
		username := ""
		if u, ok := userID.(uint); ok && u > 0 {
			username = al.getUsername(u)
		}
		
		// 确定操作状态
		status := "success"
		if wrapper.statusCode >= 400 {
			status = "failure"
		} else if wrapper.statusCode >= 300 {
			status = "warning"
		}
		
		// 创建审计日志
		auditLog := AuditLog{
			UserID:       al.uintFromInterface(userID),
			Username:     username,
			Action:       al.getAction(c.Request.Method),
			Resource:     al.getResource(c.FullPath()),
			ResourceID:   c.Param("id"),
			Method:       c.Request.Method,
			Path:         c.Request.URL.Path,
			Query:        c.Request.URL.RawQuery,
			RequestBody:  requestBody,
			ResponseCode: wrapper.statusCode,
			IP:           c.ClientIP(),
			UserAgent:    c.Request.UserAgent(),
			Duration:     duration,
			Status:       status,
			CreatedAt:    time.Now(),
		}
		
		// 异步保存到数据库
		go al.saveLog(auditLog)
	}
}

// 检查是否需要记录日志
func (al *AuditLogger) shouldLog(c *gin.Context) bool {
	// 不记录登录请求的密码
	path := c.FullPath()
	if strings.Contains(path, "/auth/login") || strings.Contains(path, "/auth/register") {
		return true
	}
	
	// 不记录静态资源
	if strings.HasPrefix(path, "/static") || strings.HasPrefix(path, "/public") {
		return false
	}
	
	// 不记录健康检查
	if path == "/health" {
		return false
	}
	
	return true
}

// 获取操作类型
func (al *AuditLogger) getAction(method string) string {
	if action, ok := al.actionPatterns[method]; ok {
		return action
	}
	return "unknown"
}

// 获取资源名称
func (al *AuditLogger) getResource(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 2 {
		return parts[len(parts)-1]
	}
	return "unknown"
}

// 清理敏感数据
func (al *AuditLogger) sanitizeData(data string) string {
	if data == "" || len(data) > 10000 {
		return data
	}
	
	// 尝试解析JSON
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(data), &jsonData); err != nil {
		return data
	}
	
	// 清理敏感字段
	al.sanitizeJSON(jsonData)
	
	// 重新编码
	cleanData, _ := json.Marshal(jsonData)
	return string(cleanData)
}

// 递归清理JSON中的敏感字段
func (al *AuditLogger) sanitizeJSON(data map[string]interface{}) {
	for key, value := range data {
		// 检查是否为敏感字段
		for _, sensitive := range sensitiveFields {
			if strings.ToLower(key) == sensitive {
				data[key] = "[REDACTED]"
				break
			}
		}
		
		// 递归处理嵌套对象
		if nested, ok := value.(map[string]interface{}); ok {
			al.sanitizeJSON(nested)
		}
		
		// 处理数组
		if arr, ok := value.([]interface{}); ok {
			for _, item := range arr {
				if nested, ok := item.(map[string]interface{}); ok {
					al.sanitizeJSON(nested)
				}
			}
		}
	}
}

// 获取用户名（需要从数据库查询）
func (al *AuditLogger) getUsername(userID uint) string {
	// 这里可以实现从缓存或数据库获取用户名
	// 暂时返回空字符串，后续完善
	return ""
}

// uintFromInterface 安全转换interface{}到uint
func (al *AuditLogger) uintFromInterface(v interface{}) uint {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case uint:
		return val
	case float64:
		return uint(val)
	case int:
		return uint(val)
	default:
		return 0
	}
}

// 异步保存日志
func (al *AuditLogger) saveLog(log AuditLog) {
	// 延迟导入以避免循环依赖
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Failed to save audit log: %v\n", r)
		}
	}()
	
	// TODO: 保存到数据库
	// config.DB.Create(&log)
	
	// 打印到控制台（开发环境）
	fmt.Printf("[AUDIT] %s | User: %d | Action: %s | Path: %s | IP: %s | Status: %s\n",
		log.CreatedAt.Format("2006-01-02 15:04:05"),
		log.UserID,
		log.Action,
		log.Path,
		log.IP,
		log.Status,
	)
}

// responseWriterWrapper 包装ResponseWriter以捕获状态码
type responseWriterWrapper struct {
	gin.ResponseWriter
	statusCode int
}

func (w *responseWriterWrapper) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// 记录特定操作的日志
func LogUserAction(userID uint, action, resource string, details map[string]interface{}) {
	log := AuditLog{
		UserID:     userID,
		Action:     action,
		Resource:   resource,
		IP:         "system",
		Status:     "success",
		CreatedAt:  time.Now(),
	}
	
	detailsJSON, _ := json.Marshal(details)
	log.ResponseMsg = string(detailsJSON)
	
	// 异步保存
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Failed to save user action log: %v\n", r)
			}
		}()
		// config.DB.Create(&log)
	}()
}