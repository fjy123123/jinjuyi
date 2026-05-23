package handlers

import (
	"chat-system-pro/models"
	"chat-system-pro/services"

	"github.com/gin-gonic/gin"
)

func GetSystemConfig(c *gin.Context) {
	config, err := services.SystemConfigService.GetConfig()
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "获取配置失败", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 200, "message": "success", "data": config})
}

func UpdateSystemConfigForSystem(c *gin.Context) {
	var update models.SystemConfig
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	config, err := services.SystemConfigService.UpdateConfig(&update)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "更新配置失败", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 200, "message": "更新成功", "data": config})
}

func UploadLogo(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "请选择文件"})
		return
	}

	if file.Size > 2*1024*1024 {
		c.JSON(400, gin.H{"code": 400, "message": "文件大小不能超过2MB"})
		return
	}

	tempPath := "/tmp/" + file.Filename
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "保存文件失败", "error": err.Error()})
		return
	}

	url, err := services.SystemConfigService.UploadLogo(tempPath)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "上传Logo失败", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 200, "message": "上传成功", "data": gin.H{"url": url}})
}

func UploadFavicon(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "请选择文件"})
		return
	}

	if file.Size > 500*1024 {
		c.JSON(400, gin.H{"code": 400, "message": "文件大小不能超过500KB"})
		return
	}

	tempPath := "/tmp/" + file.Filename
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "保存文件失败", "error": err.Error()})
		return
	}

	url, err := services.SystemConfigService.UploadFavicon(tempPath)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "上传图标失败", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 200, "message": "上传成功", "data": gin.H{"url": url}})
}

func SetMaintenanceMode(c *gin.Context) {
	var req struct {
		Mode   bool   `json:"mode"`
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	update := &models.SystemConfig{
		MaintenanceMode: req.Mode,
		MaintenanceMsg:  req.Message,
	}

	_, err := services.SystemConfigService.UpdateConfig(update)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "设置失败", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 200, "message": "设置成功", "data": gin.H{"maintenance_mode": req.Mode}})
}
