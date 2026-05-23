package models

import (
	"time"

	"gorm.io/gorm"
)

type SystemConfig struct {
	ID                    uint           `gorm:"primarykey" json:"id"`
	AppName               string         `gorm:"size:100;default:知信" json:"app_name"`
	AppVersion            string         `gorm:"size:50;default:v1.0.0" json:"app_version"`
	AppDescription        string         `gorm:"size:500" json:"app_description"`
	LogoURL               string         `gorm:"size:500" json:"logo_url"`
	FaviconURL            string         `gorm:"size:500" json:"favicon_url"`
	ThemeColor            string         `gorm:"size:50;default:#07c160" json:"theme_color"`
	ThemeSecondary        string         `gorm:"size:50;default:#576b95" json:"theme_secondary"`
	UiTemplate            string         `gorm:"size:50;default:modern" json:"ui_template"`
	MaintenanceMode       bool           `gorm:"default:false" json:"maintenance_mode"`
	MaintenanceMsg        string         `gorm:"size:500" json:"maintenance_msg"`
	ExportEnabled         bool           `gorm:"default:true" json:"export_enabled"`
	ExportMaxRecords      int            `gorm:"default:1000" json:"export_max_records"`
	
	// 消息自动清理配置（后台控制）
	AutoDeleteEnabled     bool           `gorm:"default:false" json:"auto_delete_enabled"`      // 是否开启消息自动清理
	AutoDeleteDays        int            `gorm:"default:3" json:"auto_delete_days"`              // 自动清理天数
	AutoDeleteLastRun     *time.Time     `json:"auto_delete_last_run"`                            // 上次清理时间
	
	// 消息撤回时限配置（后台控制）
	RecallEnabled         bool           `gorm:"default:true" json:"recall_enabled"`             // 是否允许消息撤回
	RecallTimeout         int            `gorm:"default:300" json:"recall_timeout"`              // 撤回时限（秒），默认5分钟=300秒
	
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`
}
