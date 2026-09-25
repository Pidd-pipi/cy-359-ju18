import { useEffect, useState } from 'react'
import { Button, Card, Space, Table, message } from 'antd'
import * as productApi from '../api/product'
import StatusBadge from '../components/StatusBadge'
import PageHeader from '../components/PageHeader'
import { useAuth } from '../hooks/useAuth'
import { isAdmin } from '../utils/auth'
import { formatDateTime } from '../utils/format'
import { RedemptionStatus } from '../constants'

export default function Redemptions() {
  useAuth()
  const [mine, setMine] = useState<productApi.Redemption[]>([])
  const [all, setAll] = useState<productApi.Redemption[]>([])
  const [loading, setLoading] = useState(false)

  const load = async () => {
    setLoading(true)
    try {
      const mineRes = await productApi.listMyRedemptions({ page: 1, page_size: 50 })
      setMine(mineRes.data.list)
      if (isAdmin()) {
        const allRes = await productApi.listAllRedemptions({ page: 1, page_size: 50 })
        setAll(allRes.data.list)
      }
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [])

  const onProcess = async (id: number, status: string) => {
    await productApi.updateRedemptionStatus(id, status)
    message.success('已更新')
    load()
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    { title: '商品', dataIndex: 'product_name' },
    { title: '数量', dataIndex: 'quantity', width: 70 },
    { title: '消耗积分', dataIndex: 'total_points', width: 100 },
    { title: '状态', dataIndex: 'status', width: 100, render: (v: string) => <StatusBadge status={v} domain="redemption" /> },
    { title: '兑换时间', dataIndex: 'redeemed_at', width: 160, render: (v: string) => formatDateTime(v) },
  ]

  return (
    <div>
      <PageHeader title="兑换记录" />
      <Card title="我的兑换记录" style={{ marginBottom: 16 }}>
        <Table<productApi.Redemption>
          rowKey="id"
          size="small"
          loading={loading}
          dataSource={mine}
          columns={columns}
          pagination={false}
        />
      </Card>
      {isAdmin() && (
        <Card title="全部兑换记录（管理员）">
          <Table<productApi.Redemption>
            rowKey="id"
            size="small"
            loading={loading}
            dataSource={all}
            pagination={false}
            columns={[
              ...columns,
              { title: '用户', dataIndex: 'username', width: 120 },
              {
                title: '操作',
                width: 180,
                render: (_, r) =>
                  r.status === RedemptionStatus.PENDING ? (
                    <Space>
                      <Button size="small" type="primary" onClick={() => onProcess(r.id, RedemptionStatus.COMPLETED)}>
                        完成
                      </Button>
                      <Button size="small" danger onClick={() => onProcess(r.id, RedemptionStatus.CANCELLED)}>
                        取消
                      </Button>
                    </Space>
                  ) : null,
              },
            ]}
          />
        </Card>
      )}
    </div>
  )
}
