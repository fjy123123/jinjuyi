package handlers

import (
	"chat-system-pro/config"
	"chat-system-pro/models"
	"chat-system-pro/utils"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebRTC 信令服务器
var webrtcUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// CallSession 通话会话
type CallSession struct {
	ID          string    `json:"id"`
	CallerID    uint      `json:"caller_id"`
	CalleeID    uint      `json:"callee_id"`
	GroupID     uint      `json:"group_id,omitempty"`
	Type        int       `json:"type"` // 1:视频 2:语音
	Status      int       `json:"status"` // 0:呼叫中 1:通话中 2:已结束 3:已拒绝
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time,omitempty"`
	Connections map[uint]*websocket.Conn
}

var callSessions = make(map[string]*CallSession)

// InitiateCall 发起通话
func InitiateCall(c *gin.Context) {
	callerID := c.GetUint("user_id")
	
	var req struct {
		TargetID uint `json:"target_id" binding:"required"`
		Type     int  `json:"type" binding:"required"` // 1:视频 2:语音
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "参数错误")
		return
	}

	// 生成会话ID
	sessionID := utils.GenerateUUID()
	
	session := &CallSession{
		ID:          sessionID,
		CallerID:    callerID,
		CalleeID:    req.TargetID,
		Type:        req.Type,
		Status:      0,
		StartTime:   time.Now(),
		Connections: make(map[uint]*websocket.Conn),
	}
	
	callSessions[sessionID] = session

	// 发送呼叫通知
	utils.SendToUser(req.TargetID, "incoming_call", map[string]interface{}{
		"session_id": sessionID,
		"caller_id":  callerID,
		"type":       req.Type,
	})

	utils.SuccessResponse(c, gin.H{"session_id": sessionID})
}

// AnswerCall 接听通话
func AnswerCall(c *gin.Context) {
	userID := c.GetUint("user_id")
	sessionID := c.Param("session_id")

	session, exists := callSessions[sessionID]
	if !exists {
		utils.ErrorResponse(c, 404, "通话不存在")
		return
	}

	if session.CalleeID != userID {
		utils.ErrorResponse(c, 403, "无权接听此通话")
		return
	}

	session.Status = 1
	
	// 通知呼叫者
	utils.SendToUser(session.CallerID, "call_answered", map[string]interface{}{
		"session_id": sessionID,
	})

	utils.SuccessResponse(c, session)
}

// RejectCall 拒绝通话
func RejectCall(c *gin.Context) {
	userID := c.GetUint("user_id")
	sessionID := c.Param("session_id")

	session, exists := callSessions[sessionID]
	if !exists {
		utils.ErrorResponse(c, 404, "通话不存在")
		return
	}

	session.Status = 3
	session.EndTime = time.Now()

	// 通知呼叫者
	utils.SendToUser(session.CallerID, "call_rejected", map[string]interface{}{
		"session_id": sessionID,
		"by":         userID,
	})

	delete(callSessions, sessionID)
	utils.SuccessResponse(c, nil)
}

// EndCall 结束通话
func EndCall(c *gin.Context) {
	userID := c.GetUint("user_id")
	sessionID := c.Param("session_id")

	session, exists := callSessions[sessionID]
	if !exists {
		utils.ErrorResponse(c, 404, "通话不存在")
		return
	}

	if session.CallerID != userID && session.CalleeID != userID {
		utils.ErrorResponse(c, 403, "无权结束此通话")
		return
	}

	session.Status = 2
	session.EndTime = time.Now()

	// 通知对方
	otherID := session.CallerID
	if userID == session.CallerID {
		otherID = session.CalleeID
	}
	
	utils.SendToUser(otherID, "call_ended", map[string]interface{}{
		"session_id": sessionID,
		"by":         userID,
		"duration":   session.EndTime.Sub(session.StartTime).Seconds(),
	})

	// 保存通话记录
	config.DB.Create(&models.CallRecord{
		SessionID:  sessionID,
		CallerID:   session.CallerID,
		CalleeID:   session.CalleeID,
		Type:       session.Type,
		Status:     session.Status,
		StartTime:  session.StartTime,
		EndTime:    session.EndTime,
		Duration:   int(session.EndTime.Sub(session.StartTime).Seconds()),
	})

	delete(callSessions, sessionID)
	utils.SuccessResponse(c, nil)
}

// WebRTCSignal WebRTC信令WebSocket
func WebRTCSignal(c *gin.Context) {
	userID := c.GetUint("user_id")
	sessionID := c.Query("session_id")

	session, exists := callSessions[sessionID]
	if !exists {
		c.JSON(404, gin.H{"error": "通话不存在"})
		return
	}

	conn, err := webrtcUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	session.Connections[userID] = conn

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var signal map[string]interface{}
		json.Unmarshal(msg, &signal)

		// 转发给对方
		otherID := session.CallerID
		if userID == session.CallerID {
			otherID = session.CalleeID
		}

		if otherConn, ok := session.Connections[otherID]; ok {
			signal["from"] = userID
			otherConn.WriteJSON(signal)
		}
	}
}

// GetCallHistory 获取通话记录
func GetCallHistory(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var records []models.CallRecord
	var total int64

	config.DB.Model(&models.CallRecord{}).
		Where("caller_id = ? OR callee_id = ?", userID, userID).
		Count(&total)

	config.DB.Where("caller_id = ? OR callee_id = ?", userID, userID).
		Order("start_time DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&records)

	utils.SuccessResponse(c, gin.H{
		"total":   total,
		"records": records,
	})
}
