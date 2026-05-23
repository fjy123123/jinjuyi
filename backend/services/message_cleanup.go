package services

import (
	"context"
	"log"
	"time"

	"chat-system-pro/config"
	"chat-system-pro/models"
	"chat-system-pro/utils"

	"go.mongodb.org/mongo-driver/bson"
)

// MessageCleanupService 消息清理服务
type MessageCleanupService struct{}

// NewMessageCleanupService 创建消息清理服务
func NewMessageCleanupService() *MessageCleanupService {
	return &MessageCleanupService{}
}

// StartCleanupScheduler 启动清理定时任务
func (s *MessageCleanupService) StartCleanupScheduler() {
	// 每30分钟检查一次是否需要清理
	ticker := time.NewTicker(30 * time.Minute)
	go func() {
		for range ticker.C {
			s.checkAndCleanup()
		}
	}()
	
	// 立即执行一次检查
	go s.checkAndCleanup()
}

// checkAndCleanup 检查并执行清理
func (s *MessageCleanupService) checkAndCleanup() {
	var sysConfig models.SystemConfig
	if err := config.DB.First(&sysConfig).Error; err != nil {
		log.Println("获取系统配置失败:", err)
		return
	}

	// 检查是否开启了自动清理
	if !sysConfig.AutoDeleteEnabled {
		return
	}

	// 检查今天是否已经清理过
	if sysConfig.AutoDeleteLastRun != nil {
		lastRun := *sysConfig.AutoDeleteLastRun
		now := time.Now()
		if lastRun.Year() == now.Year() && lastRun.YearDay() == now.YearDay() {
			// 今天已经清理过
			return
		}
	}

	// 执行清理
	s.CleanupOldMessages(sysConfig.AutoDeleteDays)
}

// CleanupOldMessages 清理过期消息
func (s *MessageCleanupService) CleanupOldMessages(days int) {
	ctx := context.Background()
	
	// 计算截止日期
	cutoffTime := time.Now().AddDate(0, 0, -days)
	
	log.Printf("开始清理 %d 天前的消息，截止日期: %v", days, cutoffTime)

	// 构建过滤器 - 删除指定时间之前且未被撤回的消息
	filter := bson.M{
		"created_at": bson.M{"$lt": cutoffTime},
		"is_recall":  false, // 不删除已撤回的消息（保留撤回记录）
	}

	// 获取要删除的消息ID列表（用于通知前端）
	cursor, err := config.MongoDBCollection.Collection("messages").Find(ctx, filter)
	if err != nil {
		log.Println("查询待删除消息失败:", err)
		return
	}
	defer cursor.Close(ctx)

	var messagesToDelete []models.MessageDoc
	if err := cursor.All(ctx, &messagesToDelete); err != nil {
		log.Println("解析待删除消息失败:", err)
		return
	}

	if len(messagesToDelete) == 0 {
		log.Println("没有需要清理的消息")
		s.updateLastRunTime()
		return
	}

	// 收集受影响的用户和群组
	affectedUsers := make(map[uint]bool)
	affectedGroups := make(map[uint]bool)
	
	for _, msg := range messagesToDelete {
		if msg.Type == 1 { // 私聊
			affectedUsers[msg.SenderID] = true
			affectedUsers[msg.ReceiverID] = true
		} else if msg.Type == 2 { // 群聊
			affectedGroups[msg.GroupID] = true
		}
	}

	// 执行删除
	deleteResult, err := config.MongoDBCollection.Collection("messages").DeleteMany(ctx, filter)
	if err != nil {
		log.Println("删除过期消息失败:", err)
		return
	}

	log.Printf("成功清理 %d 条过期消息", deleteResult.DeletedCount)

	// 广播清理通知给所有在线用户
	s.broadcastCleanupNotification(affectedUsers, affectedGroups, cutoffTime)

	// 更新上次清理时间
	s.updateLastRunTime()
}

// broadcastCleanupNotification 广播清理通知
func (s *MessageCleanupService) broadcastCleanupNotification(
	affectedUsers map[uint]bool, 
	affectedGroups map[uint]bool, 
	cutoffTime time.Time,
) {
	notification := map[string]interface{}{
		"type":             "messages_cleaned",
		"cutoff_time":       cutoffTime.Format(time.RFC3339),
		"affected_users":    affectedUsers,
		"affected_groups":   affectedGroups,
	}

	// 广播给所有在线用户
	utils.BroadcastToAll("system_notification", notification)
	
	log.Printf("已广播消息清理通知，影响 %d 个用户，%d 个群组", 
		len(affectedUsers), len(affectedGroups))
}

// updateLastRunTime 更新上次清理时间
func (s *MessageCleanupService) updateLastRunTime() {
	now := time.Now()
	config.DB.Model(&models.SystemConfig{}).Update("auto_delete_last_run", now)
}

// ForceCleanupNow 立即执行清理（管理员调用）
func (s *MessageCleanupService) ForceCleanupNow(days int) error {
	s.CleanupOldMessages(days)
	return nil
}

// GetCleanupStatus 获取清理状态
func (s *MessageCleanupService) GetCleanupStatus() map[string]interface{} {
	var sysConfig models.SystemConfig
	config.DB.First(&sysConfig)

	// 统计即将被清理的消息数量
	ctx := context.Background()
	cutoffTime := time.Now().AddDate(0, 0, -sysConfig.AutoDeleteDays)
	
	filter := bson.M{
		"created_at": bson.M{"$lt": cutoffTime},
		"is_recall":  false,
	}
	
	count, _ := config.MongoDBCollection.Collection("messages").CountDocuments(ctx, filter)

	return map[string]interface{}{
		"auto_delete_enabled": sysConfig.AutoDeleteEnabled,
		"auto_delete_days":    sysConfig.AutoDeleteDays,
		"last_run":            sysConfig.AutoDeleteLastRun,
		"pending_cleanup":     count,
	}
}

// CleanupUserMessages 清理指定用户的消息（当用户被删除时调用）
func (s *MessageCleanupService) CleanupUserMessages(userID uint) error {
	ctx := context.Background()
	
	filter := bson.M{
		"$or": []bson.M{
			{"sender_id": userID},
			{"receiver_id": userID},
		},
	}
	
	_, err := config.MongoDBCollection.Collection("messages").DeleteMany(ctx, filter)
	return err
}

// CleanupGroupMessages 清理指定群组的消息（当群组被删除时调用）
func (s *MessageCleanupService) CleanupGroupMessages(groupID uint) error {
	ctx := context.Background()
	
	filter := bson.M{"group_id": groupID}
	_, err := config.MongoDBCollection.Collection("messages").DeleteMany(ctx, filter)
	return err
}

// ArchiveOldMessages 归档旧消息（可选：将消息移动到归档集合而不是删除）
func (s *MessageCleanupService) ArchiveOldMessages(days int) error {
	ctx := context.Background()
	
	cutoffTime := time.Now().AddDate(0, 0, -days)
	
	filter := bson.M{
		"created_at": bson.M{"$lt": cutoffTime},
		"is_recall":  false,
	}
	
	// 查找要归档的消息
	cursor, err := config.MongoDBCollection.Collection("messages").Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	
	var messages []models.MessageDoc
	if err := cursor.All(ctx, &messages); err != nil {
		return err
	}
	
	if len(messages) == 0 {
		return nil
	}
	
	// 插入到归档集合
	var docs []interface{}
	for _, msg := range messages {
		docs = append(docs, msg)
	}
	
	_, err = config.MongoDBCollection.Collection("messages_archived").InsertMany(ctx, docs)
	if err != nil {
		return err
	}
	
	// 删除原消息
	_, err = config.MongoDBCollection.Collection("messages").DeleteMany(ctx, filter)
	
	log.Printf("成功归档 %d 条消息", len(messages))
	return err
}

// Global cleanup service instance
var MessageCleanup = NewMessageCleanupService()
