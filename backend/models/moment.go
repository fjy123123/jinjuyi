package models

import (
	"time"
	"gorm.io/gorm"
)

// ==================== 朋友圈数据模型 ====================

// Moment 朋友圈
type Moment struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	UserID      uint           `gorm:"index" json:"user_id"`
	Content     string         `gorm:"type:text" json:"content"`
	Images      string         `gorm:"type:text" json:"images"` // JSON array
	Location    string         `gorm:"size:200" json:"location"`
	ViewScope   int            `gorm:"default:0" json:"view_scope"` // 0:所有人,1:仅好友,2:指定可见
	VisibleTo   string         `gorm:"type:text" json:"visible_to"` // JSON array of user_ids
	HideFrom    string         `gorm:"type:text" json:"hide_from"` // JSON array of user_ids
	LikesCount  int            `gorm:"default:0" json:"likes_count"`
	CommentCount int           `gorm:"default:0" json:"comment_count"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	User        *User          `gorm:"foreignKey:UserID" json:"user"`
}

// MomentLike 朋友圈点赞
type MomentLike struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	MomentID  uint      `gorm:"index" json:"moment_id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`

	// 关联
	User      *User     `gorm:"foreignKey:UserID" json:"user"`
}

// MomentComment 朋友圈评论
type MomentComment struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	MomentID    uint      `gorm:"index" json:"moment_id"`
	UserID      uint      `gorm:"index" json:"user_id"`
	ReplyToUser uint      `gorm:"index;default:0" json:"reply_to_user"` // 回复谁
	Content     string    `gorm:"type:text" json:"content"`
	CreatedAt   time.Time `json:"created_at"`

	// 关联
	User        *User     `gorm:"foreignKey:UserID" json:"user"`
	ReplyUser   *User     `gorm:"foreignKey:ReplyToUser" json:"reply_user"`
}

// ==================== 微信小程序模型 ====================

// MiniProgramUser 小程序用户
type MiniProgramUser struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	UserID     uint           `gorm:"index" json:"user_id"`
	OpenID     string         `gorm:"uniqueIndex;size:64" json:"open_id"`
	UnionID    string         `gorm:"index;size:64" json:"union_id"`
	SessionKey string         `gorm:"size:64" json:"-"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

// ==================== 设备管理 ====================

// Device 设备
type Device struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	UserID      uint           `gorm:"index" json:"user_id"`
	DeviceID    string         `gorm:"uniqueIndex;size:64" json:"device_id"`
	DeviceType  string         `gorm:"size:20" json:"device_type"` // web, ios, android, mini
	DeviceName  string         `gorm:"size:100" json:"device_name"`
	LastActive  time.Time      `json:"last_active"`
	IsOnline    bool           `gorm:"default:false" json:"is_online"`
	PushToken   string         `gorm:"size:128" json:"push_token"` // 推送token
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
