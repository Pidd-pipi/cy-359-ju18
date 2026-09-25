import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Card, Descriptions, List, Tag, Space, Button, Table } from 'antd'
import * as teamApi from '../api/team'
import * as checkinApi from '../api/checkin'
import StatusBadge from '../components/StatusBadge'
import RegistrationStatusTag from '../components/RegistrationStatusTag'
import PageHeader from '../components/PageHeader'
import { useAuth } from '../hooks/useAuth'
import { formatDateTime, formatDuration } from '../utils/format'
import { TeamRole } from '../constants'

export default function TeamDetail() {
  useAuth()
  const { id } = useParams()
  const teamId = Number(id)
  const navigate = useNavigate()
  const [team, setTeam] = useState<teamApi.Team | null>(null)
  const [registrations, setRegistrations] = useState<teamApi.Registration[]>([])
  const [checkins, setCheckins] = useState<checkinApi.CheckinRecord[]>([])

  useEffect(() => {
    teamApi.getTeam(teamId).then((res) => setTeam(res.data)).catch(() => {})
    teamApi.listMyRegistrations().then((res) => {
      setRegistrations(res.data.filter((r) => r.team_id === teamId))
    }).catch(() => {})
  }, [teamId])

  const loadCheckins = (activityId: number) => {
    checkinApi.listTeamCheckins(teamId, activityId).then((res) => setCheckins(res.data)).catch(() => {})
  }

  return (
    <div>
      <PageHeader
        title={team?.name || '团队详情'}
        extra={<Button onClick={() => navigate('/teams')}>返回</Button>}
      />
      {team && (
        <Card style={{ marginBottom: 16 }}>
          <Descriptions column={2} bordered size="small">
            <Descriptions.Item label="口号">{team.slogan || '暂无'}</Descriptions.Item>
            <Descriptions.Item label="成员数">{team.member_count} 人</Descriptions.Item>
          </Descriptions>
          <List
            style={{ marginTop: 12 }}
            header="团队成员"
            dataSource={team.members || []}
            renderItem={(m) => (
              <List.Item>
                <Space>
                  <span>{m.nickname || m.username}</span>
                  <Tag color={m.role === TeamRole.CAPTAIN ? 'gold' : 'default'}>
                    {m.role === TeamRole.CAPTAIN ? '队长' : '队员'}
                  </Tag>
                </Space>
              </List.Item>
            )}
          />
        </Card>
      )}

      <Card title="报名记录" style={{ marginBottom: 16 }}>
        <Table<teamApi.Registration>
          rowKey="id"
          size="small"
          dataSource={registrations}
          pagination={false}
          columns={[
            { title: '活动', dataIndex: 'activity_title', render: (v: string, r: teamApi.Registration) => v || `活动 #${r.activity_id}` },
            { title: '状态', dataIndex: 'status', render: (_: string, r: teamApi.Registration) => (
              <RegistrationStatusTag status={r.status} waitlistAhead={r.waitlist_ahead} />
            ) },
            { title: '用时', dataIndex: 'duration' },
            { title: '报名时间', dataIndex: 'registered_at', render: (v: string) => formatDateTime(v) },
            {
              title: '操作',
              render: (_, r) => (
                <Space>
                  <Button size="small" onClick={() => loadCheckins(r.activity_id)}>打卡记录</Button>
                  <Button size="small" type="primary" onClick={() => navigate(`/checkin/${r.activity_id}?team=${teamId}`)}>去打卡</Button>
                </Space>
              ),
            },
          ]}
        />
      </Card>

      <Card title="打卡记录">
        <Table<checkinApi.CheckinRecord>
          rowKey="id"
          size="small"
          dataSource={checkins}
          pagination={false}
          columns={[
            { title: '打卡点 ID', dataIndex: 'checkpoint_id' },
            { title: '方式', dataIndex: 'checkin_type', render: (v: string) => (v === 'gps' ? 'GPS' : '二维码') },
            { title: '结果', dataIndex: 'result', render: (v: string) => <StatusBadge status={v} domain="checkin" fallback={v} /> },
            { title: '获得积分', dataIndex: 'points_earned' },
            { title: '时间', dataIndex: 'checked_in_at', render: (v: string) => formatDateTime(v) },
          ]}
        />
      </Card>
    </div>
  )
}
