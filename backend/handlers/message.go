package handlers

import (
	"strconv"

	"chat-system-pro/services"
	"chat-system-pro/utils"

	"github.com/gin-gonic/gin"
)

// SendMessage 发送消息
func SendMessage(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req services.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "参数错误")
		return
	}
	req.SenderID = userID
	msg, err := messageService.SendMessage(&req)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, msg)
}

// GetPrivateMessages 获取私聊消息
func GetPrivateMessages(c *gin.Context) {
	userID := c.GetUint("user_id")
	friendID, _ := strconv.ParseUint(c.Param("friend_id"), 10, 32)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	lastID := c.Query("last_id")
	msgs, err := messageService.GetPrivateMessages(userID, uint(friendID), page, pageSize, lastID)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, msgs)
}

// GetGroupMessages 获取群消息
func GetGroupMessages(c *gin.Context) {
	groupID, _ := strconv.ParseUint(c.Param("group_id"), 10, 32)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	lastID := c.Query("last_id")
	msgs, err := messageService.GetGroupMessages(uint(groupID), page, pageSize, lastID)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, msgs)
}

// MarkAsRead 标记已读
func MarkAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		TargetID uint `json:"target_id" binding:"required"`
		Type     int  `json:"type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "参数错误")
		return
	}
	if err := messageService.MarkAsRead(userID, req.TargetID, req.Type); err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, nil)
}

// RecallMessage 撤回消息
func RecallMessage(c *gin.Context) {
	userID := c.GetUint("user_id")
	messageID := c.Param("message_id")
	if err := messageService.RecallMessage(messageID, userID); err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, nil)
}
