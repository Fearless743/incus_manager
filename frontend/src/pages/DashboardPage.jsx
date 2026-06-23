import { useState, useEffect } from 'react';
import { hostAPI, instanceAPI } from '../services/api';

const DashboardPage = () => {
  const [hosts, setHosts] = useState([]);
  const [instances, setInstances] = useState([]);

  useEffect(() => {
    loadHosts();
    loadInstances();
  }, []);

  const loadHosts = async () => {
    try {
      const response = await hostAPI.getAll();
      setHosts(response.data);
    } catch (err) {
      console.error('Failed to load hosts:', err);
    }
  };

  const loadInstances = async () => {
    try {
      const response = await instanceAPI.getAll();
      setInstances(response.data);
    } catch (err) {
      console.error('Failed to load instances:', err);
    }
  };

  return (
    <div style={{ padding: 20 }}>
      <h1>Dashboard</h1>
      
      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 20 }}>
        <div>
          <h2>Hosts ({hosts.length})</h2>
          <ul>
            {hosts.map(host => (
              <li key={host.id}>{host.name} - {host.address}</li>
            ))}
          </ul>
        </div>

        <div>
          <h2>Instances ({instances.length})</h2>
          <ul>
            {instances.map(instance => (
              <li key={instance.id}>
                {instance.name} - {instance.status} (CPU: {instance.cpu}, Memory: {instance.memory}MB)
              </li>
            ))}
          </ul>
        </div>
      </div>
    </div>
  );
};

export default DashboardPage;
