package handlers

import (
	"chat-system-pro/config"
	"chat-system-pro/models"
	"chat-system-pro/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// EmojiCategory 表情包分类
type EmojiCategory struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

// EmojiItem 表情包项目
type EmojiItem struct {
	ID         uint   `json:"id"`
	CategoryID uint   `json:"category_id"`
	Name       string `json:"name"`
	URL        string `json:"url"`
	Sort       int    `json:"sort"`
}

// GetEmojiCategories 获取表情包分类列表
func GetEmojiCategories(c *gin.Context) {
	var categories []EmojiCategory
	config.DB.Model(&EmojiCategory{}).Order("sort ASC, id ASC").Find(&categories)
	utils.SuccessResponse(c, categories)
}

// GetEmojisByCategory 获取分类下的表情包
func GetEmojisByCategory(c *gin.Context) {
	categoryIDStr := c.Param("category_id")
	categoryID, _ := strconv.ParseUint(categoryIDStr, 10, 32)

	var emojis []EmojiItem
	query := config.DB.Model(&EmojiItem{})
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	query.Order("sort ASC, id ASC").Find(&emojis)

	utils.SuccessResponse(c, emojis)
}

// AddEmojiCategory 添加表情包分类（管理员）
func AddEmojiCategory(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
		Icon string `json:"icon"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "参数错误")
		return
	}

	category := EmojiCategory{
		Name: req.Name,
		Icon: req.Icon,
	}
	if err := config.DB.Create(&category).Error; err != nil {
		utils.ErrorResponse(c, 500, "创建失败")
		return
	}

	utils.SuccessResponse(c, category)
}

// AddEmojiItem 添加表情包项目（管理员）
func AddEmojiItem(c *gin.Context) {
	var req struct {
		CategoryID uint   `json:"category_id" binding:"required"`
		Name       string `json:"name" binding:"required"`
		URL        string `json:"url" binding:"required"`
		Sort       int    `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "参数错误")
		return
	}

	item := EmojiItem{
		CategoryID: req.CategoryID,
		Name:       req.Name,
		URL:        req.URL,
		Sort:       req.Sort,
	}
	if err := config.DB.Create(&item).Error; err != nil {
		utils.ErrorResponse(c, 500, "创建失败")
		return
	}

	utils.SuccessResponse(c, item)
}

// DeleteEmojiCategory 删除表情包分类（管理员）
func DeleteEmojiCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)
	
	if err := config.DB.Delete(&EmojiCategory{}, id).Error; err != nil {
		utils.ErrorResponse(c, 500, "删除失败")
		return
	}
	
	config.DB.Where("category_id = ?", id).Delete(&EmojiItem{})
	
	utils.SuccessResponse(c, nil)
}

// DeleteEmojiItem 删除表情包项目（管理员）
func DeleteEmojiItem(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)
	
	if err := config.DB.Delete(&EmojiItem{}, id).Error; err != nil {
		utils.ErrorResponse(c, 500, "删除失败")
		return
	}
	
	utils.SuccessResponse(c, nil)
}
