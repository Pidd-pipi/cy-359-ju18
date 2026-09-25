import { Button, Card, Form, Input, Typography, message } from 'antd'
import { LockOutlined, UserOutlined } from '@ant-design/icons'
import { useNavigate, Link } from 'react-router-dom'
import { useState } from 'react'
import { useAuthStore } from '../stores/authStore'

export default function Login() {
  const navigate = useNavigate()
  const login = useAuthStore((s) => s.login)
  const [loading, setLoading] = useState(false)

  const onFinish = async (values: { username: string; password: string }) => {
    setLoading(true)
    try {
      await login(values.username, values.password)
      message.success('登录成功')
      navigate('/', { replace: true })
    } catch {
      // 错误已由拦截器提示
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={{ maxWidth: 400, margin: '80px auto' }}>
      <Card>
        <Typography.Title level={3} style={{ textAlign: 'center' }}>
          城市定向越野活动平台
        </Typography.Title>
        <Typography.Paragraph type="secondary" style={{ textAlign: 'center' }}>
          登录后参与线路打卡、团队报名与积分兑换
        </Typography.Paragraph>
        <Form onFinish={onFinish} layout="vertical">
          <Form.Item name="username" rules={[{ required: true, message: '请输入用户名' }]}>
            <Input prefix={<UserOutlined />} placeholder="用户名" autoComplete="username" />
          </Form.Item>
          <Form.Item name="password" rules={[{ required: true, message: '请输入密码' }]}>
            <Input.Password prefix={<LockOutlined />} placeholder="密码" autoComplete="current-password" />
          </Form.Item>
          <Button type="primary" htmlType="submit" block loading={loading}>
            登录
          </Button>
        </Form>
        <div style={{ marginTop: 12, textAlign: 'center' }}>
          还没有账号？<Link to="/register">立即注册</Link>
        </div>
      </Card>
    </div>
  )
}
