package services

import (
	"context"
	"time"

	"chat-system-pro/config"
	"chat-system-pro/models"
	"chat-system-pro/utils"

	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gorm"
)

type DatabaseService struct {
}

func NewDatabaseService() *DatabaseService {
	return &DatabaseService{}
}

// ClearOldMessages 清除指定日期之前的消息
func (s *DatabaseService) ClearOldMessages(beforeDate time.Time) (int64, error) {
	ctx := context.TODO()
	filter := bson.M{
		"created_at": bson.M{"$lt": beforeDate},
	}

	result, err := config.MongoDBCollection.Collection("messages").DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

// ClearAllData 清空所有数据库数据（危险操作！）
func (s *DatabaseService) ClearAllData() error {
	ctx := context.TODO()

	// 清空 MongoDB 消息
	config.MongoDBCollection.Collection("messages").DeleteMany(ctx, bson.M{})
	config.MongoDBCollection.Collection("message_archived").DeleteMany(ctx, bson.M{})

	// 清空 MySQL 数据（保留表结构）
	config.DB.Exec("TRUNCATE TABLE messages")
	config.DB.Exec("TRUNCATE TABLE conversations")
	config.DB.Exec("TRUNCATE TABLE group_members")
	config.DB.Exec("TRUNCATE TABLE groups")
	config.DB.Exec("TRUNCATE TABLE friends")
	config.DB.Exec("TRUNCATE TABLE payment_orders")
	config.DB.Exec("TRUNCATE TABLE points_history")
	config.DB.Exec("TRUNCATE TABLE admin_logs")
	config.DB.Exec("TRUNCATE TABLE invite_codes")
	config.DB.Exec("TRUNCATE TABLE system_configs")
	
	// 注意：这里不删除 users 表，避免彻底清空所有账号
	// 如果需要清空 users，请单独执行

	return nil
}

// InitializeDatabase 初始化数据库（清空并重新创建表结构，插入默认数据）
func (s *DatabaseService) InitializeDatabase() error {
	// 1. 清空所有数据
	s.ClearAllData()

	// 2. 重新迁移表结构（确保）
	models.AutoMigrate(config.DB)

	// 3. 插入默认系统配置
	defaultConfigs := []models.SystemConfig{
		{Key: "site_name", Value: "Chat Pro", Description: "站点名称"},
		{Key: "ui.default_theme", Value: "modern", Description: "默认UI主题"},
		{Key: "security.invite_code_enabled", Value: "false", Description: "是否启用邀请码"},
		{Key: "security.captcha_enabled", Value: "false", Description: "是否启用验证码"},
		{Key: "security.rate_limit_enabled", Value: "true", Description: "是否启用限流"},
		{Key: "payment.wechat_pay_enabled", Value: "false", Description: "微信支付开关"},
		{Key: "payment.alipay_enabled", Value: "false", Description: "支付宝开关"},
		{Key: "payment.stripe_enabled", Value: "false", Description: "Stripe支付开关"},
	}

	for _, cfg := range defaultConfigs {
		config.DB.Create(&cfg)
	}

	// 4. 创建一个演示账号
	hashedPass, _ := utils.HashPassword("123456")
	demoUser := models.User{
		Username: "admin",
		Password: hashedPass,
		Nickname: "系统管理员",
		Phone: "",
		Email: "admin@example.com",
	}
	config.DB.Create(&demoUser)

	return nil
}

// DeleteUserAndData 删除指定用户及其所有相关数据
func (s *DatabaseService) DeleteUserAndData(userID uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 删除用户相关的好友关系
		tx.Where("user_id = ? OR friend_id = ?", userID, userID).Delete(&models.Friend{})

		// 2. 删除用户创建的群及群成员
		var groups []models.Group
		tx.Where("owner_id = ?", userID).Find(&groups)
		for _, g := range groups {
			tx.Where("group_id = ?", g.ID).Delete(&models.GroupMember{})
			tx.Where("group_id = ?", g.ID).Delete(&models.Conversation{})
		}
		tx.Where("owner_id = ?", userID).Delete(&models.Group{})

		// 3. 删除用户加入的群
		var memberGroupIDs []uint
		tx.Model(&models.GroupMember{}).Where("user_id = ?", userID).Pluck("group_id", &memberGroupIDs)
		tx.Where("user_id = ?", userID).Delete(&models.GroupMember{})
		// 群成员数量减1
		for _, gid := range memberGroupIDs {
			tx.Model(&models.Group{}).Where("id = ?", gid).UpdateColumn("member_count", gorm.Expr("member_count - 1"))
		}

		// 4. 删除用户的会话
		tx.Where("user_id = ?", userID).Delete(&models.Conversation{})

		// 5. 删除用户的支付订单
		tx.Where("user_id = ?", userID).Delete(&models.PaymentOrder{})

		// 6. 删除用户的积分记录
		tx.Where("user_id = ?", userID).Delete(&models.PointsHistory{})

		// 7. 删除用户使用的邀请码
		tx.Where("user_id = ?", userID).Delete(&models.InviteCode{})

		// 8. 从MongoDB删除用户相关的消息
		ctx := context.TODO()
		filter := bson.M{
			"$or": []bson.M{
				{"sender_id": userID},
				{"receiver_id": userID},
			},
		}
		config.MongoDBCollection.Collection("messages").DeleteMany(ctx, filter)

		// 9. 最后删除用户
		if err := tx.Delete(&models.User{}, userID).Error; err != nil {
			return err
		}

		return nil
	})
}

// DeleteGroupAndData 删除指定群及其所有相关数据
func (s *DatabaseService) DeleteGroupAndData(groupID uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 删除群成员
		tx.Where("group_id = ?", groupID).Delete(&models.GroupMember{})

		// 2. 删除群相关的会话
		tx.Where("type = 2 AND target_id = ?", groupID).Delete(&models.Conversation{})

		// 3. 从MongoDB删除群消息
		ctx := context.TODO()
		config.MongoDBCollection.Collection("messages").DeleteMany(ctx, bson.M{"group_id": groupID})

		// 4. 删除群本身
		if err := tx.Delete(&models.Group{}, groupID).Error; err != nil {
			return err
		}

		return nil
	})
}

// ArchiveOldMessages 归档指定日期前的消息到归档集合
func (s *DatabaseService) ArchiveOldMessages(beforeDate time.Time) (int64, error) {
	ctx := context.TODO()
	
	filter := bson.M{
		"created_at": bson.M{"$lt": beforeDate},
	}
	
	// 查找要归档的消息
	cursor, err := config.MongoDBCollection.Collection("messages").Find(ctx, filter)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)
	
	var messagesToArchive []models.MessageDoc
	if err = cursor.All(ctx, &messagesToArchive); err != nil {
		return 0, err
	}
	
	if len(messagesToArchive) == 0 {
		return 0, nil
	}
	
	// 插入到归档集合
	var docsToInsert []interface{}
	for _, m := range messagesToArchive {
		docsToInsert = append(docsToInsert, m)
	}
	
	if _, err = config.MongoDBCollection.Collection("message_archived").InsertMany(ctx, docsToInsert); err != nil {
		return 0, err
	}
	
	// 删除原消息
	result, err := config.MongoDBCollection.Collection("messages").DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	
	return result.DeletedCount, nil
}

// GetDatabaseStats 获取数据库统计信息
func (s *DatabaseService) GetDatabaseStats() (map[string]interface{}, error) {
	ctx := context.TODO()
	
	stats := make(map[string]interface{})
	
	// MySQL统计
	var userCount int64
	config.DB.Model(&models.User{}).Count(&userCount)
	stats["user_count"] = userCount
	
	var groupCount int64
	config.DB.Model(&models.Group{}).Count(&groupCount)
	stats["group_count"] = groupCount
	
	// MongoDB统计
	pipeline := []bson.M{
		{"$collStats": bson.M{"latencyStats": bson.M{"histograms": false}}},
	}
	var mongoStats bson.M
	cursor, _ := config.MongoDBCollection.Collection("messages").Aggregate(ctx, pipeline)
	cursor.Decode(&mongoStats)
	
	// 消息集合数量
	count, _ := config.MongoDBCollection.Collection("messages").CountDocuments(ctx, bson.M{})
	stats["mongo_message_count"] = count
	
	archivedCount, _ := config.MongoDBCollection.Collection("message_archived").CountDocuments(ctx, bson.M{})
	stats["archived_message_count"] = archivedCount
	
	return stats, nil
}
