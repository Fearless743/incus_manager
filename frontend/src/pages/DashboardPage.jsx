import { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext';

const DashboardPage = () => {
  const { user } = useAuth();
  const [stats, setStats] = useState({
    totalHosts: 0,
    totalInstances: 0,
    runningInstances: 0,
    sharedInstances: 0,
  });

  useEffect(() => {
    loadStats();
  }, []);

  const loadStats = async () => {
    // TODO: Implement stats API call
    setStats({
      totalHosts: 0,
      totalInstances: 0,
      runningInstances: 0,
      sharedInstances: 0,
    });
  };

  return (
    <div style={{ padding: 20 }}>
      <h1>Welcome, {user?.username}!</h1>
      
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 20, marginTop: 30 }}>
        <StatCard title="Total Hosts" value={stats.totalHosts} />
        <StatCard title="Total Instances" value={stats.totalInstances} />
        <StatCard title="Running" value={stats.runningInstances} />
        <StatCard title="Shared" value={stats.sharedInstances} />
      </div>

      <div style={{ marginTop: 40 }}>
        <h2>Quick Actions</h2>
        <div style={{ display: 'flex', gap: 10 }}>
          <a href="/hosts" style={{ padding: '10px 20px', backgroundColor: '#007bff', color: 'white', textDecoration: 'none', borderRadius: 5 }}>
            Manage Hosts
          </a>
          <a href="/instances" style={{ padding: '10px 20px', backgroundColor: '#28a745', color: 'white', textDecoration: 'none', borderRadius: 5 }}>
            Create Instance
          </a>
          <a href="/shared" style={{ padding: '10px 20px', backgroundColor: '#ffc107', color: 'black', textDecoration: 'none', borderRadius: 5 }}>
            Share Instances
          </a>
        </div>
      </div>
    </div>
  );
};

const StatCard = ({ title, value }) => (
  <div style={{
    padding: 20,
    backgroundColor: 'white',
    borderRadius: 8,
    boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
    textAlign: 'center'
  }}>
    <h3 style={{ margin: 0, color: '#666' }}>{title}</h3>
    <p style={{ fontSize: 32, fontWeight: 'bold', margin: '10px 0' }}>{value}</p>
  </div>
);

export default DashboardPage;
