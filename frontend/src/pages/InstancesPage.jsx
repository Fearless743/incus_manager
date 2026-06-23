import { useState, useEffect } from 'react';
import { instanceAPI } from '../services/api';

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

  useEffect(() => {
    loadInstances();
    loadHosts();
  }, []);

  const loadInstances = async () => {
    try {
      const response = await instanceAPI.getAll();
      setInstances(response.data);
    } catch (err) {
      console.error('Failed to load instances:', err);
    }
  };

  const loadHosts = async () => {
    // TODO: Implement host API call
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
      });
      setShowForm(false);
      loadInstances();
    } catch (err) {
      console.error('Failed to create instance:', err);
    }
  };

  return (
    <div style={{ padding: 20 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 20 }}>
        <h1>Instances</h1>
        <button onClick={() => setShowForm(!showForm)}>Create Instance</button>
      </div>

      {showForm && (
        <form onSubmit={handleSubmit} style={{ marginBottom: 20, padding: 20, border: '1px solid #ccc' }}>
          <div style={{ marginBottom: 10 }}>
            <label>Name:</label>
            <input type="text" value={name} onChange={(e) => setName(e.target.value)} required />
          </div>
          <div style={{ marginBottom: 10 }}>
            <label>Image:</label>
            <input type="text" value={image} onChange={(e) => setImage(e.target.value)} required />
          </div>
          <div style={{ marginBottom: 10 }}>
            <label>CPU Cores:</label>
            <input type="number" value={cpu} onChange={(e) => setCpu(parseInt(e.target.value))} required />
          </div>
          <div style={{ marginBottom: 10 }}>
            <label>Memory (MB):</label>
            <input type="number" value={memory} onChange={(e) => setMemory(parseInt(e.target.value))} required />
          </div>
          <div style={{ marginBottom: 10 }}>
            <label>Disk (GB):</label>
            <input type="number" value={disk} onChange={(e) => setDisk(parseInt(e.target.value))} required />
          </div>
          <div style={{ marginBottom: 10 }}>
            <label>Host:</label>
            <select value={hostId} onChange={(e) => setHostId(e.target.value)} required>
              <option value="">Select host</option>
              {hosts.map(host => (
                <option key={host.id} value={host.id}>{host.name}</option>
              ))}
            </select>
          </div>
          <button type="submit">Create</button>
        </form>
      )}

      <div style={{ display: 'grid', gap: 10 }}>
        {instances.map(instance => (
          <div key={instance.id} style={{ padding: 15, border: '1px solid #ddd', borderRadius: 5 }}>
            <div style={{ display: 'flex', justifyContent: 'space-between' }}>
              <h3>{instance.name}</h3>
              <span style={{ 
                padding: '5px 10px', 
                borderRadius: 3, 
                backgroundColor: instance.status === 'running' ? '#4caf50' : '#f44336',
                color: 'white'
              }}>
                {instance.status}
              </span>
            </div>
            <p>Image: {instance.image}</p>
            <p>Resources: {instance.cpu} CPU, {instance.memory}MB RAM, {instance.disk}GB Disk</p>
            <div style={{ marginTop: 10 }}>
              <button style={{ marginRight: 5 }}>Start</button>
              <button style={{ marginRight: 5 }}>Stop</button>
              <button>Delete</button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default InstancesPage;
