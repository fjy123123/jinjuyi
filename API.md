# Chat System Pro - API接口文档

## 目录

1. [基础信息](#1-基础信息)
2. [认证接口](#2-认证接口)
3. [用户接口](#3-用户接口)
4. [好友接口](#4-好友接口)
5. [会话接口](#5-会话接口)
6. [消息接口](#6-消息接口)
7. [群组接口](#7-群组接口)
8. [红包接口](#8-红包接口)
9. [朋友圈接口](#9-朋友圈接口)
10. [支付接口](#10-支付接口)
11. [管理员接口](#11-管理员接口)
12. [系统配置接口](#12-系统配置接口)
13. [WebSocket接口](#13-websocket接口)

---

## 1. 基础信息

### 基础URL

```
生产环境: https://your-domain.com/api/v1
开发环境: http://localhost:8080/api/v1
```

### 认证方式

所有需要认证的接口都需要在请求头中携带JWT Token：

```
Authorization: Bearer <token>
```

### 响应格式

所有接口统一使用JSON格式响应：

**成功响应：**
```json
{
  "code": 0,
  "msg": "success",
  "data": { ... }
}
```

**错误响应：**
```json
{
  "code": 40001,
  "msg": "错误信息",
  "data": null
}
```

### 错误码说明

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 40001 | 参数错误 |
| 40002 | 认证失败 |
| 40003 | 权限不足 |
| 40004 | 资源不存在 |
| 40005 | 操作失败 |
| 50001 | 服务器错误 |

### HTTP状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未认证或Token失效 |
| 403 | 无权限访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 2. 认证接口

### 2.1 用户注册

**请求**
```
POST /auth/register
```

**请求体：**
```json
{
  "username": "testuser",
  "password": "123456",
  "nickname": "测试用户",
  "phone": "13800138000",
  "email": "test@example.com",
  "invite_code": "ABC123"  // 可选，如果系统开启了邀请码功能
}
```

**响应示例：**
```json
{
  "code": 0,
  "msg": "注册成功",
  "data": {
    "user_id": 1,
    "username": "testuser",
    "nickname": "测试用户",
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

**验证规则：**
- username: 3-20个字符，字母数字下划线
- password: 最少6个字符
- nickname: 2-30个字符
- phone: 手机号格式（可选）
- email: 邮箱格式（可选）

---

### 2.2 用户登录

**请求**
```
POST /auth/login
```

**请求体：**
```json
{
  "username": "testuser",
  "password": "123456"
}
```

**响应示例：**
```json
{
  "code": 0,
  "msg": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_at": "2024-12-31T23:59:59Z",
    "user": {
      "id": 1,
      "username": "testuser",
      "nickname": "测试用户",
      "avatar": "https://example.com/avatar.jpg",
      "points": 1000,
      "level": 1
    }
  }
}
```

---

## 3. 用户接口

### 3.1 获取当前用户信息

**请求**
```
GET /users/me
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "username": "testuser",
    "nickname": "测试用户",
    "avatar": "https://example.com/avatar.jpg",
    "gender": 1,
    "region": "北京市",
    "sign": "这是我的个性签名",
    "phone": "13800138000",
    "email": "test@example.com",
    "points": 1000,
    "level": 1,
    "is_vip": true,
    "vip_expire_at": "2025-12-31T23:59:59Z",
    "created_at": "2024-01-01T00:00:00Z",
    "last_login_at": "2024-01-15T10:30:00Z",
    "settings": {
      "new_msg_notify": true,
      "sound_notify": true,
      "add_friend_confirm": true,
      "show_online": true,
      "show_read_receipt": true,
      "theme": "dark",
      "language": "zh-CN"
    }
  }
}
```

---

### 3.2 更新个人资料

**请求**
```
PUT /users/profile
```

**请求体：**
```json
{
  "nickname": "新昵称",
  "avatar": "https://example.com/new-avatar.jpg",
  "gender": 2,
  "region": "上海市",
  "sign": "更新后的个性签名",
  "phone": "13900139000",
  "email": "new@example.com"
}
```

**响应示例：**
```json
{
  "code": 0,
  "msg": "更新成功",
  "data": {
    "id": 1,
    "nickname": "新昵称",
    "avatar": "https://example.com/new-avatar.jpg",
    "gender": 2,
    "region": "上海市",
    "sign": "更新后的个性签名"
  }
}
```

---

### 3.3 获取用户设置

**请求**
```
GET /users/settings
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "new_msg_notify": true,
    "sound_notify": true,
    "add_friend_confirm": true,
    "show_online": true,
    "show_read_receipt": true,
    "theme": "dark",
    "language": "zh-CN",
    "font_size": "medium",
    "auto_download": true,
    "save_to_album": true
  }
}
```

---

### 3.4 更新用户设置

**请求**
```
PUT /users/settings
```

**请求体：**
```json
{
  "new_msg_notify": false,
  "sound_notify": false,
  "theme": "light",
  "language": "en-US"
}
```

**响应示例：**
```json
{
  "code": 0,
  "msg": "设置更新成功"
}
```

---

### 3.5 搜索用户

**请求**
```
GET /users/search?keyword=test&page=1&page_size=20
```

**查询参数：**
- keyword: 搜索关键词（用户名、昵称、手机号）
- page: 页码（默认1）
- page_size: 每页数量（默认20，最大100）

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 2,
        "username": "testuser2",
        "nickname": "测试用户2",
        "avatar": "https://example.com/avatar2.jpg",
        "gender": 1,
        "region": "广州市",
        "sign": "个性签名",
        "is_friend": false,
        "last_login_at": "2024-01-15T09:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 20
  }
}
```

---

## 4. 好友接口

### 4.1 获取好友列表

**请求**
```
GET /friends?page=1&page_size=20
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "friend_id": 2,
        "remark": "张三",
        "tags": ["同事", "朋友"],
        "is_star": true,
        "is_mute": false,
        "friend": {
          "id": 2,
          "username": "zhangsan",
          "nickname": "张三",
          "avatar": "https://example.com/zhangsan.jpg",
          "gender": 1,
          "region": "北京市",
          "sign": "今天很开心",
          "is_online": true,
          "last_login_at": "2024-01-15T10:30:00Z"
        },
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 20
  }
}
```

---

### 4.2 添加好友

**请求**
```
POST /friends
```

**请求体：**
```json
{
  "friend_id": 2,
  "remark": "张三",
  "tags": ["同事"],
  "message": "你好，我是XXX"
}
```

**响应示例：**
```json
{
  "code": 0,
  "msg": "好友申请已发送",
  "data": {
    "friendship_id": 1,
    "status": "pending"
  }
}
```

---

### 4.3 删除好友

**请求**
```
DELETE /friends/:friend_id
```

**响应示例：**
```json
{
  "code": 0,
  "msg": "删除成功"
}
```

---

## 5. 会话接口

### 5.1 获取会话列表

**请求**
```
GET /conversations?page=1&page_size=20
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "conversation_type": 1,
        "target_id": 2,
        "unread_count": 5,
        "last_message": {
          "id": 100,
          "sender_id": 2,
          "content": "你好",
          "message_type": 1,
          "created_at": "2024-01-15T10:30:00Z"
        },
        "target": {
          "id": 2,
          "nickname": "张三",
          "avatar": "https://example.com/zhangsan.jpg",
          "is_online": true
        },
        "updated_at": "2024-01-15T10:30:00Z"
      },
      {
        "id": 2,
        "conversation_type": 2,
        "target_id": 10,
        "unread_count": 0,
        "last_message": {
          "id": 99,
          "sender_id": 5,
          "content": "群消息测试",
          "message_type": 1,
          "created_at": "2024-01-15T10:25:00Z"
        },
        "target": {
          "id": 10,
          "nickname": "技术交流群",
          "avatar": "https://example.com/group.jpg",
          "member_count": 50
        },
        "updated_at": "2024-01-15T10:25:00Z"
      }
    ],
    "total": 2
  }
}
```

---

### 5.2 获取未读消息总数

**请求**
```
GET /conversations/unread
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "total_unread": 15,
    "private_unread": 5,
    "group_unread": 10
  }
}
```

---

## 6. 消息接口

### 6.1 发送消息

**请求**
```
POST /messages
```

**请求体：**
```json
{
  "receiver_id": 2,
  "group_id": 0,
  "content": "你好，这是一条测试消息",
  "message_type": 1,
  "red_packet_id": 0,
  "extra": {}
}
```

**字段说明：**
- receiver_id: 接收者ID（私聊时必填）
- group_id: 群ID（群聊时必填）
- content: 消息内容
- message_type: 消息类型（见下表）
- red_packet_id: 红包ID（红包消息时必填）
- extra: 扩展信息（JSON对象）

**消息类型：**
| 类型值 | 说明 |
|--------|------|
| 1 | 文本消息 |
| 2 | 图片消息 |
| 3 | 语音消息 |
| 4 | 视频消息 |
| 5 | 文件消息 |
| 6 | 红包消息 |
| 7 | 位置消息 |
| 8 | 名片消息 |

**响应示例：**
```json
{
  "code": 0,
  "msg": "发送成功",
  "data": {
    "id": 101,
    "sender_id": 1,
    "receiver_id": 2,
    "content": "你好，这是一条测试消息",
    "message_type": 1,
    "status": "sent",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

---

### 6.2 获取私聊消息

**请求**
```
GET /messages/private/:friend_id?page=1&page_size=20&last_id=0
```

**路径参数：**
- friend_id: 好友ID

**查询参数：**
- page: 页码（默认1）
- page_size: 每页数量（默认20，最大100）
- last_id: 最后一条消息ID，用于分页加载

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 100,
        "sender_id": 1,
        "receiver_id": 2,
        "content": "你好",
        "message_type": 1,
        "status": "read",
        "is_revoked": false,
        "created_at": "2024-01-15T10:30:00Z",
        "sender": {
          "id": 1,
          "nickname": "我",
          "avatar": "https://example.com/my-avatar.jpg"
        }
      },
      {
        "id": 99,
        "sender_id": 2,
        "receiver_id": 1,
        "content": "你好！",
        "message_type": 1,
        "status": "read",
        "is_revoked": false,
        "created_at": "2024-01-15T10:29:00Z",
        "sender": {
          "id": 2,
          "nickname": "张三",
          "avatar": "https://example.com/zhangsan.jpg"
        }
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20,
    "has_more": true
  }
}
```

---

### 6.3 获取群消息

**请求**
```
GET /messages/group/:group_id?page=1&page_size=20&last_id=0
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 99,
        "sender_id": 5,
        "group_id": 10,
        "content": "群消息测试",
        "message_type": 1,
        "status": "sent",
        "is_revoked": false,
        "created_at": "2024-01-15T10:25:00Z",
        "sender": {
          "id": 5,
          "nickname": "李四",
          "avatar": "https://example.com/lisi.jpg",
          "group_nickname": "管理员李四",
          "role": 2
        }
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20,
    "has_more": true
  }
}
```

---

### 6.4 标记已读

**请求**
```
POST /messages/read
```

**请求体：**
```json
{
  "target_id": 2,
  "type": 1
}
```

**字段说明：**
- target_id: 会话目标ID（好友ID或群ID）
- type: 会话类型（1:私聊, 2:群聊）

**响应示例：**
```json
{
  "code": 0,
  "msg": "标记成功"
}
```

---

### 6.5 撤回消息

**请求**
```
POST /messages/:message_id/recall
```

**响应示例：**
```json
{
  "code": 0,
  "msg": "撤回成功"
}
```

**错误响应：**
```json
{
  "code": 40005,
  "msg": "超过2分钟，无法撤回",
  "data": null
}
```

---

## 7. 群组接口

### 7.1 创建群

**请求**
```
POST /groups
```

**请求体：**
```json
{
  "name": "技术交流群",
  "avatar": "https://example.com/group.jpg",
  "description": "技术讨论交流",
  "member_ids": [2, 3, 4, 5],
  "join_mode": 1,
  "max_members": 500,
  "is_mute_all": false
}
```

**字段说明：**
- name: 群名称（必填，2-50字符）
- avatar: 群头像URL
- description: 群描述
- member_ids: 初始成员ID列表
- join_mode: 加群方式（1:直接加入, 2:需要验证, 3:禁止加入）
- max_members: 最大成员数（默认500）
- is_mute_all: 是否全员禁言

**响应示例：**
```json
{
  "code": 0,
  "msg": "创建成功",
  "data": {
    "id": 10,
    "name": "技术交流群",
    "avatar": "https://example.com/group.jpg",
    "description": "技术讨论交流",
    "member_count": 5,
    "owner_id": 1,
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

---

### 7.2 获取我的群列表

**请求**
```
GET /groups?page=1&page_size=20
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 10,
        "name": "技术交流群",
        "avatar": "https://example.com/group.jpg",
        "description": "技术讨论交流",
        "member_count": 50,
        "my_role": 2,
        "unread_count": 5,
        "last_message": {
          "id": 99,
          "content": "最新消息",
          "created_at": "2024-01-15T10:25:00Z"
        },
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 20
  }
}
```

---

### 7.3 获取群信息

**请求**
```
GET /groups/:group_id
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "id": 10,
    "name": "技术交流群",
    "avatar": "https://example.com/group.jpg",
    "description": "技术讨论交流",
    "announcement": "群公告内容",
    "member_count": 50,
    "max_members": 500,
    "join_mode": 1,
    "is_mute_all": false,
    "owner_id": 1,
    "my_role": 2,
    "my_nickname": "管理员",
    "created_at": "2024-01-01T00:00:00Z",
    "settings": {
      "show_member_nickname": true,
      "allow_member_invite": true,
      "allow_member_upload": true
    }
  }
}
```

---

### 7.4 更新群信息

**请求**
```
PUT /groups/:group_id
```

**请求体：**
```json
{
  "name": "新群名",
  "avatar": "https://example.com/new-group.jpg",
  "description": "新的描述",
  "announcement": "新的群公告",
  "join_mode": 2,
  "max_members": 1000,
  "is_mute_all": true
}
```

**响应示例：**
```json
{
  "code": 0,
  "msg": "更新成功"
}
```

---

### 7.5 获取群成员列表

**请求**
```
GET /groups/:group_id/members?page=1&page_size=50
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "user_id": 1,
        "nickname": "群主",
        "role": 3,
        "is_muted": false,
        "muted_until": null,
        "joined_at": "2024-01-01T00:00:00Z",
        "user": {
          "id": 1,
          "username": "admin",
          "nickname": "管理员",
          "avatar": "https://example.com/admin.jpg",
          "is_online": true
        }
      },
      {
        "id": 2,
        "user_id": 2,
        "nickname": "管理员",
        "role": 2,
        "is_muted": false,
        "joined_at": "2024-01-01T00:01:00Z",
        "user": {
          "id": 2,
          "username": "zhangsan",
          "nickname": "张三",
          "avatar": "https://example.com/zhangsan.jpg",
          "is_online": true
        }
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 50
  }
}
```

**角色说明：**
| 角色值 | 说明 |
|--------|------|
| 1 | 普通成员 |
| 2 | 管理员 |
| 3 | 群主 |

---

### 7.6 邀请成员入群

**请求**
```
POST /groups/:group_id/invite
```

**请求体：**
```json
{
  "user_ids": [5, 6, 7, 8]
}
```

**响应示例：**
```json
{
  "code": 0,
  "msg": "邀请成功",
  "data": {
    "invited_count": 4,
    "failed_users": []
  }
}
```

---

### 7.7 移除群成员

**请求**
```
DELETE /groups/:group_id/members/:member_id
```

**响应示例：**
```json
{
  "code": 0,
  "msg": "移除成功"
}
```

---

### 7.8 禁言群成员

**请求**
```
POST /groups/:group_id/members/:member_id/mute
```

**请求体：**
```json
{
  "duration": 3600,
  "reason": "违规发言"
}
```

**字段说明：**
- duration: 禁言时长（秒），0表示永久禁言
- reason: 禁言原因

**响应示例：**
```json
{
  "code": 0,
  "msg": "禁言成功",
  "data": {
    "muted_until": "2024-01-15T11:30:00Z"
  }
}
```

---

### 7.9 退出群

**请求**
```
POST /groups/:group_id/leave
```

**响应示例：**
```json
{
  "code": 0,
  "msg": "退出成功"
}
```

---

## 8. 红包接口

### 8.1 发送红包

**请求**
```
POST /redpackets
```

**请求体：**
```json
{
  "type": 1,
  "pay_type": 1,
  "amount": 100.00,
  "total_count": 10,
  "receiver_id": 2,
  "group_id": 0,
  "greeting": "恭喜发财，大吉大利"
}
```

**字段说明：**
- type: 红包类型（1:普通红包, 2:拼手气红包）
- pay_type: 支付类型（1:积分, 2:微信支付, 3:支付宝）
- amount: 红包总金额（必填）
- total_count: 红包个数（必填，最少1个）
- receiver_id: 接收者ID（私聊红包时必填）
- group_id: 群ID（群红包时必填）
- greeting: 祝福语（可选，默认"恭喜发财，大吉大利"）

**响应示例：**
```json
{
  "code": 0,
  "msg": "发送成功",
  "data": {
    "id": 1,
    "type": 1,
    "pay_type": 1,
    "amount": 100.00,
    "total_count": 10,
    "received_count": 0,
    "received_amount": 0,
    "greeting": "恭喜发财，大吉大利",
    "status": 0,
    "created_at": "2024-01-15T10:30:00Z",
    "expire_at": "2024-01-16T10:30:00Z"
  }
}
```

---

### 8.2 抢红包

**请求**
```
POST /redpackets/:id/grab
```

**响应示例（成功）：**
```json
{
  "code": 0,
  "msg": "恭喜你抢到红包",
  "data": {
    "amount": 10.50,
    "pay_type": 1,
    "total_amount": 100.00,
    "my_rank": 3,
    "total_people": 10
  }
}
```

**错误响应：**
```json
{
  "code": 40005,
  "msg": "红包已过期",
  "data": null
}
```

**可能的错误信息：**
- 红包已过期
- 红包已领完
- 不能领取自己的红包
- 积分不足
- 已经领取过该红包

---

### 8.3 获取红包详情

**请求**
```
GET /redpackets/:id
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "red_packet": {
      "id": 1,
      "type": 1,
      "pay_type": 1,
      "amount": 100.00,
      "total_count": 10,
      "received_count": 5,
      "received_amount": 50.00,
      "greeting": "恭喜发财，大吉大利",
      "status": 0,
      "sender": {
        "id": 1,
        "nickname": "张三",
        "avatar": "https://example.com/zhangsan.jpg"
      },
      "created_at": "2024-01-15T10:30:00Z",
      "expire_at": "2024-01-16T10:30:00Z"
    },
    "details": [
      {
        "id": 1,
        "amount": 10.00,
        "created_at": "2024-01-15T10:31:00Z",
        "receiver": {
          "id": 2,
          "nickname": "李四",
          "avatar": "https://example.com/lisi.jpg"
        }
      },
      {
        "id": 2,
        "amount": 12.50,
        "created_at": "2024-01-15T10:32:00Z",
        "receiver": {
          "id": 3,
          "nickname": "王五",
          "avatar": "https://example.com/wangwu.jpg"
        }
      }
    ]
  }
}
```

---

### 8.4 获取我发出的红包

**请求**
```
GET /redpackets/sent?page=1&page_size=20
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "type": 1,
        "pay_type": 1,
        "amount": 100.00,
        "total_count": 10,
        "received_count": 5,
        "received_amount": 50.00,
        "greeting": "恭喜发财，大吉大利",
        "status": 0,
        "created_at": "2024-01-15T10:30:00Z"
      }
    ],
    "total": 10,
    "page": 1,
    "page_size": 20
  }
}
```

---

### 8.5 获取我收到的红包

**请求**
```
GET /redpackets/received?page=1&page_size=20
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 2,
        "type": 2,
        "pay_type": 1,
        "amount": 15.50,
        "greeting": "恭喜发财",
        "total_amount": 100.00,
        "sender": {
          "id": 3,
          "nickname": "王五",
          "avatar": "https://example.com/wangwu.jpg"
        },
        "created_at": "2024-01-15T10:35:00Z"
      }
    ],
    "total": 5,
    "page": 1,
    "page_size": 20
  }
}
```

---

## 9. 朋友圈接口

### 9.1 发布朋友圈

**请求**
```
POST /moments
```

**请求体：**
```json
{
  "content": "今天天气真好",
  "images": [
    "https://example.com/photo1.jpg",
    "https://example.com/photo2.jpg"
  ],
  "location": "北京市朝阳区",
  "latitude": 39.9042,
  "longitude": 116.4074,
  "view_scope": 0,
  "visible_ids": [],
  "at_user_ids": []
}
```

**字段说明：**
- content: 文字内容（最多2000字符）
- images: 图片URL列表（最多9张）
- location: 位置名称
- latitude: 纬度
- longitude: 经度
- view_scope: 查看范围（0:公开, 1:私密, 2:部分可见, 3:不给谁看）
- visible_ids: 可见用户ID列表（view_scope=2时必填）
- at_user_ids: @的用户ID列表

**响应示例：**
```json
{
  "code": 0,
  "msg": "发布成功",
  "data": {
    "id": 1,
    "content": "今天天气真好",
    "images": ["https://example.com/photo1.jpg"],
    "location": "北京市朝阳区",
    "like_count": 0,
    "comment_count": 0,
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

---

### 9.2 获取朋友圈列表

**请求**
```
GET /moments?page=1&page_size=20
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "content": "今天天气真好",
        "images": ["https://example.com/photo1.jpg"],
        "location": "北京市朝阳区",
        "like_count": 5,
        "comment_count": 3,
        "is_liked": false,
        "created_at": "2024-01-15T10:30:00Z",
        "user": {
          "id": 1,
          "nickname": "张三",
          "avatar": "https://example.com/zhangsan.jpg"
        },
        "likes": [
          {
            "id": 1,
            "user": {
              "id": 2,
              "nickname": "李四",
              "avatar": "https://example.com/lisi.jpg"
            },
            "created_at": "2024-01-15T10:35:00Z"
          }
        ],
        "comments": [
          {
            "id": 1,
            "content": "写得真好！",
            "reply_to_user_id": 0,
            "reply_to_user": null,
            "created_at": "2024-01-15T10:40:00Z",
            "user": {
              "id": 3,
              "nickname": "王五",
              "avatar": "https://example.com/wangwu.jpg"
            }
          }
        ]
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20,
    "has_more": true
  }
}
```

---

### 9.3 点赞朋友圈

**请求**
```
POST /moments/:moment_id/like
```

**响应示例：**
```json
{
  "code": 0,
  "msg": "点赞成功"
}
```

---

### 9.4 取消点赞

**请求**
```
DELETE /moments/:moment_id/like
```

**响应示例：**
```json
{
  "code": 0,
  "msg": "取消成功"
}
```

---

### 9.5 评论朋友圈

**请求**
```
POST /moments/:moment_id/comments
```

**请求体：**
```json
{
  "content": "写得真好！",
  "reply_to_user": 0
}
```

**字段说明：**
- content: 评论内容（最多500字符）
- reply_to_user: 回复的评论ID（0表示直接评论，非0表示回复某条评论）

**响应示例：**
```json
{
  "code": 0,
  "msg": "评论成功",
  "data": {
    "id": 1,
    "content": "写得真好！",
    "reply_to_user_id": 0,
    "created_at": "2024-01-15T10:40:00Z"
  }
}
```

---

### 9.6 删除朋友圈

**请求**
```
DELETE /moments/:moment_id
```

**响应示例：**
```json
{
  "code": 0,
  "msg": "删除成功"
}
```

---

## 10. 支付接口

### 10.1 创建支付订单

**请求**
```
POST /payment/orders
```

**请求体：**
```json
{
  "amount": 100.00,
  "pay_type": 1,
  "order_type": 1,
  "subject": "积分充值",
  "description": "充值100积分"
}
```

**字段说明：**
- amount: 支付金额（必填）
- pay_type: 支付方式（1:积分, 2:微信, 3:支付宝）
- order_type: 订单类型（1:积分充值, 2:会员购买, 3:其他）
- subject: 订单标题
- description: 订单描述

**响应示例：**
```json
{
  "code": 0,
  "msg": "订单创建成功",
  "data": {
    "order_id": "ORD2024011510300001",
    "amount": 100.00,
    "pay_type": 1,
    "status": "pending",
    "created_at": "2024-01-15T10:30:00Z",
    "expire_at": "2024-01-15T11:30:00Z",
    "pay_url": ""
  }
}
```

---

### 10.2 支付订单

**请求**
```
POST /payment/orders/:order_id/pay
```

**响应示例：**
```json
{
  "code": 0,
  "msg": "支付成功",
  "data": {
    "order_id": "ORD2024011510300001",
    "status": "paid",
    "paid_at": "2024-01-15T10:35:00Z"
  }
}
```

---

### 10.3 获取订单列表

**请求**
```
GET /payment/orders?page=1&page_size=20&status=all
```

**查询参数：**
- page: 页码
- page_size: 每页数量
- status: 订单状态（all:全部, pending:待支付, paid:已支付, cancelled:已取消, refunded:已退款）

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "order_id": "ORD2024011510300001",
        "amount": 100.00,
        "pay_type": 1,
        "order_type": 1,
        "subject": "积分充值",
        "description": "充值100积分",
        "status": "paid",
        "created_at": "2024-01-15T10:30:00Z",
        "paid_at": "2024-01-15T10:35:00Z"
      }
    ],
    "total": 10,
    "page": 1,
    "page_size": 20
  }
}
```

---

### 10.4 获取积分历史

**请求**
```
GET /payment/points/history?page=1&page_size=20
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "type": 1,
        "amount": 100,
        "balance": 1100,
        "remark": "积分充值",
        "created_at": "2024-01-15T10:35:00Z"
      },
      {
        "id": 2,
        "type": 2,
        "amount": -10,
        "balance": 1090,
        "remark": "发送红包",
        "created_at": "2024-01-15T10:30:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20,
    "current_balance": 1090
  }
}
```

**积分变动类型：**
| 类型值 | 说明 |
|--------|------|
| 1 | 充值 |
| 2 | 消费 |
| 3 | 红包收入 |
| 4 | 退款 |
| 5 | 系统赠送 |
| 6 | 管理员调整 |

---

## 11. 充值提现审核接口

### 11.1 创建充值申请

**请求**
```
POST /api/v1/recharge
Authorization: Bearer <token>
```

**请求体：**
```json
{
  "amount": 100.00,
  "points": 10000,
  "recharge_type": "personal",
  "payment_image": "https://example.com/payment.jpg",
  "remark": "微信支付充值"
}
```

**参数说明：**
- `amount`: 充值金额（元）
- `points`: 兑换的积分数量
- `recharge_type`: 充值方式，`personal`（个人收款码）或 `company`（公司收款码）
- `payment_image`: 付款截图URL
- `remark`: 备注（可选）

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "user_id": 1,
    "amount": 100.00,
    "points": 10000,
    "recharge_type": "personal",
    "payment_image": "https://example.com/payment.jpg",
    "remark": "微信支付充值",
    "status": "pending",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

---

### 11.2 获取我的充值申请

**请求**
```
GET /api/v1/recharge?page=1&page_size=20
Authorization: Bearer <token>
```

**参数说明：**
- `page`: 页码，默认1
- `page_size`: 每页数量，默认20

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "amount": 100.00,
        "points": 10000,
        "status": "approved",
        "created_at": "2024-01-15T10:30:00Z",
        "reviewer": {
          "id": 2,
          "username": "admin"
        }
      }
    ],
    "total": 10,
    "page": 1,
    "page_size": 20
  }
}
```

---

### 11.3 获取充值申请详情

**请求**
```
GET /api/v1/recharge/:id
Authorization: Bearer <token>
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "amount": 100.00,
    "points": 10000,
    "recharge_type": "personal",
    "payment_image": "https://example.com/payment.jpg",
    "status": "approved",
    "review_remark": "审核通过",
    "reviewed_at": "2024-01-15T11:00:00Z",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

---

### 11.4 创建提现申请

**请求**
```
POST /api/v1/withdraw
Authorization: Bearer <token>
```

**请求体：**
```json
{
  "points": 10000,
  "amount": 100.00,
  "withdraw_type": "personal",
  "payment_code": "https://example.com/qrcode.jpg",
  "real_name": "张三",
  "phone": "13800138000",
  "remark": "提现到微信"
}
```

**参数说明：**
- `points`: 提现积分数量
- `amount`: 提现金额（元）
- `withdraw_type`: 提现方式，`personal`（个人收款码）或 `company`（公司收款码）
- `payment_code`: 收款码URL
- `real_name`: 真实姓名
- `phone`: 手机号
- `remark`: 备注（可选）

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "user_id": 1,
    "points": 10000,
    "amount": 100.00,
    "withdraw_type": "personal",
    "status": "pending",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

---

### 11.5 获取我的提现申请

**请求**
```
GET /api/v1/withdraw?page=1&page_size=20
Authorization: Bearer <token>
```

**参数说明：**
- `page`: 页码，默认1
- `page_size`: 每页数量，默认20

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "points": 10000,
        "amount": 100.00,
        "status": "pending",
        "created_at": "2024-01-15T10:30:00Z"
      }
    ],
    "total": 5,
    "page": 1,
    "page_size": 20
  }
}
```

---

### 11.6 获取提现申请详情

**请求**
```
GET /api/v1/withdraw/:id
Authorization: Bearer <token>
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "points": 10000,
    "amount": 100.00,
    "withdraw_type": "personal",
    "payment_code": "https://example.com/qrcode.jpg",
    "real_name": "张三",
    "phone": "13800138000",
    "status": "pending",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

---

## 12. 管理员接口

### 12.1 获取所有充值申请

**请求**
```
GET /api/v1/admin/recharge?status=pending&page=1&page_size=20
Authorization: Bearer <admin_token>
```

**参数说明：**
- `status`: 筛选状态，`pending`（待审核）、`approved`（已通过）、`rejected`（已拒绝）
- `page`: 页码，默认1
- `page_size`: 每页数量，默认20

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "user_id": 1,
        "amount": 100.00,
        "points": 10000,
        "status": "pending",
        "created_at": "2024-01-15T10:30:00Z"
      }
    ],
    "total": 10,
    "page": 1,
    "page_size": 20
  }
}
```

---

### 12.2 审核通过充值申请

**请求**
```
PUT /api/v1/admin/recharge/:id/approve
Authorization: Bearer <admin_token>
```

**请求体：**
```json
{
  "remark": "充值金额正确，已到账"
}
```

**响应示例：**
```json
{
  "code": 0,
  "msg": "审核成功，积分已到账"
}
```

---

### 12.3 拒绝充值申请

**请求**
```
PUT /api/v1/admin/recharge/:id/reject
Authorization: Bearer <admin_token>
```

**请求体：**
```json
{
  "remark": "付款截图不清晰，请重新上传"
}
```

**响应示例：**
```json
{
  "code": 0,
  "msg": "已拒绝"
}
```

---

### 12.4 获取所有提现申请

**请求**
```
GET /api/v1/admin/withdraw?status=pending&page=1&page_size=20
Authorization: Bearer <admin_token>
```

**参数说明：**
- `status`: 筛选状态，`pending`（待审核）、`approved`（已通过）、`rejected`（已拒绝）
- `page`: 页码，默认1
- `page_size`: 每页数量，默认20

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "user_id": 1,
        "points": 10000,
        "amount": 100.00,
        "status": "pending",
        "created_at": "2024-01-15T10:30:00Z"
      }
    ],
    "total": 5,
    "page": 1,
    "page_size": 20
  }
}
```

---

### 12.5 审核通过提现申请

**请求**
```
PUT /api/v1/admin/withdraw/:id/approve
Authorization: Bearer <admin_token>
```

**请求体：**
```json
{
  "remark": "已处理，请查收"
}
```

**响应示例：**
```json
{
  "code": 0,
  "msg": "审核成功，积分已扣除"
}
```

---

### 12.6 拒绝提现申请

**请求**
```
PUT /api/v1/admin/withdraw/:id/reject
Authorization: Bearer <admin_token>
```

**请求体：**
```json
{
  "remark": "收款码不清晰，请重新上传"
}
```

**响应示例：**
```json
{
  "code": 0,
  "msg": "已拒绝"
}
```

---

## 13. 管理员接口

### 11.1 获取数据库统计

**请求**
```
GET /admin/db/stats
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "users": {
      "total": 1000,
      "active_today": 100,
      "active_week": 500,
      "active_month": 800
    },
    "messages": {
      "total": 100000,
      "today": 5000,
      "week": 35000,
      "storage_size": "500MB"
    },
    "groups": {
      "total": 50,
      "total_members": 5000
    },
    "red_packets": {
      "total_sent": 1000,
      "total_amount": "100000.00"
    }
  }
}
```

---

### 11.2 清除旧消息

**请求**
```
POST /admin/db/clear-old-messages?date=2024-01-01
```

**查询参数：**
- date: 删除此日期之前的消息（必填，格式：YYYY-MM-DD）

**响应示例：**
```json
{
  "code": 0,
  "msg": "清理成功",
  "data": {
    "deleted_count": 5000,
    "freed_space": "100MB"
  }
}
```

---

### 11.3 清空所有数据

**请求**
```
POST /admin/db/clear-all
```

**请求体：**
```json
{
  "confirm": "YES",
  "backup": true
}
```

**响应示例：**
```json
{
  "code": 0,
  "msg": "数据已清空",
  "data": {
    "backup_path": "./backups/backup_20240115.sql"
  }
}
```

---

### 11.4 初始化数据库

**请求**
```
POST /admin/db/init
```

**请求体：**
```json
{
  "confirm": "YES"
}
```

**响应示例：**
```json
{
  "code": 0,
  "msg": "数据库初始化成功"
}
```

---

### 11.5 归档旧消息

**请求**
```
POST /admin/db/archive-old?date=2024-01-01
```

**响应示例：**
```json
{
  "code": 0,
  "msg": "归档成功",
  "data": {
    "archived_count": 5000,
    "archive_file": "./archives/messages_2024-01-01.json"
  }
}
```

---

### 11.6 删除用户及数据

**请求**
```
DELETE /admin/db/users/:user_id
```

**响应示例：**
```json
{
  "code": 0,
  "msg": "用户已删除",
  "data": {
    "deleted_messages": 1000,
    "deleted_friendships": 50,
    "deleted_groups": 3
  }
}
```

---

### 11.7 删除群及数据

**请求**
```
DELETE /admin/db/groups/:group_id
```

**响应示例：**
```json
{
  "code": 0,
  "msg": "群已删除",
  "data": {
    "deleted_messages": 5000,
    "deleted_members": 100
  }
}
```

---

### 11.8 手动调整积分

**请求**
```
POST /admin/users/points
```

**请求体：**
```json
{
  "user_id": 1,
  "points": 1000,
  "type": 1,
  "remark": "活动奖励"
}
```

**字段说明：**
- user_id: 用户ID
- points: 积分数（正数增加，负数减少）
- type: 调整类型（1:增加, 2:减少）
- remark: 备注说明

**响应示例：**
```json
{
  "code": 0,
  "msg": "积分调整成功",
  "data": {
    "user_id": 1,
    "before": 1000,
    "change": 1000,
    "after": 2000,
    "balance": 2000
  }
}
```

---

### 11.9 获取系统配置

**请求**
```
GET /admin/system/configs
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "security": {
      "invite_code_enabled": false,
      "captcha_enabled": false,
      "rate_limit_enabled": true
    },
    "payment": {
      "stripe_enabled": false,
      "wechat_pay_enabled": false,
      "alipay_enabled": false
    },
    "storage": {
      "type": "local",
      "max_file_size": 10485760
    }
  }
}
```

---

### 11.10 更新系统配置

**请求**
```
PUT /admin/system/configs
```

**请求体：**
```json
{
  "section": "security",
  "key": "invite_code_enabled",
  "value": "true"
}
```

**响应示例：**
```json
{
  "code": 0,
  "msg": "配置更新成功"
}
```

---

## 12. 系统配置接口

### 12.1 获取系统配置（公开接口）

**请求**
```
GET /api/v1/system/config
```

**响应**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "app_name": "知信",
    "app_version": "v1.0.0",
    "app_description": "知信 - 让沟通更简单",
    "logo_url": "https://example.com/logo.png",
    "favicon_url": "https://example.com/favicon.ico",
    "theme_color": "#07c160",
    "theme_secondary": "#576b95",
    "ui_template": "modern",
    "maintenance_mode": false,
    "maintenance_msg": "",
    "created_at": "2024-01-01 12:00:00",
    "updated_at": "2024-01-01 12:00:00"
  }
}
```

### 12.2 更新系统配置（管理员）

**请求**
```
PUT /api/v1/admin/system/config
Authorization: Bearer <token>
Content-Type: application/json
```

**请求体**
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| app_name | string | 否 | 系统名称 |
| app_version | string | 否 | 版本号 |
| app_description | string | 否 | 系统描述 |
| theme_color | string | 否 | 主题色 |
| theme_secondary | string | 否 | 次要颜色 |
| ui_template | string | 否 | UI模板 (modern/classic/minimal) |

```json
{
  "app_name": "知信",
  "app_version": "v1.0.1",
  "app_description": "知信 - 让沟通更简单",
  "theme_color": "#07c160",
  "theme_secondary": "#576b95",
  "ui_template": "modern"
}
```

**响应**
```json
{
  "code": 200,
  "message": "更新成功",
  "data": { ... }
}
```

### 12.3 上传Logo（管理员）

**请求**
```
POST /api/v1/admin/system/logo
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**请求体**
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file | file | 是 | Logo图片文件，支持png/jpg/jpeg，最大2MB |

**响应**
```json
{
  "code": 200,
  "message": "上传成功",
  "data": {
    "url": "https://example.com/uploads/logo/xxx.png"
  }
}
```

### 12.4 上传图标（管理员）

**请求**
```
POST /api/v1/admin/system/favicon
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**请求体**
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file | file | 是 | Favicon图标文件，支持ico/png，最大500KB |

**响应**
```json
{
  "code": 200,
  "message": "上传成功",
  "data": {
    "url": "https://example.com/uploads/favicon/xxx.ico"
  }
}
```

### 12.5 设置维护模式（管理员）

**请求**
```
POST /api/v1/admin/system/maintenance
Authorization: Bearer <token>
Content-Type: application/json
```

**请求体**
```json
{
  "mode": true,
  "message": "系统维护中，请稍后再试"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| mode | bool | 是 | 是否开启维护模式 |
| message | string | 否 | 维护提示消息 |

**响应**
```json
{
  "code": 200,
  "message": "设置成功",
  "data": {
    "maintenance_mode": true
  }
}
```

---

## 13. WebSocket接口

### 13.1 连接WebSocket

**请求**
```
GET /ws?token=<JWT_TOKEN>
```

**连接成功：**
```json
{
  "type": "connected",
  "data": {
    "user_id": 1,
    "server_time": "2024-01-15T10:30:00Z",
    "heartbeat_interval": 30000
  }
}
```

---

### 12.2 发送消息（客户端→服务器）

**消息格式：**
```json
{
  "type": "message",
  "data": {
    "receiver_id": 2,
    "group_id": 0,
    "content": "你好",
    "message_type": 1
  }
}
```

**心跳包：**
```json
{
  "type": "ping"
}
```

---

### 12.3 接收消息（服务器→客户端）

**新消息：**
```json
{
  "type": "new_message",
  "data": {
    "id": 101,
    "sender_id": 2,
    "receiver_id": 1,
    "content": "你好",
    "message_type": 1,
    "created_at": "2024-01-15T10:30:00Z",
    "sender": {
      "id": 2,
      "nickname": "张三",
      "avatar": "https://example.com/zhangsan.jpg"
    }
  }
}
```

**消息撤回：**
```json
{
  "type": "message_recalled",
  "data": {
    "message_id": "101"
  }
}
```

**消息已读：**
```json
{
  "type": "message_read",
  "data": {
    "user_id": 2,
    "message_id": 100
  }
}
```

**好友申请：**
```json
{
  "type": "friend_request",
  "data": {
    "id": 1,
    "from_user_id": 3,
    "from_user": {
      "id": 3,
      "nickname": "王五",
      "avatar": "https://example.com/wangwu.jpg"
    },
    "message": "你好，我是王五",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

**红包消息：**
```json
{
  "type": "red_packet_message",
  "data": {
    "id": 102,
    "sender_id": 2,
    "group_id": 10,
    "content": "恭喜发财，大吉大利",
    "message_type": 6,
    "red_packet_id": 1,
    "created_at": "2024-01-15T10:30:00Z",
    "sender": {
      "id": 2,
      "nickname": "张三",
      "avatar": "https://example.com/zhangsan.jpg"
    }
  }
}
```

**在线状态：**
```json
{
  "type": "user_online_status",
  "data": {
    "user_id": 2,
    "is_online": true,
    "last_login_at": "2024-01-15T10:30:00Z"
  }
}
```

**系统通知：**
```json
{
  "type": "system_notification",
  "data": {
    "title": "系统消息",
    "content": "您的账号已在新设备登录",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

**心跳响应：**
```json
{
  "type": "pong"
}
```

---

## 附录

### A. 消息类型汇总

| 类型值 | 说明 | 特殊字段 |
|--------|------|---------|
| 1 | 文本消息 | content |
| 2 | 图片消息 | content(url), extra(width, height) |
| 3 | 语音消息 | content(url), extra(duration) |
| 4 | 视频消息 | content(url), extra(thumbnail, duration) |
| 5 | 文件消息 | content(url), extra(filename, size) |
| 6 | 红包消息 | red_packet_id, content |
| 7 | 位置消息 | content(name), extra(latitude, longitude) |
| 8 | 名片消息 | extra(user_id, nickname, avatar) |

### B. 支付类型汇总

| 类型值 | 说明 |
|--------|------|
| 1 | 积分 |
| 2 | 微信支付 |
| 3 | 支付宝 |

### C. 红包类型汇总

| 类型值 | 说明 |
|--------|------|
| 1 | 普通红包 |
| 2 | 拼手气红包 |

### D. 会话类型汇总

| 类型值 | 说明 |
|--------|------|
| 1 | 私聊 |
| 2 | 群聊 |

### E. 订单状态汇总

| 状态值 | 说明 |
|--------|------|
| pending | 待支付 |
| paid | 已支付 |
| cancelled | 已取消 |
| refunded | 已退款 |

### F. 错误码汇总

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 40001 | 参数错误 |
| 40002 | 认证失败 |
| 40003 | 权限不足 |
| 40004 | 资源不存在 |
| 40005 | 操作失败 |
| 50001 | 服务器错误 |
| 50002 | 数据库错误 |
| 50003 | 缓存错误 |

---

**文档版本：** v1.0
**最后更新：** 2024-01-15
**联系支持：** support@example.com
