package main

import (
	"fmt"
	"log"

	"chat-system-pro/config"
	"chat-system-pro/handlers"
	"chat-system-pro/middleware"
	"chat-system-pro/models"
	"chat-system-pro/services"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()
	gin.SetMode(config.Cfg.Server.Mode)

	config.InitDatabase()
	config.InitMongoDB()
	config.InitRedis()

	models.AutoMigrate(config.DB)

	dbService := services.NewDatabaseService()
	dbService.EnsureAdminUser()

	// 启动消息清理定时任务
	services.MessageCleanup.StartCleanupScheduler()
	log.Println("消息清理定时任务已启动")

	r := gin.Default()

	r.Use(middleware.SecurityHeadersMiddleware())
	r.Use(middleware.CORSMiddleware())

	auditLogger := middleware.NewAuditLogger()
	r.Use(auditLogger.Middleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	api.Use(middleware.APIRateLimiter())

	auth := api.Group("/auth")
	{
		auth.Use(middleware.LoginRateLimiter())
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
		auth.POST("/refresh", handlers.RefreshToken)
		auth.POST("/logout", handlers.Logout)
	}

	authed := api.Group("")
	authed.Use(middleware.AuthMiddleware())
	{
		authed.GET("/ws", handlers.WebSocketHandler)

		authed.GET("/users/me", handlers.GetProfile)
		authed.PUT("/users/profile", handlers.UpdateProfile)
		authed.GET("/users/settings", handlers.GetUserSettings)
		authed.PUT("/users/settings", handlers.UpdateUserSettings)
		authed.GET("/users/search", handlers.SearchUsers)

		authed.GET("/friends", handlers.GetFriends)
		authed.POST("/friends", handlers.AddFriend)
		authed.DELETE("/friends/:friend_id", handlers.DeleteFriend)

		// 会话管理
		authed.GET("/conversations", handlers.GetConversations)
		authed.GET("/conversations/unread", handlers.GetUnreadCount)
		authed.POST("/conversations/pin", handlers.SetConversationPin)
		authed.POST("/conversations/mute", handlers.SetConversationMute)
		authed.POST("/conversations/archive", handlers.ArchiveConversation)
		authed.GET("/conversations/archived", handlers.GetArchivedConversations)

		// 消息功能
		authed.POST("/messages", handlers.SendMessage)
		authed.GET("/messages/private/:friend_id", handlers.GetPrivateMessages)
		authed.GET("/messages/group/:group_id", handlers.GetGroupMessages)
		authed.POST("/messages/read", handlers.MarkAsRead)
		authed.POST("/messages/:message_id/recall", handlers.RecallMessage)
		authed.GET("/messages/export", handlers.ExportChatHistory)
		authed.POST("/messages/:message_id/reaction", handlers.AddReaction)
		authed.DELETE("/messages/:message_id/reaction", handlers.RemoveReaction)
		authed.GET("/messages/:message_id/reactions", handlers.GetReactions)
		authed.POST("/messages/forward", handlers.ForwardMessage)
		authed.GET("/messages/search", handlers.SearchMessages)

		// 通话功能
		authed.POST("/calls", handlers.InitiateCall)
		authed.POST("/calls/:session_id/answer", handlers.AnswerCall)
		authed.POST("/calls/:session_id/reject", handlers.RejectCall)
		authed.POST("/calls/:session_id/end", handlers.EndCall)
		authed.GET("/calls/history", handlers.GetCallHistory)
		authed.GET("/calls/signal", handlers.WebRTCSignal)

		// 2FA双因素认证
		authed.GET("/2fa/status", handlers.Get2FAStatus)
		authed.POST("/2fa/enable", handlers.Enable2FA)
		authed.POST("/2fa/verify-enable", handlers.VerifyAndEnable2FA)
		authed.POST("/2fa/disable", handlers.Disable2FA)
		authed.POST("/2fa/verify", handlers.Verify2FA)
		authed.POST("/2fa/regenerate-codes", handlers.RegenerateBackupCodes)

		authed.POST("/groups", handlers.CreateGroup)
		authed.GET("/groups", handlers.GetMyGroups)
		authed.GET("/groups/:group_id", handlers.GetGroupInfo)
		authed.PUT("/groups/:group_id", handlers.UpdateGroup)
		authed.GET("/groups/:group_id/members", handlers.GetGroupMembers)
		authed.POST("/groups/:group_id/invite", handlers.InviteGroupMember)
		authed.DELETE("/groups/:group_id/members/:member_id", handlers.RemoveGroupMember)
		authed.POST("/groups/:group_id/members/:member_id/mute", handlers.MuteGroupMember)
		authed.POST("/groups/:group_id/leave", handlers.LeaveGroup)

		authed.POST("/redpackets", handlers.SendRedPacket)
		authed.POST("/redpackets/:id/grab", handlers.GrabRedPacket)
		authed.GET("/redpackets/:id", handlers.GetRedPacket)
		authed.GET("/redpackets/sent", handlers.GetSentRedPackets)
		authed.GET("/redpackets/received", handlers.GetReceivedRedPackets)

		authed.POST("/moments", handlers.PublishMoment)
		authed.GET("/moments", handlers.GetMoments)
		authed.POST("/moments/:moment_id/like", handlers.LikeMoment)
		authed.DELETE("/moments/:moment_id/like", handlers.UnlikeMoment)
		authed.POST("/moments/:moment_id/comments", handlers.CommentMoment)
		authed.DELETE("/moments/:moment_id", handlers.DeleteMoment)

		authed.POST("/payment/orders", handlers.CreateOrder)
		authed.POST("/payment/orders/:order_id/pay", handlers.ProcessPayment)
		authed.GET("/payment/orders", handlers.GetOrders)
		authed.GET("/payment/points/history", handlers.GetPointsHistory)

		authed.POST("/recharge", handlers.CreateRechargeRequest)
		authed.GET("/recharge", handlers.GetMyRechargeRequests)
		authed.GET("/recharge/:id", handlers.GetRechargeRequestDetail)

		authed.POST("/withdraw", handlers.CreateWithdrawRequest)
		authed.GET("/withdraw", handlers.GetMyWithdrawRequests)
		authed.GET("/withdraw/:id", handlers.GetWithdrawRequestDetail)

		authed.GET("/emoji/categories", handlers.GetEmojiCategories)
		authed.GET("/emoji/categories/:category_id", handlers.GetEmojisByCategory)
	}

	system := api.Group("/system")
	{
		system.GET("/config", handlers.GetSystemConfig)
	}

	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware())
	admin.Use(middleware.AdminMiddleware())
	{
		admin.GET("/db/stats", handlers.GetDatabaseStats)
		admin.POST("/db/clear-old-messages", handlers.ClearOldMessages)
		admin.POST("/db/archive-old", handlers.ArchiveOldMessages)
		admin.DELETE("/db/users/:user_id", handlers.DeleteUserAndData)
		admin.DELETE("/db/groups/:group_id", handlers.DeleteGroupAndData)
		admin.POST("/users/points", handlers.AdminAddPoints)
		admin.GET("/system/configs", handlers.GetSystemConfigs)
		admin.PUT("/system/configs", handlers.UpdateSystemConfig)

		admin.PUT("/system/config", handlers.UpdateSystemConfig)
		admin.POST("/system/logo", handlers.UploadLogo)
		admin.POST("/system/favicon", handlers.UploadFavicon)
		admin.POST("/system/maintenance", handlers.SetMaintenanceMode)

		admin.GET("/recharge", handlers.GetAllRechargeRequests)
		admin.PUT("/recharge/:id/approve", handlers.ApproveRecharge)
		admin.PUT("/recharge/:id/reject", handlers.RejectRecharge)

		admin.GET("/withdraw", handlers.GetAllWithdrawRequests)
		admin.PUT("/withdraw/:id/approve", handlers.ApproveWithdraw)
		admin.PUT("/withdraw/:id/reject", handlers.RejectWithdraw)

		admin.POST("/emoji/categories", handlers.AddEmojiCategory)
		admin.DELETE("/emoji/categories/:id", handlers.DeleteEmojiCategory)
		admin.POST("/emoji/items", handlers.AddEmojiItem)
		admin.DELETE("/emoji/items/:id", handlers.DeleteEmojiItem)
	}

	superAdmin := api.Group("/admin/super")
	superAdmin.Use(middleware.AuthMiddleware())
	superAdmin.Use(middleware.SuperAdminMiddleware())
	{
		superAdmin.POST("/db/clear-all", handlers.ClearAllData)
		superAdmin.POST("/db/init", handlers.InitializeDatabase)
	}

	log.Printf("Server starting on port %d", config.Cfg.Server.Port)
	if err := r.Run(fmt.Sprintf(":%d", config.Cfg.Server.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
