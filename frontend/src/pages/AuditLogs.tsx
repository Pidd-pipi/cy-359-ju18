import { useEffect, useState } from 'react'
import { Input, Table, Tag } from 'antd'
import * as auditApi from '../api/audit'
import PageHeader from '../components/PageHeader'
import { useAuth } from '../hooks/useAuth'
import { formatDateTime } from '../utils/format'

export default function AuditLogs() {
  useAuth(true)
  const [logs, setLogs] = useState<auditApi.AuditLog[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [keyword, setKeyword] = useState('')
  const [loading, setLoading] = useState(false)

  const load = async (p = page, kw = keyword) => {
    setLoading(true)
    try {
      const res = await auditApi.listAuditLogs({ page: p, page_size: 20, keyword: kw })
      setLogs(res.data.list)
      setTotal(res.data.total)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [page])

  return (
    <div>
      <PageHeader
        title="操作审计日志（管理员）"
        extra={
          <Input.Search
            placeholder="搜索用户名/详情"
            allowClear
            style={{ width: 240 }}
            onSearch={(v) => {
              setKeyword(v)
              setPage(1)
              load(1, v)
            }}
          />
        }
      />
      <Table<auditApi.AuditLog>
        rowKey="id"
        size="small"
        loading={loading}
        dataSource={logs}
        pagination={{ current: page, pageSize: 20, total, showTotal: (t) => `共 ${t} 条`, onChange: setPage }}
        scroll={{ x: 900 }}
        columns={[
          { title: 'ID', dataIndex: 'id', width: 60 },
          { title: '用户', dataIndex: 'username', width: 100 },
          { title: '动作', dataIndex: 'action', width: 200, render: (v: string) => <Tag>{v}</Tag> },
          { title: '资源', dataIndex: 'resource_type', width: 110 },
          { title: '资源 ID', dataIndex: 'resource_id', width: 90 },
          { title: '详情', dataIndex: 'detail' },
          { title: 'IP', dataIndex: 'ip', width: 120 },
          { title: '时间', dataIndex: 'created_at', width: 160, render: (v: string) => formatDateTime(v) },
        ]}
      />
    </div>
  )
}
