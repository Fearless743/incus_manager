import { useState, useEffect } from 'react';
import { instanceAPI, hostAPI } from '../services/api';
import {
  Typography,
  Button,
  Table,
  Tag,
  Modal,
  Form,
  Input,
  Select,
  InputNumber,
  Space,
  Popconfirm,
  message,
  Empty,
} from 'antd';
import { PlusOutlined, PlayCircleOutlined, PauseCircleOutlined, DeleteOutlined } from '@ant-design/icons';

const { Title } = Typography;

const getStatusTag = (status) => {
  const map = {
    running: { color: 'success', text: '运行中' },
    stopped: { color: 'error', text: '已停止' },
    creating: { color: 'warning', text: '创建中' },
    deleted: { color: 'default', text: '已删除' },
    paused: { color: 'default', text: '已暂停' },
  };
  const item = map[status] || { color: 'default', text: status?.toUpperCase() || '未知' };
  return <Tag color={item.color}>{item.text}</Tag>;
};

const InstancesPage = () => {
  const [instances, setInstances] = useState([]);
  const [modalOpen, setModalOpen] = useState(false);
  const [hosts, setHosts] = useState([]);
  const [images, setImages] = useState([]);
  const [loading, setLoading] = useState(false);
  const [tableLoading, setTableLoading] = useState(true);
  const [form] = Form.useForm();

  useEffect(() => {
    let active = true;

    (async () => {
      try {
        const instancesRes = await instanceAPI.getAll();
        if (active) setInstances(instancesRes.data);
      } catch (err) {
        if (active) {
          console.error('加载实例失败:', err);
          message.error('加载实例列表失败');
        }
      } finally {
        if (active) setTableLoading(false);
      }
    })();

    (async () => {
      try {
        const hostsRes = await hostAPI.getAll();
        if (active) setHosts(hostsRes.data);
      } catch (err) {
        if (active) console.error('加载主机失败:', err);
      }
    })();

    (async () => {
      try {
        const imagesRes = await instanceAPI.getImages();
        if (active) setImages(imagesRes.data);
      } catch (err) {
        if (active) console.error('加载镜像失败:', err);
      }
    })();

    return () => {
      active = false;
    };
  }, []);

  const loadInstances = async () => {
    setTableLoading(true);
    try {
      const response = await instanceAPI.getAll();
      setInstances(response.data);
    } catch (err) {
      console.error('加载实例失败:', err);
      message.error('加载实例列表失败');
    } finally {
      setTableLoading(false);
    }
  };

  const openCreateModal = () => {
    form.resetFields();
    form.setFieldsValue({
      image: 'ubuntu/22.04',
      cpu: 2,
      memory: 2048,
      disk: 20,
    });
    setModalOpen(true);
  };

  const closeModal = () => {
    setModalOpen(false);
    form.resetFields();
  };

  const handleSubmit = async (values) => {
    setLoading(true);
    try {
      await instanceAPI.create({
        name: values.name,
        image: values.image,
        cpu: values.cpu,
        memory: values.memory,
        disk: values.disk,
        hostId: values.hostId,
        ports: [],
        network_limit: 'unlimited',
        upload_limit: 'unlimited',
        download_limit: 'unlimited',
      });
      message.success('实例创建成功');
      closeModal();
      loadInstances();
    } catch (err) {
      message.error('创建实例失败');
      console.error('创建实例失败:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id) => {
    try {
      await instanceAPI.delete(id);
      message.success('实例已删除');
      loadInstances();
    } catch (err) {
      message.error('删除实例失败');
      console.error('删除实例失败:', err);
    }
  };

  const handleStart = async (id) => {
    try {
      await instanceAPI.start(id);
      message.success('实例已启动');
      loadInstances();
    } catch (err) {
      message.error('启动实例失败');
      console.error('启动实例失败:', err);
    }
  };

  const handleStop = async (id) => {
    try {
      await instanceAPI.stop(id);
      message.success('实例已停止');
      loadInstances();
    } catch (err) {
      message.error('停止实例失败');
      console.error('停止实例失败:', err);
    }
  };

  const columns = [
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '镜像',
      dataIndex: 'image',
      key: 'image',
    },
    {
      title: 'CPU',
      dataIndex: 'cpu',
      key: 'cpu',
      render: (cpu) => `${cpu} 核`,
    },
    {
      title: '内存',
      dataIndex: 'memory',
      key: 'memory',
      render: (memory) => `${memory} MB`,
    },
    {
      title: '磁盘',
      dataIndex: 'disk',
      key: 'disk',
      render: (disk) => `${disk} GB`,
    },
    {
      title: '映射 IP',
      dataIndex: 'mapping_ip',
      key: 'mapping_ip',
      render: (ip) => ip || '-',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status) => getStatusTag(status),
    },
    {
      title: '操作',
      key: 'actions',
      render: (_, record) => (
        <Space>
          <Button
            type="link"
            icon={<PlayCircleOutlined />}
            onClick={() => handleStart(record.id)}
          >
            启动
          </Button>
          <Button
            type="link"
            icon={<PauseCircleOutlined />}
            onClick={() => handleStop(record.id)}
          >
            停止
          </Button>
          <Popconfirm
            title="确定要删除此实例吗？"
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
        <Title level={2} style={{ margin: 0 }}>实例管理</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreateModal}>
          创建实例
        </Button>
      </Space>

      <Table
        columns={columns}
        dataSource={instances}
        rowKey="id"
        loading={tableLoading}
        locale={{ emptyText: <Empty description="尚未创建实例" /> }}
      />

      <Modal
        title="新建实例"
        open={modalOpen}
        onCancel={closeModal}
        footer={null}
        destroyOnHidden
        width={600}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          <Form.Item
            label="名称"
            name="name"
            rules={[{ required: true, message: '请输入实例名称' }]}
          >
            <Input placeholder="实例名称" />
          </Form.Item>
          <Form.Item
            label="镜像"
            name="image"
            rules={[{ required: true, message: '请选择镜像' }]}
          >
            <Select
              options={images.map((img) => ({ value: img, label: img }))}
            />
          </Form.Item>
          <Form.Item
            label="主机"
            name="hostId"
            rules={[{ required: true, message: '请选择主机' }]}
          >
            <Select
              placeholder="选择主机"
              options={hosts.map((host) => ({ value: host.id, label: host.name }))}
            />
          </Form.Item>
          <Space style={{ display: 'flex' }} align="start">
            <Form.Item
              label="CPU 核心数"
              name="cpu"
              rules={[{ required: true, message: '请输入 CPU 核心数' }]}
            >
              <InputNumber min={1} style={{ width: 120 }} />
            </Form.Item>
            <Form.Item
              label="内存（MB）"
              name="memory"
              rules={[{ required: true, message: '请输入内存大小' }]}
            >
              <InputNumber min={256} style={{ width: 120 }} />
            </Form.Item>
            <Form.Item
              label="磁盘（GB）"
              name="disk"
              rules={[{ required: true, message: '请输入磁盘大小' }]}
            >
              <InputNumber min={1} style={{ width: 120 }} />
            </Form.Item>
          </Space>
          <Form.Item style={{ marginBottom: 0 }}>
            <Space>
              <Button type="primary" htmlType="submit" loading={loading}>
                创建实例
              </Button>
              <Button onClick={closeModal}>取消</Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default InstancesPage;
