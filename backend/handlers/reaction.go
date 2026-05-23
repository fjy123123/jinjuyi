package handlers

import (
	"context"
	"strconv"
	"time"

	"chat-system-pro/config"
	"chat-system-pro/models"
	"chat-system-pro/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

// AddReaction 添加表情反应
func AddReaction(c *gin.Context) {
	userID := c.GetUint("user_id")
	messageID := c.Param("message_id")
	
	var req struct {
		Emoji string `json:"emoji" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "参数错误")
		return
	}

	ctx := context.TODO()
	oid, _ := primitive.ObjectIDFromHex(messageID)
	
	var msg models.MessageDoc
	err := config.MongoDBCollection.Collection("messages").FindOne(ctx, bson.M{"_id": oid}).Decode(&msg)
	if err != nil {
		utils.ErrorResponse(c, 404, "消息不存在")
		return
	}

	filter := bson.M{"_id": oid}
	update := bson.M{
		"$pull": bson.M{
			"reactions": bson.M{
				"user_ids": userID,
			},
		},
	}
	config.MongoDBCollection.Collection("messages").UpdateOne(ctx, filter, update)

	update = bson.M{
		"$push": bson.M{
			"reactions": bson.M{
				"emoji":    req.Emoji,
				"user_ids": userID,
				"count":    1,
			},
		},
	}
	_, err = config.MongoDBCollection.Collection("messages").UpdateOne(ctx, filter, update)
	if err != nil {
		utils.ErrorResponse(c, 500, "添加反应失败")
		return
	}

	config.DB.Create(&models.MessageReactionRecord{
		MessageID: messageID,
		UserID:    userID,
		Emoji:     req.Emoji,
		CreatedAt: time.Now(),
	})

	utils.SendToUser(msg.SenderID, "reaction_added", map[string]interface{}{
		"message_id": messageID,
		"emoji":      req.Emoji,
		"user_id":    userID,
	})

	utils.SuccessResponse(c, nil)
}

// RemoveReaction 移除表情反应
func RemoveReaction(c *gin.Context) {
	userID := c.GetUint("user_id")
	messageID := c.Param("message_id")
	
	var req struct {
		Emoji string `json:"emoji" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "参数错误")
		return
	}

	ctx := context.TODO()
	oid, _ := primitive.ObjectIDFromHex(messageID)

	filter := bson.M{"_id": oid}
	update := bson.M{
		"$pull": bson.M{
			"reactions": bson.M{
				"emoji":    req.Emoji,
				"user_ids": userID,
			},
		},
	}
	_, err := config.MongoDBCollection.Collection("messages").UpdateOne(ctx, filter, update)
	if err != nil {
		utils.ErrorResponse(c, 500, "移除反应失败")
		return
	}

	config.DB.Where("message_id = ? AND user_id = ? AND emoji = ?", messageID, userID, req.Emoji).
		Delete(&models.MessageReactionRecord{})

	utils.SuccessResponse(c, nil)
}

// GetReactions 获取消息的所有反应
func GetReactions(c *gin.Context) {
	messageID := c.Param("message_id")
	
	ctx := context.TODO()
	oid, _ := primitive.ObjectIDFromHex(messageID)
	
	var msg models.MessageDoc
	err := config.MongoDBCollection.Collection("messages").FindOne(ctx, bson.M{"_id": oid}).Decode(&msg)
	if err != nil {
		utils.ErrorResponse(c, 404, "消息不存在")
		return
	}

	utils.SuccessResponse(c, msg.Reactions)
}

// ForwardMessage 转发消息
func ForwardMessage(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var req struct {
		MessageID   string `json:"message_id" binding:"required"`
		TargetID    uint   `json:"target_id" binding:"required"`
		TargetType  int    `json:"target_type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "参数错误")
		return
	}

	ctx := context.TODO()
	oid, _ := primitive.ObjectIDFromHex(req.MessageID)
	
	var originalMsg models.MessageDoc
	err := config.MongoDBCollection.Collection("messages").FindOne(ctx, bson.M{"_id": oid}).Decode(&originalMsg)
	if err != nil {
		utils.ErrorResponse(c, 404, "原消息不存在")
		return
	}

	forwardedMsg := models.MessageDoc{
		ID:             primitive.NewObjectID().Hex(),
		SenderID:       userID,
		Content:        originalMsg.Content,
		MessageType:    9,
		MediaURL:       originalMsg.MediaURL,
		MediaSize:      originalMsg.MediaSize,
		MediaName:      originalMsg.MediaName,
		Duration:       originalMsg.Duration,
		IsForwarded:    true,
		ForwardFromID:  req.MessageID,
		OriginalSender: originalMsg.SenderID,
		ReplyToID:      originalMsg.ReplyToID,
		ReplyToContent: originalMsg.ReplyToContent,
		ReplyToSender:  originalMsg.ReplyToSender,
		IsRead:         false,
		ReadUsers:      []uint{userID},
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if req.TargetType == 1 {
		forwardedMsg.ReceiverID = req.TargetID
	} else {
		forwardedMsg.GroupID = req.TargetID
	}

	_, err = config.MongoDBCollection.Collection("messages").InsertOne(ctx, forwardedMsg)
	if err != nil {
		utils.ErrorResponse(c, 500, "转发失败")
		return
	}

	now := time.Now()
	if req.TargetType == 1 {
		config.DB.Model(&models.Conversation{}).
			Where("user_id = ? AND type = 1 AND target_id = ?", userID, req.TargetID).
			Assign(&models.Conversation{
				UserID: userID, TargetID: req.TargetID, Type: 1,
				LastMessageID: forwardedMsg.ID, LastMessage: forwardedMsg.Content, LastMessageAt: &now,
			}).FirstOrCreate()
		config.DB.Model(&models.Conversation{}).
			Where("user_id = ? AND type = 1 AND target_id = ?", req.TargetID, userID).
			Updates(map[string]interface{}{
				"last_message_at": now, "last_message_id": forwardedMsg.ID,
				"last_message": forwardedMsg.Content, "unread_count": gorm.Expr("unread_count + 1"),
			}).FirstOrCreate(&models.Conversation{
				UserID: req.TargetID, TargetID: userID, Type: 1,
				UnreadCount: 1, LastMessageAt: &now,
			})
		utils.SendToUser(req.TargetID, "new_message", forwardedMsg)
	} else {
		config.DB.Model(&models.Conversation{}).
			Where("type = 2 AND target_id = ?", req.TargetID).
			Updates(map[string]interface{}{
				"last_message_at": now, "last_message_id": forwardedMsg.ID, "last_message": forwardedMsg.Content,
			})
		var memberIDs []uint
		config.DB.Model(&models.GroupMember{}).Where("group_id = ?", req.TargetID).Pluck("user_id", &memberIDs)
		for _, uid := range memberIDs {
			utils.SendToUser(uid, "new_message", forwardedMsg)
		}
	}

	utils.SuccessResponse(c, forwardedMsg)
}

// SearchMessages 搜索消息
func SearchMessages(c *gin.Context) {
	userID := c.GetUint("user_id")
	keyword := c.Query("keyword")
	convType := c.DefaultQuery("type", "0")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if keyword == "" {
		utils.ErrorResponse(c, 400, "搜索关键词不能为空")
		return
	}

	ctx := context.TODO()
	var filter bson.M
	
	if convType == "1" {
		filter = bson.M{
			"is_recall": false,
			"$or": []bson.M{
				{"sender_id": userID, "receiver_id": bson.M{"$ne": 0}},
				{"receiver_id": userID},
			},
			"content": bson.M{"$regex": keyword},
		}
	} else if convType == "2" {
		filter = bson.M{
			"is_recall": false,
			"group_id":   bson.M{"$ne": 0},
			"content":    bson.M{"$regex": keyword},
		}
	} else {
		filter = bson.M{
			"is_recall": false,
			"$or": []bson.M{
				{"sender_id": userID},
				{"receiver_id": userID},
				{"group_id": bson.M{"$ne": 0}},
			},
			"content": bson.M{"$regex": keyword},
		}
	}

	cursor, err := config.MongoDBCollection.Collection("messages").Find(ctx, filter)
	if err != nil {
		utils.ErrorResponse(c, 500, "搜索失败")
		return
	}
	defer cursor.Close(ctx)

	var allMessages []models.MessageDoc
	cursor.All(ctx, &allMessages)

	total := len(allMessages)
	start := int64((page - 1) * pageSize)
	end := start + int64(pageSize)
	
	if int64(total) < end {
		end = int64(total)
	}
	
	if int64(total) <= start {
		utils.SuccessResponse(c, gin.H{"total": total, "messages": []models.MessageDoc{}})
		return
	}

	utils.SuccessResponse(c, gin.H{
		"total":    total,
		"messages": allMessages[start:end],
	})
}

// SetConversationPin 置顶/取消置顶会话
func SetConversationPin(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var req struct {
		TargetID uint `json:"target_id" binding:"required"`
		ConvType int  `json:"conv_type" binding:"required"`
		IsPinned bool `json:"is_pinned"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "参数错误")
		return
	}

	config.DB.Model(&models.Conversation{}).
		Where("user_id = ? AND target_id = ? AND type = ?", userID, req.TargetID, req.ConvType).
		Update("is_pinned", req.IsPinned)

	utils.SuccessResponse(c, nil)
}

// SetConversationMute 免打扰/取消免打扰
func SetConversationMute(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var req struct {
		TargetID uint `json:"target_id" binding:"required"`
		ConvType int  `json:"conv_type" binding:"required"`
		IsMuted  bool `json:"is_muted"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "参数错误")
		return
	}

	config.DB.Model(&models.Conversation{}).
		Where("user_id = ? AND target_id = ? AND type = ?", userID, req.TargetID, req.ConvType).
		Update("is_muted", req.IsMuted)

	utils.SuccessResponse(c, nil)
}

// ArchiveConversation 归档会话
func ArchiveConversation(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var req struct {
		TargetID   uint `json:"target_id" binding:"required"`
		ConvType   int  `json:"conv_type" binding:"required"`
		IsArchived bool `json:"is_archived"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "参数错误")
		return
	}

	config.DB.Model(&models.Conversation{}).
		Where("user_id = ? AND target_id = ? AND type = ?", userID, req.TargetID, req.ConvType).
		Update("is_archived", req.IsArchived)

	utils.SuccessResponse(c, nil)
}

// GetArchivedConversations 获取已归档会话
func GetArchivedConversations(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var convs []models.Conversation
	config.DB.Where("user_id = ? AND is_archived = ?", userID, true).Find(&convs)

	utils.SuccessResponse(c, convs)
}
