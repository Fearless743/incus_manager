import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import { Navigate } from 'react-router';
import { AuthProvider } from './context/AuthProvider';
import { useAuth } from './context/AuthContext';
import LoginPage from './pages/LoginPage';
import DashboardPage from './pages/DashboardPage';
import HostsPage from './pages/HostsPage';
import InstancesPage from './pages/InstancesPage';
import SharedPage from './pages/SharedPage';
import Layout from './components/Layout';

const ProtectedRoute = ({ children }) => {
  const { token } = useAuth();
  return token ? children : <Navigate to="/login" />;
};

const router = createBrowserRouter([
  {
    path: '/login',
    element: <LoginPage />,
  },
  {
    path: '/dashboard',
    element: (
      <ProtectedRoute>
        <Layout><DashboardPage /></Layout>
      </ProtectedRoute>
    ),
  },
  {
    path: '/hosts',
    element: (
      <ProtectedRoute>
        <Layout><HostsPage /></Layout>
      </ProtectedRoute>
    ),
  },
  {
    path: '/instances',
    element: (
      <ProtectedRoute>
        <Layout><InstancesPage /></Layout>
      </ProtectedRoute>
    ),
  },
  {
    path: '/shared',
    element: (
      <ProtectedRoute>
        <Layout><SharedPage /></Layout>
      </ProtectedRoute>
    ),
  },
  {
    path: '/',
    element: <Navigate to="/dashboard" />,
  },
]);

function App() {
  return (
    <AuthProvider>
      <RouterProvider router={router} />
    </AuthProvider>
  );
}

export default App;
