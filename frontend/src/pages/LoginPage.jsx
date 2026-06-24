import { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { authAPI } from '../services/api';
import { useNavigate } from 'react-router-dom';
import { Card, Form, Input, Button, Alert, Flex, Typography } from 'antd';
import { CloudServerOutlined } from '@ant-design/icons';

const LoginPage = () => {
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (values) => {
    setError('');
    setLoading(true);

    try {
      const response = await authAPI.login(values.username, values.password);
      login(response.data.user, response.data.token);
      navigate('/dashboard');
    } catch (err) {
      setError(err.response?.data?.error || '登录失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Flex
      align="center"
      justify="center"
      style={{
        minHeight: '100vh',
        background: 'linear-gradient(135deg, #1677ff 0%, #0958d9 50%, #001d66 100%)',
      }}
    >
      <Card
        style={{ width: 400, boxShadow: '0 8px 32px rgba(0, 0, 0, 0.15)' }}
        styles={{ body: { padding: 32 } }}
      >
        <Flex align="center" justify="center" gap={8} style={{ marginBottom: 32 }}>
          <CloudServerOutlined style={{ fontSize: 28, color: '#1677ff' }} />
          <Typography.Title level={3} style={{ margin: 0 }}>
            Incus 管理器
          </Typography.Title>
        </Flex>

        {error && (
          <Alert type="error" message={error} showIcon style={{ marginBottom: 24 }} />
        )}

        <Form layout="vertical" onFinish={handleSubmit} autoComplete="off">
          <Form.Item
            label="用户名"
            name="username"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input size="large" placeholder="请输入用户名" />
          </Form.Item>
          <Form.Item
            label="密码"
            name="password"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password size="large" placeholder="请输入密码" />
          </Form.Item>
          <Form.Item style={{ marginBottom: 0 }}>
            <Button type="primary" htmlType="submit" block size="large" loading={loading}>
              登录
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </Flex>
  );
};

export default LoginPage;
