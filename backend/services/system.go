package services

import (
	"chat-system-pro/backend/config"
	"chat-system-pro/backend/models"

	"gorm.io/gorm"
)

type SystemConfigService struct {
	db *gorm.DB
}

func NewSystemConfigService(db *gorm.DB) *SystemConfigService {
	return &SystemConfigService{db: db}
}

func (s *SystemConfigService) GetConfig() (*models.SystemConfig, error) {
	var config models.SystemConfig
	err := s.db.First(&config).Error
	if err == gorm.ErrRecordNotFound {
		return s.CreateDefaultConfig()
	}
	return &config, err
}

func (s *SystemConfigService) CreateDefaultConfig() (*models.SystemConfig, error) {
	defaultConfig := &models.SystemConfig{
		AppName:        "知信",
		AppVersion:     "v1.0.0",
		AppDescription: "知信 - 让沟通更简单",
		ThemeColor:     "#07c160",
		ThemeSecondary: "#576b95",
		UiTemplate:     "modern",
		ExportEnabled:  true,
		ExportMaxRecords: 1000,
	}

	err := s.db.Create(defaultConfig).Error
	return defaultConfig, err
}

func (s *SystemConfigService) UpdateConfig(update *models.SystemConfig) (*models.SystemConfig, error) {
	var config models.SystemConfig
	err := s.db.First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return s.CreateDefaultConfig()
		}
		return nil, err
	}

	if update.AppName != "" {
		config.AppName = update.AppName
	}
	if update.AppVersion != "" {
		config.AppVersion = update.AppVersion
	}
	if update.AppDescription != "" {
		config.AppDescription = update.AppDescription
	}
	if update.LogoURL != "" {
		config.LogoURL = update.LogoURL
	}
	if update.FaviconURL != "" {
		config.FaviconURL = update.FaviconURL
	}
	if update.ThemeColor != "" {
		config.ThemeColor = update.ThemeColor
	}
	if update.ThemeSecondary != "" {
		config.ThemeSecondary = update.ThemeSecondary
	}
	if update.UiTemplate != "" {
		config.UiTemplate = update.UiTemplate
	}
	config.MaintenanceMode = update.MaintenanceMode
	if update.MaintenanceMsg != "" {
		config.MaintenanceMsg = update.MaintenanceMsg
	}
	config.ExportEnabled = update.ExportEnabled
	if update.ExportMaxRecords > 0 {
		config.ExportMaxRecords = update.ExportMaxRecords
	}

	err = s.db.Save(&config).Error
	return &config, err
}

func (s *SystemConfigService) UploadLogo(filePath string) (string, error) {
	url, err := StorageService.UploadFile(filePath, "logo")
	if err != nil {
		return "", err
	}

	var config models.SystemConfig
	err = s.db.First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			config = models.SystemConfig{LogoURL: url}
			err = s.db.Create(&config).Error
		}
		return url, err
	}

	config.LogoURL = url
	err = s.db.Save(&config).Error
	return url, err
}

func (s *SystemConfigService) UploadFavicon(filePath string) (string, error) {
	url, err := StorageService.UploadFile(filePath, "favicon")
	if err != nil {
		return "", err
	}

	var config models.SystemConfig
	err = s.db.First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			config = models.SystemConfig{FaviconURL: url}
			err = s.db.Create(&config).Error
		}
		return url, err
	}

	config.FaviconURL = url
	err = s.db.Save(&config).Error
	return url, err
}

func (s *SystemConfigService) IsExportEnabled() bool {
	var config models.SystemConfig
	err := s.db.First(&config).Error
	if err != nil {
		return true
	}
	return config.ExportEnabled
}

func (s *SystemConfigService) GetExportMaxRecords() int {
	var config models.SystemConfig
	err := s.db.First(&config).Error
	if err != nil {
		return 1000
	}
	return config.ExportMaxRecords
}

var SystemConfigService = NewSystemConfigService(config.DB)
