<template>
  <view class="container">
    <view class="header">
      <text class="title">{{ title }}</text>
    </view>
    
    <view class="form">
      <view class="form-item">
        <input 
          v-model="form.username" 
          class="input" 
          placeholder="请输入用户名"
        />
      </view>
      
      <view class="form-item">
        <input 
          v-model="form.password" 
          class="input" 
          type="password"
          placeholder="请输入密码"
        />
      </view>
      
      <view v-if="isRegister" class="form-item">
        <input 
          v-model="form.email" 
          class="input" 
          type="email"
          placeholder="请输入邮箱"
        />
      </view>
      
      <button class="btn" @click="handleSubmit">
        {{ isRegister ? '注册' : '登录' }}
      </button>
      
      <view class="switch" @click="switchMode">
        <text>{{ isRegister ? '已有账号？去登录' : '没有账号？去注册' }}</text>
      </view>
    </view>
  </view>
</template>

<script>
import request from '../../utils/request.js'

export default {
  data() {
    return {
      isRegister: false,
      form: {
        username: '',
        password: '',
        email: ''
      }
    }
  },
  
  computed: {
    title() {
      return this.isRegister ? '注册' : '登录'
    }
  },
  
  methods: {
    async handleSubmit() {
      if (!this.form.username || !this.form.password) {
        uni.showToast({
          title: '请填写完整信息',
          icon: 'none'
        })
        return
      }
      
      try {
        const url = this.isRegister ? '/api/v1/auth/register' : '/api/v1/auth/login'
        const res = await request.post(url, this.form)
        
        if (res.code === 200 || res.code === 0) {
          uni.setStorageSync('token', res.data.token)
          uni.setStorageSync('userId', res.data.userId)
          
          uni.showToast({
            title: this.isRegister ? '注册成功' : '登录成功',
            icon: 'success'
          })
          
          setTimeout(() => {
            uni.switchTab({
              url: '/pages/index/index'
            })
          }, 1500)
        }
      } catch (e) {
        console.error('操作失败:', e)
      }
    },
    
    switchMode() {
      this.isRegister = !this.isRegister
      this.form = {
        username: '',
        password: '',
        email: ''
      }
    }
  }
}
</script>

<style scoped>
.container {
  padding: 40px 20px;
  height: 100vh;
  background-color: #f5f5f5;
}

.header {
  text-align: center;
  margin-bottom: 40px;
}

.title {
  font-size: 28px;
  font-weight: bold;
  color: #333;
}

.form {
  background-color: #fff;
  padding: 30px;
  border-radius: 10px;
}

.form-item {
  margin-bottom: 20px;
}

.input {
  width: 100%;
  padding: 15px;
  border: 1px solid #ddd;
  border-radius: 5px;
  font-size: 16px;
}

.btn {
  width: 100%;
  padding: 15px;
  background-color: #07c160;
  color: #fff;
  border: none;
  border-radius: 5px;
  font-size: 16px;
  margin-top: 10px;
}

.btn:active {
  background-color: #06ad56;
}

.switch {
  text-align: center;
  margin-top: 20px;
  color: #576b95;
  font-size: 14px;
}
</style>
