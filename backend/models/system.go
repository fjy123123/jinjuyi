package models

import (
	"time"

	"gorm.io/gorm"
)

type SystemConfig struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	AppName         string         `gorm:"size:100;default:知信" json:"app_name"`
	AppVersion      string         `gorm:"size:50;default:v1.0.0" json:"app_version"`
	AppDescription  string         `gorm:"size:500" json:"app_description"`
	LogoURL         string         `gorm:"size:500" json:"logo_url"`
	FaviconURL      string         `gorm:"size:500" json:"favicon_url"`
	ThemeColor      string         `gorm:"size:50;default:#07c160" json:"theme_color"`
	ThemeSecondary  string         `gorm:"size:50;default:#576b95" json:"theme_secondary"`
	UiTemplate      string         `gorm:"size:50;default:modern" json:"ui_template"`
	MaintenanceMode bool           `gorm:"default:false" json:"maintenance_mode"`
	MaintenanceMsg  string         `gorm:"size:500" json:"maintenance_msg"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
