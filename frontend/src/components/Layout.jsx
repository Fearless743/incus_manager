import { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { Link } from 'react-router-dom';

const Layout = ({ children }) => {
  const { user, logout } = useAuth();
  const [sidebarOpen, setSidebarOpen] = useState(true);

  if (!user) {
    return null;
  }

  return (
    <div style={{ display: 'flex', minHeight: '100vh' }}>
      {/* Sidebar */}
      <div style={{ 
        width: sidebarOpen ? 250 : 60, 
        backgroundColor: '#1a1a2e', 
        color: 'white',
        transition: 'width 0.3s',
        overflow: 'hidden',
        position: 'fixed',
        height: '100vh',
        left: 0,
        top: 0
      }}>
        <div style={{ padding: 15, display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          {sidebarOpen && <h2 style={{ margin: 0, fontSize: 18 }}>Incus Manager</h2>}
          <button 
            onClick={() => setSidebarOpen(!sidebarOpen)}
            style={{ background: 'none', border: 'none', color: 'white', cursor: 'pointer', fontSize: 16 }}
          >
            {sidebarOpen ? '←' : '→'}
          </button>
        </div>
        
        <nav style={{ padding: 10 }}>
          <Link to="/dashboard" style={{ display: 'block', padding: 10, color: 'white', textDecoration: 'none', borderRadius: 4 }}>
            {sidebarOpen ? '📊 Dashboard' : '📊'}
          </Link>
          <Link to="/hosts" style={{ display: 'block', padding: 10, color: 'white', textDecoration: 'none', borderRadius: 4 }}>
            {sidebarOpen ? '🖥️ Hosts' : '🖥️'}
          </Link>
          <Link to="/instances" style={{ display: 'block', padding: 10, color: 'white', textDecoration: 'none', borderRadius: 4 }}>
            {sidebarOpen ? '📦 Instances' : '📦'}
          </Link>
          <Link to="/shared" style={{ display: 'block', padding: 10, color: 'white', textDecoration: 'none', borderRadius: 4 }}>
            {sidebarOpen ? '🔗 Sharing' : '🔗'}
          </Link>
        </nav>

        <div style={{ position: 'absolute', bottom: 20, left: 10, right: 10 }}>
          {sidebarOpen && <p style={{ margin: '0 0 10px 0', fontSize: 12, opacity: 0.7 }}>{user.username}</p>}
          <button 
            onClick={logout}
            style={{ 
              width: '100%', 
              padding: '8px 0', 
              backgroundColor: '#f44336', 
              color: 'white', 
              border: 'none', 
              borderRadius: 4, 
              cursor: 'pointer',
              fontSize: 14
            }}
          >
            {sidebarOpen ? 'Logout' : '⏻'}
          </button>
        </div>
      </div>

      {/* Main Content */}
      <div style={{ flex: 1, marginLeft: sidebarOpen ? 250 : 60, padding: 20, transition: 'margin-left 0.3s' }}>
        {children}
      </div>
    </div>
  );
};

export default Layout;
