import { useEffect, useState } from 'react'
import { Button, Card, Descriptions, Form, Input, message } from 'antd'
import * as userApi from '../api/user'
import PageHeader from '../components/PageHeader'
import { useAuth } from '../hooks/useAuth'
import { useAuthStore } from '../stores/authStore'
import { formatDateTime } from '../utils/format'
import { RoleType } from '../constants'

export default function Profile() {
  useAuth()
  const [user, setUser] = useState<userApi.User | null>(null)
  const refreshProfile = useAuthStore((s) => s.refreshProfile)
  const [form] = Form.useForm()
  const [saving, setSaving] = useState(false)

  const load = async () => {
    const res = await userApi.getProfile()
    setUser(res.data)
    form.setFieldsValue(res.data)
  }

  useEffect(() => {
    load()
  }, [])

  const onSave = async (values: { nickname?: string; email?: string; phone?: string }) => {
    setSaving(true)
    try {
      await userApi.updateProfile(values)
      message.success('资料已更新')
      await refreshProfile()
      load()
    } finally {
      setSaving(false)
    }
  }

  return (
    <div>
      <PageHeader title="个人中心" />
      <Card style={{ marginBottom: 16 }}>
        <Descriptions column={2} bordered size="small">
          <Descriptions.Item label="用户名">{user?.username}</Descriptions.Item>
          <Descriptions.Item label="角色">{user?.role === RoleType.ADMIN ? '管理员' : '普通用户'}</Descriptions.Item>
          <Descriptions.Item label="当前积分"><b style={{ color: '#fa8c16' }}>{user?.points ?? 0}</b></Descriptions.Item>
          <Descriptions.Item label="注册时间">{user ? formatDateTime(user.created_at) : '-'}</Descriptions.Item>
        </Descriptions>
      </Card>
      <Card title="编辑资料">
        <Form form={form} onFinish={onSave} layout="vertical" style={{ maxWidth: 480 }}>
          <Form.Item name="nickname" label="昵称" rules={[{ required: true, min: 1, max: 50 }]}>
            <Input />
          </Form.Item>
          <Form.Item name="email" label="邮箱" rules={[{ type: 'email', message: '邮箱格式不正确' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="phone" label="手机号" rules={[{ pattern: /^1\d{10}$/, message: '手机号格式不正确' }]}>
            <Input />
          </Form.Item>
          <Button type="primary" htmlType="submit" loading={saving}>保存</Button>
        </Form>
      </Card>
    </div>
  )
}
