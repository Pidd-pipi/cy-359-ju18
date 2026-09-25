import { useEffect, useState } from 'react'
import { Button, Input, Space, Switch, Table, message } from 'antd'
import * as userApi from '../api/user'
import PageHeader from '../components/PageHeader'
import { useAuth } from '../hooks/useAuth'
import { formatDateTime, maskPhone } from '../utils/format'
import { RoleType } from '../constants'

export default function Users() {
  useAuth(true)
  const [users, setUsers] = useState<userApi.User[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [keyword, setKeyword] = useState('')
  const [loading, setLoading] = useState(false)

  const load = async (p = page, kw = keyword) => {
    setLoading(true)
    try {
      const res = await userApi.listUsers({ page: p, page_size: 10, keyword: kw })
      setUsers(res.data.list)
      setTotal(res.data.total)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [page])

  const onToggle = async (u: userApi.User, disabled: boolean) => {
    await userApi.setUserDisabled(u.id, disabled)
    message.success('已更新')
    load()
  }

  return (
    <div>
      <PageHeader
        title="用户管理（管理员）"
        extra={
          <Input.Search
            placeholder="搜索用户名/昵称"
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
      <Table<userApi.User>
        rowKey="id"
        size="small"
        loading={loading}
        dataSource={users}
        pagination={{ current: page, pageSize: 10, total, showTotal: (t) => `共 ${t} 条`, onChange: setPage }}
        columns={[
          { title: 'ID', dataIndex: 'id', width: 60 },
          { title: '用户名', dataIndex: 'username' },
          { title: '昵称', dataIndex: 'nickname' },
          { title: '角色', dataIndex: 'role', width: 90, render: (v: string) => (v === RoleType.ADMIN ? '管理员' : '用户') },
          { title: '积分', dataIndex: 'points', width: 80 },
          { title: '手机号', dataIndex: 'phone', width: 120, render: (v: string) => maskPhone(v) },
          { title: '注册时间', dataIndex: 'created_at', width: 160, render: (v: string) => formatDateTime(v) },
          {
            title: '禁用',
            width: 90,
            render: (_, u) => (
              <Switch
                size="small"
                checked={u.disabled}
                disabled={u.role === RoleType.ADMIN}
                onChange={(v) => onToggle(u, v)}
              />
            ),
          },
        ]}
      />
    </div>
  )
}
