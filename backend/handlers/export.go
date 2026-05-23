package handlers

import (
	"context"
	"encoding/csv"
	"strconv"
	"time"

	"chat-system-pro/config"
	"chat-system-pro/models"
	"chat-system-pro/services"
	"chat-system-pro/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// ExportChatHistory 导出聊天记录（CSV格式）
func ExportChatHistory(c *gin.Context) {
	// 检查导出功能是否开启
	if !services.SystemConfigService.IsExportEnabled() {
		utils.ErrorResponse(c, 403, "聊天记录导出功能已关闭")
		return
	}

	userID := c.GetUint("user_id")
	friendIDStr := c.Query("friend_id")
	groupIDStr := c.Query("group_id")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var filter bson.M = bson.M{}
	
	if friendIDStr != "" {
		friendID, _ := strconv.ParseUint(friendIDStr, 10, 32)
		filter["$or"] = []bson.M{
			{"sender_id": userID, "receiver_id": uint(friendID)},
			{"sender_id": uint(friendID), "receiver_id": userID},
		}
		filter["type"] = 1
	} else if groupIDStr != "" {
		groupID, _ := strconv.ParseUint(groupIDStr, 10, 32)
		filter["group_id"] = uint(groupID)
		filter["type"] = 2
	} else {
		filter["$or"] = []bson.M{
			{"sender_id": userID},
			{"receiver_id": userID},
		}
	}

	if startDateStr != "" {
		startDate, _ := time.Parse("2006-01-02", startDateStr)
		if !startDate.IsZero() {
			if filter["created_at"] == nil {
				filter["created_at"] = bson.M{}
			}
			filter["created_at"].(bson.M)["$gte"] = startDate
		}
	}

	if endDateStr != "" {
		endDate, _ := time.Parse("2006-01-02", endDateStr)
		if !endDate.IsZero() {
			endDate = endDate.Add(24 * time.Hour)
			if filter["created_at"] == nil {
				filter["created_at"] = bson.M{}
			}
			filter["created_at"].(bson.M)["$lt"] = endDate
		}
	}

	ctx := context.TODO()
	
	// 获取最大导出记录数
	maxRecords := services.SystemConfigService.GetExportMaxRecords()
	
	cursor, err := config.MongoDBCollection.Collection("messages").Find(ctx, filter)
	if err != nil {
		utils.ErrorResponse(c, 500, "获取消息失败")
		return
	}
	defer cursor.Close(ctx)

	var messages []models.MessageDoc
	if err = cursor.All(ctx, &messages); err != nil {
		utils.ErrorResponse(c, 500, "解析消息失败")
		return
	}

	// 限制导出记录数
	exportMessages := messages
	if len(messages) > maxRecords {
		exportMessages = messages[:maxRecords]
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=chat_history_"+time.Now().Format("20060102")+".csv")
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	writer.Write([]string{"时间", "发送者ID", "接收者/群ID", "消息类型", "内容", "是否已读"})

	for _, msg := range exportMessages {
		var msgType = "文本"
		if msg.MessageType == 2 {
			msgType = "图片"
		} else if msg.MessageType == 3 {
			msgType = "文件"
		} else if msg.MessageType == 4 {
			msgType = "语音"
		} else if msg.MessageType == 5 {
			msgType = "表情包"
		}

		readStatus := "否"
		if msg.IsRead {
			readStatus = "是"
		}

		targetID := ""
		if msg.Type == 1 {
			targetID = strconv.FormatUint(uint64(msg.ReceiverID), 10)
		} else {
			targetID = strconv.FormatUint(uint64(msg.GroupID), 10)
		}

		writer.Write([]string{
			msg.CreatedAt.Format("2006-01-02 15:04:05"),
			strconv.FormatUint(uint64(msg.SenderID), 10),
			targetID,
			msgType,
			msg.Content,
			readStatus,
		})
	}

	writer.Flush()
}
