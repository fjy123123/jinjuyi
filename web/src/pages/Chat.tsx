import React, { useEffect, useRef, useState } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { 
  Layout, 
  List, 
  Input, 
  Button, 
  Avatar, 
  Badge, 
  Dropdown, 
  Space, 
  MenuProps, 
  message as antdMessage,
  Modal,
  Select,
  Switch
} from 'antd';
import { 
  SendOutlined, 
  PlusOutlined, 
  UserOutlined, 
  TeamOutlined, 
  SearchOutlined, 
  LogoutOutlined,
  SettingOutlined,
  SmileOutlined,
  GiftOutlined,
  HistoryOutlined,
  ThunderboltOutlined,
  MoreOutlined,
  PictureOutlined
} from '@ant-design/icons';
import { RootState, setCurrentConversation, setConversations, addMessage } from '@/store';
import { api } from '@/services/api';
import { useWebSocket } from '@/hooks/useWebSocket';
import './Chat.css';

const { Sider, Content } = Layout;
const { TextArea } = Input;
const { Option } = Select;

interface ThemeConfig {
  name: string;
  primaryColor: string;
  bgColor: string;
  sidebarBg: string;
  messageBg: string;
  myMessageBg: string;
  borderColor: string;
}

const THEMES: Record<string, ThemeConfig> = {
  modern: {
    name: 'Modern',
    primaryColor: '#1890ff',
    bgColor: '#f0f2f5',
    sidebarBg: '#ffffff',
    messageBg: '#ffffff',
    myMessageBg: '#95ec69',
    borderColor: '#e8e8e8'
  },
  dark: {
    name: 'Dark',
    primaryColor: '#177ddc',
    bgColor: '#141414',
    sidebarBg: '#1f1f1f',
    messageBg: '#2a2a2a',
    myMessageBg: '#07c160',
    borderColor: '#333'
  }
};

const Chat: React.FC = () => {
  const dispatch = useDispatch();
  const { user } = useSelector((state: RootState) => state.auth);
  const { conversations, currentConversation, messages } = useSelector((state: RootState) => state.chat);
  const [currentTheme, setCurrentTheme] = useState('modern');
  const [inputValue, setInputValue] = useState('');
  const [showSettings, setShowSettings] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  
  const theme = THEMES[currentTheme];
  
  // WebSocket hook
  const { isConnected } = useWebSocket();

  useEffect(() => {
    loadConversations();
  }, []);

  useEffect(() => {
    if (currentConversation) {
      loadMessages();
    }
  }, [currentConversation]);

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const loadConversations = async () => {
    try {
      const res = await api.conversation.getList();
      if (res.code === 0) {
        dispatch(setConversations(res.data));
      }
    } catch (error) {
      console.error(error);
    }
  };

  const loadMessages = async () => {
    if (!currentConversation) return;
    
    try {
      const key = currentConversation.type === 1 
        ? `private_${currentConversation.target_id}`
        : `group_${currentConversation.target_id}`;
      
      let res;
      if (currentConversation.type === 1) {
        res = await api.message.getPrivate(currentConversation.target_id);
      } else {
        res = await api.message.getGroup(currentConversation.target_id);
      }
      
      if (res.code === 0) {
        dispatch(addMessage({
          type: currentConversation.type,
          targetId: currentConversation.target_id,
          message: res.data
        }));
      }
    } catch (error) {
      console.error(error);
    }
  };

  const handleSend = async () => {
    if (!inputValue.trim() || !currentConversation) return;
    
    const content = inputValue.trim();
    setInputValue('');
    
    try {
      const res = await api.message.send({
        receiver_id: currentConversation.type === 1 ? currentConversation.target_id : undefined,
        group_id: currentConversation.type === 2 ? currentConversation.target_id : undefined,
        content,
        message_type: 1
      });
      
      if (res.code === 0) {
        // 成功，消息会通过WebSocket推送回来
      }
    } catch (error) {
      antdMessage.error('发送失败');
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const userMenuItems: MenuProps['items'] = [
    {
      key: 'settings',
      label: '设置',
      icon: <SettingOutlined />,
      onClick: () => setShowSettings(true)
    },
    {
      type: 'divider'
    },
    {
      key: 'logout',
      label: '退出登录',
      icon: <LogoutOutlined />,
      danger: true,
      onClick: () => {
        Modal.confirm({
          title: '确定要退出登录吗？',
          onOk: () => {
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            window.location.href = '/login';
          }
        });
      }
    }
  ];

  const renderMessage = (msg: any) => {
    const isMyMessage = msg.sender_id === user?.id;
    
    return (
      <div 
        key={msg.id} 
        className={`message-item ${isMyMessage ? 'my' : ''}`}
      >
        <Avatar src={msg.sender?.avatar} icon={<UserOutlined />} />
        <div className="message-bubble">
          {msg.is_recall ? (
            <div className="recall-message">
              <span style={{ color: '#999' }}>对方撤回了一条消息</span>
            </div>
          ) : (
            <div className="message-content">{msg.content}</div>
          )}
          <div className="message-time">
            {new Date(msg.created_at).toLocaleTimeString()}
          </div>
        </div>
      </div>
    );
  };

  const currentMessages = currentConversation ? 
    messages[currentConversation.type === 1 
      ? `private_${currentConversation.target_id}` 
      : `group_${currentConversation.target_id}`] || [] 
    : [];

  return (
    <Layout className="chat-layout" style={{ backgroundColor: theme.bgColor }}>
      <Sider width={320} className="chat-sider" style={{ backgroundColor: theme.sidebarBg }}>
        <div className="sider-header">
          <div className="user-info">
            <Dropdown menu={{ items: userMenuItems }} placement="bottomLeft">
              <Avatar src={user?.avatar} icon={<UserOutlined />} size={48} />
            </Dropdown>
            <div className="user-details">
              <div className="username">{user?.nickname || user?.username}</div>
              <div className="user-status">
                {isConnected ? <span className="status-dot online"></span> : <span className="status-dot"></span>}
                {isConnected ? '在线' : '连接中...'}
              </div>
            </div>
          </div>
          <div className="header-actions">
            <Button type="text" icon={<PlusOutlined />} />
            <Button type="text" icon={<TeamOutlined />} />
          </div>
        </div>

        <div className="search-bar">
          <Input placeholder="搜索" prefix={<SearchOutlined />} />
        </div>

        <div className="tab-bar">
          <div className="tab-item active">聊天</div>
          <div className="tab-item">联系人</div>
          <div className="tab-item">发现</div>
        </div>

        <List
          dataSource={conversations}
          renderItem={(conv) => {
            const isActive = currentConversation?.target_id === conv.target_id && currentConversation?.type === conv.type;
            
            return (
              <List.Item
                className={`conversation-item ${isActive ? 'active' : ''}`}
                onClick={() => dispatch(setCurrentConversation(conv))}
              >
                <Avatar 
                  src={conv.target_avatar} 
                  icon={conv.type === 1 ? <UserOutlined /> : <TeamOutlined />} 
                  size={48} 
                />
                <div className="conversation-info">
                  <div className="conv-header">
                    <span className="conv-name">{conv.target_name || 'Chat'}</span>
                    <span className="conv-time">{conv.last_message_at && new Date(conv.last_message_at).toLocaleTimeString()}</span>
                  </div>
                  <div className="conv-preview">
                    <span className="preview-text">{conv.last_message?.content || '暂无消息'}</span>
                    {conv.unread_count > 0 && (
                      <Badge count={conv.unread_count} overflowCount={99} />
                    )}
                  </div>
                </div>
              </List.Item>
            );
          }}
        />
      </Sider>

      <Layout>
        <Content className="chat-content">
          {currentConversation ? (
            <>
              <div className="chat-window" style={{ backgroundColor: theme.bgColor }}>
                <div className="chat-header">
                  <div className="header-info">
                    <Avatar 
                      src={currentConversation.target_avatar} 
                      icon={currentConversation.type === 1 ? <UserOutlined /> : <TeamOutlined />}
                    />
                    <span className="chat-name">
                      {currentConversation.target_name || 'Chat'}
                    </span>
                  </div>
                  <div className="header-actions">
                    <Button type="text" icon={<SearchOutlined />} />
                    <Button type="text" icon={<UserOutlined />} />
                    <Button type="text" icon={<MoreOutlined />} />
                  </div>
                </div>

                <div className="messages-area">
                  {currentMessages.length === 0 ? (
                    <div className="no-messages">
                      <p>开始聊天吧！</p>
                    </div>
                  ) : (
                    currentMessages.map(renderMessage)
                  )}
                  <div ref={messagesEndRef} />
                </div>

                <div className="input-area">
                  <div className="input-toolbar">
                    <Button type="text" icon={<SmileOutlined />} />
                    <Button type="text" icon={<PictureOutlined />} />
                    <Button type="text" icon={<GiftOutlined />} />
                    <Button type="text" icon={<HistoryOutlined />} />
                  </div>
                  <div className="input-wrapper">
                    <TextArea
                      value={inputValue}
                      onChange={(e) => setInputValue(e.target.value)}
                      onKeyPress={handleKeyPress}
                      placeholder="输入消息..."
                      autoSize={{ minRows: 1, maxRows: 4 }}
                      style={{ border: 'none', boxShadow: 'none' }}
                    />
                    <Button
                      type="primary"
                      icon={<SendOutlined />}
                      onClick={handleSend}
                      disabled={!inputValue.trim()}
                    >
                      发送
                    </Button>
                  </div>
                </div>
              </div>
            </>
          ) : (
            <div className="chat-empty">
              <div className="empty-content">
                <ThunderboltOutlined style={{ fontSize: 64, color: '#d9d9d9' }} />
                <p>选择一个聊天开始对话</p>
              </div>
            </div>
          )}
        </Content>
      </Layout>

      {/* 设置弹窗 */}
      <Modal
        title="系统设置"
        open={showSettings}
        onCancel={() => setShowSettings(false)}
        footer={null}
      >
        <div className="settings-content">
          <div className="setting-item">
            <label>UI主题</label>
            <Select 
              value={currentTheme} 
              onChange={setCurrentTheme}
              style={{ width: '100%' }}
            >
              {Object.entries(THEMES).map(([key, config]) => (
                <Option key={key} value={key}>
                  {config.name}
                </Option>
              ))}
            </Select>
          </div>
          
          <div className="setting-section-title">安全设置</div>
          
          <div className="setting-item">
            <label>邀请码注册</label>
            <Switch checked={false} />
          </div>
          
          <div className="setting-item">
            <label>验证码登录</label>
            <Switch checked={false} />
          </div>
        </div>
      </Modal>
    </Layout>
  );
};

export default Chat;