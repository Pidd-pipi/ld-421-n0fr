import { useMemo } from 'react'
import { Layout, Menu, Dropdown, Avatar, Space, Typography } from 'antd'
import {
  DashboardOutlined,
  ExperimentOutlined,
  SwapOutlined,
  CalendarOutlined,
  ToolOutlined,
  AuditOutlined,
  LogoutOutlined,
  UserOutlined
} from '@ant-design/icons'
import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import { useAuthStore } from '../../stores/authStore'

const { Header, Sider, Content } = Layout

interface MenuItem {
  key: string
  icon: JSX.Element
  label: string
  roles?: string[]
}

const MENU_ITEMS: MenuItem[] = [
  { key: '/dashboard', icon: <DashboardOutlined />, label: '设备总览' },
  { key: '/equipment', icon: <ExperimentOutlined />, label: '设备管理' },
  { key: '/borrow', icon: <SwapOutlined />, label: '借用管理' },
  { key: '/reservations', icon: <CalendarOutlined />, label: '预约管理' },
  { key: '/maintenance', icon: <ToolOutlined />, label: '维护管理' },
  { key: '/audit-logs', icon: <AuditOutlined />, label: '操作日志', roles: ['Admin'] }
]

export function AppLayout() {
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout } = useAuthStore()

  const items = useMemo(() => {
    const role = user?.roleCode || ''
    return MENU_ITEMS.filter((item) => !item.roles || item.roles.includes(role)).map((item) => ({
      key: item.key,
      icon: item.icon,
      label: item.label
    }))
  }, [user])

  const selectedKey = '/' + (location.pathname.split('/')[1] || 'dashboard')

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider theme="dark" width={210} collapsible>
        <div
          style={{
            height: 56,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            color: '#fff',
            fontWeight: 600,
            letterSpacing: 1
          }}
        >
          实验室设备管理
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[selectedKey]}
          items={items}
          onClick={({ key }) => navigate(key)}
        />
      </Sider>
      <Layout>
        <Header
          style={{
            background: '#fff',
            padding: '0 24px',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            boxShadow: '0 1px 4px rgba(0,21,41,.08)'
          }}
        >
          <Typography.Title level={4} style={{ margin: 0 }}>
            实验室设备管理系统
          </Typography.Title>
          <Dropdown
            menu={{
              items: [
                {
                  key: 'logout',
                  icon: <LogoutOutlined />,
                  label: '退出登录',
                  onClick: () => {
                    logout()
                    navigate('/login')
                  }
                }
              ]
            }}
          >
            <Space style={{ cursor: 'pointer' }}>
              <Avatar size="small" icon={<UserOutlined />} />
              <span>{user?.name || user?.username || '未登录'}</span>
              <span style={{ color: '#999' }}>{user?.roleName || ''}</span>
            </Space>
          </Dropdown>
        </Header>
        <Content style={{ margin: 16 }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}

export default AppLayout
