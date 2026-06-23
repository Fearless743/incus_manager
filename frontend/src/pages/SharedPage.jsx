import { useState, useEffect } from 'react';
import { instanceAPI, shareAPI } from '../services/api';

const SharedPage = () => {
  const [sharedInstances, setSharedInstances] = useState([]);
  const [myInstances, setMyInstances] = useState([]);
  const [selectedInstance, setSelectedInstance] = useState('');
  const [sharedWithUser, setSharedWithUser] = useState('');
  const [expiresAt, setExpiresAt] = useState('');

  useEffect(() => {
    loadSharedInstances();
    loadMyInstances();
  }, []);

  const loadSharedInstances = async () => {
    try {
      const response = await instanceAPI.getAll();
      setSharedInstances(response.data.filter(inst => inst.shared_with && inst.shared_with.length > 0));
    } catch (err) {
      console.error('Failed to load shared instances:', err);
    }
  };

  const loadMyInstances = async () => {
    try {
      const response = await instanceAPI.getAll();
      setMyInstances(response.data);
    } catch (err) {
      console.error('Failed to load instances:', err);
    }
  };

  const handleShare = async (e) => {
    e.preventDefault();
    try {
      await shareAPI.share(parseInt(selectedInstance), parseInt(sharedWithUser), expiresAt);
      loadSharedInstances();
    } catch (err) {
      console.error('Failed to share instance:', err);
    }
  };

  return (
    <div style={{ padding: 20 }}>
      <h1>Share Instances</h1>
      
      <form onSubmit={handleShare} style={{ marginBottom: 30, padding: 20, border: '1px solid #ccc' }}>
        <h3>Share an Instance</h3>
        <div style={{ marginBottom: 10 }}>
          <label>Select Instance:</label>
          <select value={selectedInstance} onChange={(e) => setSelectedInstance(e.target.value)} required>
            <option value="">Select instance</option>
            {myInstances.map(inst => (
              <option key={inst.id} value={inst.id}>{inst.name}</option>
            ))}
          </select>
        </div>
        <div style={{ marginBottom: 10 }}>
          <label>Share with User ID:</label>
          <input 
            type="number" 
            value={sharedWithUser} 
            onChange={(e) => setSharedWithUser(e.target.value)} 
            required 
          />
        </div>
        <div style={{ marginBottom: 10 }}>
          <label>Expires At:</label>
          <input 
            type="datetime-local" 
            value={expiresAt} 
            onChange={(e) => setExpiresAt(e.target.value)} 
            required 
          />
        </div>
        <button type="submit">Share</button>
      </form>

      <h2>Shared Instances</h2>
      <div style={{ display: 'grid', gap: 10 }}>
        {sharedInstances.map(instance => (
          <div key={instance.id} style={{ padding: 15, border: '1px solid #ddd', borderRadius: 5 }}>
            <h3>{instance.name}</h3>
            <p>Shared with: {instance.shared_with?.join(', ')}</p>
            <p>Expires: {new Date(instance.expiry_date).toLocaleString()}</p>
          </div>
        ))}
      </div>
    </div>
  );
};

export default SharedPage;
