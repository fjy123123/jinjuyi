package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": data})
}

func ErrorResponse(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, gin.H{"code": code, "msg": msg, "data": nil})
}

func PaginatedResponse(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0, "msg": "success",
		"data": gin.H{
			"list":      list,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
