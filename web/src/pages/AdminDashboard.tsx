import React, { useState, useEffect } from 'react';
import { Layout, Menu, Card, Table, Button, Modal, Form, Input, message, Tag, Space, Switch, InputNumber, Divider, Alert } from 'antd';
import { 
  DashboardOutlined, 
  UserOutlined, 
  WalletOutlined, 
  SettingOutlined, 
  SmileOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  DownloadOutlined,
  DeleteOutlined,
  ClockCircleOutlined
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { RootState } from '../store';
import api from '../services/api';

const { Header, Content, Sider } = Layout;
const { TextArea } = Input;

interface SystemConfig {
  id: number;
  app_name: string;
  app_version: string;
  maintenance_mode: boolean;
  maintenance_msg: string;
  export_enabled: boolean;
  export_max_records: number;
  auto_delete_enabled: boolean;
  auto_delete_days: number;
  auto_delete_last_run: string;
  recall_enabled: boolean;
  recall_timeout: number;
}

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
  const [systemConfig, setSystemConfig] = useState<SystemConfig | null>(null);
  const [cleanupStatus, setCleanupStatus] = useState<any>(null);
  
  const [categoryModalVisible, setCategoryModalVisible] = useState(false);
  const [emojiModalVisible, setEmojiModalVisible] = useState(false);
  const [configModalVisible, setConfigModalVisible] = useState(false);
  const [cleanupModalVisible, setCleanupModalVisible] = useState(false);
  const [form] = Form.useForm();
  const [emojiForm] = Form.useForm();
  const [configForm] = Form.useForm();

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
    } else if (selectedKey === 'settings') {
      loadSystemConfig();
    } else if (selectedKey === 'cleanup') {
      loadCleanupStatus();
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

  const loadSystemConfig = async () => {
    try {
      const response = await api.get('/system/config');
      setSystemConfig(response.data.data);
      configForm.setFieldsValue({
        app_name: response.data.data.app_name,
        maintenance_mode: response.data.data.maintenance_mode,
        maintenance_msg: response.data.data.maintenance_msg,
        export_enabled: response.data.data.export_enabled,
        export_max_records: response.data.data.export_max_records,
        auto_delete_enabled: response.data.data.auto_delete_enabled,
        auto_delete_days: response.data.data.auto_delete_days,
        recall_enabled: response.data.data.recall_enabled,
        recall_timeout: response.data.data.recall_timeout
      });
      setConfigModalVisible(true);
    } catch (error) {
      message.error('加载系统配置失败');
    }
  };

  const loadCleanupStatus = async () => {
    try {
      const response = await api.get('/admin/cleanup/status');
      setCleanupStatus(response.data.data);
      setCleanupModalVisible(true);
    } catch (error) {
      message.error('加载清理状态失败');
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

  const handleUpdateConfig = async (values: any) => {
    try {
      await api.put('/admin/system/configs', {
        ...values,
        maintenance_mode: values.maintenance_mode,
        export_enabled: values.export_enabled,
        export_max_records: values.export_max_records,
        auto_delete_enabled: values.auto_delete_enabled,
        auto_delete_days: values.auto_delete_days,
        recall_enabled: values.recall_enabled,
        recall_timeout: values.recall_timeout
      });
      message.success('配置更新成功');
      setConfigModalVisible(false);
    } catch (error) {
      message.error('更新配置失败');
    }
  };

  const handleForceCleanup = async () => {
    try {
      await api.post('/admin/cleanup/force');
      message.success('清理任务已启动');
      loadCleanupStatus();
    } catch (error) {
      message.error('启动清理失败');
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

      case 'cleanup':
        return (
          <div>
            <h2>消息清理管理</h2>
            {cleanupStatus && (
              <div style={{ marginTop: 20 }}>
                <Alert
                  message="自动清理状态"
                  description={
                    <div>
                      <p>自动清理: {cleanupStatus.auto_delete_enabled ? '已开启' : '已关闭'}</p>
                      <p>清理天数: {cleanupStatus.auto_delete_days} 天</p>
                      <p>上次清理: {cleanupStatus.last_run || '从未'}</p>
                      <p>待清理消息: {cleanupStatus.pending_cleanup} 条</p>
                    </div>
                  }
                  type={cleanupStatus.auto_delete_enabled ? 'info' : 'warning'}
                  showIcon
                />
                <div style={{ marginTop: 16 }}>
                  <Button type="primary" danger icon={<DeleteOutlined />} onClick={handleForceCleanup}>
                    立即执行清理
                  </Button>
                </div>
              </div>
            )}
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
              { key: 'cleanup', icon: <DeleteOutlined />, label: '消息清理' },
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

      <Modal
        title="系统设置"
        open={configModalVisible}
        onCancel={() => setConfigModalVisible(false)}
        footer={null}
        width={700}
      >
        <Form
          form={configForm}
          onFinish={handleUpdateConfig}
          layout="vertical"
          initialValues={systemConfig || {}}
        >
          <Form.Item label="应用名称" name="app_name">
            <Input />
          </Form.Item>

          <Form.Item 
            label="维护模式" 
            name="maintenance_mode" 
            valuePropName="checked"
          >
            <Switch />
          </Form.Item>

          <Form.Item label="维护消息" name="maintenance_msg">
            <TextArea rows={3} />
          </Form.Item>

          <Divider />

          <h4>聊天记录导出设置</h4>
          
          <Form.Item 
            label="允许导出聊天记录" 
            name="export_enabled"
            valuePropName="checked"
            tooltip="关闭后用户将无法导出聊天记录"
          >
            <Switch />
          </Form.Item>

          <Form.Item 
            label="最大导出记录数" 
            name="export_max_records"
            tooltip="单次导出的最大消息记录数量"
          >
            <InputNumber min={100} max={10000} step={100} />
          </Form.Item>

          <Divider />

          <h4>消息自动清理设置</h4>
          
          <Form.Item 
            label="开启消息自动清理" 
            name="auto_delete_enabled"
            valuePropName="checked"
            tooltip="开启后系统将自动清理超过指定天数的旧消息"
          >
            <Switch />
          </Form.Item>

          <Form.Item 
            label="自动清理天数" 
            name="auto_delete_days"
            tooltip="消息保留天数，超过此天数的旧消息将被自动删除"
          >
            <InputNumber min={1} max={365} step={1} addonAfter="天" />
          </Form.Item>

          <Alert
            message="⚠️ 警告"
            description="开启自动清理后，超过指定天数的旧消息将被永久删除，且会通知所有在线用户清理本地消息。此操作不可恢复！"
            type="warning"
            showIcon
            style={{ marginBottom: 16 }}
          />

          <Divider />

          <h4>消息撤回设置</h4>
          
          <Form.Item 
            label="允许消息撤回" 
            name="recall_enabled"
            valuePropName="checked"
            tooltip="关闭后用户将无法撤回已发送的消息"
          >
            <Switch />
          </Form.Item>

          <Form.Item 
            label="撤回时限" 
            name="recall_timeout"
            tooltip="发送消息后多长时间内可以撤回"
          >
            <InputNumber min={10} max={3600} step={10} addonAfter="秒" />
          </Form.Item>

          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit">保存配置</Button>
              <Button onClick={() => setConfigModalVisible(false)}>取消</Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </Layout>
  );
};

export default AdminDashboard;
