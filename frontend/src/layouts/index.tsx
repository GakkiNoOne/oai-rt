import React, { useState, useEffect } from 'react';
import { Outlet, useNavigate, useLocation } from 'umi';
import { 
  ApiOutlined,
  MenuFoldOutlined, 
  MenuUnfoldOutlined,
  UserOutlined,
  SettingOutlined,
  ControlOutlined,
  LogoutOutlined,
  FileTextOutlined
} from '@ant-design/icons';
import { Layout, Menu, Button, Avatar, Dropdown, Space, Typography, message, Modal } from 'antd';

const { Header, Sider, Content } = Layout;
const { Text } = Typography;

const AdminLayout: React.FC = () => {
  const [collapsed, setCollapsed] = useState(false);
  const [username, setUsername] = useState('管理员');
  const navigate = useNavigate();
  const location = useLocation();

  // 检查登录状态
  useEffect(() => {
    const token = localStorage.getItem('token');
    const savedUsername = localStorage.getItem('username');
    
    if (!token) {
      message.warning('请先登录');
      navigate('/login');
      return;
    }

    if (savedUsername) {
      setUsername(savedUsername);
    }
  }, [navigate]);

  const menuItems = [
    {
      key: '/rts',
      icon: <ApiOutlined />,
      label: 'RT管理',
    },
    {
      key: '/configs',
      icon: <ControlOutlined />,
      label: '配置管理',
    },
    {
      key: '/api-docs',
      icon: <FileTextOutlined />,
      label: 'API文档',
    },
  ];

  const userMenuItems = [
    {
      key: 'settings',
      icon: <SettingOutlined />,
      label: '设置',
    },
    {
      type: 'divider',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      danger: true,
    },
  ];

  const handleMenuClick = ({ key }: { key: string }) => {
    navigate(key);
  };

  const handleLogout = () => {
    // 清除本地存储的token和用户信息
    localStorage.removeItem('token');
    localStorage.removeItem('username');
    message.success('已退出登录');
    navigate('/login');
  };

  const handleUserMenuClick = ({ key }: { key: string }) => {
    if (key === 'settings') {
      console.log('打开设置');
    } else if (key === 'logout') {
      Modal.confirm({
        title: '确认退出',
        content: '确定要退出登录吗？',
        okText: '确定',
        cancelText: '取消',
        onOk: handleLogout,
      });
    }
  };

  return (
    <Layout style={{ minHeight: '100vh' }}>
      {/* 顶部导航栏 */}
      <Header 
        style={{ 
          background: '#fff', 
          padding: '0 24px', 
          borderBottom: '1px solid #f0f0f0',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between'
        }}
      >
        <div style={{ display: 'flex', alignItems: 'center' }}>
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => setCollapsed(!collapsed)}
            style={{ fontSize: '16px', width: 64, height: 64 }}
          />
          <div style={{ marginLeft: 16, fontSize: '18px', fontWeight: 'bold', color: '#1890ff' }}>
            RT管理系统
          </div>
        </div>
        
        <Dropdown
          menu={{ 
            items: userMenuItems,
            onClick: handleUserMenuClick 
          }}
          placement="bottomRight"
        >
          <Space style={{ cursor: 'pointer' }}>
            <Avatar icon={<UserOutlined />} size="small" />
            <Text>{username}</Text>
          </Space>
        </Dropdown>
      </Header>

      <Layout>
        {/* 左侧菜单 */}
        <Sider 
          trigger={null} 
          collapsible 
          collapsed={collapsed}
          style={{ 
            background: '#fff',
            borderRight: '1px solid #f0f0f0'
          }}
        >
          <Menu
            mode="inline"
            selectedKeys={[location.pathname]}
            items={menuItems}
            onClick={handleMenuClick}
            style={{ 
              height: '100%', 
              borderRight: 0,
              marginTop: 16
            }}
          />
        </Sider>

        {/* 主内容区域 */}
        <Layout>
          <Content
            style={{
              padding: '24px',
              background: '#f0f2f5',
              minHeight: 'calc(100vh - 64px)',
            }}
          >
            <Outlet />
          </Content>
        </Layout>
      </Layout>
    </Layout>
  );
};

export default AdminLayout;
