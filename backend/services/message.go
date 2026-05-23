package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"chat-system-pro/config"
	"chat-system-pro/models"
	"chat-system-pro/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type MessageService struct {
}

func NewMessageService() *MessageService {
	return &MessageService{}
}

// SendMessage 发送消息（高性能版本）
func (s *MessageService) SendMessage(req *SendMessageRequest) (*models.MessageDoc, error) {
	ctx := context.TODO()
	
	now := time.Now()
	
	// 1. 构建消息文档
	message := &models.MessageDoc{
		ID:           primitive.NewObjectID().Hex(),
		SenderID:     req.SenderID,
		ReceiverID:   req.ReceiverID,
		GroupID:      req.GroupID,
		Content:      req.Content,
		MessageType:  req.MessageType,
		MediaURL:     req.MediaURL,
		MediaSize:    req.MediaSize,
		Duration:     req.Duration,
		RedPacketID:  req.RedPacketID,
		IsRecall:     false,
		IsRead:       false,
		ReadUsers:    []uint{req.SenderID},
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// 2. 先保存到MongoDB（高性能写入）
	_, err := config.MongoDBCollection.Collection("messages").InsertOne(ctx, message)
	if err != nil {
		return nil, fmt.Errorf("save message failed: %v", err)
	}

	// 3. 异步处理后续业务（通过消息队列或 goroutine）
	go func() {
		if req.GroupID > 0 {
			// 群聊
			s.updateGroupConversations(req.GroupID, message.ID)
		} else {
			// 私聊
			s.updatePrivateConversations(req.SenderID, req.ReceiverID, message.ID)
		}
	}()

	// 4. 推送到 WebSocket（快速响应）
	broadcastNewMessage(message)

	return message, nil
}

// updatePrivateConversations 更新私聊会话
func (s *MessageService) updatePrivateConversations(senderID, receiverID uint, msgID string) {
	now := time.Now()
	
	// 更新发送者会话
	config.DB.Model(&models.Conversation{}).
		Where("user_id = ? AND type = 1 AND target_id = ?", senderID, receiverID).
		FirstOrCreate(&models.Conversation{
			UserID:        senderID,
			TargetID:      receiverID,
			Type:          1,
			UnreadCount:   0,
			LastMessageAt: &now,
		})

	// 更新接收者会话（增加未读）
	var convReceiver models.Conversation
	config.DB.Where("user_id = ? AND type = 1 AND target_id = ?", receiverID, senderID).First(&convReceiver)
	
	if convReceiver.ID == 0 {
		config.DB.Create(&models.Conversation{
			UserID:        receiverID,
			TargetID:      senderID,
			Type:          1,
			UnreadCount:   1,
			LastMessageAt: &now,
		})
	} else {
		config.DB.Model(&convReceiver).Updates(map[string]interface{}{
			"last_message_at": now,
			"unread_count":    gorm.Expr("unread_count + 1"),
		})
	}
}

// updateGroupConversations 更新群聊会话
func (s *MessageService) updateGroupConversations(groupID uint, msgID string) {
	now := time.Now()
	config.DB.Model(&models.Conversation{}).
		Where("type = 2 AND target_id = ?", groupID).
		Updates(map[string]interface{}{
			"last_message_at": now,
		})
	
	// 增加群成员未读（排除发送者）
	config.DB.Exec(`
		UPDATE conversations 
		SET unread_count = unread_count + 1 
		WHERE type = 2 AND target_id = ? AND user_id != ?`, 
		groupID, msgID)
}

// GetPrivateMessages 获取私聊消息（带分页）
func (s *MessageService) GetPrivateMessages(userID, friendID uint, page, pageSize int, lastID string) ([]models.MessageDoc, error) {
	ctx := context.TODO()
	
	filter := bson.M{
		"$and": []bson.M{
			{"is_recall": false},
			{"$or": []bson.M{
				{"sender_id": userID, "receiver_id": friendID},
				{"sender_id": friendID, "receiver_id": userID},
			}},
		},
	}
	
	// 如果有lastID，只获取比它新的消息
	if lastID != "" {
		oid, _ := primitive.ObjectIDFromHex(lastID)
		filter["_id"] = bson.M{"$gt": oid}
	}
	
	options := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize))
	
	cursor, err := config.MongoDBCollection.Collection("messages").Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var messages []models.MessageDoc
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, err
	}
	
	// 反转顺序
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	
	return messages, nil
}

// GetGroupMessages 获取群消息
func (s *MessageService) GetGroupMessages(groupID uint, page, pageSize int, lastID string) ([]models.MessageDoc, error) {
	ctx := context.TODO()
	
	filter := bson.M{
		"group_id":   groupID,
		"is_recall": false,
	}
	
	if lastID != "" {
		oid, _ := primitive.ObjectIDFromHex(lastID)
		filter["_id"] = bson.M{"$gt": oid}
	}
	
	options := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize))
	
	cursor, err := config.MongoDBCollection.Collection("messages").Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var messages []models.MessageDoc
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, err
	}
	
	// 反转顺序
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	
	return messages, nil
}

// MarkAsRead 标记消息已读
func (s *MessageService) MarkAsRead(userID, targetID uint, convType int) error {
	ctx := context.TODO()
	
	now := time.Now()
	var updatedIDs []string
	
	if convType == 1 {
		// 私聊 - 先查找需要更新的消息
		filter := bson.M{
			"sender_id":   targetID,
			"receiver_id": userID,
			"is_read":     false,
		}
		
		// 查找要更新的消息ID
		cursor, _ := config.MongoDBCollection.Collection("messages").Find(ctx, filter)
		var msgsToUpdate []models.MessageDoc
		cursor.All(ctx, &msgsToUpdate)
		
		for _, msg := range msgsToUpdate {
			updatedIDs = append(updatedIDs, msg.ID)
		}
		
		if len(updatedIDs) > 0 {
			// 执行更新
			update := bson.M{
				"$set": bson.M{
					"is_read":    true,
					"read_at":    now,
					"updated_at": now,
				},
				"$addToSet": bson.M{"read_users": userID},
			}
			
			_, err := config.MongoDBCollection.Collection("messages").UpdateMany(ctx, filter, update)
			if err != nil {
				return err
			}
			
			// 广播已读回执
			go utils.BroadcastReadReceipt(userID, targetID, convType, updatedIDs)
		}
		
		// 清零会话未读数
		config.DB.Model(&models.Conversation{}).
			Where("user_id = ? AND type = 1 AND target_id = ?", userID, targetID).
			Update("unread_count", 0)
		
	} else if convType == 2 {
		// 群聊
		filter := bson.M{
			"group_id":   targetID,
			"read_users": bson.M{"$ne": userID},
		}
		
		// 查找要更新的消息ID
		cursor, _ := config.MongoDBCollection.Collection("messages").Find(ctx, filter)
		var msgsToUpdate []models.MessageDoc
		cursor.All(ctx, &msgsToUpdate)
		
		for _, msg := range msgsToUpdate {
			updatedIDs = append(updatedIDs, msg.ID)
		}
		
		if len(updatedIDs) > 0 {
			// 更新群消息的已读列表
			oidList := make([]primitive.ObjectID, len(updatedIDs))
			for i, id := range updatedIDs {
				oid, _ := primitive.ObjectIDFromHex(id)
				oidList[i] = oid
			}
			
			updateFilter := bson.M{"_id": bson.M{"$in": oidList}}
			update := bson.M{
				"$addToSet": bson.M{"read_users": userID},
				"$set":      bson.M{"updated_at": now},
			}
			
			_, err := config.MongoDBCollection.Collection("messages").UpdateMany(ctx, updateFilter, update)
			if err != nil {
				return err
			}
			
			// 广播已读回执
			go utils.BroadcastReadReceipt(userID, targetID, convType, updatedIDs)
		}
		
		config.DB.Model(&models.Conversation{}).
			Where("user_id = ? AND type = 2 AND target_id = ?", userID, targetID).
			Update("unread_count", 0)
	}
	
	return nil
}

// RecallMessage 撤回消息
func (s *MessageService) RecallMessage(messageID string, userID uint) error {
	ctx := context.TODO()
	
	oid, _ := primitive.ObjectIDFromHex(messageID)
	filter := bson.M{
		"_id":      oid,
		"sender_id": userID,
	}
	
	// 检查是否超过时间（2分钟）
	var msg models.MessageDoc
	config.MongoDBCollection.Collection("messages").FindOne(ctx, filter).Decode(&msg)
	
	if msg.ID == "" {
		return errors.New("message not found")
	}
	
	if time.Since(msg.CreatedAt) > 2*time.Minute {
		return errors.New("can only recall messages within 2 minutes")
	}
	
	update := bson.M{
		"$set": bson.M{
			"is_recall": true,
			"updated_at": time.Now(),
		},
	}
	
	_, err := config.MongoDBCollection.Collection("messages").UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	
	// 通知撤回
	broadcastRecallMessage(messageID)
	
	return nil
}

type SendMessageRequest struct {
	SenderID     uint   `json:"sender_id" binding:"required"`
	ReceiverID   uint   `json:"receiver_id"`
	GroupID      uint   `json:"group_id"`
	Content      string `json:"content" binding:"required"`
	MessageType  int    `json:"message_type"`
	MediaURL     string `json:"media_url"`
	MediaSize    int64  `json:"media_size"`
	Duration     int    `json:"duration"`
	RedPacketID  uint   `json:"red_packet_id"`
}

// broadcastNewMessage 广播新消息（通过WebSocket）
func broadcastNewMessage(msg *models.MessageDoc) {
	// 调用 utils 模块的广播函数
	utils.BroadcastNewMessage(msg)
}

// broadcastRecallMessage 广播撤回消息
func broadcastRecallMessage(msgID string) {
	// 调用 utils 模块的广播函数
	utils.BroadcastRecallMessage(msgID)
}
