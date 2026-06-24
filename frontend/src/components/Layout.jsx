import { useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { Layout as AntLayout, Menu, Button, Dropdown, Switch, Space } from 'antd';
import {
  DashboardOutlined,
  CloudServerOutlined,
  AppstoreOutlined,
  ShareAltOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  UserOutlined,
  LogoutOutlined,
  BulbOutlined,
} from '@ant-design/icons';
import { useAuth } from '../context/AuthContext';
import { useTheme } from '../theme/ThemeContext';

const { Header, Sider, Content } = AntLayout;

const menuItems = [
  { key: '/dashboard', icon: <DashboardOutlined />, label: <Link to="/dashboard">仪表盘</Link> },
  { key: '/hosts', icon: <CloudServerOutlined />, label: <Link to="/hosts">主机</Link> },
  { key: '/instances', icon: <AppstoreOutlined />, label: <Link to="/instances">实例</Link> },
  { key: '/shared', icon: <ShareAltOutlined />, label: <Link to="/shared">共享</Link> },
];

const Layout = ({ children }) => {
  const { user, logout } = useAuth();
  const { themeMode, toggleTheme } = useTheme();
  const location = useLocation();
  const [collapsed, setCollapsed] = useState(false);

  if (!user) {
    return null;
  }

  const userMenuItems = [
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: logout,
    },
  ];

  return (
    <AntLayout style={{ minHeight: '100vh' }}>
      <Sider
        collapsible
        collapsed={collapsed}
        onCollapse={setCollapsed}
        trigger={null}
        width={220}
        theme={themeMode === 'dark' ? 'dark' : 'light'}
      >
        <div
          style={{
            height: 64,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            fontWeight: 600,
            fontSize: collapsed ? 14 : 16,
            overflow: 'hidden',
            whiteSpace: 'nowrap',
          }}
        >
          {collapsed ? 'Incus' : 'Incus 管理器'}
        </div>
        <Menu
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          theme={themeMode === 'dark' ? 'dark' : 'light'}
        />
      </Sider>
      <AntLayout>
        <Header
          style={{
            padding: '0 24px',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            borderBottom: '1px solid rgba(5, 5, 5, 0.06)',
            backgroundColor: themeMode === 'dark' ? '#001529' : '#ffffff',
          }}
          theme={themeMode === 'dark' ? 'dark' : 'light'}
        >
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => setCollapsed(!collapsed)}
          />
          <Space size="middle">
            <Space size="small">
              <BulbOutlined />
              <Switch
                checked={themeMode === 'dark'}
                onChange={toggleTheme}
                checkedChildren="暗"
                unCheckedChildren="亮"
              />
            </Space>
            <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
              <Button type="text" icon={<UserOutlined />}>
                {user.username}
              </Button>
            </Dropdown>
          </Space>
        </Header>
        <Content style={{ padding: 24 }}>
          {children}
        </Content>
      </AntLayout>
    </AntLayout>
  );
};

export default Layout;
