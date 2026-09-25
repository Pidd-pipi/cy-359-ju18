import { useEffect, useState } from 'react'
import { Button, Card, Form, Input, List, Modal, Space, Typography, message } from 'antd'
import { TeamOutlined, PlusOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import * as teamApi from '../api/team'
import { useTeamStore } from '../stores/teamStore'
import PageHeader from '../components/PageHeader'
import EmptyState from '../components/EmptyState'
import { useAuth } from '../hooks/useAuth'
import { confirmAction } from '../components/ConfirmDialog'

export default function TeamList() {
  useAuth()
  const navigate = useNavigate()
  const { myTeams, fetchMyTeams, fetchMyRegistrations } = useTeamStore()
  const [open, setOpen] = useState(false)
  const [joinOpen, setJoinOpen] = useState(false)
  const [form] = Form.useForm()
  const [joinForm] = Form.useForm()
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    fetchMyTeams()
    fetchMyRegistrations()
  }, [fetchMyTeams, fetchMyRegistrations])

  const onCreate = async (values: { name: string; slogan?: string }) => {
    setLoading(true)
    try {
      await teamApi.createTeam(values)
      message.success('创建成功')
      setOpen(false)
      form.resetFields()
      fetchMyTeams()
    } catch {
      // 已提示
    } finally {
      setLoading(false)
    }
  }

  const onJoin = async (values: { team_id: number }) => {
    setLoading(true)
    try {
      await teamApi.joinTeam(values.team_id)
      message.success('加入成功')
      setJoinOpen(false)
      joinForm.resetFields()
      fetchMyTeams()
    } catch {
      // 已提示
    } finally {
      setLoading(false)
    }
  }

  const onLeave = (teamId: number, name: string) => {
    confirmAction('退出团队', `确定退出团队「${name}」吗？`, async () => {
      await teamApi.leaveTeam(teamId)
      message.success('已退出')
      fetchMyTeams()
    })
  }

  return (
    <div>
      <PageHeader
        title="我的团队"
        extra={
          <Space>
            <Button icon={<PlusOutlined />} type="primary" onClick={() => setOpen(true)}>
              创建团队
            </Button>
            <Button onClick={() => setJoinOpen(true)}>加入团队</Button>
          </Space>
        }
      />
      {myTeams.length === 0 ? (
        <EmptyState description="还没有团队，创建一个或加入他人团队开始定向越野" />
      ) : (
        <List
          grid={{ gutter: 16, xs: 1, sm: 2, lg: 3 }}
          dataSource={myTeams}
          renderItem={(team) => (
            <List.Item>
              <Card
                title={<Space><TeamOutlined />{team.name}</Space>}
                extra={
                  <Space>
                    <Button size="small" type="link" onClick={() => navigate(`/teams/${team.id}`)}>详情</Button>
                    {team.captain_id !== JSON.parse(localStorage.getItem('orienteering_user') || '{}').id && (
                      <Button size="small" danger onClick={() => onLeave(team.id, team.name)}>退出</Button>
                    )}
                  </Space>
                }
              >
                <Typography.Paragraph type="secondary">{team.slogan || '暂无口号'}</Typography.Paragraph>
                <Typography.Text>成员 {team.member_count} 人</Typography.Text>
              </Card>
            </List.Item>
          )}
        />
      )}

      <Modal title="创建团队" open={open} onCancel={() => setOpen(false)} footer={null}>
        <Form form={form} onFinish={onCreate} layout="vertical">
          <Form.Item name="name" label="团队名称" rules={[{ required: true, min: 2, max: 80, message: '2-80 个字符' }]}>
            <Input placeholder="团队名称" />
          </Form.Item>
          <Form.Item name="slogan" label="团队口号">
            <Input placeholder="口号（可选）" />
          </Form.Item>
          <Button type="primary" htmlType="submit" block loading={loading}>
            创建
          </Button>
        </Form>
      </Modal>

      <Modal title="加入团队" open={joinOpen} onCancel={() => setJoinOpen(false)} footer={null}>
        <Form form={joinForm} onFinish={onJoin} layout="vertical">
          <Form.Item name="team_id" label="团队 ID" rules={[{ required: true, message: '请输入团队 ID' }]}>
            <Input type="number" placeholder="团队 ID" />
          </Form.Item>
          <Button type="primary" htmlType="submit" block loading={loading}>
            加入
          </Button>
        </Form>
      </Modal>
    </div>
  )
}
