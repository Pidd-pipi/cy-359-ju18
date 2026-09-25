import { useEffect, useState } from 'react'
import { useParams, useSearchParams, useNavigate } from 'react-router-dom'
import {
  Card, Button, List, Space, Typography, Input, Radio, Tag, message, Col, Row, Modal,
} from 'antd'
import * as activityApi from '../api/activity'
import * as checkinApi from '../api/checkin'
import Leaderboard from '../components/Leaderboard'
import PageHeader from '../components/PageHeader'
import { useAuth } from '../hooks/useAuth'
import { TaskType, CheckinType, ActivityStatus } from '../constants'

export default function Checkin() {
  useAuth()
  const { id } = useParams()
  const [searchParams] = useSearchParams()
  const activityId = Number(id)
  const teamId = Number(searchParams.get('team') || 0)
  const navigate = useNavigate()

  const [activity, setActivity] = useState<activityApi.Activity | null>(null)
  const [answer, setAnswer] = useState('')
  const [checkinType, setCheckinType] = useState<string>(CheckinType.GPS)
  const [selectedCp, setSelectedCp] = useState<activityApi.Checkpoint | null>(null)
  const [doneCps, setDoneCps] = useState<Set<number>>(new Set())
  const [records, setRecords] = useState<checkinApi.CheckinRecord[]>([])
  const [loading, setLoading] = useState(false)

  const load = async () => {
    const res = await activityApi.getActivity(activityId)
    setActivity(res.data)
    if (teamId) {
      const rec = await checkinApi.listTeamCheckins(teamId, activityId)
      setRecords(rec.data)
      setDoneCps(new Set(rec.data.map((r) => r.checkpoint_id)))
    }
  }

  useEffect(() => {
    load()
  }, [activityId, teamId])

  const onCheckin = async () => {
    if (!selectedCp) return
    if (!teamId) {
      message.warning('请先选择团队')
      navigate('/teams')
      return
    }
    setLoading(true)
    try {
      await checkinApi.checkin(teamId, {
        checkpoint_id: selectedCp.id,
        checkin_type: checkinType,
        latitude: checkinType === CheckinType.GPS ? selectedCp.lat + 0.0001 : undefined,
        longitude: checkinType === CheckinType.GPS ? selectedCp.lng + 0.0001 : undefined,
        answer: selectedCp.task_type === TaskType.QUIZ ? answer : undefined,
        photo_url: selectedCp.task_type === TaskType.PHOTO ? 'https://example.com/photo.jpg' : undefined,
      })
      message.success('打卡成功！')
      setAnswer('')
      setSelectedCp(null)
      load()
    } catch {
      // 已提示
    } finally {
      setLoading(false)
    }
  }

  const openCp = (cp: activityApi.Checkpoint) => {
    setSelectedCp(cp)
    setAnswer('')
  }

  if (!activity) return null

  if (activity.status !== ActivityStatus.ONGOING) {
    return (
      <div>
        <PageHeader title={`打卡 · ${activity.title}`} />
        <Card>
          <Typography.Paragraph type="secondary">
            当前活动状态为「{activity.status}」，仅进行中的活动可打卡。
          </Typography.Paragraph>
          <Button type="primary" onClick={() => navigate(`/activities/${activityId}`)}>返回活动详情</Button>
        </Card>
      </div>
    )
  }

  return (
    <div>
      <PageHeader
        title={`打卡 · ${activity.title}`}
        extra={<Tag color="blue">团队 #{teamId || '-'}</Tag>}
      />
      <Row gutter={16}>
        <Col xs={24} lg={12}>
          <Card title="打卡点列表" style={{ marginBottom: 16 }}>
            <List
              dataSource={activity.checkpoints || []}
              renderItem={(cp) => {
                const done = doneCps.has(cp.id)
                return (
                  <List.Item
                    actions={[
                      <Button
                        key="btn"
                        type={done ? 'default' : 'primary'}
                        disabled={done}
                        onClick={() => openCp(cp)}
                      >
                        {done ? '已打卡' : '打卡'}
                      </Button>,
                    ]}
                  >
                    <List.Item.Meta
                      title={
                        <Space>
                          <span>CP{cp.sequence} · {cp.name}</span>
                          {done && <Tag color="success">已完成</Tag>}
                        </Space>
                      }
                      description={
                        <Space direction="vertical" size={0}>
                          <Typography.Text type="secondary">线索：{cp.clue || '无'}</Typography.Text>
                          <Typography.Text type="secondary">任务：{cp.task_content || '无任务'}</Typography.Text>
                        </Space>
                      }
                    />
                  </List.Item>
                )
              }}
            />
          </Card>
        </Col>
        <Col xs={24} lg={12}>
          <Card title="排行榜（实时）" style={{ marginBottom: 16 }}>
            <Leaderboard activityId={activityId} height={300} />
          </Card>
        </Col>
      </Row>

      <Modal
        title={`打卡 CP${selectedCp?.sequence || ''} · ${selectedCp?.name || ''}`}
        open={!!selectedCp}
        onCancel={() => setSelectedCp(null)}
        onOk={onCheckin}
        confirmLoading={loading}
        okText="确认打卡"
      >
        {selectedCp && (
          <Space direction="vertical" style={{ width: '100%' }}>
            <Radio.Group value={checkinType} onChange={(e) => setCheckinType(e.target.value)}>
              <Radio.Button value={CheckinType.GPS}>GPS 定位</Radio.Button>
              <Radio.Button value={CheckinType.QRCODE}>二维码</Radio.Button>
            </Radio.Group>
            {selectedCp.task_type === TaskType.QUIZ && (
              <div>
                <Typography.Paragraph strong>任务：{selectedCp.task_content}</Typography.Paragraph>
                <Input.TextArea
                  placeholder="请输入答案"
                  value={answer}
                  onChange={(e) => setAnswer(e.target.value)}
                />
              </div>
            )}
            {selectedCp.task_type === TaskType.PHOTO && (
              <Typography.Text type="secondary">请在现场拍照完成任务（已模拟提交示例照片）</Typography.Text>
            )}
          </Space>
        )}
      </Modal>
    </div>
  )
}
