package handlers

import (
	"strconv"

	"chat-system-pro/services"
	"chat-system-pro/utils"

	"github.com/gin-gonic/gin"
)

// SendRedPacket 发送红包
func SendRedPacket(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req services.SendRedPacketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "参数错误: "+err.Error())
		return
	}

	redPacket, err := redPacketService.SendRedPacket(userID, &req)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	// 如果是私聊红包，自动发送一条红包消息
	if req.ReceiverID > 0 {
		messageReq := &services.SendMessageRequest{
			SenderID:    userID,
			ReceiverID: req.ReceiverID,
			Content:     "红包",
			MessageType: 6,  // 红包类型
			RedPacketID: redPacket.ID,
		}
		messageService.SendMessage(messageReq)
	}

	// 如果是群红包，自动发送一条红包消息
	if req.GroupID > 0 {
		messageReq := &services.SendMessageRequest{
			SenderID:    userID,
			GroupID:     req.GroupID,
			Content:     "红包",
			MessageType: 6,  // 红包类型
			RedPacketID: redPacket.ID,
		}
		messageService.SendMessage(messageReq)
	}

	utils.SuccessResponse(c, redPacket)
}

// GrabRedPacket 抢红包
func GrabRedPacket(c *gin.Context) {
	userID := c.GetUint("user_id")
	redPacketID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	detail, err := redPacketService.GrabRedPacket(userID, uint(redPacketID))
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, detail)
}

// GetRedPacket 获取红包详情
func GetRedPacket(c *gin.Context) {
	userID := c.GetUint("user_id")
	redPacketID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	redPacket, details, err := redPacketService.GetRedPacket(userID, uint(redPacketID))
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"red_packet": redPacket,
		"details":    details,
	})
}

// GetSentRedPackets 获取我发出的红包
func GetSentRedPackets(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	list, total, err := redPacketService.GetRedPacketList(userID, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, 500, "获取失败")
		return
	}

	utils.PaginatedResponse(c, list, total, page, pageSize)
}

// GetReceivedRedPackets 获取我收到的红包
func GetReceivedRedPackets(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	list, total, err := redPacketService.GetMyRedPacketRecords(userID, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, 500, "获取失败")
		return
	}

	utils.PaginatedResponse(c, list, total, page, pageSize)
}
