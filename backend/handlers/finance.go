package handlers

import (
	"chat-system-pro/models"
	"chat-system-pro/services"
	"chat-system-pro/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ==================== 充值接口 ====================

// CreateRechargeRequest 创建充值申请
func CreateRechargeRequest(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Amount       float64               `json:"amount" binding:"required"`
		Points       int64                 `json:"points" binding:"required"`
		RechargeType models.RechargeType   `json:"recharge_type" binding:"required"`
		PaymentImage string                `json:"payment_image" binding:"required"`
		Remark       string                `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "参数错误", err.Error())
		return
	}

	rechargeReq, err := services.RechargeService.CreateRechargeRequest(
		userID,
		req.Amount,
		req.Points,
		req.RechargeType,
		req.PaymentImage,
		req.Remark,
	)

	if err != nil {
		utils.ErrorResponse(c, 500, "创建充值申请失败", err.Error())
		return
	}

	utils.SuccessResponse(c, rechargeReq)
}

// GetMyRechargeRequests 获取我的充值申请
func GetMyRechargeRequests(c *gin.Context) {
	userID := c.GetUint("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	requests, total, err := services.RechargeService.GetMyRechargeRequests(userID, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, 500, "获取充值申请列表失败", err.Error())
		return
	}

	utils.PaginatedResponse(c, requests, page, pageSize, total)
}

// GetRechargeRequestDetail 获取充值申请详情
func GetRechargeRequestDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	req, err := services.RechargeService.GetRechargeRequest(uint(id))
	if err != nil {
		utils.ErrorResponse(c, 404, "充值申请不存在", err.Error())
		return
	}

	utils.SuccessResponse(c, req)
}

// ==================== 提现接口 ====================

// CreateWithdrawRequest 创建提现申请
func CreateWithdrawRequest(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Points       int64                 `json:"points" binding:"required"`
		Amount       float64               `json:"amount" binding:"required"`
		WithdrawType models.WithdrawType   `json:"withdraw_type" binding:"required"`
		PaymentCode  string                `json:"payment_code" binding:"required"`
		RealName     string                `json:"real_name" binding:"required"`
		Phone        string                `json:"phone" binding:"required"`
		Remark       string                `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "参数错误", err.Error())
		return
	}

	withdrawReq, err := services.WithdrawService.CreateWithdrawRequest(
		userID,
		req.Points,
		req.Amount,
		req.WithdrawType,
		req.PaymentCode,
		req.RealName,
		req.Phone,
		req.Remark,
	)

	if err != nil {
		utils.ErrorResponse(c, 500, "创建提现申请失败", err.Error())
		return
	}

	utils.SuccessResponse(c, withdrawReq)
}

// GetMyWithdrawRequests 获取我的提现申请
func GetMyWithdrawRequests(c *gin.Context) {
	userID := c.GetUint("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	requests, total, err := services.WithdrawService.GetMyWithdrawRequests(userID, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, 500, "获取提现申请列表失败", err.Error())
		return
	}

	utils.PaginatedResponse(c, requests, page, pageSize, total)
}

// GetWithdrawRequestDetail 获取提现申请详情
func GetWithdrawRequestDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	req, err := services.WithdrawService.GetWithdrawRequest(uint(id))
	if err != nil {
		utils.ErrorResponse(c, 404, "提现申请不存在", err.Error())
		return
	}

	utils.SuccessResponse(c, req)
}

// ==================== 管理员审核接口 ====================

// GetAllRechargeRequests 获取所有充值申请
func GetAllRechargeRequests(c *gin.Context) {
	status := models.RechargeStatus(c.Query("status"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	requests, total, err := services.RechargeService.GetAllRechargeRequests(status, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, 500, "获取充值申请列表失败", err.Error())
		return
	}

	utils.PaginatedResponse(c, requests, page, pageSize, total)
}

// ApproveRecharge 审核通过充值
func ApproveRecharge(c *gin.Context) {
	reviewerID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		Remark string `json:"remark"`
	}
	c.ShouldBindJSON(&req)

	rechargeReq, err := services.RechargeService.ApproveRecharge(reviewerID, uint(id), req.Remark)
	if err != nil {
		utils.ErrorResponse(c, 500, "审核失败", err.Error())
		return
	}

	utils.SuccessResponse(c, rechargeReq)
}

// RejectRecharge 拒绝充值申请
func RejectRecharge(c *gin.Context) {
	reviewerID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		Remark string `json:"remark"`
	}
	c.ShouldBindJSON(&req)

	rechargeReq, err := services.RechargeService.RejectRecharge(reviewerID, uint(id), req.Remark)
	if err != nil {
		utils.ErrorResponse(c, 500, "拒绝失败", err.Error())
		return
	}

	utils.SuccessResponse(c, rechargeReq)
}

// GetAllWithdrawRequests 获取所有提现申请
func GetAllWithdrawRequests(c *gin.Context) {
	status := models.WithdrawStatus(c.Query("status"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	requests, total, err := services.WithdrawService.GetAllWithdrawRequests(status, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, 500, "获取提现申请列表失败", err.Error())
		return
	}

	utils.PaginatedResponse(c, requests, page, pageSize, total)
}

// ApproveWithdraw 审核通过提现
func ApproveWithdraw(c *gin.Context) {
	reviewerID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		Remark string `json:"remark"`
	}
	c.ShouldBindJSON(&req)

	withdrawReq, err := services.WithdrawService.ApproveWithdraw(reviewerID, uint(id), req.Remark)
	if err != nil {
		utils.ErrorResponse(c, 500, "审核失败", err.Error())
		return
	}

	utils.SuccessResponse(c, withdrawReq)
}

// RejectWithdraw 拒绝提现申请
func RejectWithdraw(c *gin.Context) {
	reviewerID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		Remark string `json:"remark"`
	}
	c.ShouldBindJSON(&req)

	withdrawReq, err := services.WithdrawService.RejectWithdraw(reviewerID, uint(id), req.Remark)
	if err != nil {
		utils.ErrorResponse(c, 500, "拒绝失败", err.Error())
		return
	}

	utils.SuccessResponse(c, withdrawReq)
}
