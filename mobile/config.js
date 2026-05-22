const config = {
  development: {
    API_BASE_URL: 'http://localhost:8080',
    WS_BASE_URL: 'ws://localhost:8080/ws',
    UPLOAD_URL: 'http://localhost:8080/api/v1/upload'
  },
  production: {
    API_BASE_URL: 'https://your-domain.com',
    WS_BASE_URL: 'wss://your-domain.com/ws',
    UPLOAD_URL: 'https://your-domain.com/api/v1/upload'
  }
}

const env = process.env.NODE_ENV || 'development'

export default {
  ...config[env],
  
  // 通用配置
  APP_NAME: '即时聊天',
  VERSION: '1.0.0',
  TIMEOUT: 15000,
  DEBUG: env === 'development'
}
