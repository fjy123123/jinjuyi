package models

import (
	"time"

	"gorm.io/gorm"
)

// ==================== MySQL Models ====================

// User 用户模型
type User struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Username    string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password    string         `gorm:"size:255;not null" json:"-"`
	Nickname    string         `gorm:"size:100" json:"nickname"`
	Avatar      string         `gorm:"size:255" json:"avatar"`
	Phone       string         `gorm:"size:20" json:"phone"`
	Email       string         `gorm:"size:100" json:"email"`
	Gender      int            `gorm:"default:0" json:"gender"` // 0:未知 1:男 2:女
	Birthday    *time.Time     `json:"birthday"`
	Region      string         `gorm:"size:200" json:"region"`
	Sign        string         `gorm:"size:500" json:"sign"`
	Status      int            `gorm:"default:0" json:"status"`
	Role        int            `gorm:"default:0" json:"role"` // 0:普通用户 1:管理员 2:超级管理员
	Balance     float64        `gorm:"default:0" json:"balance"`
	Points      int64          `gorm:"default:0" json:"points"`
	VIPLevel    int            `gorm:"default:0" json:"vip_level"`
	LastLoginAt *time.Time     `json:"last_login_at"`
	LastLoginIP string         `gorm:"size:50" json:"last_login_ip"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// UserSettings 用户设置
type UserSettings struct {
	ID              uint      `gorm:"primarykey" json:"id"`
	UserID          uint      `gorm:"uniqueIndex" json:"user_id"`
	NewMsgNotify    bool      `gorm:"default:true" json:"new_msg_notify"`
	SoundNotify     bool      `gorm:"default:true" json:"sound_notify"`
	AddFriendConfirm bool     `gorm:"default:true" json:"add_friend_confirm"`
	ShowOnline      bool      `gorm:"default:true" json:"show_online"`
	ShowReadReceipt bool      `gorm:"default:true" json:"show_read_receipt"`
	Theme           string    `gorm:"size:50;default:modern" json:"theme"`
	Language        string    `gorm:"size:20;default:zh-CN" json:"language"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Friend 好友关系
type Friend struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	FriendID  uint      `gorm:"index;not null" json:"friend_id"`
	Remark    string    `gorm:"size:100" json:"remark"`
	Status    int       `gorm:"default:0" json:"status"` // 0:正常 1:待验证
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User   User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Friend User `gorm:"foreignKey:FriendID" json:"friend,omitempty"`
}

// Group 群组
type Group struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Avatar      string         `gorm:"size:255" json:"avatar"`
	Description string         `gorm:"size:500" json:"description"`
	Announcement string        `gorm:"type:text" json:"announcement"`
	OwnerID     uint           `gorm:"not null" json:"owner_id"`
	MemberCount int            `gorm:"default:0" json:"member_count"`
	MaxMembers  int            `gorm:"default:500" json:"max_members"`
	JoinMode    int            `gorm:"default:0" json:"join_mode"` // 0:自由 1:验证 2:禁止
	IsMuteAll   bool           `gorm:"default:false" json:"is_mute_all"`
	AllowInvite bool           `gorm:"default:true" json:"allow_invite"`
	ShowMember  bool           `gorm:"default:true" json:"show_member"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Owner User `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
}

// GroupMember 群成员
type GroupMember struct {
	ID         uint       `gorm:"primarykey" json:"id"`
	GroupID    uint       `gorm:"index;not null" json:"group_id"`
	UserID     uint       `gorm:"index;not null" json:"user_id"`
	Nickname   string     `gorm:"size:100" json:"nickname"`
	Role       int        `gorm:"default:0" json:"role"` // 0:成员 1:管理员 2:群主
	IsMute     bool       `gorm:"default:false" json:"is_mute"`
	JoinAt     time.Time  `json:"join_at"`
	LastReadAt *time.Time `json:"last_read_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`

	Group Group `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// GroupInvite 群邀请
type GroupInvite struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	GroupID   uint      `gorm:"index" json:"group_id"`
	InviterID uint      `gorm:"index" json:"inviter_id"`
	InviteeID uint      `gorm:"index" json:"invitee_id"`
	Status    int       `gorm:"default:0" json:"status"` // 0:待处理 1:已接受 2:已拒绝
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GroupJoinRequest 入群申请
type GroupJoinRequest struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	GroupID   uint      `gorm:"index" json:"group_id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Remark    string    `gorm:"size:200" json:"remark"`
	Status    int       `gorm:"default:0" json:"status"` // 0:待处理 1:通过 2:拒绝
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Conversation 会话
type Conversation struct {
	ID            uint       `gorm:"primarykey" json:"id"`
	UserID        uint       `gorm:"index;not null" json:"user_id"`
	Type          int        `gorm:"not null" json:"type"` // 1:私聊 2:群聊
	TargetID      uint       `gorm:"not null" json:"target_id"`
	UnreadCount   int        `gorm:"default:0" json:"unread_count"`
	LastMessageAt *time.Time `json:"last_message_at"`
	TopAt         *time.Time `json:"top_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// InviteCode 邀请码
type InviteCode struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Code      string    `gorm:"uniqueIndex;size:50" json:"code"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Status    int       `gorm:"default:0" json:"status"`
	UsedCount int       `gorm:"default:0" json:"used_count"`
	MaxCount  int       `gorm:"default:1" json:"max_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PaymentOrder 支付订单
type PaymentOrder struct {
	ID            uint       `gorm:"primarykey" json:"id"`
	UserID        uint       `gorm:"index" json:"user_id"`
	OrderNo       string     `gorm:"uniqueIndex;size:50" json:"order_no"`
	Amount        float64    `json:"amount"`
	PointsAwarded int64      `json:"points_awarded"`
	Type          int        `json:"type"`
	Status        int        `gorm:"default:0" json:"status"` // 0:待支付 1:已支付 2:已退款
	PayType       int        `json:"pay_type"`
	PayTime       *time.Time `json:"pay_time"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// PointsHistory 积分记录
// PointsType 积分变动类型
type PointsType string

const (
	PointsTypeRecharge PointsType = "recharge" // 充值
	PointsTypeWithdraw PointsType = "withdraw" // 提现
	PointsTypeConsume  PointsType = "consume"  // 消费
	PointsTypeReward   PointsType = "reward"   // 奖励
)

// PointsHistory 积分流水记录
type PointsHistory struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	UserID      uint       `gorm:"index" json:"user_id"`
	Amount      int64      `json:"amount"` // 变动数量，正数增加，负数减少
	Type        PointsType `gorm:"size:20" json:"type"`
	Description string     `gorm:"size:500" json:"description"`
	RelatedID   uint       `json:"related_id"`
	CreatedAt   time.Time  `json:"created_at"`
}

// SystemConfig 系统配置
type SystemConfig struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Key         string    `gorm:"uniqueIndex;size:100" json:"key"`
	Value       string    `gorm:"type:text" json:"value"`
	Description string    `gorm:"size:255" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AdminLog 管理员日志
type AdminLog struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	AdminID   uint      `gorm:"index" json:"admin_id"`
	Action    string    `gorm:"size:100" json:"action"`
	IP        string    `gorm:"size:50" json:"ip"`
	UserAgent string    `gorm:"size:255" json:"user_agent"`
	Params    string    `gorm:"type:text" json:"params"`
	CreatedAt time.Time `json:"created_at"`
}

// ==================== 红包相关 ====================

// RedPacket 红包
type RedPacket struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	SenderID     uint           `gorm:"index" json:"sender_id"`
	ReceiverID   uint           `gorm:"index" json:"receiver_id,omitempty"`  // 0 表示群红包
	GroupID      uint           `gorm:"index" json:"group_id,omitempty"`
	Type           int            `gorm:"default:1" json:"type"`  // 1:普通红包 2:拼手气红包
	PayType      int            `gorm:"default:1" json:"pay_type"`  // 1:积分 2:微信 3:支付宝
	Amount        float64         `json:"amount"`
	TotalCount  int            `gorm:"default:1" json:"total_count"`
	ReceivedCount int            `gorm:"default:0" json:"received_count"`
	ReceivedAmount float64        `gorm:"default:0" json:"received_amount"`
	Greeting    string         `gorm:"size:200" json:"greeting"`
	Status       int            `gorm:"default:0" json:"status"`  // 0:进行中 1:已抢完 2:已过期
	ExpireAt    *time.Time     `json:"expire_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Sender User `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
}

// RedPacketDetail 红包领取记录
type RedPacketDetail struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	RedPacketID uint       `gorm:"index" json:"red_packet_id"`
	UserID      uint       `gorm:"index" json:"user_id"`
	Amount      float64    `json:"amount"`
	IsBest      bool       `gorm:"default:false" json:"is_best"`
	CreatedAt   time.Time  `json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// ==================== MongoDB Models ====================

// MessageDoc 消息文档
type MessageDoc struct {
	ID             string                  `bson:"_id,omitempty" json:"id"`
	SenderID       uint                    `bson:"sender_id" json:"sender_id"`
	ReceiverID     uint                    `bson:"receiver_id,omitempty" json:"receiver_id"`
	GroupID        uint                    `bson:"group_id,omitempty" json:"group_id"`
	ConversationID uint                    `bson:"conversation_id,omitempty" json:"conversation_id"`
	Content        string                  `bson:"content" json:"content"`
	MessageType    int                     `bson:"message_type" json:"message_type"` // 1:文本 2:图片 3:文件 4:语音 5:视频 6:红包
	MediaURL       string                  `bson:"media_url,omitempty" json:"media_url"`
	MediaSize      int64                   `bson:"media_size,omitempty" json:"media_size"`
	Duration       int                     `bson:"duration,omitempty" json:"duration"`
	RedPacketID    uint                    `bson:"red_packet_id,omitempty" json:"red_packet_id"`
	IsRecall       bool                    `bson:"is_recall" json:"is_recall"`
	IsRead         bool                    `bson:"is_read" json:"is_read"`
	ReadUsers      []uint                  `bson:"read_users,omitempty" json:"read_users"`
	ReadAt         *time.Time              `bson:"read_at,omitempty" json:"read_at"`
	ExtInfo        map[string]interface{}  `bson:"ext_info,omitempty" json:"ext_info"`
	CreatedAt      time.Time               `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time               `bson:"updated_at" json:"updated_at"`
}

// AutoMigrate 自动迁移
func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&User{},
		&UserSettings{},
		&Friend{},
		&Group{},
		&GroupMember{},
		&GroupInvite{},
		&GroupJoinRequest{},
		&Conversation{},
		&InviteCode{},
		&PaymentOrder{},
		&PointsHistory{},
		&SystemConfig{},
		&AdminLog{},
		&RedPacket{},
		&RedPacketDetail{},
		&RechargeRequest{},
		&WithdrawRequest{},
	)
}
