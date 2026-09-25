import { useEffect, useState } from 'react'
import { Input, Select, Space, Row, Col, Spin } from 'antd'
import { useNavigate } from 'react-router-dom'
import ActivityCard from '../components/ActivityCard'
import EmptyState from '../components/EmptyState'
import { useActivityStore } from '../stores/activityStore'
import { Difficulty, ActivityStatus } from '../constants'

export default function ActivityList() {
  const navigate = useNavigate()
  const { list, total, loading, fetchList } = useActivityStore()
  const [difficulty, setDifficulty] = useState<string>()
  const [status, setStatus] = useState<string>()
  const [keyword, setKeyword] = useState<string>()

  useEffect(() => {
    fetchList({ difficulty, status, keyword })
  }, [fetchList, difficulty, status, keyword])

  return (
    <div>
      <Space wrap style={{ marginBottom: 16 }}>
        <Input.Search
          placeholder="搜索线路名称/地点"
          allowClear
          style={{ width: 260 }}
          onSearch={(v) => setKeyword(v || undefined)}
        />
        <Select
          placeholder="难度"
          allowClear
          style={{ width: 120 }}
          onChange={(v) => setDifficulty(v || undefined)}
          options={Object.entries(Difficulty).map(([k, v]) => ({ label: k, value: v }))}
        />
        <Select
          placeholder="状态"
          allowClear
          style={{ width: 140 }}
          onChange={(v) => setStatus(v || undefined)}
          options={Object.entries(ActivityStatus).map(([k, v]) => ({ label: k, value: v }))}
        />
      </Space>
      {loading ? (
        <Spin style={{ display: 'block', margin: '60px auto' }} />
      ) : list.length === 0 ? (
        <EmptyState description="暂无活动线路" />
      ) : (
        <Row gutter={16}>
          {list.map((a) => (
            <Col xs={24} md={12} lg={8} key={a.id}>
              <ActivityCard activity={a} onClick={() => navigate(`/activities/${a.id}`)} />
            </Col>
          ))}
        </Row>
      )}
      <div style={{ textAlign: 'right', color: '#999' }}>共 {total} 条线路</div>
    </div>
  )
}
