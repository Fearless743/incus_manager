import { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { authAPI } from '../services/api';
import { useNavigate } from 'react-router-dom';

const LoginPage = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const response = await authAPI.login(username, password);
      login(response.data.user, response.data.token);
      navigate('/dashboard');
    } catch (err) {
      setError(err.response?.data?.error || '登录失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ 
      maxWidth: 400, 
      margin: '100px auto', 
      padding: 30, 
      borderRadius: 8, 
      boxShadow: '0 2px 10px rgba(0,0,0,0.1)',
      backgroundColor: 'white'
    }}>
      <h1 style={{ textAlign: 'center', marginBottom: 30 }}>Incus 管理器</h1>
      {error && (
        <div style={{ 
          padding: 10, 
          backgroundColor: '#ffebee', 
          color: '#c62828', 
          borderRadius: 4,
          marginBottom: 20
        }}>
          {error}
        </div>
      )}
      <form onSubmit={handleSubmit}>
        <div style={{ marginBottom: 15 }}>
          <label style={{ display: 'block', marginBottom: 5, fontWeight: 'bold' }}>用户名</label>
          <input
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
            style={{ 
              width: '100%', 
              padding: 10, 
              border: '1px solid #ddd', 
              borderRadius: 4,
              boxSizing: 'border-box'
            }}
          />
        </div>
        <div style={{ marginBottom: 20 }}>
          <label style={{ display: 'block', marginBottom: 5, fontWeight: 'bold' }}>密码</label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            style={{ 
              width: '100%', 
              padding: 10, 
              border: '1px solid #ddd', 
              borderRadius: 4,
              boxSizing: 'border-box'
            }}
          />
        </div>
        <button 
          type="submit" 
          disabled={loading}
          style={{ 
            width: '100%', 
            padding: 12, 
            backgroundColor: '#007bff', 
            color: 'white', 
            border: 'none',
            borderRadius: 4,
            cursor: loading ? 'not-allowed' : 'pointer',
            fontSize: 16
          }}
        >
          {loading ? '登录中...' : '登录'}
        </button>
      </form>
    </div>
  );
};

export default LoginPage;
