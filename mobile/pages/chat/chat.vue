<template>
  <view class="container">
    <view class="header">
      <text class="title">{{ title }}</text>
    </view>
    
    <view class="message-list" v-if="messages.length > 0">
      <view 
        v-for="(msg, index) in messages" 
        :key="index" 
        :class="['message-item', msg.self ? 'self' : '']"
      >
        <image v-if="!msg.self" :src="msg.avatar" class="avatar"></image>
        <view :class="['bubble', msg.type === 'text' ? 'text' : msg.type]">
          <text v-if="msg.type === 'text'" class="text">{{ msg.content }}</text>
          <image v-if="msg.type === 'image'" :src="msg.content" class="image"></image>
        </view>
        <image v-if="msg.self" :src="msg.avatar" class="avatar"></image>
      </view>
    </view>
    
    <view v-else class="empty">
      <text>开始聊天吧</text>
    </view>
    
    <view class="input-area">
      <input 
        v-model="inputText" 
        class="input" 
        placeholder="请输入消息"
        @confirm="sendMessage"
      />
      <button class="send-btn" @click="sendMessage">发送</button>
    </view>
  </view>
</template>

<script>
import request from '../../utils/request.js'
import ws from '../../utils/websocket.js'

export default {
  data() {
    return {
      title: '聊天',
      friendId: '',
      inputText: '',
      messages: []
    }
  },
  
  onLoad(options) {
    this.title = options.name || '聊天'
    this.friendId = options.id
    this.loadMessages()
    this.initWebSocket()
  },
  
  onUnload() {
    ws.disconnect()
  },
  
  methods: {
    async loadMessages() {
      try {
        const res = await request.get(`/api/v1/messages/private/${this.friendId}`)
        if (res.code === 200 || res.code === 0) {
          this.messages = res.data || []
        }
      } catch (e) {
        console.error('加载消息失败:', e)
      }
    },
    
    initWebSocket() {
      const userId = uni.getStorageSync('userId')
      if (userId) {
        ws.connect(userId)
        
        ws.on('new_message', (data) => {
          if (data.friendId == this.friendId || data.senderId == this.friendId) {
            this.messages.push(data)
          }
        })
      }
    },
    
    async sendMessage() {
      if (!this.inputText.trim()) return
      
      try {
        const res = await request.post('/api/v1/messages', {
          receiverId: this.friendId,
          type: 'text',
          content: this.inputText
        })
        
        if (res.code === 200 || res.code === 0) {
          this.messages.push({
            self: true,
            content: this.inputText,
            type: 'text',
            avatar: 'https://via.placeholder.com/40',
            time: new Date().toLocaleTimeString()
          })
          
          ws.send({
            type: 'message',
            data: {
              receiverId: this.friendId,
              content: this.inputText
            }
          })
          
          this.inputText = ''
        }
      } catch (e) {
        console.error('发送失败:', e)
      }
    }
  }
}
</script>

<style scoped>
.container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: #f5f5f5;
}

.header {
  background-color: #fff;
  padding: 20px;
  text-align: center;
  border-bottom: 1px solid #eee;
}

.title {
  font-size: 18px;
  font-weight: bold;
}

.message-list {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.message-item {
  display: flex;
  margin-bottom: 20px;
}

.message-item.self {
  flex-direction: row-reverse;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 5px;
}

.bubble {
  max-width: 70%;
  padding: 10px 15px;
  margin: 0 10px;
  border-radius: 5px;
}

.message-item:not(.self) .bubble {
  background-color: #fff;
}

.message-item.self .bubble {
  background-color: #95ec69;
}

.text {
  font-size: 16px;
  line-height: 1.6;
  word-wrap: break-word;
}

.image {
  max-width: 200px;
  border-radius: 5px;
}

.input-area {
  display: flex;
  padding: 10px;
  background-color: #fff;
  border-top: 1px solid #eee;
}

.input {
  flex: 1;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 5px;
  margin-right: 10px;
}

.send-btn {
  background-color: #07c160;
  color: #fff;
  border: none;
  border-radius: 5px;
  padding: 10px 20px;
}

.empty {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
}
</style>
