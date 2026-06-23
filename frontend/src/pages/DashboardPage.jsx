import { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import { statsAPI } from '../services/api';

const DashboardPage = () => {
  const { user } = useAuth();
  const [stats, setStats] = useState({
    total_hosts: 0,
    total_instances: 0,
    running_instances: 0,
    shared_instances: 0,
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadStats();
  }, []);

  const loadStats = async () => {
    try {
      const response = await statsAPI.get();
      setStats(response.data);
    } catch (err) {
      console.error('Failed to load stats:', err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div style={{ padding: 20 }}>Loading...</div>;
  }

  return (
    <div style={{ padding: 20 }}>
      <h1>Welcome, {user?.username}!</h1>
      <p style={{ color: '#666' }}>Here is an overview of your Incus management panel.</p>
      
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 20, marginTop: 30 }}>
        <StatCard title="Total Hosts" value={stats.total_hosts} icon="🖥️" color="#2196f3" />
        <StatCard title="Total Instances" value={stats.total_instances} icon="📦" color="#4caf50" />
        <StatCard title="Running" value={stats.running_instances} icon="⚡" color="#ff9800" />
        <StatCard title="Shared" value={stats.shared_instances} icon="🔗" color="#9c27b0" />
      </div>

      <div style={{ marginTop: 40 }}>
        <h2>Quick Actions</h2>
        <div style={{ display: 'flex', gap: 15, marginTop: 15 }}>
          <a href="/hosts" style={{ padding: '15px 25px', backgroundColor: '#2196f3', color: 'white', textDecoration: 'none', borderRadius: 8, boxShadow: '0 2px 4px rgba(0,0,0,0.1)' }}>
            🖥️ Manage Hosts
          </a>
          <a href="/instances" style={{ padding: '15px 25px', backgroundColor: '#4caf50', color: 'white', textDecoration: 'none', borderRadius: 8, boxShadow: '0 2px 4px rgba(0,0,0,0.1)' }}>
            📦 Create Instance
          </a>
          <a href="/shared" style={{ padding: '15px 25px', backgroundColor: '#9c27b0', color: 'white', textDecoration: 'none', borderRadius: 8, boxShadow: '0 2px 4px rgba(0,0,0,0.1)' }}>
            🔗 Share Instances
          </a>
        </div>
      </div>

      <div style={{ marginTop: 40, padding: 20, backgroundColor: '#f5f5f5', borderRadius: 8 }}>
        <h3>System Information</h3>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 15, marginTop: 15 }}>
          <div><strong>User:</strong> {user?.username}</div>
          <div><strong>Email:</strong> {user?.email}</div>
          <div><strong>Role:</strong> {user?.role}</div>
          <div><strong>Server Time:</strong> {new Date().toLocaleString()}</div>
        </div>
      </div>
    </div>
  );
};

const StatCard = ({ title, value, icon, color }) => (
  <div style={{
    padding: 25,
    backgroundColor: 'white',
    borderRadius: 12,
    boxShadow: '0 2px 8px rgba(0,0,0,0.08)',
    textAlign: 'center',
    borderLeft: `4px solid ${color}`
  }}>
    <div style={{ fontSize: 32, marginBottom: 10 }}>{icon}</div>
    <h3 style={{ margin: 0, color: '#666', fontSize: 14, textTransform: 'uppercase' }}>{title}</h3>
    <p style={{ fontSize: 36, fontWeight: 'bold', margin: '10px 0', color }}>{value}</p>
  </div>
);

export default DashboardPage;
