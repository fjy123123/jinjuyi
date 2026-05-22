package services

import (
	"chat-system-pro/config"
	"chat-system-pro/models"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type PaymentService struct {
}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

// CreateOrder 创建支付订单
func (s *PaymentService) CreateOrder(userID uint, amount float64, payType int, orderType int) (*models.PaymentOrder, error) {
	// 生成订单号
	orderNo := fmt.Sprintf("CHAT%s%d", time.Now().Format("20060102150405"), userID)
	
	// 计算积分奖励
	pointsAwarded := int64(amount * 100) // 1元 = 100积分
	
	order := &models.PaymentOrder{
		UserID:        userID,
		OrderNo:       orderNo,
		Amount:        amount,
		PointsAwarded: pointsAwarded,
		Type:          orderType,
		PayType:       payType,
		Status:        0,
		CreatedAt:     time.Now(),
	}
	
	if err := config.DB.Create(order).Error; err != nil {
		return nil, err
	}
	
	return order, nil
}

// ProcessStripePayment 处理Stripe支付
func (s *PaymentService) ProcessStripePayment(orderID uint, token string) error {
	// TODO: 实现实际支付处理
	return nil
}

// completeOrder 完成订单，发放积分
func (s *PaymentService) completeOrder(orderID uint) error {
	var order models.PaymentOrder
	if err := config.DB.First(&order, orderID).Error; err != nil {
		return err
	}
	
	if order.Status != 0 {
		return errors.New("order already processed")
	}
	
	return config.DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		
		// 1. 更新订单状态
		tx.Model(&order).Updates(map[string]interface{}{
			"status":    1,
			"pay_time": now,
		})
		
		// 2. 增加用户余额和积分
		var user models.User
		tx.First(&user, order.UserID)
		tx.Model(&user).Updates(map[string]interface{}{
			"balance": gorm.Expr("balance + ?", order.Amount),
			"points": gorm.Expr("points + ?", order.PointsAwarded),
		})
		
		// 3. 记录积分历史
		pointsHistory := &models.PointsHistory{
			UserID:    order.UserID,
			Change:    order.PointsAwarded,
			Balance:   user.Points + order.PointsAwarded,
			Type:      1, // 充值获取
			Remark:    fmt.Sprintf("充值%s元", order.Amount),
			RelatedID: order.ID,
			CreatedAt: now,
		}
		tx.Create(pointsHistory)
		
		return nil
	})
}

// AddPoints 手动添加积分
func (s *PaymentService) AddPoints(userID uint, points int64, remark string) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 更新用户积分
		var user models.User
		tx.First(&user, userID)
		newBalance := user.Points + points
		tx.Model(&user).Update("points", newBalance)
		
		// 记录积分历史
		history := &models.PointsHistory{
			UserID: userID,
			Change: points,
			Balance: newBalance,
			Type: 2,
			Remark: remark,
			CreatedAt: time.Now(),
		}
		tx.Create(history)
		
		return nil
	})
}

// DeductPoints 扣除积分
func (s *PaymentService) DeductPoints(userID uint, points int64, remark string) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var user models.User
		tx.First(&user, userID)
		
		if user.Points < points {
			return errors.New("insufficient points")
		}
		
		newBalance := user.Points - points
		tx.Model(&user).Update("points", newBalance)
		
		history := &models.PointsHistory{
			UserID: userID,
			Change: -points,
			Balance: newBalance,
			Type: 3,
			Remark: remark,
			CreatedAt: time.Now(),
		}
		tx.Create(history)
		
		return nil
	})
}

// GetPointsHistory 获取积分历史
func (s *PaymentService) GetPointsHistory(userID uint, page, pageSize int) ([]models.PointsHistory, int64, error) {
	var history []models.PointsHistory
	var total int64
	
	offset := (page - 1) * pageSize
	
	config.DB.Model(&models.PointsHistory{}).Where("user_id = ?", userID).Count(&total)
	config.DB.Where("user_id = ?", userID).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&history)
	
	return history, total, nil
}

// GetOrders 获取用户订单
func (s *PaymentService) GetOrders(userID uint, page, pageSize int) ([]models.PaymentOrder, int64, error) {
	var orders []models.PaymentOrder
	var total int64
	
	offset := (page - 1) * pageSize
	
	config.DB.Model(&models.PaymentOrder{}).Where("user_id = ?", userID).Count(&total)
	config.DB.Where("user_id = ?", userID).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders)
	
	return orders, total, nil
}
