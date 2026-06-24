import { useState, useEffect } from 'react';
import { hostAPI } from '../services/api';

const HostsPage = () => {
  const [hosts, setHosts] = useState([]);
  const [showForm, setShowForm] = useState(false);
  const [name, setName] = useState('');
  const [address, setAddress] = useState('');
  const [certificate, setCertificate] = useState('');

  useEffect(() => {
    loadHosts();
  }, []);

  const loadHosts = async () => {
    try {
      const response = await hostAPI.getAll();
      setHosts(response.data);
    } catch (err) {
      console.error('加载主机失败:', err);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await hostAPI.add(name, address, certificate);
      setShowForm(false);
      setName('');
      setAddress('');
      setCertificate('');
      loadHosts();
    } catch (err) {
      console.error('添加主机失败:', err);
    }
  };

  return (
    <div style={{ padding: 20 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 20 }}>
        <h1>主机管理</h1>
        <button 
          onClick={() => setShowForm(!showForm)}
          style={{ padding: '10px 20px', backgroundColor: '#4caf50', color: 'white', border: 'none', borderRadius: 4, cursor: 'pointer' }}
        >
          添加主机
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleSubmit} style={{ marginBottom: 20, padding: 20, border: '1px solid #ccc', borderRadius: 8, backgroundColor: '#f9f9f9' }}>
          <h3 style={{ marginTop: 0 }}>新建主机</h3>
          <div style={{ marginBottom: 10 }}>
            <label>名称：</label>
            <input type="text" value={name} onChange={(e) => setName(e.target.value)} required style={{ width: '100%', padding: 8, border: '1px solid #ddd', borderRadius: 4, boxSizing: 'border-box' }} />
          </div>
          <div style={{ marginBottom: 10 }}>
            <label>地址（例如：https://192.168.1.100:8443）：</label>
            <input type="text" value={address} onChange={(e) => setAddress(e.target.value)} required style={{ width: '100%', padding: 8, border: '1px solid #ddd', borderRadius: 4, boxSizing: 'border-box' }} />
          </div>
          <div style={{ marginBottom: 10 }}>
            <label>客户端证书（PEM 格式）：</label>
            <textarea value={certificate} onChange={(e) => setCertificate(e.target.value)} required rows={4} style={{ width: '100%', padding: 8, border: '1px solid #ddd', borderRadius: 4, boxSizing: 'border-box', fontFamily: 'monospace', fontSize: 12 }} />
          </div>
          <button type="submit" style={{ padding: '10px 20px', backgroundColor: '#4caf50', color: 'white', border: 'none', borderRadius: 4, cursor: 'pointer' }}>添加主机</button>
        </form>
      )}

      {hosts.length === 0 ? (
        <div style={{ textAlign: 'center', padding: 40, color: '#999' }}>
          <p style={{ fontSize: 48 }}>🖥️</p>
          <p>尚未添加任何主机。点击"添加主机"开始使用。</p>
        </div>
      ) : (
        <div style={{ display: 'grid', gap: 15 }}>
          {hosts.map(host => (
            <div key={host.id} style={{ padding: 20, border: '1px solid #ddd', borderRadius: 8, backgroundColor: 'white', boxShadow: '0 2px 4px rgba(0,0,0,0.05)' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 10 }}>
                <h3 style={{ margin: 0 }}>{host.name}</h3>
                <span style={{ 
                  padding: '4px 12px', 
                  borderRadius: 12, 
                  backgroundColor: host.status === 'active' ? '#4caf50' : '#9e9e9e',
                  color: 'white',
                  fontSize: 12,
                  fontWeight: 'bold'
                }}>
                  {host.status === 'active' ? '在线' : host.status.toUpperCase()}
                </span>
              </div>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 10 }}>
                <div><strong>地址：</strong> {host.address}</div>
                <div><strong>项目：</strong> {host.project}</div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default HostsPage;
