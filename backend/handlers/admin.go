package handlers

import (
	"strconv"
	"time"

	"chat-system-pro/config"
	"chat-system-pro/models"
	"chat-system-pro/services"
	"chat-system-pro/utils"

	"github.com/gin-gonic/gin"
)

var databaseService = services.NewDatabaseService()
var paymentService = services.NewPaymentService()

// ==================== 数据库管理接口 ====================

func ClearOldMessages(c *gin.Context) {
	dateStr := c.Query("date")
	beforeDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid date format")
		return
	}
	count, err := databaseService.ClearOldMessages(beforeDate)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, gin.H{"deleted_count": count})
}

func ClearAllData(c *gin.Context) {
	confirm := c.PostForm("confirm")
	if confirm != "YES" {
		utils.ErrorResponse(c, 400, "Confirm required")
		return
	}
	if err := databaseService.ClearAllData(); err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, nil)
}

func InitializeDatabase(c *gin.Context) {
	confirm := c.PostForm("confirm")
	if confirm != "YES" {
		utils.ErrorResponse(c, 400, "Confirm required")
		return
	}
	if err := databaseService.InitializeDatabase(); err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, nil)
}

func DeleteUserAndData(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err := databaseService.DeleteUserAndData(uint(userID)); err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, nil)
}

func DeleteGroupAndData(c *gin.Context) {
	groupID, _ := strconv.ParseUint(c.Param("group_id"), 10, 32)
	if err := databaseService.DeleteGroupAndData(uint(groupID)); err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, nil)
}

func ArchiveOldMessages(c *gin.Context) {
	dateStr := c.Query("date")
	beforeDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid date format")
		return
	}
	count, err := databaseService.ArchiveOldMessages(beforeDate)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, gin.H{"archived_count": count})
}

func GetDatabaseStats(c *gin.Context) {
	stats, err := databaseService.GetDatabaseStats()
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, stats)
}

// ==================== 支付相关接口 ====================

func CreateOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		Amount    float64 `json:"amount" binding:"required"`
		PayType   int     `json:"pay_type" binding:"required"`
		OrderType int     `json:"order_type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "Invalid params")
		return
	}
	order, err := paymentService.CreateOrder(userID, req.Amount, req.PayType, req.OrderType)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, order)
}

func ProcessPayment(c *gin.Context) {
	orderID, _ := strconv.ParseUint(c.Param("order_id"), 10, 32)
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "Invalid params")
		return
	}
	if err := paymentService.ProcessStripePayment(uint(orderID), req.Token); err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, nil)
}

func GetOrders(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	orders, total, err := paymentService.GetOrders(userID, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.PaginatedResponse(c, orders, total, page, pageSize)
}

func GetPointsHistory(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	history, total, err := paymentService.GetPointsHistory(userID, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.PaginatedResponse(c, history, total, page, pageSize)
}

func AdminAddPoints(c *gin.Context) {
	var req struct {
		UserID uint   `json:"user_id" binding:"required"`
		Points int64  `json:"points" binding:"required"`
		Remark string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "Invalid params")
		return
	}
	if err := paymentService.AddPoints(req.UserID, req.Points, req.Remark); err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, nil)
}

// ==================== 系统配置接口 ====================

func UpdateSystemConfig(c *gin.Context) {
	var req struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "Invalid params")
		return
	}
	var cfg models.SystemConfig
	if err := config.DB.Where("`key` = ?", req.Key).First(&cfg).Error; err != nil {
		cfg = models.SystemConfig{Key: req.Key}
	}
	cfg.Value = req.Value
	config.DB.Save(&cfg)
	utils.SuccessResponse(c, cfg)
}

func GetSystemConfigs(c *gin.Context) {
	var configs []models.SystemConfig
	config.DB.Find(&configs)
	configMap := make(map[string]string)
	for _, c := range configs {
		configMap[c.Key] = c.Value
	}
	utils.SuccessResponse(c, configMap)
}
