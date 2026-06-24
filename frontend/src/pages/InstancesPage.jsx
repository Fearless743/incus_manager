import { useState, useEffect } from 'react';
import { instanceAPI, hostAPI } from '../services/api';

const InstancesPage = () => {
  const [instances, setInstances] = useState([]);
  const [showForm, setShowForm] = useState(false);
  const [name, setName] = useState('');
  const [image, setImage] = useState('ubuntu/22.04');
  const [cpu, setCpu] = useState(2);
  const [memory, setMemory] = useState(2048);
  const [disk, setDisk] = useState(20);
  const [hostId, setHostId] = useState('');
  const [hosts, setHosts] = useState([]);
  const [images, setImages] = useState([]);

  useEffect(() => {
    loadInstances();
    loadHosts();
    loadImages();
  }, []);

  const loadInstances = async () => {
    try {
      const response = await instanceAPI.getAll();
      setInstances(response.data);
    } catch (err) {
      console.error('加载实例失败:', err);
    }
  };

  const loadHosts = async () => {
    try {
      const response = await hostAPI.getAll();
      setHosts(response.data);
    } catch (err) {
      console.error('加载主机失败:', err);
    }
  };

  const loadImages = async () => {
    try {
      const response = await instanceAPI.getImages();
      setImages(response.data);
    } catch (err) {
      console.error('加载镜像失败:', err);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await instanceAPI.create({
        name,
        image,
        cpu,
        memory,
        disk,
        hostId: parseInt(hostId),
        ports: [],
        network_limit: 'unlimited',
        upload_limit: 'unlimited',
        download_limit: 'unlimited',
      });
      setShowForm(false);
      setName('');
      setImage('ubuntu/22.04');
      setCpu(2);
      setMemory(2048);
      setDisk(20);
      setHostId('');
      loadInstances();
    } catch (err) {
      console.error('创建实例失败:', err);
    }
  };

  const handleDelete = async (id) => {
    if (!confirm('确定要删除此实例吗？')) return;
    try {
      await instanceAPI.delete(id);
      loadInstances();
    } catch (err) {
      console.error('删除实例失败:', err);
    }
  };

  const handleStart = async (id) => {
    try {
      await instanceAPI.start(id);
      loadInstances();
    } catch (err) {
      console.error('启动实例失败:', err);
    }
  };

  const handleStop = async (id) => {
    try {
      await instanceAPI.stop(id);
      loadInstances();
    } catch (err) {
      console.error('停止实例失败:', err);
    }
  };

  const getStatusColor = (status) => {
    switch(status) {
      case 'running': return '#4caf50';
      case 'stopped': return '#f44336';
      case 'creating': return '#ff9800';
      default: return '#9e9e9e';
    }
  };

  const getStatusText = (status) => {
    switch(status) {
      case 'running': return '运行中';
      case 'stopped': return '已停止';
      case 'creating': return '创建中';
      case 'deleted': return '已删除';
      case 'paused': return '已暂停';
      default: return status?.toUpperCase() || '未知';
    }
  };

  return (
    <div style={{ padding: 20 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 20 }}>
        <h1>实例管理</h1>
        <button onClick={() => setShowForm(!showForm)} style={{ padding: '10px 20px', backgroundColor: '#4caf50', color: 'white', border: 'none', borderRadius: 4, cursor: 'pointer' }}>
          创建实例
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleSubmit} style={{ marginBottom: 20, padding: 20, border: '1px solid #ccc', borderRadius: 8, backgroundColor: '#f9f9f9' }}>
          <h3>新建实例</h3>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 15 }}>
            <div>
              <label>名称：</label>
              <input type="text" value={name} onChange={(e) => setName(e.target.value)} required style={{ width: '100%', padding: 8, border: '1px solid #ddd', borderRadius: 4 }} />
            </div>
            <div>
              <label>镜像：</label>
              <select value={image} onChange={(e) => setImage(e.target.value)} required style={{ width: '100%', padding: 8, border: '1px solid #ddd', borderRadius: 4 }}>
                {images.map(img => (
                  <option key={img} value={img}>{img}</option>
                ))}
              </select>
            </div>
            <div>
              <label>主机：</label>
              <select value={hostId} onChange={(e) => setHostId(e.target.value)} required style={{ width: '100%', padding: 8, border: '1px solid #ddd', borderRadius: 4 }}>
                <option value="">选择主机</option>
                {hosts.map(host => (
                  <option key={host.id} value={host.id}>{host.name}</option>
                ))}
              </select>
            </div>
            <div>
              <label>CPU 核心数：</label>
              <input type="number" value={cpu} onChange={(e) => setCpu(parseInt(e.target.value))} min="1" required style={{ width: '100%', padding: 8, border: '1px solid #ddd', borderRadius: 4 }} />
            </div>
            <div>
              <label>内存（MB）：</label>
              <input type="number" value={memory} onChange={(e) => setMemory(parseInt(e.target.value))} min="256" required style={{ width: '100%', padding: 8, border: '1px solid #ddd', borderRadius: 4 }} />
            </div>
            <div>
              <label>磁盘（GB）：</label>
              <input type="number" value={disk} onChange={(e) => setDisk(parseInt(e.target.value))} min="1" required style={{ width: '100%', padding: 8, border: '1px solid #ddd', borderRadius: 4 }} />
            </div>
          </div>
          <div style={{ marginTop: 15 }}>
            <button type="submit" style={{ padding: '10px 20px', backgroundColor: '#4caf50', color: 'white', border: 'none', borderRadius: 4, cursor: 'pointer' }}>创建实例</button>
            <button type="button" onClick={() => setShowForm(false)} style={{ marginLeft: 10, padding: '10px 20px', backgroundColor: '#9e9e9e', color: 'white', border: 'none', borderRadius: 4, cursor: 'pointer' }}>取消</button>
          </div>
        </form>
      )}

      <div style={{ display: 'grid', gap: 15 }}>
        {instances.map(instance => (
          <div key={instance.id} style={{ padding: 20, border: '1px solid #ddd', borderRadius: 8, backgroundColor: 'white', boxShadow: '0 2px 4px rgba(0,0,0,0.05)' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 15 }}>
              <h3 style={{ margin: 0 }}>{instance.name}</h3>
              <span style={{ 
                padding: '5px 12px', 
                borderRadius: 12, 
                backgroundColor: getStatusColor(instance.status),
                color: 'white',
                fontSize: 12,
                fontWeight: 'bold'
              }}>
                {getStatusText(instance.status)}
              </span>
            </div>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 10, marginBottom: 15 }}>
              <div><strong>镜像：</strong> {instance.image}</div>
              <div><strong>CPU：</strong> {instance.cpu} 核</div>
              <div><strong>内存：</strong> {instance.memory} MB</div>
              <div><strong>磁盘：</strong> {instance.disk} GB</div>
            </div>
            {instance.mapping_ip && (
              <div style={{ marginBottom: 15, padding: 10, backgroundColor: '#e3f2fd', borderRadius: 4 }}>
                <strong>映射 IP：</strong> {instance.mapping_ip}
              </div>
            )}
            <div style={{ display: 'flex', gap: 10 }}>
              <button onClick={() => handleStart(instance.id)} style={{ padding: '8px 16px', backgroundColor: '#4caf50', color: 'white', border: 'none', borderRadius: 4, cursor: 'pointer' }}>启动</button>
              <button onClick={() => handleStop(instance.id)} style={{ padding: '8px 16px', backgroundColor: '#f44336', color: 'white', border: 'none', borderRadius: 4, cursor: 'pointer' }}>停止</button>
              <button onClick={() => handleDelete(instance.id)} style={{ padding: '8px 16px', backgroundColor: '#ff9800', color: 'white', border: 'none', borderRadius: 4, cursor: 'pointer' }}>删除</button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default InstancesPage;
