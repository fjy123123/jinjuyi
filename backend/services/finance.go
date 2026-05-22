package services

import (
	"chat-system-pro/backend/config"
	"chat-system-pro/backend/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

// ==================== 充值服务 ====================

// RechargeService 充值服务
type RechargeService struct {
	db *gorm.DB
}

// NewRechargeService 创建充值服务
func NewRechargeService(db *gorm.DB) *RechargeService {
	return &RechargeService{db: db}
}

// CreateRechargeRequest 创建充值申请
func (s *RechargeService) CreateRechargeRequest(userID uint, amount float64, points int64, rechargeType models.RechargeType, paymentImage string, remark string) (*models.RechargeRequest, error) {
	req := &models.RechargeRequest{
		UserID:       userID,
		Amount:       amount,
		Points:       points,
		RechargeType: rechargeType,
		PaymentImage: paymentImage,
		Remark:       remark,
		Status:       models.RechargeStatusPending,
	}

	if err := s.db.Create(req).Error; err != nil {
		return nil, err
	}

	return req, nil
}

// GetMyRechargeRequests 获取我的充值申请
func (s *RechargeService) GetMyRechargeRequests(userID uint, page, pageSize int) ([]models.RechargeRequest, int64, error) {
	var requests []models.RechargeRequest
	var total int64

	offset := (page - 1) * pageSize
	s.db.Model(&models.RechargeRequest{}).Where("user_id = ?", userID).Count(&total)
	s.db.Where("user_id = ?", userID).Preload("Reviewer").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&requests)

	return requests, total, nil
}

// GetRechargeRequest 获取充值申请详情
func (s *RechargeService) GetRechargeRequest(id uint) (*models.RechargeRequest, error) {
	var req models.RechargeRequest
	if err := s.db.Preload("Reviewer").First(&req, id).Error; err != nil {
		return nil, err
	}
	return &req, nil
}

// GetAllRechargeRequests 获取所有充值申请（管理员）
func (s *RechargeService) GetAllRechargeRequests(status models.RechargeStatus, page, pageSize int) ([]models.RechargeRequest, int64, error) {
	var requests []models.RechargeRequest
	var total int64

	query := s.db.Model(&models.RechargeRequest{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	query.Preload("Reviewer").Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&requests)

	return requests, total, nil
}

// ApproveRecharge 审核通过充值
func (s *RechargeService) ApproveRecharge(reviewerID uint, id uint, remark string) (*models.RechargeRequest, error) {
	var req models.RechargeRequest
	if err := s.db.First(&req, id).Error; err != nil {
		return nil, err
	}

	if req.Status != models.RechargeStatusPending {
		return nil, errors.New("该申请已处理")
	}

	now := time.Now()

	// 开启事务
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 更新申请状态
		req.Status = models.RechargeStatusApproved
		req.ReviewerID = &reviewerID
		req.ReviewRemark = remark
		req.ReviewedAt = &now

		if err := tx.Save(&req).Error; err != nil {
			return err
		}

		// 增加用户积分
		if err := tx.Model(&models.User{}).Where("id = ?", req.UserID).UpdateColumn("points", gorm.Expr("points + ?", req.Points)).Error; err != nil {
			return err
		}

		// 记录积分流水
		pointsHistory := &models.PointsHistory{
			UserID:      req.UserID,
			Type:        models.PointsTypeRecharge,
			Amount:      req.Points,
			Description: "充值积分",
			RelatedID:   req.ID,
		}
		if err := tx.Create(pointsHistory).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &req, nil
}

// RejectRecharge 拒绝充值申请
func (s *RechargeService) RejectRecharge(reviewerID uint, id uint, remark string) (*models.RechargeRequest, error) {
	var req models.RechargeRequest
	if err := s.db.First(&req, id).Error; err != nil {
		return nil, err
	}

	if req.Status != models.RechargeStatusPending {
		return nil, errors.New("该申请已处理")
	}

	now := time.Now()

	req.Status = models.RechargeStatusRejected
	req.ReviewerID = &reviewerID
	req.ReviewRemark = remark
	req.ReviewedAt = &now

	if err := s.db.Save(&req).Error; err != nil {
		return nil, err
	}

	return &req, nil
}

// ==================== 提现服务 ====================

// WithdrawService 提现服务
type WithdrawService struct {
	db *gorm.DB
}

// NewWithdrawService 创建提现服务
func NewWithdrawService(db *gorm.DB) *WithdrawService {
	return &WithdrawService{db: db}
}

// CreateWithdrawRequest 创建提现申请
func (s *WithdrawService) CreateWithdrawRequest(userID uint, points int64, amount float64, withdrawType models.WithdrawType, paymentCode string, realName, phone, remark string) (*models.WithdrawRequest, error) {
	// 检查用户积分是否足够
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	if user.Points < points {
		return nil, errors.New("积分不足")
	}

	req := &models.WithdrawRequest{
		UserID:       userID,
		Points:       points,
		Amount:       amount,
		WithdrawType: withdrawType,
		PaymentCode:  paymentCode,
		RealName:     realName,
		Phone:        phone,
		Remark:       remark,
		Status:       models.WithdrawStatusPending,
	}

	// 先冻结积分
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 创建申请
		if err := tx.Create(req).Error; err != nil {
			return err
		}

		// 冻结积分（这里只是记录，实际上等审核通过再扣减）
		// 为了防止重复提现，可以先扣减并记录在流水里，审核失败时加回去

		return nil
	})

	if err != nil {
		return nil, err
	}

	return req, nil
}

// GetMyWithdrawRequests 获取我的提现申请
func (s *WithdrawService) GetMyWithdrawRequests(userID uint, page, pageSize int) ([]models.WithdrawRequest, int64, error) {
	var requests []models.WithdrawRequest
	var total int64

	offset := (page - 1) * pageSize
	s.db.Model(&models.WithdrawRequest{}).Where("user_id = ?", userID).Count(&total)
	s.db.Where("user_id = ?", userID).Preload("Reviewer").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&requests)

	return requests, total, nil
}

// GetWithdrawRequest 获取提现申请详情
func (s *WithdrawService) GetWithdrawRequest(id uint) (*models.WithdrawRequest, error) {
	var req models.WithdrawRequest
	if err := s.db.Preload("Reviewer").First(&req, id).Error; err != nil {
		return nil, err
	}
	return &req, nil
}

// GetAllWithdrawRequests 获取所有提现申请（管理员）
func (s *WithdrawService) GetAllWithdrawRequests(status models.WithdrawStatus, page, pageSize int) ([]models.WithdrawRequest, int64, error) {
	var requests []models.WithdrawRequest
	var total int64

	query := s.db.Model(&models.WithdrawRequest{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	query.Preload("Reviewer").Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&requests)

	return requests, total, nil
}

// ApproveWithdraw 审核通过提现
func (s *WithdrawService) ApproveWithdraw(reviewerID uint, id uint, remark string) (*models.WithdrawRequest, error) {
	var req models.WithdrawRequest
	if err := s.db.First(&req, id).Error; err != nil {
		return nil, err
	}

	if req.Status != models.WithdrawStatusPending {
		return nil, errors.New("该申请已处理")
	}

	now := time.Now()

	// 开启事务
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 更新申请状态
		req.Status = models.WithdrawStatusApproved
		req.ReviewerID = &reviewerID
		req.ReviewRemark = remark
		req.ReviewedAt = &now

		if err := tx.Save(&req).Error; err != nil {
			return err
		}

		// 扣减用户积分
		if err := tx.Model(&models.User{}).Where("id = ?", req.UserID).UpdateColumn("points", gorm.Expr("points - ?", req.Points)).Error; err != nil {
			return err
		}

		// 记录积分流水
		pointsHistory := &models.PointsHistory{
			UserID:      req.UserID,
			Type:        models.PointsTypeWithdraw,
			Amount:      -req.Points, // 负数表示扣减
			Description: "提现扣除积分",
			RelatedID:   req.ID,
		}
		if err := tx.Create(pointsHistory).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &req, nil
}

// RejectWithdraw 拒绝提现申请
func (s *WithdrawService) RejectWithdraw(reviewerID uint, id uint, remark string) (*models.WithdrawRequest, error) {
	var req models.WithdrawRequest
	if err := s.db.First(&req, id).Error; err != nil {
		return nil, err
	}

	if req.Status != models.WithdrawStatusPending {
		return nil, errors.New("该申请已处理")
	}

	now := time.Now()

	req.Status = models.WithdrawStatusRejected
	req.ReviewerID = &reviewerID
	req.ReviewRemark = remark
	req.ReviewedAt = &now

	if err := s.db.Save(&req).Error; err != nil {
		return nil, err
	}

	return &req, nil
}

// ==================== 全局服务实例 ====================

var RechargeService = NewRechargeService(config.DB)
var WithdrawService = NewWithdrawService(config.DB)
