package models

import (
	"time"

	"gorm.io/gorm"
)

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
	OnlineStatus int           `gorm:"default:0" json:"online_status"` // 0:离线 1:在线 2:忙碌 3:离开
	IsTyping    bool           `gorm:"-" json:"is_typing,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// UserSettings 用户设置
type UserSettings struct {
	ID              uint   `gorm:"primarykey" json:"id"`
	UserID          uint   `gorm:"uniqueIndex;not null" json:"user_id"`
	NewMsgNotify    bool   `gorm:"default:true" json:"new_msg_notify"`
	SoundNotify     bool   `gorm:"default:true" json:"sound_notify"`
	AddFriendConfirm bool   `gorm:"default:true" json:"add_friend_confirm"`
	ShowOnline      bool   `gorm:"default:true" json:"show_online"`
	ShowReadReceipt bool   `gorm:"default:true" json:"show_read_receipt"`
	Theme           string `gorm:"size:50;default:modern" json:"theme"`
	Language        string `gorm:"size:10;default:zh-CN" json:"language"`
	DndStart        string `gorm:"size:10" json:"dnd_start"` // 免打扰开始时间
	DndEnd          string `gorm:"size:10" json:"dnd_end"`   // 免打扰结束时间
	AutoDeleteDays  int    `gorm:"default:0" json:"auto_delete_days"` // 消息自动删除天数，0表示不删除
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Friend 好友关系
type Friend struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	FriendID  uint      `gorm:"index;not null" json:"friend_id"`
	Remark    string    `gorm:"size:100" json:"remark"`     // 备注名
	Tag       string    `gorm:"size:100" json:"tag"`        // 标签
	IsPinned  bool      `gorm:"default:false" json:"is_pinned"`
	IsMuted   bool      `gorm:"default:false" json:"is_muted"`
	IsBlocked bool      `gorm:"default:false" json:"is_blocked"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Friend    User      `gorm:"foreignKey:FriendID" json:"friend,omitempty"`
}

// Group 群组
type Group struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	Name         string         `gorm:"size:100;not null" json:"name"`
	Avatar       string         `gorm:"size:500" json:"avatar"`
	Description  string         `gorm:"size:500" json:"description"`
	OwnerID      uint           `gorm:"index" json:"owner_id"`
	Notice       string         `gorm:"size:500" json:"notice"`        // 群公告
	IsMuted      bool           `gorm:"default:false" json:"is_muted"`   // 群免打扰
	AllowInvite  bool           `gorm:"default:true" json:"allow_invite"`
	NeedApprove  bool           `gorm:"default:false" json:"need_approve"` // 加入需审核
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Members      []GroupMember  `gorm:"-" json:"members,omitempty"`
}

// GroupMember 群成员
type GroupMember struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	GroupID   uint      `gorm:"index;not null" json:"group_id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Role      int       `gorm:"default:0" json:"role"` // 0:普通成员 1:管理员 2:群主
	Nickname  string    `gorm:"size:100" json:"nickname"` // 群内昵称
	IsMuted   bool      `gorm:"default:false" json:"is_muted"`
	JoinedAt  time.Time `json:"joined_at"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// GroupInvite 群邀请
type GroupInvite struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	GroupID   uint      `gorm:"index" json:"group_id"`
	InviterID uint      `gorm:"index" json:"inviter_id"`
	InviteeID uint      `gorm:"index" json:"invitee_id"`
	Status    int       `gorm:"default:0" json:"status"` // 0:待处理 1:已接受 2:已拒绝
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// GroupJoinRequest 加群申请
type GroupJoinRequest struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	GroupID   uint      `gorm:"index" json:"group_id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Message   string    `gorm:"size:200" json:"message"` // 申请留言
	Status    int       `gorm:"default:0" json:"status"` // 0:待审核 1:已通过 2:已拒绝
	HandledBy uint      `json:"handled_by"`
	HandleAt  *time.Time `json:"handle_at"`
	CreatedAt time.Time `json:"created_at"`
}

// Conversation 会话
type Conversation struct {
	ID            uint       `gorm:"primarykey" json:"id"`
	UserID        uint       `gorm:"index;not null" json:"user_id"`
	TargetID      uint       `gorm:"index" json:"target_id"`
	Type          int        `gorm:"not null" json:"type"` // 1:私聊 2:群聊
	UnreadCount   int        `gorm:"default:0" json:"unread_count"`
	LastMessageID string     `gorm:"size:50" json:"last_message_id"`
	LastMessage   string     `gorm:"size:200" json:"last_message"`
	LastMessageAt *time.Time `json:"last_message_at"`
	IsPinned      bool       `gorm:"default:false" json:"is_pinned"`
	IsMuted       bool       `gorm:"default:false" json:"is_muted"`
	IsArchived    bool       `gorm:"default:false" json:"is_archived"`
	TopMessageID  string     `gorm:"size:50" json:"top_message_id"` // 置顶消息ID
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// InviteCode 邀请码
type InviteCode struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Code      string    `gorm:"uniqueIndex;size:20;not null" json:"code"`
	CreatorID uint      `gorm:"index" json:"creator_id"`
	UsedCount int       `gorm:"default:0" json:"used_count"`
	MaxUses   int       `gorm:"default:10" json:"max_uses"`
	ExpiresAt *time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
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
	ID                uint           `gorm:"primarykey" json:"id"`
	AppName           string         `gorm:"size:100;default:知信" json:"app_name"`
	AppVersion        string         `gorm:"size:50;default:v1.0.0" json:"app_version"`
	AppDescription    string         `gorm:"size:500" json:"app_description"`
	LogoURL           string         `gorm:"size:500" json:"logo_url"`
	FaviconURL        string         `gorm:"size:500" json:"favicon_url"`
	ThemeColor        string         `gorm:"size:50;default:#07c160" json:"theme_color"`
	ThemeSecondary    string         `gorm:"size:50;default:#576b95" json:"theme_secondary"`
	UiTemplate        string         `gorm:"size:50;default:modern" json:"ui_template"`
	MaintenanceMode   bool           `gorm:"default:false" json:"maintenance_mode"`
	MaintenanceMsg    string         `gorm:"size:500" json:"maintenance_msg"`
	ExportEnabled     bool           `gorm:"default:true" json:"export_enabled"`
	ExportMaxRecords  int            `gorm:"default:1000" json:"export_max_records"`
	RecallTimeout     int            `gorm:"default:300" json:"recall_timeout"` // 消息撤回超时时间（秒），默认5分钟
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
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
	ID             uint           `gorm:"primarykey" json:"id"`
	SenderID       uint           `gorm:"index" json:"sender_id"`
	ReceiverID     uint           `gorm:"index" json:"receiver_id,omitempty"`  // 0 表示群红包
	GroupID        uint           `gorm:"index" json:"group_id,omitempty"`
	Type           int            `gorm:"default:1" json:"type"`  // 1:普通红包 2:拼手气红包
	PayType        int            `gorm:"default:1" json:"pay_type"`  // 1:积分 2:微信 3:支付宝
	Amount         float64        `json:"amount"`
	TotalCount     int            `gorm:"default:1" json:"total_count"`
	ReceivedCount  int            `gorm:"default:0" json:"received_count"`
	ReceivedAmount float64        `gorm:"default:0" json:"received_amount"`
	Greeting       string         `gorm:"size:200" json:"greeting"`
	Status         int            `gorm:"default:0" json:"status"`  // 0:进行中 1:已抢完 2:已过期
	ExpireAt       *time.Time     `json:"expire_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	Sender User `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
}

// RedPacketDetail 红包领取记录
type RedPacketDetail struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	RedPacketID uint      `gorm:"index" json:"red_packet_id"`
	UserID      uint      `gorm:"index" json:"user_id"`
	Amount      float64   `json:"amount"`
	IsBest      bool      `gorm:"default:false" json:"is_best"`
	CreatedAt   time.Time `json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// ==================== 充值提现相关 ====================

// RechargeRequest 充值申请
type RechargeRequest struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	UserID      uint       `gorm:"index" json:"user_id"`
	Amount      float64    `json:"amount"`
	Points      int64      `json:"points"`
	ProofImage  string     `gorm:"size:500" json:"proof_image"` // 支付凭证图片
	Status      int        `gorm:"default:0" json:"status"`     // 0:待审核 1:已通过 2:已拒绝
	HandlerID   uint       `json:"handler_id"`
	HandleNote  string     `gorm:"size:500" json:"handle_note"`
	HandledAt   *time.Time `json:"handled_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	User        User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// WithdrawRequest 提现申请
type WithdrawRequest struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	UserID      uint       `gorm:"index" json:"user_id"`
	Amount      float64    `json:"amount"`
	Points      int64      `json:"points"`
	PayType     int        `json:"pay_type"` // 1:微信 2:支付宝
	AccountInfo string     `gorm:"size:500" json:"account_info"` // 收款账户信息
	Status      int        `gorm:"default:0" json:"status"`      // 0:待审核 1:已通过 2:已拒绝
	HandlerID   uint       `json:"handler_id"`
	HandleNote  string     `gorm:"size:500" json:"handle_note"`
	HandledAt   *time.Time `json:"handled_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	User        User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
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
	MessageType    int                     `bson:"message_type" json:"message_type"` // 1:文本 2:图片 3:文件 4:语音 5:视频 6:红包 7:系统消息 8:引用消息 9:转发消息
	MediaURL       string                  `bson:"media_url,omitempty" json:"media_url"`
	MediaSize      int64                   `bson:"media_size,omitempty" json:"media_size"`
	MediaName      string                  `bson:"media_name,omitempty" json:"media_name"` // 文件名
	Duration       int                     `bson:"duration,omitempty" json:"duration"`     // 语音/视频时长(秒)
	RedPacketID    uint                    `bson:"red_packet_id,omitempty" json:"red_packet_id"`
	
	// 回复/引用功能
	ReplyToID     string                  `bson:"reply_to_id,omitempty" json:"reply_to_id"`       // 回复的消息ID
	ReplyToContent string                  `bson:"reply_to_content,omitempty" json:"reply_to_content"` // 被回复的消息内容预览
	ReplyToSender  uint                    `bson:"reply_to_sender,omitempty" json:"reply_to_sender"`  // 被回复消息的发送者
	
	// 转发功能
	IsForwarded    bool                    `bson:"is_forwarded,omitempty" json:"is_forwarded"`      // 是否是转发消息
	ForwardFromID  string                  `bson:"forward_from_id,omitempty" json:"forward_from_id"` // 转发来源消息ID
	OriginalSender uint                    `bson:"original_sender,omitempty" json:"original_sender"` // 原消息发送者
	
	// 已读功能
	IsRecall       bool                    `bson:"is_recall" json:"is_recall"`
	IsRead         bool                    `bson:"is_read" json:"is_read"`
	ReadUsers      []uint                  `bson:"read_users,omitempty" json:"read_users"`
	ReadAt         *time.Time              `bson:"read_at,omitempty" json:"read_at"`
	
	// 表情反应 (类Telegram)
	Reactions      []MessageReaction       `bson:"reactions,omitempty" json:"reactions"`
	
	// 阅后即焚
	SelfDestruct   int                     `bson:"self_destruct,omitempty" json:"self_destruct"` // 秒数，0表示不销毁
	
	IsBurned       bool                    `bson:"is_burned,omitempty" json:"is_burned"` // 是否已销毁
	
	ExtInfo        map[string]interface{}  `bson:"ext_info,omitempty" json:"ext_info"`
	CreatedAt      time.Time               `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time               `bson:"updated_at" json:"updated_at"`
}

// MessageReaction 消息表情反应
type MessageReaction struct {
	Emoji    string   `bson:"emoji" json:"emoji"`       // 表情符号
	UserIDs  []uint   `bson:"user_ids" json:"user_ids"` // 点了这个表情的用户
	Count    int      `bson:"count" json:"count"`
}

// MessageReactionRecord 用户反应记录 (用于追踪谁反应了什么)
type MessageReactionRecord struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	MessageID  string    `gorm:"size:50;index" json:"message_id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	Emoji      string    `gorm:"size:50" json:"emoji"`
	CreatedAt  time.Time `json:"created_at"`
}

// EmojiCategory 表情包分类
type EmojiCategory struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Icon      string    `gorm:"size:500" json:"icon"`
	Sort      int       `gorm:"default:0" json:"sort"`
	CreatedAt time.Time `json:"created_at"`
}

// EmojiItem 表情包项目
type EmojiItem struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	CategoryID uint      `gorm:"index" json:"category_id"`
	Name       string    `gorm:"size:100" json:"name"`
	URL        string    `gorm:"size:500" json:"url"`
	Sort       int       `gorm:"default:0" json:"sort"`
	CreatedAt  time.Time `json:"created_at"`
}

// ==================== 通话相关 ====================

// CallRecord 通话记录
type CallRecord struct {
	ID         uint       `gorm:"primarykey" json:"id"`
	SessionID  string     `gorm:"size:100;index" json:"session_id"`
	CallerID   uint       `gorm:"index" json:"caller_id"`
	CalleeID   uint       `gorm:"index" json:"callee_id"`
	GroupID    uint       `gorm:"index" json:"group_id,omitempty"`
	Type       int        `json:"type"` // 1:视频 2:语音
	Status     int        `json:"status"` // 0:呼叫中 1:通话中 2:已结束 3:已拒绝
	StartTime  time.Time  `json:"start_time"`
	EndTime    time.Time  `json:"end_time,omitempty"`
	Duration   int        `json:"duration"` // 通话时长（秒）
	CreatedAt  time.Time  `json:"created_at"`
	
	Caller User `gorm:"foreignKey:CallerID" json:"caller,omitempty"`
	Callee User `gorm:"foreignKey:CalleeID" json:"callee,omitempty"`
}

// ==================== 2FA双因素认证 ====================

// TwoFactorAuth 双因素认证
type TwoFactorAuth struct {
	ID           uint       `gorm:"primarykey" json:"id"`
	UserID       uint       `gorm:"uniqueIndex" json:"user_id"`
	Secret       string     `gorm:"size:100;not null" json:"-"` // TOTP密钥
	Enabled      bool       `gorm:"default:false" json:"enabled"`
	BackupCodes  string     `gorm:"size:500" json:"-"` // 备用验证码（JSON数组）
	VerifiedAt   *time.Time `json:"verified_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
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
		&MessageReactionRecord{},
		&EmojiCategory{},
		&EmojiItem{},
		&CallRecord{},
		&TwoFactorAuth{},
	)
}
