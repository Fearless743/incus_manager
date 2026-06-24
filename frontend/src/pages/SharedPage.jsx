import { useState, useEffect } from 'react';
import { instanceAPI, shareAPI } from '../services/api';
import {
  Typography,
  Button,
  Table,
  Tag,
  Card,
  Form,
  Select,
  InputNumber,
  DatePicker,
  Space,
  Popconfirm,
  message,
  Empty,
} from 'antd';

const { Title } = Typography;

const getExpiryStatus = (expiryDate) => {
  if (!expiryDate) return { color: 'default', text: '永不过期' };
  const now = new Date();
  const expiry = new Date(expiryDate);
  if (expiry < now) return { color: 'error', text: '已过期' };
  const daysLeft = Math.ceil((expiry - now) / (1000 * 60 * 60 * 24));
  return { color: 'processing', text: `剩余 ${daysLeft} 天` };
};

const SharedPage = () => {
  const [myInstances, setMyInstances] = useState([]);
  const [sharedInstances, setSharedInstances] = useState([]);
  const [loading, setLoading] = useState(false);
  const [tableLoading, setTableLoading] = useState(true);
  const [form] = Form.useForm();

  const loadInstances = async () => {
    try {
      const response = await instanceAPI.getAll();
      setMyInstances(response.data);
    } catch (err) {
      console.error('加载实例失败:', err);
      message.error('加载实例列表失败');
    }
  };

  const loadSharedInstances = async () => {
    setTableLoading(true);
    try {
      const response = await instanceAPI.getAll();
      const shared = response.data.filter(
        (inst) => inst.shared_with && inst.shared_with.length > 0
      );
      setSharedInstances(shared);
    } catch (err) {
      console.error('加载共享实例失败:', err);
      message.error('加载共享实例失败');
    } finally {
      setTableLoading(false);
    }
  };

  useEffect(() => {
    let active = true;

    (async () => {
      try {
        const response = await instanceAPI.getAll();
        if (active) setMyInstances(response.data);
      } catch (err) {
        if (active) {
          console.error('加载实例失败:', err);
          message.error('加载实例列表失败');
        }
      }
    })();

    (async () => {
      try {
        const response = await instanceAPI.getAll();
        if (active) {
          const shared = response.data.filter(
            (inst) => inst.shared_with && inst.shared_with.length > 0
          );
          setSharedInstances(shared);
        }
      } catch (err) {
        if (active) {
          console.error('加载共享实例失败:', err);
          message.error('加载共享实例失败');
        }
      } finally {
        if (active) setTableLoading(false);
      }
    })();

    return () => {
      active = false;
    };
  }, []);

  const handleShare = async (values) => {
    setLoading(true);
    try {
      const expiresDateTime = values.expiresAt.toDate().toISOString();
      await shareAPI.share(values.instanceId, values.userId, expiresDateTime);
      message.success('实例共享成功');
      form.resetFields();
      loadInstances();
      loadSharedInstances();
    } catch (err) {
      message.error('共享实例失败');
      console.error('共享实例失败:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleRevoke = async (instanceId, userId) => {
    try {
      await shareAPI.revoke(instanceId, userId);
      message.success('已撤销共享权限');
      loadInstances();
      loadSharedInstances();
    } catch (err) {
      message.error('撤销共享失败');
      console.error('撤销共享失败:', err);
    }
  };

  const columns = [
    {
      title: '实例名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status) => (
        <Tag color={status === 'running' ? 'success' : 'error'}>
          {status === 'running' ? '运行中' : '已停止'}
        </Tag>
      ),
    },
    {
      title: '共享给',
      dataIndex: 'shared_with',
      key: 'shared_with',
      render: (sharedWith) => sharedWith?.join(', ') || '无',
    },
    {
      title: '过期状态',
      dataIndex: 'expiry_date',
      key: 'expiry_date',
      render: (expiryDate) => {
        const status = getExpiryStatus(expiryDate);
        return <Tag color={status.color}>{status.text}</Tag>;
      },
    },
    {
      title: '操作',
      key: 'actions',
      render: (_, record) => (
        <Space wrap>
          {record.shared_with?.map((userId) => (
            <Popconfirm
              key={userId}
              title={`确定要撤销用户 #${userId} 的权限吗？`}
              onConfirm={() => handleRevoke(record.id, userId)}
              okText="确定"
              cancelText="取消"
            >
              <Button danger size="small">
                撤销用户 #{userId}
              </Button>
            </Popconfirm>
          ))}
        </Space>
      ),
    },
  ];

  return (
    <div>
      <Title level={2}>实例共享</Title>

      <Card title="共享实例" style={{ marginBottom: 24 }}>
        <Form form={form} layout="vertical" onFinish={handleShare}>
          <Space wrap style={{ width: '100%' }} size="middle">
            <Form.Item
              label="选择实例"
              name="instanceId"
              rules={[{ required: true, message: '请选择实例' }]}
              style={{ minWidth: 200 }}
            >
              <Select
                placeholder="选择实例"
                options={myInstances.map((inst) => ({
                  value: inst.id,
                  label: inst.name,
                }))}
              />
            </Form.Item>
            <Form.Item
              label="共享给用户 ID"
              name="userId"
              rules={[{ required: true, message: '请输入用户 ID' }]}
              style={{ minWidth: 160 }}
            >
              <InputNumber placeholder="用户 ID" style={{ width: '100%' }} />
            </Form.Item>
            <Form.Item
              label="过期时间"
              name="expiresAt"
              rules={[{ required: true, message: '请选择过期时间' }]}
              style={{ minWidth: 220 }}
            >
              <DatePicker showTime format="YYYY-MM-DD HH:mm:ss" style={{ width: '100%' }} />
            </Form.Item>
            <Form.Item label=" " style={{ marginBottom: 0 }}>
              <Button type="primary" htmlType="submit" loading={loading}>
                共享实例
              </Button>
            </Form.Item>
          </Space>
        </Form>
      </Card>

      <Title level={4}>当前共享的实例</Title>
      <Table
        columns={columns}
        dataSource={sharedInstances}
        rowKey="id"
        loading={tableLoading}
        locale={{ emptyText: <Empty description="暂无共享实例" /> }}
        style={{ marginTop: 8 }}
      />
    </div>
  );
};

export default SharedPage;
