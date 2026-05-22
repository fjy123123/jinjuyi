package handlers

import (
	"strconv"

	"chat-system-pro/config"
	"chat-system-pro/models"
	"chat-system-pro/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateGroup 创建群
func CreateGroup(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		Name        string `json:"name" binding:"required"`
		Avatar      string `json:"avatar"`
		Description string `json:"description"`
		MemberIDs   []uint `json:"member_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "参数错误")
		return
	}
	group := models.Group{
		Name: req.Name, Avatar: req.Avatar, Description: req.Description,
		OwnerID: userID, MemberCount: 1,
	}
	if err := config.DB.Create(&group).Error; err != nil {
		utils.ErrorResponse(c, 500, "创建群失败")
		return
	}
	// 群主加入
	config.DB.Create(&models.GroupMember{GroupID: group.ID, UserID: userID, Role: 2, JoinAt: group.CreatedAt})
	// 其他成员
	for _, mid := range req.MemberIDs {
		if mid != userID {
			config.DB.Create(&models.GroupMember{GroupID: group.ID, UserID: mid, Role: 0})
			config.DB.Model(&group).Update("member_count", gorm.Expr("member_count + 1"))
		}
	}
	utils.SuccessResponse(c, group)
}

// GetMyGroups 获取我的群
func GetMyGroups(c *gin.Context) {
	userID := c.GetUint("user_id")
	var memberIDs []uint
	config.DB.Model(&models.GroupMember{}).Where("user_id = ?", userID).Pluck("group_id", &memberIDs)
	var groups []models.Group
	config.DB.Preload("Owner").Where("id IN ?", memberIDs).Find(&groups)
	utils.SuccessResponse(c, groups)
}

// GetGroupInfo 获取群信息
func GetGroupInfo(c *gin.Context) {
	groupID, _ := strconv.ParseUint(c.Param("group_id"), 10, 32)
	var group models.Group
	if err := config.DB.Preload("Owner").First(&group, groupID).Error; err != nil {
		utils.ErrorResponse(c, 404, "群不存在")
		return
	}
	utils.SuccessResponse(c, group)
}

// UpdateGroup 更新群信息
func UpdateGroup(c *gin.Context) {
	userID := c.GetUint("user_id")
	groupID, _ := strconv.ParseUint(c.Param("group_id"), 10, 32)
	var group models.Group
	config.DB.First(&group, groupID)
	if group.OwnerID != userID {
		utils.ErrorResponse(c, 403, "只有群主可以修改")
		return
	}
	var req struct {
		Name         string `json:"name"`
		Avatar       string `json:"avatar"`
		Description  string `json:"description"`
		Announcement string `json:"announcement"`
		JoinMode     int    `json:"join_mode"`
		MaxMembers   int    `json:"max_members"`
		IsMuteAll    *bool  `json:"is_mute_all"`
	}
	c.ShouldBindJSON(&req)
	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Announcement != "" {
		updates["announcement"] = req.Announcement
	}
	if req.JoinMode > 0 {
		updates["join_mode"] = req.JoinMode
	}
	if req.MaxMembers > 0 {
		updates["max_members"] = req.MaxMembers
	}
	if req.IsMuteAll != nil {
		updates["is_mute_all"] = *req.IsMuteAll
	}
	config.DB.Model(&group).Updates(updates)
	utils.SuccessResponse(c, nil)
}

// GetGroupMembers 获取群成员
func GetGroupMembers(c *gin.Context) {
	groupID, _ := strconv.ParseUint(c.Param("group_id"), 10, 32)
	var members []models.GroupMember
	config.DB.Preload("User").Where("group_id = ?", groupID).Find(&members)
	utils.SuccessResponse(c, members)
}

// InviteGroupMember 邀请入群
func InviteGroupMember(c *gin.Context) {
	userID := c.GetUint("user_id")
	groupID, _ := strconv.ParseUint(c.Param("group_id"), 10, 32)
	var req struct {
		UserIDs []uint `json:"user_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "参数错误")
		return
	}
	var group models.Group
	config.DB.First(&group, groupID)
	if !group.AllowInvite && group.OwnerID != userID {
		utils.ErrorResponse(c, 403, "当前不允许邀请")
		return
	}
	for _, uid := range req.UserIDs {
		var existing models.GroupMember
		if config.DB.Where("group_id = ? AND user_id = ?", groupID, uid).First(&existing).Error != nil {
			config.DB.Create(&models.GroupMember{GroupID: uint(groupID), UserID: uid, Role: 0})
			config.DB.Model(&group).Update("member_count", gorm.Expr("member_count + 1"))
		}
	}
	utils.SuccessResponse(c, nil)
}

// RemoveGroupMember 踢出成员
func RemoveGroupMember(c *gin.Context) {
	userID := c.GetUint("user_id")
	groupID, _ := strconv.ParseUint(c.Param("group_id"), 10, 32)
	memberID, _ := strconv.ParseUint(c.Param("member_id"), 10, 32)
	var group models.Group
	config.DB.First(&group, groupID)
	if group.OwnerID != userID {
		utils.ErrorResponse(c, 403, "只有群主可以踢人")
		return
	}
	if uint(memberID) == group.OwnerID {
		utils.ErrorResponse(c, 400, "不能踢出群主")
		return
	}
	config.DB.Where("group_id = ? AND user_id = ?", groupID, memberID).Delete(&models.GroupMember{})
	config.DB.Model(&group).Update("member_count", gorm.Expr("member_count - 1"))
	utils.SuccessResponse(c, nil)
}

// MuteGroupMember 禁言成员
func MuteGroupMember(c *gin.Context) {
	userID := c.GetUint("user_id")
	groupID, _ := strconv.ParseUint(c.Param("group_id"), 10, 32)
	memberID, _ := strconv.ParseUint(c.Param("member_id"), 10, 32)
	var group models.Group
	config.DB.First(&group, groupID)
	if group.OwnerID != userID {
		utils.ErrorResponse(c, 403, "只有群主可以禁言")
		return
	}
	config.DB.Model(&models.GroupMember{}).Where("group_id = ? AND user_id = ?", groupID, memberID).Update("is_mute", true)
	utils.SuccessResponse(c, nil)
}

// LeaveGroup 退出群
func LeaveGroup(c *gin.Context) {
	userID := c.GetUint("user_id")
	groupID, _ := strconv.ParseUint(c.Param("group_id"), 10, 32)
	config.DB.Where("group_id = ? AND user_id = ?", groupID, userID).Delete(&models.GroupMember{})
	config.DB.Model(&models.Group{}).Where("id = ?", groupID).Update("member_count", gorm.Expr("member_count - 1"))
	utils.SuccessResponse(c, nil)
}
