package handlers

import (
	"chat-system-pro/config"
	"chat-system-pro/models"
	"chat-system-pro/utils"
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
)

// Enable2FA 开启2FA
func Enable2FA(c *gin.Context) {
	userID := c.GetUint("user_id")

	// 检查是否已开启
	var existing models.TwoFactorAuth
	if err := config.DB.Where("user_id = ?", userID).First(&existing).Error; err == nil && existing.Enabled {
		utils.ErrorResponse(c, 400, "2FA已开启")
		return
	}

	// 生成密钥
	key := make([]byte, 20)
	rand.Read(key)
	secret := base32.StdEncoding.EncodeToString(key)

	// 生成TOTP密钥
	totpKey, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      "知信",
		AccountName: "",
		Secret:      key,
	})

	// 生成备用验证码
	backupCodes := generateBackupCodes()
	backupCodesJSON, _ := json.Marshal(backupCodes)

	// 保存到数据库
	tfa := models.TwoFactorAuth{
		UserID:      userID,
		Secret:      secret,
		Enabled:     false, // 需要验证后才真正开启
		BackupCodes: string(backupCodesJSON),
	}
	config.DB.Save(&tfa)

	utils.SuccessResponse(c, gin.H{
		"secret":       secret,
		"qr_code":      totpKey.URL(),
		"backup_codes": backupCodes,
	})
}

// VerifyAndEnable2FA 验证并启用2FA
func VerifyAndEnable2FA(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "参数错误")
		return
	}

	var tfa models.TwoFactorAuth
	if err := config.DB.Where("user_id = ?", userID).First(&tfa).Error; err != nil {
		utils.ErrorResponse(c, 404, "请先开启2FA")
		return
	}

	// 验证TOTP码
	valid := totp.Validate(req.Code, tfa.Secret)
	if !valid {
		utils.ErrorResponse(c, 400, "验证码错误")
		return
	}

	// 启用2FA
	now := time.Now()
	tfa.Enabled = true
	tfa.VerifiedAt = &now
	config.DB.Save(&tfa)

	utils.SuccessResponse(c, gin.H{"message": "2FA已启用"})
}

// Disable2FA 关闭2FA
func Disable2FA(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Password string `json:"password" binding:"required"`
		Code     string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "参数错误")
		return
	}

	// 验证密码
	var user models.User
	config.DB.First(&user, userID)
	if !utils.CheckPassword(req.Password, user.Password) {
		utils.ErrorResponse(c, 400, "密码错误")
		return
	}

	// 验证2FA码
	var tfa models.TwoFactorAuth
	if err := config.DB.Where("user_id = ?", userID).First(&tfa).Error; err != nil {
		utils.ErrorResponse(c, 404, "2FA未开启")
		return
	}

	valid := totp.Validate(req.Code, tfa.Secret)
	if !valid {
		// 尝试备用验证码
		var backupCodes []string
		json.Unmarshal([]byte(tfa.BackupCodes), &backupCodes)
		valid = false
		for i, code := range backupCodes {
			if code == req.Code {
				valid = true
				backupCodes = append(backupCodes[:i], backupCodes[i+1:]...)
				newCodes, _ := json.Marshal(backupCodes)
				tfa.BackupCodes = string(newCodes)
				break
			}
		}
	}

	if !valid {
		utils.ErrorResponse(c, 400, "验证码错误")
		return
	}

	// 删除2FA
	config.DB.Delete(&tfa)

	utils.SuccessResponse(c, gin.H{"message": "2FA已关闭"})
}

// Verify2FA 验证2FA码
func Verify2FA(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "参数错误")
		return
	}

	var tfa models.TwoFactorAuth
	if err := config.DB.Where("user_id = ? AND enabled = ?", userID, true).First(&tfa).Error; err != nil {
		utils.ErrorResponse(c, 404, "2FA未开启")
		return
	}

	valid := totp.Validate(req.Code, tfa.Secret)
	if !valid {
		utils.ErrorResponse(c, 400, "验证码错误")
		return
	}

	utils.SuccessResponse(c, gin.H{"valid": true})
}

// Get2FAStatus 获取2FA状态
func Get2FAStatus(c *gin.Context) {
	userID := c.GetUint("user_id")

	var tfa models.TwoFactorAuth
	err := config.DB.Where("user_id = ?", userID).First(&tfa).Error

	if err != nil {
		utils.SuccessResponse(c, gin.H{"enabled": false})
		return
	}

	utils.SuccessResponse(c, gin.H{
		"enabled":     tfa.Enabled,
		"verified_at": tfa.VerifiedAt,
	})
}

// RegenerateBackupCodes 重新生成备用验证码
func RegenerateBackupCodes(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "参数错误")
		return
	}

	var tfa models.TwoFactorAuth
	if err := config.DB.Where("user_id = ? AND enabled = ?", userID, true).First(&tfa).Error; err != nil {
		utils.ErrorResponse(c, 404, "2FA未开启")
		return
	}

	// 验证当前2FA码
	valid := totp.Validate(req.Code, tfa.Secret)
	if !valid {
		utils.ErrorResponse(c, 400, "验证码错误")
		return
	}

	// 生成新的备用验证码
	backupCodes := generateBackupCodes()
	backupCodesJSON, _ := json.Marshal(backupCodes)
	tfa.BackupCodes = string(backupCodesJSON)
	config.DB.Save(&tfa)

	utils.SuccessResponse(c, gin.H{"backup_codes": backupCodes})
}

// generateBackupCodes 生成备用验证码
func generateBackupCodes() []string {
	codes := make([]string, 10)
	for i := 0; i < 10; i++ {
		code := make([]byte, 8)
		rand.Read(code)
		codes[i] = base32.StdEncoding.EncodeToString(code)[:8]
	}
	return codes
}
