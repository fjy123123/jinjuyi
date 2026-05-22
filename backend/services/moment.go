package services

import (
	"chat-system-pro/config"
	"chat-system-pro/models"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// MomentService 朋友圈服务
type MomentService struct {
}

// NewMomentService 创建朋友圈服务
func NewMomentService() *MomentService {
	return &MomentService{}
}

// PublishMoment 发布朋友圈
func (ms *MomentService) PublishMoment(userID uint, content string, images []string, location string, viewScope int, visibleTo, hideFrom []uint) (*models.Moment, error) {
	// 序列化
	imagesJSON, _ := json.Marshal(images)
	visibleJSON, _ := json.Marshal(visibleTo)
	hideJSON, _ := json.Marshal(hideFrom)
	
	moment := &models.Moment{
		UserID:    userID,
		Content:   content,
		Images:    string(imagesJSON),
		Location:  location,
		ViewScope: viewScope,
		VisibleTo: string(visibleJSON),
		HideFrom:  string(hideJSON),
	}
	
	if err := config.DB.Create(moment).Error; err != nil {
		return nil, err
	}
	
	// 预加载用户信息
	config.DB.Preload("User").First(moment, moment.ID)
	
	return moment, nil
}

// GetMoments 获取朋友圈列表
func (ms *MomentService) GetMoments(userID uint, page, pageSize int) ([]*models.Moment, int64, error) {
	// 获取我的好友
	var friendIDs []uint
	config.DB.Model(&models.Friend{}).Where("user_id = ?", userID).Pluck("friend_id", &friendIDs)
	
	// 也包含我自己的
	friendIDs = append(friendIDs, userID)
	
	// 查询
	var moments []*models.Moment
	var total int64
	
	query := config.DB.Model(&models.Moment{}).Where("user_id IN ?", friendIDs)
	
	query.Count(&total)
	query.Preload("User").Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&moments)
	
	// 加载点赞和评论
	for _, moment := range moments {
		ms.loadMomentDetails(moment)
	}
	
	return moments, total, nil
}

// GetMyMoments 获取我的朋友圈
func (ms *MomentService) GetMyMoments(userID uint, page, pageSize int) ([]*models.Moment, int64, error) {
	var moments []*models.Moment
	var total int64
	
	config.DB.Model(&models.Moment{}).Where("user_id = ?", userID).Count(&total)
	config.DB.Where("user_id = ?", userID).Preload("User").Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&moments)
	
	for _, moment := range moments {
		ms.loadMomentDetails(moment)
	}
	
	return moments, total, nil
}

func (ms *MomentService) loadMomentDetails(moment *models.Moment) {
	// TODO: 加载点赞和评论详情
}

// LikeMoment 点赞朋友圈
func (ms *MomentService) LikeMoment(userID, momentID uint) error {
	// 检查是否已经点赞过
	var existingLike models.MomentLike
	result := config.DB.Where("moment_id = ? AND user_id = ?", momentID, userID).First(&existingLike)
	if result.Error == nil {
		// 已经点赞过
		return errors.New("already liked")
	}
	
	// 创建点赞
	like := &models.MomentLike{
		MomentID: momentID,
		UserID:   userID,
	}
	
	if err := config.DB.Create(like).Error; err != nil {
		return err
	}
	
	// 更新点赞数
	config.DB.Model(&models.Moment{}).Where("id = ?", momentID).Update("likes_count", gorm.Expr("likes_count + 1"))
	
	return nil
}

// UnlikeMoment 取消点赞
func (ms *MomentService) UnlikeMoment(userID, momentID uint) error {
	result := config.DB.Where("moment_id = ? AND user_id = ?", momentID, userID).Delete(&models.MomentLike{})
	
	if result.RowsAffected > 0 {
		config.DB.Model(&models.Moment{}).Where("id = ?", momentID).Update("likes_count", gorm.Expr("likes_count - 1"))
	}
	
	return nil
}

// CommentMoment 评论朋友圈
func (ms *MomentService) CommentMoment(userID, momentID, replyToUser uint, content string) (*models.MomentComment, error) {
	comment := &models.MomentComment{
		MomentID:    momentID,
		UserID:      userID,
		ReplyToUser: replyToUser,
		Content:     content,
	}
	
	if err := config.DB.Create(comment).Error; err != nil {
		return nil, err
	}
	
	// 更新评论数
	config.DB.Model(&models.Moment{}).Where("id = ?", momentID).Update("comment_count", gorm.Expr("comment_count + 1"))
	
	// 预加载
	config.DB.Preload("User").Preload("ReplyUser").First(comment, comment.ID)
	
	return comment, nil
}

// DeleteComment 删除评论
func (ms *MomentService) DeleteComment(userID, commentID uint) error {
	var comment models.MomentComment
	if err := config.DB.First(&comment, commentID).Error; err != nil {
		return err
	}
	
	// 只能删除自己的评论
	if comment.UserID != userID {
		return errors.New("permission denied")
	}
	
	if err := config.DB.Delete(&comment).Error; err != nil {
		return err
	}
	
	config.DB.Model(&models.Moment{}).Where("id = ?", comment.MomentID).Update("comment_count", gorm.Expr("comment_count - 1"))
	
	return nil
}

// DeleteMoment 删除朋友圈
func (ms *MomentService) DeleteMoment(userID, momentID uint) error {
	var moment models.Moment
	if err := config.DB.First(&moment, momentID).Error; err != nil {
		return err
	}
	
	if moment.UserID != userID {
		return errors.New("permission denied")
	}
	
	// 级联删除点赞和评论
	config.DB.Where("moment_id = ?", momentID).Delete(&models.MomentLike{})
	config.DB.Where("moment_id = ?", momentID).Delete(&models.MomentComment{})
	config.DB.Delete(&moment)
	
	return nil
}

// ==================== 小程序服务 ====================

type MiniProgramService struct {
}

func NewMiniProgramService() *MiniProgramService {
	return &MiniProgramService{}
}

// LoginWechatMiniProgram 微信小程序登录
func (mps *MiniProgramService) LoginWechatMiniProgram(code, encryptedData, iv string) (*models.User, error) {
	// 1. 使用 code 换取 openid 和 session_key
	// 2. 解密用户信息（如果有）
	// 3. 查找或创建用户
	
	// TODO: 调用微信官方API
	
	demoOpenID := "demo_open_id_" + code
	
	// 查找是否已有
	var mpUser models.MiniProgramUser
	result := config.DB.Where("open_id = ?", demoOpenID).First(&mpUser)
	
	if result.Error != nil {
		// 创建新用户
		user := &models.User{
			Username:  "wx_" + config.GenerateShortID(8),
			Nickname:  "微信用户",
			Status:    0,
		}
		
		if err := config.DB.Create(user).Error; err != nil {
			return nil, err
		}
		
		// 创建小程序用户关联
		mpUser = models.MiniProgramUser{
			UserID:     user.ID,
			OpenID:     demoOpenID,
		}
		config.DB.Create(&mpUser)
		
		return user, nil
	}
	
	// 获取已关联用户
	var user models.User
	config.DB.First(&user, mpUser.UserID)
	
	return &user, nil
}

// ==================== 设备管理服务 ====================

type DeviceService struct {
}

func NewDeviceService() *DeviceService {
	return &DeviceService{}
}

// RegisterDevice 注册设备
func (ds *DeviceService) RegisterDevice(userID uint, deviceID, deviceType, deviceName, pushToken string) error {
	var device models.Device
	
	result := config.DB.Where("user_id = ? AND device_id = ?", userID, deviceID).First(&device)
	
	if result.Error != nil {
		// 创建新设备
		device = models.Device{
			UserID:     userID,
			DeviceID:   deviceID,
			DeviceType: deviceType,
			DeviceName: deviceName,
			PushToken:  pushToken,
			IsOnline:   true,
			LastActive: time.Now(),
		}
		return config.DB.Create(&device).Error
	}
	
	// 更新设备信息
	device.PushToken = pushToken
	device.IsOnline = true
	device.LastActive = time.Now()
	
	return config.DB.Save(&device).Error
}

// GetUserDevices 获取用户设备
func (ds *DeviceService) GetUserDevices(userID uint) ([]*models.Device, error) {
	var devices []*models.Device
	err := config.DB.Where("user_id = ?", userID).Find(&devices).Error
	return devices, err
}

// SetDeviceOffline 设置设备离线
func (ds *DeviceService) SetDeviceOffline(userID uint, deviceID string) {
	config.DB.Model(&models.Device{}).Where("user_id = ? AND device_id = ?", userID, deviceID).Updates(map[string]interface{}{
		"is_online":    false,
		"last_active":  time.Now(),
	})
}
