import { Layout, Menu, Dropdown, Space, Avatar, Typography, Tag } from 'antd'
import {
  HomeOutlined, TeamOutlined, ShopOutlined, StarOutlined, UserOutlined,
  SettingOutlined, SafetyCertificateOutlined, LogoutOutlined, FlagOutlined,
} from '@ant-design/icons'
import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import { useAuthStore } from './stores/authStore'
import { isAdmin } from './utils/auth'

const { Header, Sider, Content } = Layout

export default function App() {
  const navigate = useNavigate()
  const location = useLocation()
  const user = useAuthStore((s) => s.user)
  const logout = useAuthStore((s) => s.logout)

  const admin = isAdmin()

  const menuItems = [
    { key: '/', icon: <HomeOutlined />, label: '活动线路' },
    { key: '/teams', icon: <TeamOutlined />, label: '我的团队' },
    { key: '/mall', icon: <ShopOutlined />, label: '积分商城' },
    { key: '/favorites', icon: <StarOutlined />, label: '我的收藏' },
    { key: '/profile', icon: <UserOutlined />, label: '个人中心' },
    ...(admin
      ? [
          { key: 'admin', type: 'group' as const, label: '管理后台' },
          { key: '/admin/activities', icon: <FlagOutlined />, label: '活动管理' },
          { key: '/admin/products', icon: <ShopOutlined />, label: '商品管理' },
          { key: '/admin/users', icon: <UserOutlined />, label: '用户管理' },
          { key: '/admin/audit', icon: <SafetyCertificateOutlined />, label: '审计日志' },
        ]
      : []),
  ]

  const selectedKey = (() => {
    if (location.pathname.startsWith('/admin')) return location.pathname
    if (location.pathname.startsWith('/activities') || location.pathname.startsWith('/leaderboard') || location.pathname.startsWith('/checkin')) return '/'
    return location.pathname
  })()

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider breakpoint="lg" collapsedWidth={0} theme="dark">
        <div style={{ height: 56, display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#fff', fontWeight: 600, fontSize: 16 }}>
          🧭 城市定向越野
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[selectedKey]}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
        />
      </Sider>
      <Layout>
        <Header style={{ background: '#fff', padding: '0 24px', display: 'flex', justifyContent: 'flex-end', alignItems: 'center', boxShadow: '0 1px 4px rgba(0,0,0,.08)' }}>
          {user ? (
            <Dropdown
              menu={{
                items: [
                  { key: 'profile', icon: <UserOutlined />, label: '个人中心', onClick: () => navigate('/profile') },
                  { key: 'logout', icon: <LogoutOutlined />, label: '退出登录', onClick: () => { logout(); navigate('/login') } },
                ],
              }}
            >
              <Space style={{ cursor: 'pointer' }}>
                <Avatar size="small" icon={<UserOutlined />} />
                <Typography.Text>{user.nickname || user.username}</Typography.Text>
                <Tag color="gold">{user.points} 分</Tag>
                {admin && <Tag color="blue">管理员</Tag>}
              </Space>
            </Dropdown>
          ) : (
            <Space>
              <SettingOutlined />
              <Typography.Text>未登录</Typography.Text>
            </Space>
          )}
        </Header>
        <Content style={{ margin: 16 }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}
