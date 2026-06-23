import { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import LoginPage from './LoginPage';

const Layout = ({ children }) => {
  const { user, logout } = useAuth();
  const [sidebarOpen, setSidebarOpen] = useState(true);

  if (!user) {
    return <LoginPage />;
  }

  return (
    <div style={{ display: 'flex', minHeight: '100vh' }}>
      {/* Sidebar */}
      <div style={{ 
        width: sidebarOpen ? 250 : 60, 
        backgroundColor: '#1a1a2e', 
        color: 'white',
        transition: 'width 0.3s',
        overflow: 'hidden'
      }}>
        <div style={{ padding: 15, display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          {sidebarOpen && <h2>Incus Manager</h2>}
          <button onClick={() => setSidebarOpen(!sidebarOpen)}>
            {sidebarOpen ? '←' : '→'}
          </button>
        </div>
        
        <nav style={{ padding: 10 }}>
          <a href="/dashboard" style={{ display: 'block', padding: 10, color: 'white', textDecoration: 'none' }}>
            {sidebarOpen ? 'Dashboard' : '📊'}
          </a>
          <a href="/hosts" style={{ display: 'block', padding: 10, color: 'white', textDecoration: 'none' }}>
            {sidebarOpen ? 'Hosts' : '🖥️'}
          </a>
          <a href="/instances" style={{ display: 'block', padding: 10, color: 'white', textDecoration: 'none' }}>
            {sidebarOpen ? 'Instances' : '📦'}
          </a>
        </nav>

        <div style={{ position: 'absolute', bottom: 20, left: 10 }}>
          <p>{sidebarOpen && user.username}</p>
          <button onClick={logout}>Logout</button>
        </div>
      </div>

      {/* Main Content */}
      <div style={{ flex: 1, padding: 20 }}>
        {children}
      </div>
    </div>
  );
};

export default Layout;
