// 共享组件：活动卡片（首页/收藏页复用）
import { Card, Tag, Space, Typography } from 'antd'
import { EnvironmentOutlined, ClockCircleOutlined, TeamOutlined } from '@ant-design/icons'
import type { Activity } from '../api/activity'
import { difficultyConfig } from '../constants'
import { formatDateTime } from '../utils/format'
import StatusBadge from './StatusBadge'

interface Props {
  activity: Activity
  onClick?: () => void
  extra?: React.ReactNode
}

export default function ActivityCard({ activity, onClick, extra }: Props) {
  return (
    <Card
      hoverable={!!onClick}
      onClick={onClick}
      title={activity.title}
      extra={
        <Space>
          <StatusBadge status={activity.status} />
          {extra}
        </Space>
      }
      style={{ marginBottom: 16 }}
    >
      <Space direction="vertical" size={4} style={{ width: '100%' }}>
        <Typography.Text type="secondary">{activity.description || '暂无描述'}</Typography.Text>
        <Space wrap>
          <Tag color="blue">{difficultyConfig[activity.difficulty] || activity.difficulty}</Tag>
          <Tag icon={<ClockCircleOutlined />}>{activity.duration_minutes} 分钟</Tag>
          <Tag icon={<TeamOutlined />}>已报名 {activity.team_count}/{activity.max_teams}</Tag>
          <Tag icon={<EnvironmentOutlined />}>{activity.address}</Tag>
        </Space>
        <Typography.Text type="secondary" style={{ fontSize: 12 }}>
          时间：{formatDateTime(activity.start_time)} ~ {formatDateTime(activity.end_time)}
        </Typography.Text>
      </Space>
    </Card>
  )
}
