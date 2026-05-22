package handlers

import (
	"chat-system-pro/config"
	"chat-system-pro/models"
	"chat-system-pro/utils"

	"github.com/gin-gonic/gin"
)

// GetConversations 获取会话列表
func GetConversations(c *gin.Context) {
	userID := c.GetUint("user_id")
	var conversations []models.Conversation
	config.DB.Where("user_id = ?", userID).Order("last_message_at DESC").Find(&conversations)
	utils.SuccessResponse(c, conversations)
}

// GetUnreadCount 获取未读消息总数
func GetUnreadCount(c *gin.Context) {
	userID := c.GetUint("user_id")
	var total int64
	config.DB.Model(&models.Conversation{}).Where("user_id = ?", userID).Select("COALESCE(SUM(unread_count), 0)").Scan(&total)
	utils.SuccessResponse(c, gin.H{"unread_count": total})
}
