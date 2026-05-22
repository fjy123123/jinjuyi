import config from '../config.js'

class HttpRequest {
  constructor() {
    this.baseURL = config.API_BASE_URL
    this.timeout = config.TIMEOUT
  }

  request(options) {
    const { url, method = 'GET', data = {}, header = {} } = options
    
    return new Promise((resolve, reject) => {
      uni.request({
        url: this.baseURL + url,
        method,
        data,
        header: {
          'Content-Type': 'application/json',
          'Authorization': uni.getStorageSync('token') ? `Bearer ${uni.getStorageSync('token')}` : '',
          ...header
        },
        timeout: this.timeout,
        success: (res) => {
          if (res.statusCode === 200) {
            if (res.data.code === 0 || res.data.code === 200) {
              resolve(res.data)
            } else {
              uni.showToast({
                title: res.data.message || '请求失败',
                icon: 'none'
              })
              reject(res.data)
            }
          } else {
            uni.showToast({
              title: '网络错误',
              icon: 'none'
            })
            reject(res)
          }
        },
        fail: (err) => {
          uni.showToast({
            title: '网络连接失败',
            icon: 'none'
          })
          reject(err)
        }
      })
    })
  }

  get(url, data, header) {
    return this.request({ url, method: 'GET', data, header })
  }

  post(url, data, header) {
    return this.request({ url, method: 'POST', data, header })
  }

  put(url, data, header) {
    return this.request({ url, method: 'PUT', data, header })
  }

  delete(url, data, header) {
    return this.request({ url, method: 'DELETE', data, header })
  }

  upload(url, filePath, formData = {}) {
    return new Promise((resolve, reject) => {
      uni.uploadFile({
        url: this.baseURL + url,
        filePath,
        name: 'file',
        formData,
        header: {
          'Authorization': uni.getStorageSync('token') ? `Bearer ${uni.getStorageSync('token')}` : ''
        },
        success: (res) => {
          const data = JSON.parse(res.data)
          if (data.code === 0 || data.code === 200) {
            resolve(data)
          } else {
            reject(data)
          }
        },
        fail: reject
      })
    })
  }
}

export default new HttpRequest()
