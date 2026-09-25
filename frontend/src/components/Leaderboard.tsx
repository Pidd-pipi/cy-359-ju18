// 共享组件：实时排行榜（活动详情页 / 打卡页 / 排行榜页复用）
import { Table, Tag, Typography } from 'antd'
import { useEffect } from 'react'
import type { LeaderboardRow } from '../api/checkin'
import { useWebSocket } from '../hooks/useWebSocket'
import { useCheckinStore } from '../stores/checkinStore'
import { RegistrationStatus } from '../constants'

interface Props {
  activityId: number
  height?: number
}

export default function Leaderboard({ activityId, height }: Props) {
  const { rows: wsRows } = useWebSocket(activityId)
  const leaderboard = useCheckinStore((s) => s.leaderboard)
  const fetchLeaderboard = useCheckinStore((s) => s.fetchLeaderboard)

  useEffect(() => {
    fetchLeaderboard(activityId)
  }, [activityId, fetchLeaderboard])

  const data = wsRows || leaderboard[activityId] || []

  const columns = [
    { title: '排名', dataIndex: 'rank', width: 60, render: (v: number) => <Tag color={v <= 3 ? 'gold' : 'default'}>#{v}</Tag> },
    { title: '队伍', dataIndex: 'team_name', render: (_: string, r: LeaderboardRow) => r.team_name || `队伍 #${r.team_id}` },
    { title: '打卡数', dataIndex: 'checkpoint_count', width: 80 },
    {
      title: '用时',
      dataIndex: 'duration',
      width: 100,
      render: (v: string, r: LeaderboardRow) =>
        r.status === RegistrationStatus.FINISHED ? <Typography.Text strong>{v}</Typography.Text> : '-',
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 90,
      render: (v: string) =>
        v === RegistrationStatus.FINISHED ? '已完成' : '进行中',
    },
  ]

  return (
    <Table<LeaderboardRow>
      rowKey="team_id"
      size="small"
      columns={columns}
      dataSource={data}
      pagination={false}
      scroll={{ y: height }}
    />
  )
}
