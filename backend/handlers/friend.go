package handlers

import (
	"strconv"

	"chat-system-pro/config"
	"chat-system-pro/models"
	"chat-system-pro/utils"

	"github.com/gin-gonic/gin"
)

// GetFriends 获取好友列表
func GetFriends(c *gin.Context) {
	userID := c.GetUint("user_id")
	var friends []models.Friend
	config.DB.Preload("Friend").Where("user_id = ? AND status = 0", userID).Find(&friends)
	utils.SuccessResponse(c, friends)
}

// AddFriend 添加好友
func AddFriend(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		FriendID uint   `json:"friend_id" binding:"required"`
		Remark   string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "参数错误")
		return
	}
	if userID == req.FriendID {
		utils.ErrorResponse(c, 400, "不能添加自己为好友")
		return
	}
	var target models.User
	if config.DB.First(&target, req.FriendID).Error != nil {
		utils.ErrorResponse(c, 404, "用户不存在")
		return
	}
	// 检查是否已是好友
	var existing models.Friend
	if config.DB.Where("user_id = ? AND friend_id = ?", userID, req.FriendID).First(&existing).Error == nil {
		utils.ErrorResponse(c, 400, "已经是好友")
		return
	}
	friend := models.Friend{UserID: userID, FriendID: req.FriendID, Remark: req.Remark, Status: 0}
	config.DB.Create(&friend)
	// 双向好友
	friend2 := models.Friend{UserID: req.FriendID, FriendID: userID, Status: 0}
	config.DB.Create(&friend2)
	utils.SuccessResponse(c, nil)
}

// DeleteFriend 删除好友
func DeleteFriend(c *gin.Context) {
	userID := c.GetUint("user_id")
	friendID, _ := strconv.ParseUint(c.Param("friend_id"), 10, 32)
	config.DB.Where("user_id = ? AND friend_id = ?", userID, friendID).Delete(&models.Friend{})
	config.DB.Where("user_id = ? AND friend_id = ?", friendID, userID).Delete(&models.Friend{})
	utils.SuccessResponse(c, nil)
}
