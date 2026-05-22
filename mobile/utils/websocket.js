import config from '../config.js'

class WebSocketClient {
  constructor() {
    this.ws = null
    this.reconnectAttempts = 0
    this.maxReconnectAttempts = 5
    this.heartbeatInterval = null
    this.messageHandlers = new Map()
    this.isConnected = false
  }

  connect(userId) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      return
    }

    const wsUrl = `${config.WS_BASE_URL}?user_id=${userId}&token=${uni.getStorageSync('token')}`
    
    this.ws = new WebSocket(wsUrl)
    
    this.ws.onopen = () => {
      console.log('WebSocket连接成功')
      this.isConnected = true
      this.reconnectAttempts = 0
      this.startHeartbeat()
    }

    this.ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data)
        this.handleMessage(message)
      } catch (e) {
        console.error('消息解析失败:', e)
      }
    }

    this.ws.onclose = () => {
      console.log('WebSocket连接关闭')
      this.isConnected = false
      this.stopHeartbeat()
      this.reconnect(userId)
    }

    this.ws.onerror = (error) => {
      console.error('WebSocket错误:', error)
      this.isConnected = false
    }
  }

  handleMessage(message) {
    const handler = this.messageHandlers.get(message.type)
    if (handler) {
      handler(message.data)
    }
  }

  on(eventType, handler) {
    this.messageHandlers.set(eventType, handler)
  }

  send(message) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message))
    } else {
      console.error('WebSocket未连接')
    }
  }

  startHeartbeat() {
    this.heartbeatInterval = setInterval(() => {
      this.send({ type: 'ping' })
    }, 30000)
  }

  stopHeartbeat() {
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval)
      this.heartbeatInterval = null
    }
  }

  reconnect(userId) {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++
      console.log(`尝试重连 (${this.reconnectAttempts}/${this.maxReconnectAttempts})...`)
      setTimeout(() => {
        this.connect(userId)
      }, 3000 * this.reconnectAttempts)
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
    this.stopHeartbeat()
    this.isConnected = false
  }
}

export default new WebSocketClient()
