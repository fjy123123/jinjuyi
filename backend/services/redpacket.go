package services

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"chat-system-pro/config"
	"chat-system-pro/models"
	"gorm.io/gorm"
)

// RedPacketService 红包服务
type RedPacketService struct {
	db  *gorm.DB
	ctx context.Context
}

// NewRedPacketService 创建红包服务
func NewRedPacketService() *RedPacketService {
	return &RedPacketService{
		db:  config.DB,
		ctx: context.Background(),
	}
}

// SendRedPacketRequest 发送红包请求
type SendRedPacketRequest struct {
	ReceiverID uint    `json:"receiver_id"`  // 私聊红包接收者
	GroupID    uint    `json:"group_id"`     // 群聊红包
	Type       int     `json:"type"`         // 1:普通 2:拼手气
	PayType    int     `json:"pay_type"`     // 1:积分 2:微信 3:支付宝
	Amount     float64 `json:"amount"`       // 总金额
	Count      int     `json:"count"`        // 红包个数
	Greeting   string  `json:"greeting"`     // 祝福语
}

// SendRedPacket 发送红包
func (s *RedPacketService) SendRedPacket(senderID uint, req *SendRedPacketRequest) (*models.RedPacket, error) {
	// 验证参数
	if req.Amount <= 0 {
		return nil, errors.New("金额必须大于0")
	}
	if req.Count <= 0 {
		return nil, errors.New("红包个数必须大于0")
	}
	if req.PayType != 1 && req.PayType != 2 && req.PayType != 3 {
		return nil, errors.New("不支持的支付方式")
	}
	if req.ReceiverID > 0 && req.GroupID > 0 {
		return nil, errors.New("不能同时指定私聊和群聊")
	}
	if req.ReceiverID == 0 && req.GroupID == 0 {
		return nil, errors.New("请指定接收者")
	}
	if req.GroupID > 0 && req.Count <= 1 {
		return nil, errors.New("群红包个数必须大于1")
	}
	if req.ReceiverID > 0 && req.Count > 1 {
		return nil, errors.New("私聊红包只能发1个")
	}

	// 验证发送者余额
	var sender models.User
	if err := s.db.First(&sender, senderID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 扣款
	if req.PayType == 1 {
		// 积分支付
		if sender.Points < int64(req.Amount) {
			tx.Rollback()
			return nil, errors.New("积分不足")
		}
		// 扣除积分
		if err := tx.Model(&sender).UpdateColumn("points", gorm.Expr("points - ?", int64(req.Amount))).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		// 记录积分变动
		pointsHistory := models.PointsHistory{
			UserID:    senderID,
			Change:    -int64(req.Amount),
			Balance:   sender.Points - int64(req.Amount),
			Type:      3, // 发红包
			Remark:    "发红包",
		}
		if err := tx.Create(&pointsHistory).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	} else {
		// 微信/支付宝，这里只做占位，实际需要对接第三方支付
		// 实际场景中应该先创建支付订单，用户完成支付后再创建红包
		// 这里简化处理，假设已支付
	}

	// 创建红包
	expireAt := time.Now().Add(24 * time.Hour) // 24小时过期
	redPacket := models.RedPacket{
		SenderID:     senderID,
		ReceiverID:   req.ReceiverID,
		GroupID:      req.GroupID,
		Type:           req.Type,
		PayType:      req.PayType,
		Amount:        req.Amount,
		TotalCount:  req.Count,
		Greeting:    req.Greeting,
		ExpireAt:    &expireAt,
		Status:       0,
	}

	if err := tx.Create(&redPacket).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 如果是拼手气红包，预分配金额
	if req.Type == 2 && req.GroupID > 0 && req.Count > 1 {
		if err := s.preallocateRedPacket(tx, &redPacket); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()

	return &redPacket, nil
}

// preallocateRedPacket 预分配拼手气红包金额
func (s *RedPacketService) preallocateRedPacket(tx *gorm.DB, redPacket *models.RedPacket) error {
	// 使用二倍均值法分配红包
	remainingAmount := redPacket.Amount
	remainingCount := redPacket.TotalCount

	amounts := make([]float64, redPacket.TotalCount)
	for i := 0; i < redPacket.TotalCount-1; i++ {
		// 二倍均值：[0.01, (剩余金额/剩余人数)*2]
		maxAmount := (remainingAmount / float64(remainingCount)) * 2
		amount := rand.Float64()*(maxAmount-0.01) + 0.01
		amount = float64(int(amount*100)) / 100 // 保留2位小数
		amounts[i] = amount
		remainingAmount -= amount
		remainingCount--
	}
	amounts[redPacket.TotalCount-1] = remainingAmount

	// 存入Redis，便于快速领取
	key := fmt.Sprintf("red_packet:prealloc:%d", redPacket.ID)
	for _, amount := range amounts {
		if err := config.RDB.LPush(s.ctx, key, amount).Err(); err != nil {
			return err
		}
	}
	// 设置过期时间24小时
	config.RDB.Expire(s.ctx, key, 24*time.Hour)

	return nil
}

// GrabRedPacket 抢红包
func (s *RedPacketService) GrabRedPacket(userID uint, redPacketID uint) (*models.RedPacketDetail, error) {
	// 获取红包信息
	var redPacket models.RedPacket
	if err := s.db.Preload("Sender").First(&redPacket, redPacketID).Error; err != nil {
		return nil, errors.New("红包不存在")
	}

	// 验证红包状态
	if redPacket.Status != 0 {
		return nil, errors.New("红包已抢完或已过期")
	}
	if redPacket.ExpireAt != nil && time.Now().After(*redPacket.ExpireAt) {
		// 标记过期
		s.db.Model(&redPacket).Update("status", 2)
		return nil, errors.New("红包已过期")
	}

	// 验证是否已抢过
	var existingDetail models.RedPacketDetail
	if s.db.Where("red_packet_id = ? AND user_id = ?", redPacketID, userID).First(&existingDetail).Error == nil {
		return &existingDetail, nil // 已经抢过了
	}

	// 如果是群红包，验证是否在群里
	if redPacket.GroupID > 0 {
		var groupMember models.GroupMember
		if s.db.Where("group_id = ? AND user_id = ?", redPacket.GroupID, userID).First(&groupMember).Error != nil {
			return nil, errors.New("不在该群，无法领取")
		}
	}

	// 如果是私聊红包，验证是否是接收者
	if redPacket.ReceiverID > 0 && redPacket.ReceiverID != userID {
		return nil, errors.New("这个红包不是发给你的")
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var amount float64
	if redPacket.Type == 2 && redPacket.GroupID > 0 && redPacket.TotalCount > 1 {
		// 拼手气红包，从Redis取预分配的金额
		key := fmt.Sprintf("red_packet:prealloc:%d", redPacketID)
		result, err := config.RDB.RPop(s.ctx, key).Result()
		if err != nil {
			tx.Rollback()
			return nil, errors.New("红包已抢完")
		}
		_, err = fmt.Sscanf(result, "%f", &amount)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	} else {
		// 普通红包，平均分配
		amount = redPacket.Amount / float64(redPacket.TotalCount)
		amount = float64(int(amount*100)) / 100
	}

	// 创建领取记录
	detail := models.RedPacketDetail{
		RedPacketID: redPacketID,
		UserID:      userID,
		Amount:      amount,
	}
	if err := tx.Create(&detail).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 更新红包状态
	newReceivedCount := redPacket.ReceivedCount + 1
	newReceivedAmount := redPacket.ReceivedAmount + amount

	updates := map[string]interface{}{
		"received_count": newReceivedCount,
		"received_amount": newReceivedAmount,
	}

	if newReceivedCount >= redPacket.TotalCount {
		updates["status"] = 1 // 已抢完
	}

	if err := tx.Model(&redPacket).Updates(updates).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 给用户加钱/积分
	var receiver models.User
	if err := tx.First(&receiver, userID).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("用户不存在")
	}

	if redPacket.PayType == 1 {
		// 积分
		if err := tx.Model(&receiver).UpdateColumn("points", gorm.Expr("points + ?", int64(amount))).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		// 记录积分变动
		pointsHistory := models.PointsHistory{
			UserID:    userID,
			Change:    int64(amount),
			Balance:   receiver.Points + int64(amount),
			Type:      4, // 抢红包
			Remark:    "抢红包",
		}
		if err := tx.Create(&pointsHistory).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	} else {
		// 现金（这里只做占位）
	}

	// 检查是否是手气最佳
	if redPacket.Type == 2 && newReceivedCount >= redPacket.TotalCount {
		s.checkAndUpdateBestLucky(tx, redPacketID)
	}

	tx.Commit()

	// 预加载用户信息
	if err := s.db.Preload("User").First(&detail, detail.ID).Error; err == nil {
	}

	return &detail, nil
}

// checkAndUpdateBestLucky 检查并更新手气最佳
func (s *RedPacketService) checkAndUpdateBestLucky(tx *gorm.DB, redPacketID uint) {
	var details []models.RedPacketDetail
	if err := tx.Where("red_packet_id = ?", redPacketID).Order("amount desc").Find(&details).Error; err == nil {
		if len(details) > 0 {
			tx.Model(&details[0]).Update("is_best", true)
		}
	}
}

// GetRedPacket 获取红包信息
func (s *RedPacketService) GetRedPacket(userID uint, redPacketID uint) (*models.RedPacket, []models.RedPacketDetail, error) {
	var redPacket models.RedPacket
	if err := s.db.Preload("Sender").First(&redPacket, redPacketID).Error; err != nil {
		return nil, nil, errors.New("红包不存在")
	}

	// 获取领取记录
	var details []models.RedPacketDetail
	if err := s.db.Preload("User").Where("red_packet_id = ?", redPacketID).Order("created_at asc").Find(&details).Error; err != nil {
		return &redPacket, nil, err
	}

	return &redPacket, details, nil
}

// GetRedPacketList 获取我发出的红包列表
func (s *RedPacketService) GetRedPacketList(userID uint, page, pageSize int) ([]models.RedPacket, int64, error) {
	var list []models.RedPacket
	var total int64

	offset := (page - 1) * pageSize
	if err := s.db.Model(&models.RedPacket{}).Where("sender_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := s.db.Preload("Sender").Where("sender_id = ?", userID).Order("created_at desc").Limit(pageSize).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// GetMyRedPacketRecords 获取我收到的红包记录
func (s *RedPacketService) GetMyRedPacketRecords(userID uint, page, pageSize int) ([]models.RedPacketDetail, int64, error) {
	var list []models.RedPacketDetail
	var total int64

	offset := (page - 1) * pageSize
	if err := s.db.Model(&models.RedPacketDetail{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := s.db.Preload("User").Where("user_id = ?", userID).Order("created_at desc").Limit(pageSize).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
