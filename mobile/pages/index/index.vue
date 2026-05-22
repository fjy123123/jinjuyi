<template>
  <view class="container">
    <view class="header">
      <text class="title">即时聊天</text>
    </view>
    
    <view class="chat-list">
      <view 
        v-for="(item, index) in chatList" 
        :key="index" 
        class="chat-item"
        @click="goToChat(item)"
      >
        <image :src="item.avatar" class="avatar"></image>
        <view class="content">
          <view class="name">{{ item.name }}</view>
          <view class="last-message">{{ item.lastMessage }}</view>
        </view>
        <view class="meta">
          <text class="time">{{ item.time }}</text>
          <text v-if="item.unread > 0" class="unread">{{ item.unread }}</text>
        </view>
      </view>
    </view>
    
    <view v-if="chatList.length === 0" class="empty">
      <text>暂无消息</text>
    </view>
  </view>
</template>

<script>
import request from '../../utils/request.js'
import ws from '../../utils/websocket.js'

export default {
  data() {
    return {
      chatList: []
    }
  },
  
  onLoad() {
    this.loadChatList()
    this.initWebSocket()
  },
  
  onShow() {
    this.loadChatList()
  },
  
  onUnload() {
    ws.disconnect()
  },
  
  methods: {
    async loadChatList() {
      try {
        const res = await request.get('/api/v1/conversations')
        if (res.code === 200 || res.code === 0) {
          this.chatList = res.data || []
        }
      } catch (e) {
        console.error('加载聊天列表失败:', e)
      }
    },
    
    initWebSocket() {
      const userId = uni.getStorageSync('userId')
      if (userId) {
        ws.connect(userId)
        
        ws.on('new_message', (data) => {
          this.loadChatList()
        })
      }
    },
    
    goToChat(item) {
      uni.navigateTo({
        url: `/pages/chat/chat?id=${item.id}&name=${item.name}`
      })
    }
  }
}
</script>

<style scoped>
.container {
  padding: 0;
  height: 100vh;
  background-color: #f5f5f5;
}

.header {
  background-color: #fff;
  padding: 20px;
  text-align: center;
  border-bottom: 1px solid #eee;
}

.title {
  font-size: 20px;
  font-weight: bold;
  color: #333;
}

.chat-list {
  background-color: #fff;
}

.chat-item {
  display: flex;
  padding: 15px;
  border-bottom: 1px solid #f0f0f0;
}

.avatar {
  width: 50px;
  height: 50px;
  border-radius: 5px;
  margin-right: 12px;
}

.content {
  flex: 1;
  overflow: hidden;
}

.name {
  font-size: 16px;
  color: #333;
  margin-bottom: 5px;
}

.last-message {
  font-size: 14px;
  color: #999;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  justify-content: space-between;
}

.time {
  font-size: 12px;
  color: #999;
}

.unread {
  background-color: #fa5151;
  color: #fff;
  font-size: 12px;
  padding: 2px 6px;
  border-radius: 10px;
  min-width: 16px;
  text-align: center;
}

.empty {
  text-align: center;
  padding: 50px 20px;
  color: #999;
}
</style>
