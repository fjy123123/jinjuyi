package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"chat-system-pro/config"
	"chat-system-pro/models"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client 代表一个 WebSocket 连接
type Client struct {
	ID   uint
	Conn *websocket.Conn
	Send chan []byte
}

// Hub 管理所有客户端和广播消息
type Hub struct {
	Clients    map[uint]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	mu         sync.RWMutex
}

var hub = &Hub{
	Clients:    make(map[uint]*Client),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Broadcast:  make(chan []byte, 1000),
}

func init() {
	go hub.Run()
}

// Run 启动 Hub 处理
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.ID] = client
			h.mu.Unlock()
			log.Printf("Client %d connected", client.ID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)
				close(client.Send)
				log.Printf("Client %d disconnected", client.ID)
			}
			h.mu.Unlock()

		case message := <-h.Broadcast:
			h.mu.RLock()
			for _, client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client.ID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastNewMessage 广播新消息
func BroadcastNewMessage(msg *models.MessageDoc) {
	wsMsg := WSMessage{
		Type: "new_message",
		Data: msg,
	}
	bytes, _ := json.Marshal(wsMsg)

	if msg.GroupID > 0 {
		broadcastToGroup(msg.GroupID, msg.SenderID, bytes)
	} else {
		broadcastToUser(msg.ReceiverID, bytes)
	}
}

// BroadcastRecallMessage 广播撤回消息
func BroadcastRecallMessage(messageID string) {
	wsMsg := WSMessage{
		Type: "recall_message",
		Data: map[string]interface{}{"message_id": messageID},
	}
	bytes, _ := json.Marshal(wsMsg)
	hub.Broadcast <- bytes
}

// WSMessage WebSocket 消息格式
type WSMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// RegisterClient 注册一个客户端
func RegisterClient(client *Client) {
	hub.Register <- client
}

// UnregisterClient 注销一个客户端
func UnregisterClient(client *Client) {
	hub.Unregister <- client
}

// SendToUser 给用户发送消息
func SendToUser(userID uint, msgType string, data interface{}) {
	wsMsg := WSMessage{
		Type: msgType,
		Data: data,
	}
	bytes, _ := json.Marshal(wsMsg)
	broadcastToUser(userID, bytes)
}

func broadcastToUser(userID uint, bytes []byte) {
	hub.mu.RLock()
	if client, ok := hub.Clients[userID]; ok {
		select {
		case client.Send <- bytes:
		default:
		}
	}
	hub.mu.RUnlock()
}

func broadcastToGroup(groupID uint, excludeUserID uint, bytes []byte) {
	var memberIDs []uint
	config.DB.Model(&models.GroupMember{}).Where("group_id = ?", groupID).Pluck("user_id", &memberIDs)

	hub.mu.RLock()
	for _, userID := range memberIDs {
		if userID == excludeUserID {
			continue
		}
		if client, ok := hub.Clients[userID]; ok {
			select {
			case client.Send <- bytes:
			default:
			}
		}
	}
	hub.mu.RUnlock()
}
