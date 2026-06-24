import { useState, useEffect } from 'react';
import { hostAPI } from '../services/api';
import {
  Typography,
  Button,
  Table,
  Tag,
  Modal,
  Form,
  Input,
  Space,
  Popconfirm,
  message,
  Empty,
} from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';

const { Title } = Typography;

const HostsPage = () => {
  const [hosts, setHosts] = useState([]);
  const [modalOpen, setModalOpen] = useState(false);
  const [editingHost, setEditingHost] = useState(null);
  const [loading, setLoading] = useState(false);
  const [connecting, setConnecting] = useState(false);
  const [tableLoading, setTableLoading] = useState(true);
  const [form] = Form.useForm();

  const loadHosts = async () => {
    setTableLoading(true);
    try {
      const response = await hostAPI.getAll();
      setHosts(response.data);
    } catch (err) {
      console.error('加载主机失败:', err);
      message.error('加载主机列表失败');
    } finally {
      setTableLoading(false);
    }
  };

  useEffect(() => {
    let active = true;

    (async () => {
      try {
        const response = await hostAPI.getAll();
        if (active) setHosts(response.data);
      } catch (err) {
        if (active) {
          console.error('加载主机失败:', err);
          message.error('加载主机列表失败');
        }
      } finally {
        if (active) setTableLoading(false);
      }
    })();

    return () => {
      active = false;
    };
  }, []);

  const openCreateModal = () => {
    setEditingHost(null);
    form.resetFields();
    setModalOpen(true);
  };

  const openEditModal = (host) => {
    setEditingHost(host);
    form.setFieldsValue({
      name: host.name,
      address: host.address,
      certificate: host.certificate || '',
    });
    setModalOpen(true);
  };

  const closeModal = () => {
    setModalOpen(false);
    setEditingHost(null);
    form.resetFields();
  };

  const handleSubmit = async (values) => {
    setLoading(true);
    try {
      setConnecting(true);
      const response = await hostAPI.test(values.address, values.certificate);
      setConnecting(false);

      if (response.data.success) {
        if (editingHost) {
          await hostAPI.update(editingHost.id, values);
          message.success('主机已更新');
        } else {
          await hostAPI.add(values.name, values.address, values.certificate);
          message.success('主机已添加');
        }
        closeModal();
        loadHosts();
      } else {
        message.error('连接测试失败：' + response.data.message);
      }
    } catch (err) {
      setConnecting(false);
      message.error('连接测试失败：' + (err.response?.data?.error || '无法连接到主机'));
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id) => {
    try {
      await hostAPI.delete(id);
      message.success('主机已删除');
      loadHosts();
    } catch (err) {
      message.error('删除主机失败');
      console.error('删除主机失败:', err);
    }
  };

  const columns = [
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '地址',
      dataIndex: 'address',
      key: 'address',
    },
    {
      title: '项目',
      dataIndex: 'project',
      key: 'project',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status) => (
        <Tag color={status === 'active' ? 'success' : 'default'}>
          {status === 'active' ? '在线' : status?.toUpperCase()}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'actions',
      render: (_, record) => (
        <Space>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => openEditModal(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除此主机吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <Space style={{ marginBottom: 16, width: '100%', justifyContent: 'space-between' }}>
        <Title level={2} style={{ margin: 0 }}>主机管理</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreateModal}>
          添加主机
        </Button>
      </Space>

      <Table
        columns={columns}
        dataSource={hosts}
        rowKey="id"
        loading={tableLoading}
        locale={{ emptyText: <Empty description="尚未添加主机" /> }}
      />

      <Modal
        title={editingHost ? '编辑主机' : '新建主机'}
        open={modalOpen}
        onCancel={closeModal}
        footer={null}
        destroyOnHidden
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          <Form.Item
            label="名称"
            name="name"
            rules={[{ required: true, message: '请输入主机名称' }]}
          >
            <Input placeholder="主机名称" />
          </Form.Item>
          <Form.Item
            label="地址（IP:端口，如 192.168.1.100:8443）"
            name="address"
            rules={[{ required: true, message: '请输入主机地址' }]}
          >
            <Input placeholder="192.168.1.100:8443" />
          </Form.Item>
          <Form.Item
            label="凭证"
            name="certificate"
            rules={[{ required: true, message: '请输入凭证' }]}
          >
            <Input.TextArea rows={4} style={{ fontFamily: 'monospace', fontSize: 12 }} />
          </Form.Item>
          <Form.Item style={{ marginBottom: 0 }}>
            <Space>
              <Button
                type="primary"
                htmlType="submit"
                loading={loading || connecting}
              >
                {connecting ? '测试连接中...' : editingHost ? '保存' : '添加主机'}
              </Button>
              <Button onClick={closeModal}>取消</Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default HostsPage;
