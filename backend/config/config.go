package config

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	MongoDB   MongoDBConfig
	Redis     RedisConfig
	JWT       JWTConfig
	Upload    UploadConfig
	WebSocket WebSocketConfig
	Payment   PaymentConfig
	Security  SecurityConfig
	System    SystemConfig
}

type ServerConfig struct {
	Port int
	Mode string
	Name string
}

type DatabaseConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string
	MaxOpenConns int
	MaxIdleConns int
}

type MongoDBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

type UploadConfig struct {
	Path    string
	MaxSize int64
}

type WebSocketConfig struct {
	ReadBufferSize  int
	WriteBufferSize int
	PingInterval    int
}

type PaymentConfig struct {
	StripeSecretKey   string
	WechatPayEnabled  bool
	AlipayEnabled     bool
	WechatMchID       string
	WechatAPIv3Key    string
	AlipayAppID       string
	AlipayPrivateKey  string
	AlipayPublicKey   string
}

type SecurityConfig struct {
	InviteCodeEnabled bool
	CaptchaEnabled    bool
	RateLimitEnabled  bool
}

type SystemConfig struct {
	UIDefault string
}

var (
	DB                *gorm.DB
	MongoClient       *mongo.Client
	MongoDBCollection *mongo.Database
	RDB               *redis.Client
	Cfg               *Config
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	Cfg = &Config{
		Server: ServerConfig{
			Port: viper.GetInt("server.port"),
			Mode: viper.GetString("server.mode"),
			Name: viper.GetString("server.name"),
		},
		Database: DatabaseConfig{
			Host:         viper.GetString("database.host"),
			Port:         viper.GetInt("database.port"),
			User:         viper.GetString("database.user"),
			Password:     viper.GetString("database.password"),
			DBName:       viper.GetString("database.dbname"),
			MaxOpenConns: viper.GetInt("database.max_open_conns"),
			MaxIdleConns: viper.GetInt("database.max_idle_conns"),
		},
		MongoDB: MongoDBConfig{
			Host:     viper.GetString("mongodb.host"),
			Port:     viper.GetInt("mongodb.port"),
			User:     viper.GetString("mongodb.user"),
			Password: viper.GetString("mongodb.password"),
			DBName:   viper.GetString("mongodb.dbname"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("redis.host"),
			Port:     viper.GetInt("redis.port"),
			Password: viper.GetString("redis.password"),
			DB:       viper.GetInt("redis.db"),
		},
		JWT: JWTConfig{
			Secret:      viper.GetString("jwt.secret"),
			ExpireHours: viper.GetInt("jwt.expire_hours"),
		},
		Upload: UploadConfig{
			Path:    viper.GetString("upload.path"),
			MaxSize: viper.GetInt64("upload.max_size"),
		},
		WebSocket: WebSocketConfig{
			ReadBufferSize:  viper.GetInt("websocket.read_buffer_size"),
			WriteBufferSize: viper.GetInt("websocket.write_buffer_size"),
			PingInterval:   viper.GetInt("websocket.ping_interval"),
		},
		Payment: PaymentConfig{
			StripeSecretKey:  viper.GetString("payment.stripe_secret_key"),
			WechatPayEnabled: viper.GetBool("payment.wechat_pay_enabled"),
			AlipayEnabled:    viper.GetBool("payment.alipay_enabled"),
			WechatMchID:      viper.GetString("payment.wechat_mch_id"),
			WechatAPIv3Key:   viper.GetString("payment.wechat_apiv3_key"),
			AlipayAppID:      viper.GetString("payment.alipay_app_id"),
			AlipayPrivateKey: viper.GetString("payment.alipay_private_key"),
			AlipayPublicKey:  viper.GetString("payment.alipay_public_key"),
		},
		Security: SecurityConfig{
			InviteCodeEnabled: viper.GetBool("security.invite_code_enabled"),
			CaptchaEnabled:    viper.GetBool("security.captcha_enabled"),
			RateLimitEnabled:  viper.GetBool("security.rate_limit_enabled"),
		},
		System: SystemConfig{
			UIDefault: viper.GetString("system.ui_default"),
		},
	}
}

func InitDatabase() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		Cfg.Database.User, Cfg.Database.Password,
		Cfg.Database.Host, Cfg.Database.Port, Cfg.Database.DBName,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	sqlDB.SetMaxOpenConns(Cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(Cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("MySQL connection established")
}

func InitMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var uri string
	if Cfg.MongoDB.User != "" && Cfg.MongoDB.Password != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin",
			Cfg.MongoDB.User, Cfg.MongoDB.Password,
			Cfg.MongoDB.Host, Cfg.MongoDB.Port, Cfg.MongoDB.DBName)
	} else {
		uri = fmt.Sprintf("mongodb://%s:%d", Cfg.MongoDB.Host, Cfg.MongoDB.Port)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	MongoClient = client
	MongoDBCollection = client.Database(Cfg.MongoDB.DBName)
	log.Println("MongoDB connection established")
}

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", Cfg.Redis.Host, Cfg.Redis.Port),
		Password: Cfg.Redis.Password,
		DB:       Cfg.Redis.DB,
	})
	log.Println("Redis connection established")
}

// ==================== 辅助函数 ====================

// GetString 获取字符串配置
func GetString(key string, defaultVal string) string {
	return viper.GetString(key)
}

// GetInt 获取整数配置
func GetInt(key string, defaultVal int) int {
	v := viper.GetInt(key)
	if v == 0 {
		return defaultVal
	}
	return v
}

// GenerateShortID 生成短随机ID
func GenerateShortID(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		result[i] = chars[n.Int64()]
	}
	return string(result)
}

// JSONMarshal JSON序列化
func JSONMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// RedisPublish Redis发布消息
func RedisPublish(ctx context.Context, channel string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return RDB.Publish(ctx, channel, data).Err()
}
