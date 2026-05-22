package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"chat-system-pro/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocketHandler WebSocket 连接处理
func WebSocketHandler(c *gin.Context) {
	userID := c.GetUint("user_id")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &utils.Client{
		ID:   userID,
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	utils.RegisterClient(client)

	wsMsg := utils.WSMessage{
		Type: "connected",
		Data: map[string]interface{}{"user_id": userID},
	}
	bytes, _ := json.Marshal(wsMsg)
	client.Send <- bytes

	go readPump(client)
	go writePump(client)
}

func readPump(c *utils.Client) {
	defer func() {
		utils.UnregisterClient(c)
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		var wsMsg utils.WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			log.Printf("Parse message error: %v", err)
			continue
		}

		switch wsMsg.Type {
		case "ping":
			resp := utils.WSMessage{Type: "pong"}
			bytes, _ := json.Marshal(resp)
			c.Send <- bytes
		}
	}
}

func writePump(c *utils.Client) {
	defer c.Conn.Close()

	for message := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("WebSocket write error: %v", err)
			break
		}
	}
}
