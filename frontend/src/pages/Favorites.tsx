import { useEffect } from 'react'
import { List, Button, Space, Tag } from 'antd'
import { useNavigate } from 'react-router-dom'
import { useFavoriteStore } from '../stores/favoriteStore'
import { useAuthStore } from '../stores/authStore'
import PageHeader from '../components/PageHeader'
import EmptyState from '../components/EmptyState'
import StatusBadge from '../components/StatusBadge'
import { useAuth } from '../hooks/useAuth'
import * as favoriteApi from '../api/favorite'
import { formatDateTime } from '../utils/format'

export default function Favorites() {
  useAuth()
  const navigate = useNavigate()
  const { favorites, fetchFavorites } = useFavoriteStore()
  const user = useAuthStore((s) => s.user)

  useEffect(() => {
    if (user) fetchFavorites()
  }, [user, fetchFavorites])

  const onRemove = async (activityId: number) => {
    await favoriteApi.removeFavorite(activityId)
    fetchFavorites()
  }

  return (
    <div>
      <PageHeader title="历史线路收藏" />
      {favorites.length === 0 ? (
        <EmptyState description="还没有收藏线路，去已结束的活动详情页收藏吧" />
      ) : (
        <List
          dataSource={favorites}
          renderItem={(fav) => (
            <List.Item
              actions={[
                <Button key="detail" type="link" onClick={() => navigate(`/activities/${fav.activity_id}`)}>查看</Button>,
                <Button key="remove" type="link" danger onClick={() => onRemove(fav.activity_id)}>取消收藏</Button>,
              ]}
            >
              <List.Item.Meta
                title={
                  <Space>
                    <span>{fav.activity.title}</span>
                    <StatusBadge status={fav.activity.status} />
                    <Tag>{fav.activity.address}</Tag>
                  </Space>
                }
                description={`收藏时间：${formatDateTime(fav.created_at)}`}
              />
            </List.Item>
          )}
        />
      )}
    </div>
  )
}
