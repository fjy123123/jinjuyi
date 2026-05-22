package handlers

import (
	"strconv"

	"chat-system-pro/services"
	"chat-system-pro/utils"

	"github.com/gin-gonic/gin"
)

var momentService = services.NewMomentService()

// PublishMoment 发布朋友圈
func PublishMoment(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		Content   string   `json:"content"`
		Images    []string `json:"images"`
		Location  string   `json:"location"`
		ViewScope int      `json:"view_scope"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "参数错误")
		return
	}
	moment, err := momentService.PublishMoment(userID, req.Content, req.Images, req.Location, req.ViewScope, nil, nil)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, moment)
}

// GetMoments 获取朋友圈列表
func GetMoments(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	moments, total, err := momentService.GetMoments(userID, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.PaginatedResponse(c, moments, total, page, pageSize)
}

// LikeMoment 点赞朋友圈
func LikeMoment(c *gin.Context) {
	userID := c.GetUint("user_id")
	momentID, _ := strconv.ParseUint(c.Param("moment_id"), 10, 32)
	if err := momentService.LikeMoment(userID, uint(momentID)); err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, nil)
}

// UnlikeMoment 取消点赞
func UnlikeMoment(c *gin.Context) {
	userID := c.GetUint("user_id")
	momentID, _ := strconv.ParseUint(c.Param("moment_id"), 10, 32)
	momentService.UnlikeMoment(userID, uint(momentID))
	utils.SuccessResponse(c, nil)
}

// CommentMoment 评论朋友圈
func CommentMoment(c *gin.Context) {
	userID := c.GetUint("user_id")
	momentID, _ := strconv.ParseUint(c.Param("moment_id"), 10, 32)
	var req struct {
		ReplyToUser uint   `json:"reply_to_user"`
		Content     string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "参数错误")
		return
	}
	comment, err := momentService.CommentMoment(userID, uint(momentID), req.ReplyToUser, req.Content)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, comment)
}

// DeleteMoment 删除朋友圈
func DeleteMoment(c *gin.Context) {
	userID := c.GetUint("user_id")
	momentID, _ := strconv.ParseUint(c.Param("moment_id"), 10, 32)
	if err := momentService.DeleteMoment(userID, uint(momentID)); err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, nil)
}
