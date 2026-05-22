package main

import (
	"fmt"
	"log"

	"chat-system-pro/config"
	"chat-system-pro/handlers"
	"chat-system-pro/middleware"
	"chat-system-pro/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.InitConfig()

	// 设置Gin模式
	gin.SetMode(config.Cfg.Server.Mode)

	// 初始化数据库
	config.InitDatabase()
	config.InitMongoDB()
	config.InitRedis()

	// 自动迁移
	models.AutoMigrate(config.DB)

	r := gin.Default()

	// 安全中间件
	r.Use(middleware.SecurityHeadersMiddleware())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")

	// 认证接口（无需登录）
	auth := api.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	// 需要登录的接口
	authed := api.Group("")
	authed.Use(middleware.AuthMiddleware())
	{
		// WebSocket
		authed.GET("/ws", handlers.WebSocketHandler)

		// 用户
		authed.GET("/users/me", handlers.GetProfile)
		authed.PUT("/users/profile", handlers.UpdateProfile)
		authed.GET("/users/settings", handlers.GetUserSettings)
		authed.PUT("/users/settings", handlers.UpdateUserSettings)
		authed.GET("/users/search", handlers.SearchUsers)

		// 好友
		authed.GET("/friends", handlers.GetFriends)
		authed.POST("/friends", handlers.AddFriend)
		authed.DELETE("/friends/:friend_id", handlers.DeleteFriend)

		// 会话
		authed.GET("/conversations", handlers.GetConversations)
		authed.GET("/conversations/unread", handlers.GetUnreadCount)

		// 消息
		authed.POST("/messages", handlers.SendMessage)
		authed.GET("/messages/private/:friend_id", handlers.GetPrivateMessages)
		authed.GET("/messages/group/:group_id", handlers.GetGroupMessages)
		authed.POST("/messages/read", handlers.MarkAsRead)
		authed.POST("/messages/:message_id/recall", handlers.RecallMessage)

		// 群组
		authed.POST("/groups", handlers.CreateGroup)
		authed.GET("/groups", handlers.GetMyGroups)
		authed.GET("/groups/:group_id", handlers.GetGroupInfo)
		authed.PUT("/groups/:group_id", handlers.UpdateGroup)
		authed.GET("/groups/:group_id/members", handlers.GetGroupMembers)
		authed.POST("/groups/:group_id/invite", handlers.InviteGroupMember)
		authed.DELETE("/groups/:group_id/members/:member_id", handlers.RemoveGroupMember)
		authed.POST("/groups/:group_id/members/:member_id/mute", handlers.MuteGroupMember)
		authed.POST("/groups/:group_id/leave", handlers.LeaveGroup)

		// 红包
		authed.POST("/redpackets", handlers.SendRedPacket)
		authed.POST("/redpackets/:id/grab", handlers.GrabRedPacket)
		authed.GET("/redpackets/:id", handlers.GetRedPacket)
		authed.GET("/redpackets/sent", handlers.GetSentRedPackets)
		authed.GET("/redpackets/received", handlers.GetReceivedRedPackets)

		// 朋友圈
		authed.POST("/moments", handlers.PublishMoment)
		authed.GET("/moments", handlers.GetMoments)
		authed.POST("/moments/:moment_id/like", handlers.LikeMoment)
		authed.DELETE("/moments/:moment_id/like", handlers.UnlikeMoment)
		authed.POST("/moments/:moment_id/comments", handlers.CommentMoment)
		authed.DELETE("/moments/:moment_id", handlers.DeleteMoment)

		// 支付
		authed.POST("/payment/orders", handlers.CreateOrder)
		authed.POST("/payment/orders/:order_id/pay", handlers.ProcessPayment)
		authed.GET("/payment/orders", handlers.GetOrders)
		authed.GET("/payment/points/history", handlers.GetPointsHistory)
	}

	// 管理员接口
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware())
	{
		admin.GET("/db/stats", handlers.GetDatabaseStats)
		admin.POST("/db/clear-old-messages", handlers.ClearOldMessages)
		admin.POST("/db/clear-all", handlers.ClearAllData)
		admin.POST("/db/init", handlers.InitializeDatabase)
		admin.POST("/db/archive-old", handlers.ArchiveOldMessages)
		admin.DELETE("/db/users/:user_id", handlers.DeleteUserAndData)
		admin.DELETE("/db/groups/:group_id", handlers.DeleteGroupAndData)
		admin.POST("/users/points", handlers.AdminAddPoints)
		admin.GET("/system/configs", handlers.GetSystemConfigs)
		admin.PUT("/system/configs", handlers.UpdateSystemConfig)
	}

	log.Printf("Server starting on port %d", config.Cfg.Server.Port)
	if err := r.Run(fmt.Sprintf(":%d", config.Cfg.Server.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
