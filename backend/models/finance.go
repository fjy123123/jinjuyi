package models

import (
	"time"

	"gorm.io/gorm"
)

// ==================== 充值申请模型 ====================

// RechargeType 充值方式
type RechargeType string

const (
	RechargeTypePersonal RechargeType = "personal" // 个人收款码
	RechargeTypeCompany  RechargeType = "company"  // 公司收款码
)

// RechargeStatus 充值状态
type RechargeStatus string

const (
	RechargeStatusPending  RechargeStatus = "pending"  // 待审核
	RechargeStatusApproved RechargeStatus = "approved" // 已通过
	RechargeStatusRejected RechargeStatus = "rejected" // 已拒绝
)

// RechargeRequest 充值申请
type RechargeRequest struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	UserID         uint           `gorm:"index;not null" json:"user_id"`
	Amount         float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	Points         int64          `gorm:"not null" json:"points"` // 兑换的积分数量
	RechargeType   RechargeType   `gorm:"size:20;default:'personal'" json:"recharge_type"`
	PaymentImage   string         `gorm:"size:500;not null" json:"payment_image"` // 付款截图
	Remark         string         `gorm:"size:500" json:"remark"`                // 用户备注
	Status         RechargeStatus `gorm:"size:20;default:'pending'" json:"status"`
	ReviewerID     *uint          `gorm:"index" json:"reviewer_id,omitempty"`     // 审核人ID
	Reviewer       *User          `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
	ReviewRemark   string         `gorm:"size:500" json:"review_remark"` // 审核备注
	ReviewedAt     *time.Time     `json:"reviewed_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// ==================== 提现申请模型 ====================

// WithdrawType 提现方式
type WithdrawType string

const (
	WithdrawTypePersonal WithdrawType = "personal" // 个人收款码
	WithdrawTypeCompany  WithdrawType = "company"  // 公司收款码
)

// WithdrawStatus 提现状态
type WithdrawStatus string

const (
	WithdrawStatusPending  WithdrawStatus = "pending"  // 待审核
	WithdrawStatusApproved WithdrawStatus = "approved" // 已通过
	WithdrawStatusRejected WithdrawStatus = "rejected" // 已拒绝
)

// WithdrawRequest 提现申请
type WithdrawRequest struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	UserID         uint           `gorm:"index;not null" json:"user_id"`
	Points         int64          `gorm:"not null" json:"points"` // 提现的积分数量
	Amount         float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	WithdrawType   WithdrawType   `gorm:"size:20;default:'personal'" json:"withdraw_type"`
	PaymentCode    string         `gorm:"size:500;not null" json:"payment_code"` // 收款码
	RealName       string         `gorm:"size:50" json:"real_name"`           // 真实姓名
	Phone          string         `gorm:"size:20" json:"phone"`               // 手机号
	Remark         string         `gorm:"size:500" json:"remark"`             // 用户备注
	Status         WithdrawStatus `gorm:"size:20;default:'pending'" json:"status"`
	ReviewerID     *uint          `gorm:"index" json:"reviewer_id,omitempty"`     // 审核人ID
	Reviewer       *User          `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
	ReviewRemark   string         `gorm:"size:500" json:"review_remark"` // 审核备注
	ReviewedAt     *time.Time     `json:"reviewed_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
