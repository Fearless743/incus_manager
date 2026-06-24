import { useState, useEffect } from 'react';
import { instanceAPI, shareAPI } from '../services/api';

const SharedPage = () => {
  const [myInstances, setMyInstances] = useState([]);
  const [selectedInstance, setSelectedInstance] = useState('');
  const [sharedWithUser, setSharedWithUser] = useState('');
  const [expiresAt, setExpiresAt] = useState('');
  const [sharedInstances, setSharedInstances] = useState([]);

  useEffect(() => {
    loadInstances();
    loadSharedInstances();
  }, []);

  const loadInstances = async () => {
    try {
      const response = await instanceAPI.getAll();
      setMyInstances(response.data);
    } catch (err) {
      console.error('加载实例失败:', err);
    }
  };

  const loadSharedInstances = async () => {
    try {
      const response = await instanceAPI.getAll();
      const shared = response.data.filter(inst => inst.shared_with && inst.shared_with.length > 0);
      setSharedInstances(shared);
    } catch (err) {
      console.error('加载共享实例失败:', err);
    }
  };

  const handleShare = async (e) => {
    e.preventDefault();
    try {
      const expiresDateTime = new Date(expiresAt).toISOString();
      await shareAPI.share(parseInt(selectedInstance), parseInt(sharedWithUser), expiresDateTime);
      setSelectedInstance('');
      setSharedWithUser('');
      setExpiresAt('');
      loadInstances();
      loadSharedInstances();
    } catch (err) {
      console.error('共享实例失败:', err);
    }
  };

  const handleRevoke = async (instanceId, userId) => {
    if (!confirm('确定要撤销对此用户的共享吗？')) return;
    try {
      await shareAPI.revoke(instanceId, userId);
      loadInstances();
      loadSharedInstances();
    } catch (err) {
      console.error('撤销共享失败:', err);
    }
  };

  const getExpiryStatus = (expiryDate) => {
    if (!expiryDate) return '永不过期';
    const now = new Date();
    const expiry = new Date(expiryDate);
    if (expiry < now) return '已过期';
    const daysLeft = Math.ceil((expiry - now) / (1000 * 60 * 60 * 24));
    return `剩余 ${daysLeft} 天`;
  };

  return (
    <div style={{ padding: 20 }}>
      <h1>实例共享</h1>
      
      <form onSubmit={handleShare} style={{ marginBottom: 30, padding: 20, border: '1px solid #ccc', borderRadius: 8, backgroundColor: '#f9f9f9' }}>
        <h3>共享实例</h3>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: 15 }}>
          <div>
            <label>选择实例：</label>
            <select value={selectedInstance} onChange={(e) => setSelectedInstance(e.target.value)} required style={{ width: '100%', padding: 8, border: '1px solid #ddd', borderRadius: 4, marginTop: 5 }}>
              <option value="">选择实例</option>
              {myInstances.map(inst => (
                <option key={inst.id} value={inst.id}>{inst.name}</option>
              ))}
            </select>
          </div>
          <div>
            <label>共享给用户 ID：</label>
            <input 
              type="number" 
              value={sharedWithUser} 
              onChange={(e) => setSharedWithUser(e.target.value)} 
              placeholder="输入用户 ID"
              required 
              style={{ width: '100%', padding: 8, border: '1px solid #ddd', borderRadius: 4, marginTop: 5 }}
            />
          </div>
          <div>
            <label>过期时间：</label>
            <input 
              type="datetime-local" 
              value={expiresAt} 
              onChange={(e) => setExpiresAt(e.target.value)} 
              required 
              style={{ width: '100%', padding: 8, border: '1px solid #ddd', borderRadius: 4, marginTop: 5 }}
            />
          </div>
        </div>
        <div style={{ marginTop: 15 }}>
          <button type="submit" style={{ padding: '10px 20px', backgroundColor: '#2196f3', color: 'white', border: 'none', borderRadius: 4, cursor: 'pointer' }}>共享实例</button>
        </div>
      </form>

      <h3>当前共享的实例</h3>
      <div style={{ display: 'grid', gap: 15 }}>
        {sharedInstances.map(instance => (
          <div key={instance.id} style={{ padding: 20, border: '1px solid #ddd', borderRadius: 8, backgroundColor: 'white' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 10 }}>
              <h4 style={{ margin: 0 }}>{instance.name}</h4>
              <span style={{ 
                padding: '4px 10px', 
                borderRadius: 12, 
                backgroundColor: instance.status === 'running' ? '#4caf50' : '#f44336',
                color: 'white',
                fontSize: 12
              }}>
                {instance.status === 'running' ? '运行中' : '已停止'}
              </span>
            </div>
            <div style={{ display: 'flex', gap: 20, marginBottom: 10 }}>
              <div><strong>共享给：</strong> {instance.shared_with?.join(', ') || '无'}</div>
              <div><strong>过期时间：</strong> {getExpiryStatus(instance.expiry_date)}</div>
            </div>
            {instance.shared_with && instance.shared_with.map(userId => (
              <button 
                key={userId}
                onClick={() => handleRevoke(instance.id, userId)}
                style={{ padding: '5px 10px', backgroundColor: '#ff5722', color: 'white', border: 'none', borderRadius: 4, cursor: 'pointer', marginRight: 5 }}
              >
                撤销用户 #{userId} 的权限
              </button>
            ))}
          </div>
        ))}
      </div>
    </div>
  );
};

export default SharedPage;
