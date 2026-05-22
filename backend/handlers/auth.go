package handlers

import (
	"chat-system-pro/config"
	"chat-system-pro/models"
	"chat-system-pro/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Register 用户注册
func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Password string `json:"password" binding:"required,min=6"`
		Nickname string `json:"nickname"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "参数错误: "+err.Error())
		return
	}
	// 检查用户名是否已存在
	var existing models.User
	if config.DB.Where("username = ?", req.Username).First(&existing).Error == nil {
		utils.ErrorResponse(c, 400, "用户名已存在")
		return
	}
	// 加密密码
	hashedPass, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.ErrorResponse(c, 500, "系统错误")
		return
	}
	nickname := req.Nickname
	if nickname == "" {
		nickname = req.Username
	}
	user := models.User{
		Username: req.Username,
		Password: hashedPass,
		Nickname: nickname,
		Phone:    req.Phone,
		Email:    req.Email,
	}
	if err := config.DB.Create(&user).Error; err != nil {
		utils.ErrorResponse(c, 500, "创建用户失败")
		return
	}
	token, _ := generateToken(user.ID)
	utils.SuccessResponse(c, gin.H{"token": token, "user": user})
}

// Login 用户登录
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "参数错误")
		return
	}
	var user models.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.ErrorResponse(c, 400, "用户名或密码错误")
		return
	}
	if !utils.CheckPassword(req.Password, user.Password) {
		utils.ErrorResponse(c, 400, "用户名或密码错误")
		return
	}
	now := time.Now()
	config.DB.Model(&user).Updates(map[string]interface{}{
		"last_login_at": now,
		"last_login_ip": c.ClientIP(),
	})
	token, _ := generateToken(user.ID)
	utils.SuccessResponse(c, gin.H{"token": token, "user": user})
}

// GetProfile 获取当前用户信息
func GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		utils.ErrorResponse(c, 404, "用户不存在")
		return
	}
	utils.SuccessResponse(c, user)
}

// UpdateProfile 更新个人资料
func UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
		Gender   int    `json:"gender"`
		Region   string `json:"region"`
		Sign     string `json:"sign"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "参数错误")
		return
	}
	config.DB.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"nickname": req.Nickname,
		"avatar":   req.Avatar,
		"gender":   req.Gender,
		"region":   req.Region,
		"sign":     req.Sign,
		"phone":    req.Phone,
		"email":    req.Email,
	})
	utils.SuccessResponse(c, nil)
}

// GetUserSettings 获取用户设置
func GetUserSettings(c *gin.Context) {
	userID := c.GetUint("user_id")
	var settings models.UserSettings
	if err := config.DB.Where("user_id = ?", userID).First(&settings).Error; err != nil {
		// 创建默认设置
		settings = models.UserSettings{UserID: userID}
		config.DB.Create(&settings)
	}
	utils.SuccessResponse(c, settings)
}

// UpdateUserSettings 更新用户设置
func UpdateUserSettings(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req models.UserSettings
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "参数错误")
		return
	}
	config.DB.Model(&models.UserSettings{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"new_msg_notify":     req.NewMsgNotify,
		"sound_notify":       req.SoundNotify,
		"add_friend_confirm": req.AddFriendConfirm,
		"show_online":        req.ShowOnline,
		"show_read_receipt":  req.ShowReadReceipt,
		"theme":              req.Theme,
		"language":           req.Language,
	})
	utils.SuccessResponse(c, nil)
}

// SearchUsers 搜索用户
func SearchUsers(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		utils.ErrorResponse(c, 400, "搜索关键词不能为空")
		return
	}
	userID := c.GetUint("user_id")
	var users []models.User
	config.DB.Where("id != ? AND (username LIKE ? OR nickname LIKE ?)", userID, "%"+keyword+"%", "%"+keyword+"%").
		Limit(20).Find(&users)
	utils.SuccessResponse(c, users)
}

func generateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Duration(config.Cfg.JWT.ExpireHours) * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Cfg.JWT.Secret))
}
