import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import {
  Card, Descriptions, Tag, Space, Button, List, Typography, Select, message, Timeline,
} from 'antd'
import { EnvironmentOutlined } from '@ant-design/icons'
import * as activityApi from '../api/activity'
import * as teamApi from '../api/team'
import * as favoriteApi from '../api/favorite'
import StatusBadge from '../components/StatusBadge'
import Leaderboard from '../components/Leaderboard'
import { difficultyConfig, taskTypeConfig, ActivityStatus } from '../constants'
import { formatDateTime } from '../utils/format'
import { getToken } from '../utils/auth'

export default function ActivityDetail() {
  const { id } = useParams()
  const activityId = Number(id)
  const navigate = useNavigate()
  const [activity, setActivity] = useState<activityApi.Activity | null>(null)
  const [myTeams, setMyTeams] = useState<teamApi.Team[]>([])
  const [selectedTeam, setSelectedTeam] = useState<number>()
  const [loading, setLoading] = useState(false)

  const load = async () => {
    const res = await activityApi.getActivity(activityId)
    setActivity(res.data)
  }

  useEffect(() => {
    load()
  }, [activityId])

  useEffect(() => {
    if (getToken()) {
      teamApi.listMyTeams().then((res) => {
        setMyTeams(res.data)
        if (res.data.length > 0) setSelectedTeam(res.data[0].id)
      }).catch(() => {})
    }
  }, [])

  const onApply = async () => {
    if (!selectedTeam) {
      message.warning('请先创建或加入一个团队')
      navigate('/teams')
      return
    }
    setLoading(true)
    try {
      const res = await teamApi.applyActivity({ team_id: selectedTeam, activity_id: activityId })
      // 后端在名额满时返回 waitlist 状态及前面队伍数
      if (res.data?.status === 'waitlist') {
        message.success(`名额已满，已进入候补，前面还有 ${res.data.waitlist_ahead ?? 0} 队`)
      } else {
        message.success('报名成功，等待管理员审核')
      }
      load()
    } catch {
      // 已提示
    } finally {
      setLoading(false)
    }
  }

  const onFavorite = async () => {
    try {
      await favoriteApi.addFavorite(activityId)
      message.success('收藏成功')
    } catch {
      // 已提示
    }
  }

  if (!activity) return null

  const canApply = activity.status === ActivityStatus.PUBLISHED

  return (
    <div>
      <Card
        title={activity.title}
        extra={<StatusBadge status={activity.status} domain="activity" />}
        style={{ marginBottom: 16 }}
      >
        <Descriptions column={{ xs: 1, md: 2 }} bordered size="small">
          <Descriptions.Item label="难度">
            <Tag color="blue">{difficultyConfig[activity.difficulty] || activity.difficulty}</Tag>
          </Descriptions.Item>
          <Descriptions.Item label="时长">{activity.duration_minutes} 分钟</Descriptions.Item>
          <Descriptions.Item label="地点">
            <Space>
              <EnvironmentOutlined />
              {activity.address}
            </Space>
          </Descriptions.Item>
          <Descriptions.Item label="报名情况">
            {activity.team_count}/{activity.max_teams} 队
          </Descriptions.Item>
          <Descriptions.Item label="开始时间">{formatDateTime(activity.start_time)}</Descriptions.Item>
          <Descriptions.Item label="结束时间">{formatDateTime(activity.end_time)}</Descriptions.Item>
          <Descriptions.Item label="装备要求" span={2}>
            {activity.equipment_requirement || '无'}
          </Descriptions.Item>
          <Descriptions.Item label="线路描述" span={2}>
            {activity.description || '暂无'}
          </Descriptions.Item>
        </Descriptions>
        <Space style={{ marginTop: 16 }}>
          {canApply && (
            <Space>
              <Select
                placeholder="选择报名团队"
                style={{ width: 200 }}
                value={selectedTeam}
                onChange={setSelectedTeam}
                options={myTeams.map((t) => ({ label: t.name, value: t.id }))}
              />
              <Button type="primary" loading={loading} onClick={onApply}>
                团队报名
              </Button>
            </Space>
          )}
          {activity.status === ActivityStatus.FINISHED && (
            <Button onClick={onFavorite}>收藏线路</Button>
          )}
          <Button onClick={() => navigate('/leaderboard/' + activityId)}>查看排行榜</Button>
        </Space>
      </Card>

      <Card title="打卡点（CP 点）" style={{ marginBottom: 16 }}>
        {activity.checkpoints && activity.checkpoints.length > 0 ? (
          <Timeline
            items={activity.checkpoints.map((cp) => ({
              children: (
                <div>
                  <Typography.Text strong>
                    CP{cp.sequence} · {cp.name}
                  </Typography.Text>
                  <div>
                    <Tag>{taskTypeConfig[cp.task_type] || cp.task_type}</Tag>
                    <Tag>半径 {cp.radius_meters} 米</Tag>
                    <Typography.Text type="secondary" style={{ marginLeft: 8 }}>
                      坐标：{cp.lat.toFixed(6)}, {cp.lng.toFixed(6)}
                    </Typography.Text>
                  </div>
                  {cp.clue && <Typography.Paragraph type="secondary" style={{ marginBottom: 0 }}>线索：{cp.clue}</Typography.Paragraph>}
                  {cp.task_content && <Typography.Paragraph style={{ marginBottom: 0 }}>任务：{cp.task_content}</Typography.Paragraph>}
                </div>
              ),
            }))}
          />
        ) : (
          <Typography.Text type="secondary">暂无打卡点</Typography.Text>
        )}
      </Card>
    </div>
  )
}
