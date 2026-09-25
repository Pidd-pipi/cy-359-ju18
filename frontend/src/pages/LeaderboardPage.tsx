import { useParams, useNavigate } from 'react-router-dom'
import { Button, Card } from 'antd'
import Leaderboard from '../components/Leaderboard'
import PageHeader from '../components/PageHeader'

export default function LeaderboardPage() {
  const { id } = useParams()
  const activityId = Number(id)
  const navigate = useNavigate()
  return (
    <div>
      <PageHeader
        title="实时排行榜"
        extra={<Button onClick={() => navigate(`/activities/${activityId}`)}>返回活动</Button>}
      />
      <Card>
        <Leaderboard activityId={activityId} height={480} />
      </Card>
    </div>
  )
}
