package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"chat-system-pro/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// ==================== 文件存储服务 ====================

// FileStorage 文件存储服务
type FileStorage struct {
	storageType string // local, oss, s3
	s3Client    *s3.S3
	ossClient   *oss.Client
}

// NewFileStorage 创建文件存储服务
func NewFileStorage() *FileStorage {
	storageType := config.GetString("storage.type", "local")
	
	fs := &FileStorage{
		storageType: storageType,
	}
	
	if storageType == "s3" {
		// AWS S3 配置
		accessKey := config.GetString("storage.s3.access_key", "")
		secretKey := config.GetString("storage.s3.secret_key", "")
		region := config.GetString("storage.s3.region", "us-east-1")
		endpoint := config.GetString("storage.s3.endpoint", "")
		
		sess, err := session.NewSession(&aws.Config{
			Region:           aws.String(region),
			Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
			Endpoint:         aws.String(endpoint),
			S3ForcePathStyle: aws.Bool(true),
		})
		if err == nil {
			fs.s3Client = s3.New(sess)
		}
	} else if storageType == "oss" {
		// 阿里云 OSS 配置
		endpoint := config.GetString("storage.oss.endpoint", "")
		accessKey := config.GetString("storage.oss.access_key", "")
		secretKey := config.GetString("storage.oss.secret_key", "")
		
		client, err := oss.New(endpoint, accessKey, secretKey)
		if err == nil {
			fs.ossClient = client
		}
	}
	
	return fs
}

// UploadFile 上传文件
func (fs *FileStorage) UploadFile(file *multipart.FileHeader, fileName string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	
	switch fs.storageType {
	case "local":
		return fs.uploadLocal(src, fileName, file.Size)
	case "s3":
		return fs.uploadS3(src, fileName)
	case "oss":
		return fs.uploadOSS(src, fileName)
	default:
		return fs.uploadLocal(src, fileName, file.Size)
	}
}

// uploadLocal 本地上传
func (fs *FileStorage) uploadLocal(src io.Reader, fileName string, size int64) (string, error) {
	uploadPath := config.GetString("upload.path", "./uploads")
	filePath := filepath.Join(uploadPath, fileName)
	
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}
	
	return fmt.Sprintf("/static/uploads/%s", fileName), nil
}

// uploadS3 S3上传
func (fs *FileStorage) uploadS3(src io.ReadSeeker, fileName string) (string, error) {
	if fs.s3Client == nil {
		return "", errors.New("S3 client not initialized")
	}
	
	bucket := config.GetString("storage.s3.bucket", "")
	fileKey := fmt.Sprintf("chat/%s", fileName)
	
	_, err := fs.s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
		Body:   src,
	})
	
	if err != nil {
		return "", err
	}
	
	// 返回URL
	domain := config.GetString("storage.s3.domain", "")
	if domain != "" {
		return fmt.Sprintf("%s/%s", domain, fileKey), nil
	}
	
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, fileKey), nil
}

// uploadOSS OSS上传
func (fs *FileStorage) uploadOSS(src io.ReadSeeker, fileName string) (string, error) {
	if fs.ossClient == nil {
		return "", errors.New("OSS client not initialized")
	}
	
	bucketName := config.GetString("storage.oss.bucket", "")
	fileKey := fmt.Sprintf("chat/%s", fileName)
	
	bucket, err := fs.ossClient.Bucket(bucketName)
	if err != nil {
		return "", err
	}
	
	err = bucket.PutObject(fileKey, src)
	if err != nil {
		return "", err
	}
	
	// 返回URL
	domain := config.GetString("storage.oss.domain", "")
	if domain != "" {
		return fmt.Sprintf("%s/%s", domain, fileKey), nil
	}
	
	return fmt.Sprintf("https://%s.oss-cn-hangzhou.aliyuncs.com/%s", bucketName, fileKey), nil
}

// DeleteFile 删除文件
func (fs *FileStorage) DeleteFile(fileURL string) error {
	// 解析文件Key
	fileKey := fs.getFileKey(fileURL)
	if fileKey == "" {
		return errors.New("invalid file URL")
	}
	
	switch fs.storageType {
	case "s3":
		bucket := config.GetString("storage.s3.bucket", "")
		_, err := fs.s3Client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(fileKey),
		})
		return err
	case "oss":
		bucketName := config.GetString("storage.oss.bucket", "")
		bucket, err := fs.ossClient.Bucket(bucketName)
		if err != nil {
			return err
		}
		return bucket.DeleteObject(fileKey)
	default:
		// 本地删除
		filePath := filepath.Join(config.GetString("upload.path", "./uploads"), filepath.Base(fileKey))
		return os.Remove(filePath)
	}
}

func (fs *FileStorage) getFileKey(url string) string {
	// 从URL提取文件Key
	if idx := strings.Index(url, "/chat/"); idx != -1 {
		return url[idx+1:]
	}
	return ""
}

// ==================== 推送服务 ====================

// PushService 推送服务
type PushService struct {
	provider string
}

// NewPushService 创建推送服务
func NewPushService() *PushService {
	return &PushService{
		provider: config.GetString("push.provider", "jpush"),
	}
}

// PushToUser 推送消息给用户
func (ps *PushService) PushToUser(userID uint, title, content string, extra map[string]interface{}) error {
	switch ps.provider {
	case "jpush":
		return ps.pushJPush(userID, title, content, extra)
	case "getui":
		return ps.pushGeTui(userID, title, content, extra)
	default:
		return fmt.Errorf("push provider not supported")
	}
}

// PushToAll 推送给所有用户
func (ps *PushService) PushToAll(title, content string, extra map[string]interface{}) error {
	switch ps.provider {
	case "jpush":
		return ps.pushJPushAll(title, content, extra)
	case "getui":
		return ps.pushGeTuiAll(title, content, extra)
	default:
		return fmt.Errorf("push provider not supported")
	}
}

func (ps *PushService) pushJPush(userID uint, title, content string, extra map[string]interface{}) error {
	// 极光推送实现
	jpushKey := config.GetString("push.jpush.app_key", "")
	jpushSecret := config.GetString("push.jpush.app_secret", "")
	
	if jpushKey == "" || jpushSecret == "" {
		return errors.New("jpush config not set")
	}
	
	// TODO: 实现实际的推送调用
	fmt.Printf("JPush to user %d: %s\n", userID, content)
	return nil
}

func (ps *PushService) pushJPushAll(title, content string, extra map[string]interface{}) error {
	fmt.Printf("JPush all: %s\n", content)
	return nil
}

func (ps *PushService) pushGeTui(userID uint, title, content string, extra map[string]interface{}) error {
	// 个推实现
	fmt.Printf("GeTui to user %d: %s\n", userID, content)
	return nil
}

func (ps *PushService) pushGeTuiAll(title, content string, extra map[string]interface{}) error {
	fmt.Printf("GeTui all: %s\n", content)
	return nil
}

// ==================== 多端同步 ====================

// SyncService 多端同步服务
type SyncService struct {
}

// NewSyncService 创建同步服务
func NewSyncService() *SyncService {
	return &SyncService{}
}

// SyncMessage 同步消息
func (ss *SyncService) SyncMessage(msg interface{}, excludeDeviceID string) {
	// 通过Redis发布消息同步
	ctx := context.Background()
	
	syncData := map[string]interface{}{
		"type":      "message",
		"data":      msg,
		"timestamp": time.Now().UnixMilli(),
		"exclude_device": excludeDeviceID,
	}
	
	syncJSON, _ := config.JSONMarshal(syncData)
	config.RedisPublish(ctx, "chat:sync", syncJSON)
}

// GetSyncMessages 获取同步消息
func (ss *SyncService) GetSyncMessages(deviceID string, lastSync time.Time) []interface{} {
	// TODO: 实现获取离线同步消息
	return []interface{}{}
}

// ==================== 支付服务扩展 ====================

// WechatPayService 微信支付
type WechatPayService struct {
}

// NewWechatPayService 创建微信支付
func NewWechatPayService() *WechatPayService {
	return &WechatPayService{}
}

// CreateOrder 创建微信订单
func (wp *WechatPayService) CreateOrder(orderNo string, amount int, desc string) (map[string]interface{}, error) {
	// 实现微信支付
	result := map[string]interface{}{
		"prepay_id": "wx_prepay_id",
		"timestamp": time.Now().Unix(),
		"nonce":     "random_nonce",
		"sign":      "generated_sign",
	}
	
	return result, nil
}

// AlipayService 支付宝
type AlipayService struct {
}

// NewAlipayService 创建支付宝
func NewAlipayService() *AlipayService {
	return &AlipayService{}
}

// CreateOrder 创建支付宝订单
func (ap *AlipayService) CreateOrder(orderNo string, amount int, desc string) (string, error) {
	// 实现支付宝支付
	return "alipay_order_form", nil
}
