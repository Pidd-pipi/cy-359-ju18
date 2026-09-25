import { useEffect, useState } from 'react'
import { Card, Col, Row, Button, Tag, Space, Typography, Modal, InputNumber, message, Spin } from 'antd'
import { GiftOutlined, ShoppingOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import * as productApi from '../api/product'
import { useProductStore } from '../stores/productStore'
import { useAuthStore } from '../stores/authStore'
import PageHeader from '../components/PageHeader'
import EmptyState from '../components/EmptyState'
import { useAuth } from '../hooks/useAuth'
import { productTypeConfig } from '../constants'

export default function Mall() {
  useAuth()
  const navigate = useNavigate()
  const { products, total, loading, fetchProducts } = useProductStore()
  const user = useAuthStore((s) => s.user)
  const refreshProfile = useAuthStore((s) => s.refreshProfile)
  const [selected, setSelected] = useState<productApi.Product | null>(null)
  const [quantity, setQuantity] = useState(1)
  const [redeeming, setRedeeming] = useState(false)

  useEffect(() => {
    fetchProducts({ status: 'on' })
  }, [fetchProducts])

  const onRedeem = async () => {
    if (!selected) return
    setRedeeming(true)
    try {
      await productApi.redeem({ product_id: selected.id, quantity })
      message.success('兑换成功，等待处理')
      setSelected(null)
      fetchProducts({ status: 'on' })
      refreshProfile()
    } catch {
      // 已提示
    } finally {
      setRedeeming(false)
    }
  }

  return (
    <div>
      <PageHeader
        title="积分商城"
        extra={
          <Space>
            <Tag icon={<GiftOutlined />} color="gold">我的积分：{user?.points ?? 0}</Tag>
            <Button onClick={() => navigate('/redemptions')}>我的兑换记录</Button>
          </Space>
        }
      />
      {loading ? (
        <Spin style={{ display: 'block', margin: '60px auto' }} />
      ) : products.length === 0 ? (
        <EmptyState description="商城暂无商品" />
      ) : (
        <Row gutter={16}>
          {products.map((p) => (
            <Col xs={24} sm={12} lg={8} key={p.id}>
              <Card
                style={{ marginBottom: 16 }}
                cover={
                  p.image_url ? (
                    <img src={p.image_url} alt={p.name} style={{ height: 160, objectFit: 'cover' }} />
                  ) : (
                    <div style={{ height: 160, background: '#f5f5f5', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                      <ShoppingOutlined style={{ fontSize: 48, color: '#bbb' }} />
                    </div>
                  )
                }
                actions={[
                  <Button key="redeem" type="primary" disabled={p.stock <= 0} onClick={() => setSelected(p)}>
                    {p.stock > 0 ? '兑换' : '已售罄'}
                  </Button>,
                ]}
              >
                <Card.Meta
                  title={p.name}
                  description={
                    <Space direction="vertical" size={4}>
                      <Tag color="blue">{productTypeConfig[p.type] || p.type}</Tag>
                      <Typography.Text strong style={{ color: '#fa8c16' }}>
                        {p.points_cost} 积分
                      </Typography.Text>
                      <Typography.Text type="secondary">库存 {p.stock}</Typography.Text>
                    </Space>
                  }
                />
              </Card>
            </Col>
          ))}
        </Row>
      )}
      <Typography.Text type="secondary">共 {total} 件商品</Typography.Text>

      <Modal
        title={`兑换「${selected?.name || ''}」`}
        open={!!selected}
        onCancel={() => setSelected(null)}
        onOk={onRedeem}
        confirmLoading={redeeming}
        okText="确认兑换"
      >
        {selected && (
          <Space direction="vertical" style={{ width: '100%' }}>
            <Typography.Text>
              单价 {selected.points_cost} 积分，共需 {(selected.points_cost * quantity)} 积分
            </Typography.Text>
            <InputNumber min={1} max={Math.min(selected.stock, 99)} value={quantity} onChange={(v) => setQuantity(v || 1)} />
          </Space>
        )}
      </Modal>
    </div>
  )
}
