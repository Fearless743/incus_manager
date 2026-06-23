import { useState } from 'react';
import { hostAPI } from '../services/api';

const HostsPage = () => {
  const [hosts, setHosts] = useState([]);
  const [showForm, setShowForm] = useState(false);
  const [name, setName] = useState('');
  const [address, setAddress] = useState('');
  const [certificate, setCertificate] = useState('');

  const loadHosts = async () => {
    try {
      const response = await hostAPI.getAll();
      setHosts(response.data);
    } catch (err) {
      console.error('Failed to load hosts:', err);
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
      console.error('Failed to add host:', err);
    }
  };

  return (
    <div style={{ padding: 20 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 20 }}>
        <h1>Hosts</h1>
        <button onClick={() => setShowForm(!showForm)}>Add Host</button>
      </div>

      {showForm && (
        <form onSubmit={handleSubmit} style={{ marginBottom: 20, padding: 20, border: '1px solid #ccc' }}>
          <div style={{ marginBottom: 10 }}>
            <label>Name:</label>
            <input type="text" value={name} onChange={(e) => setName(e.target.value)} required />
          </div>
          <div style={{ marginBottom: 10 }}>
            <label>Address:</label>
            <input type="text" value={address} onChange={(e) => setAddress(e.target.value)} required />
          </div>
          <div style={{ marginBottom: 10 }}>
            <label>Certificate:</label>
            <textarea value={certificate} onChange={(e) => setCertificate(e.target.value)} required />
          </div>
          <button type="submit">Save</button>
        </form>
      )}

      <div style={{ display: 'grid', gap: 10 }}>
        {hosts.map(host => (
          <div key={host.id} style={{ padding: 15, border: '1px solid #ddd', borderRadius: 5 }}>
            <h3>{host.name}</h3>
            <p>Address: {host.address}</p>
            <p>Status: {host.status}</p>
            <p>Project: {host.project}</p>
          </div>
        ))}
      </div>
    </div>
  );
};

export default HostsPage;
