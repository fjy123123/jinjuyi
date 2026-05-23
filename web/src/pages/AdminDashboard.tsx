import React, { useState, useEffect } from 'react';
import { Layout, Menu, Card, Table, Button, Modal, Form, Input, message, Tag, Space } from 'antd';
import { 
  DashboardOutlined, 
  UserOutlined, 
  WalletOutlined, 
  SettingOutlined, 
  SmileOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { RootState } from '../store';
import api from '../services/api';

const { Header, Content, Sider } = Layout;
const { TextArea } = Input;

interface RechargeRequest {
  id: number;
  user_id: number;
  amount: number;
  proof_image: string;
  status: string;
  created_at: string;
}

interface WithdrawRequest {
  id: number;
  user_id: number;
  amount: number;
  account_info: string;
  status: string;
  created_at: string;
}

interface EmojiCategory {
  id: number;
  name: string;
  icon: string;
}

interface EmojiItem {
  id: number;
  category_id: number;
  name: string;
  url: string;
}

const AdminDashboard: React.FC = () => {
  const navigate = useNavigate();
  const { isAuthenticated, user } = useSelector((state: RootState) => state.auth);
  const [selectedKey, setSelectedKey] = useState('dashboard');
  const [stats, setStats] = useState<any>({});
  const [rechargeRequests, setRechargeRequests] = useState<RechargeRequest[]>([]);
  const [withdrawRequests, setWithdrawRequests] = useState<WithdrawRequest[]>([]);
  const [emojiCategories, setEmojiCategories] = useState<EmojiCategory[]>([]);
  const [emojis, setEmojis] = useState<EmojiItem[]>([]);
  
  const [categoryModalVisible, setCategoryModalVisible] = useState(false);
  const [emojiModalVisible, setEmojiModalVisible] = useState(false);
  const [form] = Form.useForm();
  const [emojiForm] = Form.useForm();

  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/login');
      return;
    }
    if (selectedKey === 'dashboard') {
      loadStats();
    } else if (selectedKey === 'recharge') {
      loadRechargeRequests();
    } else if (selectedKey === 'withdraw') {
      loadWithdrawRequests();
    } else if (selectedKey === 'emoji') {
      loadEmojiCategories();
    }
  }, [selectedKey, isAuthenticated, navigate]);

  const loadStats = async () => {
    try {
      const response = await api.get('/admin/db/stats');
      setStats(response.data.data || response.data);
    } catch (error) {
      message.error('加载统计数据失败');
    }
  };

  const loadRechargeRequests = async () => {
    try {
      const response = await api.get('/admin/recharge');
      setRechargeRequests(response.data.data || []);
    } catch (error) {
      message.error('加载充值申请失败');
    }
  };

  const loadWithdrawRequests = async () => {
    try {
      const response = await api.get('/admin/withdraw');
      setWithdrawRequests(response.data.data || []);
    } catch (error) {
      message.error('加载提现申请失败');
    }
  };

  const loadEmojiCategories = async () => {
    try {
      const response = await api.get('/emoji/categories');
      setEmojiCategories(response.data.data || []);
    } catch (error) {
      message.error('加载表情包分类失败');
    }
  };

  const handleRechargeAction = async (id: number, action: 'approve' | 'reject') => {
    try {
      await api.put(`/admin/recharge/${id}/${action}`);
      message.success(action === 'approve' ? '审核通过' : '已拒绝');
      loadRechargeRequests();
    } catch (error) {
      message.error('操作失败');
    }
  };

  const handleWithdrawAction = async (id: number, action: 'approve' | 'reject') => {
    try {
      await api.put(`/admin/withdraw/${id}/${action}`);
      message.success(action === 'approve' ? '审核通过' : '已拒绝');
      loadWithdrawRequests();
    } catch (error) {
      message.error('操作失败');
    }
  };

  const handleAddCategory = async (values: any) => {
    try {
      await api.post('/admin/emoji/categories', values);
      message.success('添加分类成功');
      setCategoryModalVisible(false);
      form.resetFields();
      loadEmojiCategories();
    } catch (error) {
      message.error('添加分类失败');
    }
  };

  const handleDeleteCategory = async (id: number) => {
    try {
      await api.delete(`/admin/emoji/categories/${id}`);
      message.success('删除成功');
      loadEmojiCategories();
    } catch (error) {
      message.error('删除失败');
    }
  };

  const getStatusTag = (status: string) => {
    const colors: Record<string, string> = {
      pending: 'orange',
      approved: 'green',
      rejected: 'red'
    };
    const labels: Record<string, string> = {
      pending: '待审核',
      approved: '已通过',
      rejected: '已拒绝'
    };
    return <Tag color={colors[status] || 'default'}>{labels[status] || status}</Tag>;
  };

  const renderContent = () => {
    switch (selectedKey) {
      case 'dashboard':
        return (
          <div>
            <h2>系统概览</h2>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 16, marginTop: 20 }}>
              <Card title="用户数" bordered={false}>{stats.user_count || 0}</Card>
              <Card title="群组数" bordered={false}>{stats.group_count || 0}</Card>
              <Card title="消息数" bordered={false}>{stats.mongo_message_count || 0}</Card>
              <Card title="归档消息" bordered={false}>{stats.archived_message_count || 0}</Card>
            </div>
          </div>
        );
      
      case 'recharge':
        return (
          <div>
            <h2>充值审核</h2>
            <Table 
              dataSource={rechargeRequests}
              rowKey="id"
              columns={[
                { title: 'ID', dataIndex: 'id', key: 'id' },
                { title: '用户ID', dataIndex: 'user_id', key: 'user_id' },
                { title: '金额', dataIndex: 'amount', key: 'amount', render: (v: number) => `¥${v}` },
                { title: '状态', dataIndex: 'status', key: 'status', render: getStatusTag },
                { title: '创建时间', dataIndex: 'created_at', key: 'created_at' },
                { 
                  title: '操作', 
                  key: 'action', 
                  render: (_, record: RechargeRequest) => (
                    <Space>
                      {record.status === 'pending' && (
                        <>
                          <Button 
                            type="primary" 
                            icon={<CheckCircleOutlined />} 
                            onClick={() => handleRechargeAction(record.id, 'approve')}
                          >通过</Button>
                          <Button 
                            danger 
                            icon={<CloseCircleOutlined />} 
                            onClick={() => handleRechargeAction(record.id, 'reject')}
                          >拒绝</Button>
                        </>
                      )}
                    </Space>
                  )
                }
              ]}
            />
          </div>
        );
      
      case 'withdraw':
        return (
          <div>
            <h2>提现审核</h2>
            <Table 
              dataSource={withdrawRequests}
              rowKey="id"
              columns={[
                { title: 'ID', dataIndex: 'id', key: 'id' },
                { title: '用户ID', dataIndex: 'user_id', key: 'user_id' },
                { title: '金额', dataIndex: 'amount', key: 'amount', render: (v: number) => `¥${v}` },
                { title: '账户信息', dataIndex: 'account_info', key: 'account_info' },
                { title: '状态', dataIndex: 'status', key: 'status', render: getStatusTag },
                { title: '创建时间', dataIndex: 'created_at', key: 'created_at' },
                { 
                  title: '操作', 
                  key: 'action', 
                  render: (_, record: WithdrawRequest) => (
                    <Space>
                      {record.status === 'pending' && (
                        <>
                          <Button 
                            type="primary" 
                            icon={<CheckCircleOutlined />} 
                            onClick={() => handleWithdrawAction(record.id, 'approve')}
                          >通过</Button>
                          <Button 
                            danger 
                            icon={<CloseCircleOutlined />} 
                            onClick={() => handleWithdrawAction(record.id, 'reject')}
                          >拒绝</Button>
                        </>
                      )}
                    </Space>
                  )
                }
              ]}
            />
          </div>
        );
      
      case 'emoji':
        return (
          <div>
            <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
              <h2>表情包管理</h2>
              <Button type="primary" onClick={() => setCategoryModalVisible(true)}>添加分类</Button>
            </div>
            <Table 
              dataSource={emojiCategories}
              rowKey="id"
              columns={[
                { title: 'ID', dataIndex: 'id', key: 'id' },
                { title: '分类名称', dataIndex: 'name', key: 'name' },
                { title: '图标', dataIndex: 'icon', key: 'icon' },
                { 
                  title: '操作', 
                  key: 'action', 
                  render: (_, record: EmojiCategory) => (
                    <Space>
                      <Button danger onClick={() => handleDeleteCategory(record.id)}>删除</Button>
                    </Space>
                  )
                }
              ]}
            />
          </div>
        );
      
      default:
        return <div>请选择功能模块</div>;
    }
  };

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ background: '#001529', color: '#fff', display: 'flex', alignItems: 'center' }}>
        <h2 style={{ color: '#fff', margin: 0 }}>知信管理后台</h2>
        <div style={{ marginLeft: 'auto' }}>
          <span>欢迎，管理员</span>
          <Button type="link" style={{ color: '#fff' }} onClick={() => navigate('/')}>返回聊天</Button>
        </div>
      </Header>
      <Layout>
        <Sider width={200} theme="dark">
          <Menu
            mode="inline"
            selectedKeys={[selectedKey]}
            style={{ height: '100%', borderRight: 0 }}
            onClick={({ key }) => setSelectedKey(key)}
            items={[
              { key: 'dashboard', icon: <DashboardOutlined />, label: '系统概览' },
              { key: 'recharge', icon: <WalletOutlined />, label: '充值审核' },
              { key: 'withdraw', icon: <WalletOutlined />, label: '提现审核' },
              { key: 'emoji', icon: <SmileOutlined />, label: '表情包管理' },
              { key: 'settings', icon: <SettingOutlined />, label: '系统设置' }
            ]}
          />
        </Sider>
        <Layout style={{ padding: '24px' }}>
          <Content>
            {renderContent()}
          </Content>
        </Layout>
      </Layout>

      <Modal 
        title="添加表情包分类" 
        open={categoryModalVisible} 
        onCancel={() => setCategoryModalVisible(false)}
        onOk={() => form.submit()}
      >
        <Form form={form} onFinish={handleAddCategory} layout="vertical">
          <Form.Item label="分类名称" name="name" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item label="图标URL" name="icon">
            <Input />
          </Form.Item>
        </Form>
      </Modal>
    </Layout>
  );
};

export default AdminDashboard;
