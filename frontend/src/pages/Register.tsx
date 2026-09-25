import { Button, Card, Form, Input, Typography, message } from 'antd'
import { useNavigate, Link } from 'react-router-dom'
import { useState } from 'react'
import { useAuthStore } from '../stores/authStore'

export default function Register() {
  const navigate = useNavigate()
  const register = useAuthStore((s) => s.register)
  const [loading, setLoading] = useState(false)

  const onFinish = async (values: { username: string; password: string; confirm: string; nickname: string; email?: string; phone?: string }) => {
    if (values.password !== values.confirm) {
      message.error('两次输入的密码不一致')
      return
    }
    setLoading(true)
    try {
      await register({
        username: values.username,
        password: values.password,
        nickname: values.nickname,
        email: values.email,
        phone: values.phone,
      })
      message.success('注册成功，请登录')
      navigate('/login', { replace: true })
    } catch {
      // 错误已由拦截器提示
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={{ maxWidth: 420, margin: '60px auto' }}>
      <Card>
        <Typography.Title level={3} style={{ textAlign: 'center' }}>
          注册账号
        </Typography.Title>
        <Form onFinish={onFinish} layout="vertical">
          <Form.Item name="username" label="用户名" rules={[{ required: true, min: 3, max: 50, message: '3-50 个字符' }]}>
            <Input placeholder="用户名" />
          </Form.Item>
          <Form.Item name="nickname" label="昵称" rules={[{ required: true, message: '请输入昵称' }]}>
            <Input placeholder="昵称" />
          </Form.Item>
          <Form.Item name="password" label="密码" rules={[{ required: true, min: 6, message: '至少 6 位' }]}>
            <Input.Password placeholder="密码" />
          </Form.Item>
          <Form.Item name="confirm" label="确认密码" rules={[{ required: true, message: '请再次输入密码' }]}>
            <Input.Password placeholder="确认密码" />
          </Form.Item>
          <Form.Item name="email" label="邮箱" rules={[{ type: 'email', message: '邮箱格式不正确' }]}>
            <Input placeholder="邮箱（可选）" />
          </Form.Item>
          <Form.Item name="phone" label="手机号" rules={[{ pattern: /^1\d{10}$/, message: '手机号格式不正确' }]}>
            <Input placeholder="手机号（可选）" />
          </Form.Item>
          <Button type="primary" htmlType="submit" block loading={loading}>
            注册
          </Button>
        </Form>
        <div style={{ marginTop: 12, textAlign: 'center' }}>
          已有账号？<Link to="/login">去登录</Link>
        </div>
      </Card>
    </div>
  )
}
