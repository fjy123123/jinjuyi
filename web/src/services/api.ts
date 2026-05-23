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

  async uploadFile<T>(path: string, formData: FormData): Promise<ApiResponse<T>> {
    const url = `${this.baseUrl}${path}`
    const token = this.getToken()

    const headers: Record<string, string> = {}
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }

    const response = await fetch(url, {
      method: 'POST',
      headers,
      body: formData
    })

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
    refresh: (data: { refresh_token: string }) => apiClient.post('/api/v1/auth/refresh', data),
    logout: () => apiClient.post('/api/v1/auth/logout')
  },
  user: {
    getProfile: () => apiClient.get('/api/v1/users/me'),
    updateProfile: (data: { nickname?: string; avatar?: string; gender?: number; region?: string; sign?: string }) => apiClient.put('/api/v1/users/profile', data),
    getSettings: () => apiClient.get('/api/v1/users/settings'),
    updateSettings: (data: { new_msg_notify?: boolean; sound_notify?: boolean; theme?: string; language?: string }) => apiClient.put('/api/v1/users/settings', data),
    search: (keyword: string, page = 1, pageSize = 20) => apiClient.get(`/api/v1/users/search?keyword=${keyword}&page=${page}&page_size=${pageSize}`)
  },
  friend: {
    getList: () => apiClient.get('/api/v1/friends'),
    add: (data: { friend_id: number; remark?: string; message?: string }) => apiClient.post('/api/v1/friends', data),
    delete: (friendId: number) => apiClient.delete(`/api/v1/friends/${friendId}`)
  },
  conversation: {
    getList: (page = 1, pageSize = 20) => apiClient.get(`/api/v1/conversations?page=${page}&page_size=${pageSize}`),
    getUnreadCount: () => apiClient.get('/api/v1/conversations/unread'),
    pin: (data: { target_id: number; type: number; is_pinned: boolean }) => apiClient.post('/api/v1/conversations/pin', data),
    mute: (data: { target_id: number; type: number; is_muted: boolean }) => apiClient.post('/api/v1/conversations/mute', data),
    archive: (data: { target_id: number; type: number; is_archived: boolean }) => apiClient.post('/api/v1/conversations/archive', data),
    getArchived: () => apiClient.get('/api/v1/conversations/archived')
  },
  message: {
    getPrivate: (friendId: number, page = 1, pageSize = 20) => apiClient.get(`/api/v1/messages/private/${friendId}?page=${page}&page_size=${pageSize}`),
    getGroup: (groupId: number, page = 1, pageSize = 20) => apiClient.get(`/api/v1/messages/group/${groupId}?page=${page}&page_size=${pageSize}`),
    send: (data: { receiver_id?: number; group_id?: number; content: string; message_type: number; red_packet_id?: number }) => apiClient.post('/api/v1/messages', data),
    markAsRead: (data: { target_id: number; type: number }) => apiClient.post('/api/v1/messages/read', data),
    recall: (messageId: string) => apiClient.post(`/api/v1/messages/${messageId}/recall`),
    export: (format = 'json') => apiClient.get(`/api/v1/messages/export?format=${format}`),
    addReaction: (messageId: string, emoji: string) => apiClient.post(`/api/v1/messages/${messageId}/reaction`, { emoji }),
    removeReaction: (messageId: string) => apiClient.delete(`/api/v1/messages/${messageId}/reaction`),
    getReactions: (messageId: string) => apiClient.get(`/api/v1/messages/${messageId}/reactions`),
    forward: (data: { message_id: string; target_ids: number[]; target_type: number }) => apiClient.post('/api/v1/messages/forward', data),
    search: (keyword: string, page = 1, pageSize = 20) => apiClient.get(`/api/v1/messages/search?keyword=${keyword}&page=${page}&page_size=${pageSize}`)
  },
  group: {
    getList: (page = 1, pageSize = 20) => apiClient.get(`/api/v1/groups?page=${page}&page_size=${pageSize}`),
    create: (data: { name: string; avatar?: string; description?: string; member_ids?: number[]; join_mode?: number }) => apiClient.post('/api/v1/groups', data),
    getInfo: (groupId: number) => apiClient.get(`/api/v1/groups/${groupId}`),
    update: (groupId: number, data: { name?: string; avatar?: string; description?: string; announcement?: string; join_mode?: number; is_mute_all?: boolean }) => apiClient.put(`/api/v1/groups/${groupId}`, data),
    getMembers: (groupId: number, page = 1, pageSize = 50) => apiClient.get(`/api/v1/groups/${groupId}/members?page=${page}&page_size=${pageSize}`),
    invite: (groupId: number, userIds: number[]) => apiClient.post(`/api/v1/groups/${groupId}/invite`, { user_ids: userIds }),
    removeMember: (groupId: number, memberId: number) => apiClient.delete(`/api/v1/groups/${groupId}/members/${memberId}`),
    muteMember: (groupId: number, memberId: number, data: { duration: number; reason?: string }) => apiClient.post(`/api/v1/groups/${groupId}/members/${memberId}/mute`, data),
    leave: (groupId: number) => apiClient.post(`/api/v1/groups/${groupId}/leave`)
  },
  redpacket: {
    send: (data: { type: number; pay_type: number; amount: number; total_count: number; receiver_id?: number; group_id?: number; greeting?: string }) => apiClient.post('/api/v1/redpackets', data),
    grab: (id: number) => apiClient.post(`/api/v1/redpackets/${id}/grab`),
    getDetail: (id: number) => apiClient.get(`/api/v1/redpackets/${id}`),
    getSent: (page = 1, pageSize = 20) => apiClient.get(`/api/v1/redpackets/sent?page=${page}&page_size=${pageSize}`),
    getReceived: (page = 1, pageSize = 20) => apiClient.get(`/api/v1/redpackets/received?page=${page}&page_size=${pageSize}`)
  },
  moment: {
    publish: (data: { content: string; images?: string[]; location?: string; latitude?: number; longitude?: number; view_scope?: number }) => apiClient.post('/api/v1/moments', data),
    getList: (page = 1, pageSize = 20) => apiClient.get(`/api/v1/moments?page=${page}&page_size=${pageSize}`),
    like: (momentId: number) => apiClient.post(`/api/v1/moments/${momentId}/like`),
    unlike: (momentId: number) => apiClient.delete(`/api/v1/moments/${momentId}/like`),
    comment: (momentId: number, data: { content: string; reply_to_user?: number }) => apiClient.post(`/api/v1/moments/${momentId}/comments`, data),
    delete: (momentId: number) => apiClient.delete(`/api/v1/moments/${momentId}`)
  },
  payment: {
    createOrder: (data: { amount: number; pay_type: number; order_type: number; subject?: string; description?: string }) => apiClient.post('/api/v1/payment/orders', data),
    pay: (orderId: string) => apiClient.post(`/api/v1/payment/orders/${orderId}/pay`),
    getOrders: (page = 1, pageSize = 20, status = 'all') => apiClient.get(`/api/v1/payment/orders?page=${page}&page_size=${pageSize}&status=${status}`),
    getPointsHistory: (page = 1, pageSize = 20) => apiClient.get(`/api/v1/payment/points/history?page=${page}&page_size=${pageSize}`)
  },
  recharge: {
    create: (data: { amount: number; points: number; recharge_type: string; payment_image: string; remark?: string }) => apiClient.post('/api/v1/recharge', data),
    getList: (page = 1, pageSize = 20) => apiClient.get(`/api/v1/recharge?page=${page}&page_size=${pageSize}`),
    getDetail: (id: number) => apiClient.get(`/api/v1/recharge/${id}`)
  },
  withdraw: {
    create: (data: { points: number; amount: number; withdraw_type: string; payment_code: string; real_name: string; phone: string; remark?: string }) => apiClient.post('/api/v1/withdraw', data),
    getList: (page = 1, pageSize = 20) => apiClient.get(`/api/v1/withdraw?page=${page}&page_size=${pageSize}`),
    getDetail: (id: number) => apiClient.get(`/api/v1/withdraw/${id}`)
  },
  emoji: {
    getCategories: () => apiClient.get('/api/v1/emoji/categories'),
    getByCategory: (categoryId: number) => apiClient.get(`/api/v1/emoji/categories/${categoryId}`)
  },
  call: {
    initiate: (data: { target_id: number; type: number }) => apiClient.post('/api/v1/calls', data),
    answer: (sessionId: string) => apiClient.post(`/api/v1/calls/${sessionId}/answer`),
    reject: (sessionId: string) => apiClient.post(`/api/v1/calls/${sessionId}/reject`),
    end: (sessionId: string) => apiClient.post(`/api/v1/calls/${sessionId}/end`),
    getHistory: (page = 1, pageSize = 20) => apiClient.get(`/api/v1/calls/history?page=${page}&page_size=${pageSize}`)
  },
  twoFA: {
    getStatus: () => apiClient.get('/api/v1/2fa/status'),
    enable: () => apiClient.post('/api/v1/2fa/enable'),
    verifyAndEnable: (data: { code: string }) => apiClient.post('/api/v1/2fa/verify-enable', data),
    disable: (data: { code: string }) => apiClient.post('/api/v1/2fa/disable', data),
    verify: (data: { code: string }) => apiClient.post('/api/v1/2fa/verify', data),
    regenerateCodes: () => apiClient.post('/api/v1/2fa/regenerate-codes')
  },
  admin: {
    getDatabaseStats: () => apiClient.get('/api/v1/admin/db/stats'),
    clearOldMessages: (date: string) => apiClient.post(`/api/v1/admin/db/clear-old-messages?date=${date}`),
    archiveOld: (date: string) => apiClient.post(`/api/v1/admin/db/archive-old?date=${date}`),
    deleteUser: (userId: number) => apiClient.delete(`/api/v1/admin/db/users/${userId}`),
    deleteGroup: (groupId: number) => apiClient.delete(`/api/v1/admin/db/groups/${groupId}`),
    addPoints: (data: { user_id: number; points: number; type: number; remark: string }) => apiClient.post('/api/v1/admin/users/points', data),
    getSystemConfigs: () => apiClient.get('/api/v1/admin/system/configs'),
    updateSystemConfig: (data: { section: string; key: string; value: string }) => apiClient.put('/api/v1/admin/system/configs', data),
    updateConfig: (data: { app_name?: string; app_version?: string; app_description?: string; theme_color?: string; theme_secondary?: string; ui_template?: string }) => apiClient.put('/api/v1/admin/system/config', data),
    uploadLogo: (formData: FormData) => apiClient.uploadFile('/api/v1/admin/system/logo', formData),
    uploadFavicon: (formData: FormData) => apiClient.uploadFile('/api/v1/admin/system/favicon', formData),
    setMaintenance: (data: { mode: boolean; message?: string }) => apiClient.post('/api/v1/admin/system/maintenance', data),
    getAllRecharge: (status = 'pending', page = 1, pageSize = 20) => apiClient.get(`/api/v1/admin/recharge?status=${status}&page=${page}&page_size=${pageSize}`),
    approveRecharge: (id: number, remark: string) => apiClient.put(`/api/v1/admin/recharge/${id}/approve`, { remark }),
    rejectRecharge: (id: number, remark: string) => apiClient.put(`/api/v1/admin/recharge/${id}/reject`, { remark }),
    getAllWithdraw: (status = 'pending', page = 1, pageSize = 20) => apiClient.get(`/api/v1/admin/withdraw?status=${status}&page=${page}&page_size=${pageSize}`),
    approveWithdraw: (id: number, remark: string) => apiClient.put(`/api/v1/admin/withdraw/${id}/approve`, { remark }),
    rejectWithdraw: (id: number, remark: string) => apiClient.put(`/api/v1/admin/withdraw/${id}/reject`, { remark }),
    addEmojiCategory: (data: { name: string; icon: string }) => apiClient.post('/api/v1/admin/emoji/categories', data),
    deleteEmojiCategory: (id: number) => apiClient.delete(`/api/v1/admin/emoji/categories/${id}`),
    addEmojiItem: (data: { category_id: number; emoji: string; name: string }) => apiClient.post('/api/v1/admin/emoji/items', data),
    deleteEmojiItem: (id: number) => apiClient.delete(`/api/v1/admin/emoji/items/${id}`),
    clearAll: (data: { confirm: string; backup?: boolean }) => apiClient.post('/api/v1/admin/super/db/clear-all', data),
    initDb: (data: { confirm: string }) => apiClient.post('/api/v1/admin/super/db/init', data)
  },
  system: {
    getConfig: () => apiClient.get('/api/v1/system/config')
  }
}

export default api
