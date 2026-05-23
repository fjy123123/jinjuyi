package utils

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"chat-system-pro/config"
)

// ClusterHub 集群Hub - 支持多节点消息同步
type ClusterHub struct {
	NodeID     string
	Clients    map[uint]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	mu         sync.RWMutex // 添加互斥锁保证并发安全
}

// Redis Pub/Sub channels
const (
	ChannelNewMessage    = "chat:new_message"
	ChannelRecall       = "chat:recall"
	ChannelReadReceipt  = "chat:read_receipt"
	ChannelSystemNotify = "chat:system_notify"
	ChannelCall         = "chat:call"
)

var clusterHub *ClusterHub

// InitClusterHub 初始化集群Hub
func InitClusterHub(nodeID string) {
	clusterHub = &ClusterHub{
		NodeID:     nodeID,
		Clients:    make(map[uint]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte, 10000),
	}

	// 启动本地Hub处理
	go clusterHub.Run()

	// 启动Redis订阅监听
	go clusterHub.SubscribeFromRedis()

	log.Printf("Cluster hub initialized for node: %s", nodeID)
}

// Run 启动Hub
func (h *ClusterHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.ID] = client
			h.mu.Unlock()
			log.Printf("Client %d registered on node %s", client.ID, h.NodeID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)
				close(client.Send)
				log.Printf("Client %d unregistered from node %s", client.ID, h.NodeID)
			}
			h.mu.Unlock()

		case message := <-h.Broadcast:
			h.muBroadcast(message)
		}
	}
}

// muBroadcast 发送消息给所有本地客户端
func (h *ClusterHub) muBroadcast(message []byte) {
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

// BroadcastToCluster 广播到集群所有节点
func BroadcastToCluster(channel string, msgType string, data interface{}) {
	if clusterHub == nil {
		return
	}
	
	msg := ClusterMessage{
		FromNode: clusterHub.NodeID,
		Type:     msgType,
		Data:     data,
	}

	bytes, _ := json.Marshal(msg)
	
	ctx := context.Background()
	config.RedisClient.Publish(ctx, channel, string(bytes))
}

// SubscribeFromRedis 从Redis订阅消息
func (h *ClusterHub) SubscribeFromRedis() {
	ctx := context.Background()
	
	pubsub := config.RedisClient.Subscribe(ctx,
		ChannelNewMessage,
		ChannelRecall,
		ChannelReadReceipt,
		ChannelSystemNotify,
		ChannelCall,
	)

	defer pubsub.Close()

	ch := pubsub.Channel()
	
	for msg := range ch {
		var clusterMsg ClusterMessage
		if err := json.Unmarshal([]byte(msg.Payload), &clusterMsg); err != nil {
			continue
		}

		// 忽略自己发送的消息
		if clusterMsg.FromNode == h.NodeID {
			continue
		}

		// 广播到本地客户端
		bytes, _ := json.Marshal(clusterMsg)
		h.muBroadcast(bytes)
	}
}

// ClusterMessage 集群消息
type ClusterMessage struct {
	FromNode string      `json:"from_node"`
	Type     string      `json:"type"`
	Data     interface{} `json:"data"`
}

// RegisterClientToCluster 注册客户端到集群Hub
func RegisterClientToCluster(client *Client) {
	if clusterHub != nil {
		clusterHub.Register <- client
	}
}

// UnregisterClientFromCluster 从集群Hub注销客户端
func UnregisterClientFromCluster(client *Client) {
	if clusterHub != nil {
		clusterHub.Unregister <- client
	}
}

// BroadcastNewMessageCluster 集群广播新消息
func BroadcastNewMessageCluster(msg interface{}) {
	BroadcastToCluster(ChannelNewMessage, "new_message", msg)
}

// BroadcastRecallCluster 集群广播撤回消息
func BroadcastRecallCluster(messageID string) {
	BroadcastToCluster(ChannelRecall, "recall_message", map[string]interface{}{"message_id": messageID})
}

// BroadcastReadReceiptCluster 集群广播已读回执
func BroadcastReadReceiptCluster(userID uint, targetID uint, convType int, messageIDs []string) {
	BroadcastToCluster(ChannelReadReceipt, "read_receipt", map[string]interface{}{
		"user_id":     userID,
		"target_id":   targetID,
		"type":        convType,
		"message_ids": messageIDs,
	})
}

// BroadcastSystemNotificationCluster 集群广播系统通知
func BroadcastSystemNotificationCluster(data interface{}) {
	BroadcastToCluster(ChannelSystemNotify, "system_notification", data)
}
