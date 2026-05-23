const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

class ApiClient {
  private baseUrl: string

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl
  }

  private getToken(): string | null {
    return localStorage.getItem('token')
  }

  private async request<T>(
    method: 'GET' | 'POST' | 'PUT' | 'DELETE',
    path: string,
    data?: any
  ): Promise<ApiResponse<T>> {
    const url = `${this.baseUrl}${path}`
    const token = this.getToken()

    const headers: Record<string, string> = {
      'Content-Type': 'application/json'
    }

    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }

    const options: RequestInit = {
      method,
      headers,
      credentials: 'include'
    }

    if (data && method !== 'GET') {
      options.body = JSON.stringify(data)
    }

    const response = await fetch(url, options)

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`)
    }

    return response.json()
  }

  get<T>(path: string): Promise<ApiResponse<T>> {
    return this.request<T>('GET', path)
  }

  post<T>(path: string, data?: any): Promise<ApiResponse<T>> {
    return this.request<T>('POST', path, data)
  }

  put<T>(path: string, data?: any): Promise<ApiResponse<T>> {
    return this.request<T>('PUT', path, data)
  }

  delete<T>(path: string): Promise<ApiResponse<T>> {
    return this.request<T>('DELETE', path)
  }
}

const apiClient = new ApiClient(API_BASE_URL)

export const api = {
  auth: {
    login: (data: { username: string; password: string }) => apiClient.post('/api/v1/auth/login', data),
    register: (data: { username: string; password: string; nickname: string }) => apiClient.post('/api/v1/auth/register', data),
    getInfo: () => apiClient.get('/api/v1/auth/me')
  },
  conversation: {
    getList: () => apiClient.get('/api/v1/conversation/list')
  },
  message: {
    getPrivate: (userId: number) => apiClient.get(`/api/v1/message/private/${userId}`),
    getGroup: (groupId: number) => apiClient.get(`/api/v1/message/group/${groupId}`),
    send: (data: { receiver_id?: number; group_id?: number; content: string; message_type: number }) => apiClient.post('/api/v1/message/send', data),
    markAsRead: (data: { target_id: number; type: number }) => apiClient.post('/api/v1/message/read', data)
  },
  friend: {
    getList: () => apiClient.get('/api/v1/friend/list'),
    add: (username: string) => apiClient.post('/api/v1/friend/add', { username })
  },
  group: {
    getList: () => apiClient.get('/api/v1/group/list'),
    create: (data: { name: string; member_ids: number[] }) => apiClient.post('/api/v1/group/create', data)
  },
  user: {
    search: (keyword: string) => apiClient.get(`/api/v1/user/search?keyword=${keyword}`)
  },
  redpacket: {
    send: (data: { receiver_id?: number; group_id?: number; amount: number; count: number; type: 1 | 2; greeting?: string }) => apiClient.post('/api/v1/redpacket/send', data),
    open: (redPacketId: number) => apiClient.post(`/api/v1/redpacket/open/${redPacketId}`)
  },
  finance: {
    recharge: (data: { amount: number; recharge_type: 1 | 2; payment_image: string }) => apiClient.post('/api/v1/finance/recharge', data),
    withdraw: (data: { points: number; withdraw_type: 1 | 2; payment_code: string; real_name: string; phone?: string }) => apiClient.post('/api/v1/finance/withdraw', data)
  }
}