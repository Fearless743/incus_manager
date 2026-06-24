import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { statsAPI } from '../services/api';
import {
  Typography,
  Row,
  Col,
  Card,
  Statistic,
  Button,
  Descriptions,
  Spin,
  Space,
} from 'antd';
import {
  CloudServerOutlined,
  AppstoreOutlined,
  ThunderboltOutlined,
  ShareAltOutlined,
} from '@ant-design/icons';

const { Title, Text } = Typography;

const statCards = [
  { key: 'total_hosts', title: '主机总数', icon: <CloudServerOutlined />, color: '#1677ff' },
  { key: 'total_instances', title: '实例总数', icon: <AppstoreOutlined />, color: '#52c41a' },
  { key: 'running_instances', title: '运行中', icon: <ThunderboltOutlined />, color: '#fa8c16' },
  { key: 'shared_instances', title: '已共享', icon: <ShareAltOutlined />, color: '#722ed1' },
];

const DashboardPage = () => {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [stats, setStats] = useState({
    total_hosts: 0,
    total_instances: 0,
    running_instances: 0,
    shared_instances: 0,
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let active = true;

    (async () => {
      try {
        const response = await statsAPI.get();
        if (active) setStats(response.data);
      } catch (err) {
        console.error('加载统计失败:', err);
      } finally {
        if (active) setLoading(false);
      }
    })();

    return () => {
      active = false;
    };
  }, []);

  return (
    <Spin spinning={loading}>
      <Title level={2}>欢迎，{user?.username}！</Title>
      <Text type="secondary">这里是您的 Incus 管理面板概览。</Text>

      <Row gutter={[16, 16]} style={{ marginTop: 24 }}>
        {statCards.map((card) => (
          <Col key={card.key} xs={24} sm={12} lg={6}>
            <Card>
              <Statistic
                title={card.title}
                value={stats[card.key]}
                prefix={card.icon}
                valueStyle={{ color: card.color }}
              />
            </Card>
          </Col>
        ))}
      </Row>

      <Title level={4} style={{ marginTop: 32 }}>快捷操作</Title>
      <Space wrap style={{ marginTop: 8 }}>
        <Button type="primary" icon={<CloudServerOutlined />} onClick={() => navigate('/hosts')}>
          管理主机
        </Button>
        <Button type="primary" icon={<AppstoreOutlined />} onClick={() => navigate('/instances')}>
          创建实例
        </Button>
        <Button icon={<ShareAltOutlined />} onClick={() => navigate('/shared')}>
          共享实例
        </Button>
      </Space>

      <Card title="系统信息" style={{ marginTop: 32 }}>
        <Descriptions bordered column={2}>
          <Descriptions.Item label="用户">{user?.username}</Descriptions.Item>
          <Descriptions.Item label="邮箱">{user?.email}</Descriptions.Item>
          <Descriptions.Item label="角色">{user?.role}</Descriptions.Item>
          <Descriptions.Item label="服务器时间">
            {new Date().toLocaleString('zh-CN')}
          </Descriptions.Item>
        </Descriptions>
      </Card>
    </Spin>
  );
};

export default DashboardPage;
