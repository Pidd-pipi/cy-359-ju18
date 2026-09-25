import { useEffect, useState } from 'react'
import {
  Button, Card, Col, Form, Input, InputNumber, Modal, Row, Select, Space, DatePicker, Table, message, Popconfirm,
} from 'antd'
import dayjs from 'dayjs'
import * as activityApi from '../api/activity'
import * as checkpointApi from '../api/checkpoint'
import * as teamApi from '../api/team'
import PageHeader from '../components/PageHeader'
import StatusBadge from '../components/StatusBadge'
import RegistrationStatusTag from '../components/RegistrationStatusTag'
import { useAuth } from '../hooks/useAuth'
import { ActivityStatus, Difficulty, RegistrationStatus, TaskType } from '../constants'
import { formatDateTime } from '../utils/format'

export default function ActivityManage() {
  useAuth(true)
  const [list, setList] = useState<activityApi.Activity[]>([])
  const [editing, setEditing] = useState<activityApi.Activity | null>(null)
  const [open, setOpen] = useState(false)
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)
  const [cpActivity, setCpActivity] = useState<activityApi.Activity | null>(null)

  const load = async () => {
    const res = await activityApi.listActivities({ page: 1, page_size: 100 })
    setList(res.data.list)
  }

  useEffect(() => {
    load()
  }, [])

  const openCreate = () => {
    setEditing(null)
    form.resetFields()
    form.setFieldsValue({
      difficulty: Difficulty.ADULT,
      duration_minutes: 120,
      max_teams: 50,
      start_lat: 31.2304,
      start_lng: 121.4737,
      end_lat: 31.2404,
      end_lng: 121.4937,
    })
    setOpen(true)
  }

  const openEdit = (a: activityApi.Activity) => {
    setEditing(a)
    form.setFieldsValue({
      ...a,
      start_time: dayjs(a.start_time),
      end_time: dayjs(a.end_time),
    })
    setOpen(true)
  }

  const onSave = async () => {
    const values = await form.validateFields()
    const payload: activityApi.ActivityPayload = {
      ...values,
      start_time: values.start_time.toISOString(),
      end_time: values.end_time.toISOString(),
    }
    setLoading(true)
    try {
      if (editing) {
        await activityApi.updateActivity(editing.id, payload)
      } else {
        await activityApi.createActivity(payload)
      }
      message.success('保存成功')
      setOpen(false)
      load()
    } finally {
      setLoading(false)
    }
  }

  const onTransition = async (a: activityApi.Activity, status: string) => {
    await activityApi.transitionActivity(a.id, status)
    message.success('状态已更新')
    load()
  }

  const onDelete = async (id: number) => {
    await activityApi.deleteActivity(id)
    message.success('已删除')
    load()
  }

  const nextAction = (a: activityApi.Activity) => {
    switch (a.status) {
      case ActivityStatus.DRAFT:
        return { status: ActivityStatus.PUBLISHED, label: '发布' }
      case ActivityStatus.PUBLISHED:
        return { status: ActivityStatus.ONGOING, label: '开始' }
      case ActivityStatus.ONGOING:
        return { status: ActivityStatus.FINISHED, label: '结束' }
      default:
        return null
    }
  }

  return (
    <div>
      <PageHeader
        title="活动管理（管理员）"
        extra={<Button type="primary" onClick={openCreate}>新建活动</Button>}
      />
      <Table<activityApi.Activity>
        rowKey="id"
        size="small"
        loading={loading}
        dataSource={list}
        pagination={{ pageSize: 10, showTotal: (t) => `共 ${t} 条` }}
        scroll={{ x: 1100 }}
        columns={[
          { title: 'ID', dataIndex: 'id', width: 60 },
          { title: '标题', dataIndex: 'title' },
          { title: '难度', dataIndex: 'difficulty' },
          { title: '状态', dataIndex: 'status', width: 90, render: (v: string) => <StatusBadge status={v} domain="activity" /> },
          { title: '报名', dataIndex: 'team_count', width: 80 },
          { title: '时间', dataIndex: 'start_time', width: 160, render: (v: string) => formatDateTime(v) },
          {
            title: '操作',
            width: 320,
            render: (_, a) => {
              const action = nextAction(a)
              return (
                <Space>
                  {a.status === ActivityStatus.DRAFT && (
                    <Button size="small" onClick={() => openEdit(a)}>编辑</Button>
                  )}
                  <Button size="small" onClick={() => setCpActivity(a)}>打卡点</Button>
                  {action && (
                    <Button size="small" type="primary" onClick={() => onTransition(a, action.status)}>
                      {action.label}
                    </Button>
                  )}
                  {(a.status === ActivityStatus.DRAFT || a.status === ActivityStatus.CANCELLED) && (
                    <Popconfirm title="确定删除该活动？" onConfirm={() => onDelete(a.id)}>
                      <Button size="small" danger>删除</Button>
                    </Popconfirm>
                  )}
                </Space>
              )
            },
          },
        ]}
      />

      <Modal title={editing ? '编辑活动' : '新建活动'} open={open} onCancel={() => setOpen(false)} onOk={onSave} confirmLoading={loading} width={720}>
        <Form form={form} layout="vertical">
          <Row gutter={12}>
            <Col span={12}>
              <Form.Item name="title" label="标题" rules={[{ required: true, min: 2 }]}>
                <Input />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="difficulty" label="难度" rules={[{ required: true }]}>
                <Select options={Object.entries(Difficulty).map(([k, v]) => ({ label: k, value: v }))} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="duration_minutes" label="时长(分钟)" rules={[{ required: true }]}>
                <InputNumber min={15} max={1440} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="max_teams" label="最大团队数" rules={[{ required: true }]}>
                <InputNumber min={1} max={500} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="start_time" label="开始时间" rules={[{ required: true }]}>
                <DatePicker showTime style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="end_time" label="结束时间" rules={[{ required: true }]}>
                <DatePicker showTime style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="start_lat" label="起点纬度" rules={[{ required: true }]}>
                <InputNumber step={0.0001} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="start_lng" label="起点经度" rules={[{ required: true }]}>
                <InputNumber step={0.0001} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="end_lat" label="终点纬度" rules={[{ required: true }]}>
                <InputNumber step={0.0001} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="end_lng" label="终点经度" rules={[{ required: true }]}>
                <InputNumber step={0.0001} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={24}>
              <Form.Item name="address" label="地点" rules={[{ required: true }]}>
                <Input />
              </Form.Item>
            </Col>
            <Col span={24}>
              <Form.Item name="equipment_requirement" label="装备要求">
                <Input.TextArea rows={2} />
              </Form.Item>
            </Col>
            <Col span={24}>
              <Form.Item name="description" label="描述">
                <Input.TextArea rows={3} />
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </Modal>

      <CheckpointManager activity={cpActivity} onClose={() => setCpActivity(null)} onChanged={load} />
    </div>
  )
}

function CheckpointManager({ activity, onClose, onChanged }: { activity: activityApi.Activity | null; onClose: () => void; onChanged: () => void }) {
  const [form] = Form.useForm()
  const [cps, setCps] = useState<activityApi.Checkpoint[]>([])
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (activity) {
      setLoading(true)
      activityApi.getActivity(activity.id).then((res) => setCps(res.data.checkpoints || [])).finally(() => setLoading(false))
      form.setFieldsValue({ task_type: TaskType.NONE, radius_meters: 200, sequence: 1 })
    }
  }, [activity, form])

  const onCreate = async () => {
    if (!activity) return
    const values = await form.validateFields()
    await checkpointApi.createCheckpoint(activity.id, values)
    message.success('打卡点已创建')
    form.setFieldsValue({ sequence: values.sequence + 1 })
    const res = await activityApi.getActivity(activity.id)
    setCps(res.data.checkpoints || [])
    onChanged()
  }

  const onDelete = async (id: number) => {
    if (!activity) return
    await checkpointApi.deleteCheckpoint(id)
    message.success('已删除')
    const res = await activityApi.getActivity(activity.id)
    setCps(res.data.checkpoints || [])
    onChanged()
  }

  const onRegistrations = async () => {
    if (!activity) return
    const res = await teamApi.listRegistrationsByActivity(activity.id)
    const regs = res.data
    Modal.info({
      title: `报名列表（${activity.title}） 名额 ${activity.max_teams} 队`,
      width: 720,
      content: (
        <Table
          size="small"
          rowKey="id"
          dataSource={regs}
          pagination={false}
          rowClassName={(r) => (r.status === RegistrationStatus.WAITLIST ? 'ant-table-row-warning' : '')}
          columns={[
            { title: '报名 ID', dataIndex: 'id', width: 80 },
            { title: '团队 ID', dataIndex: 'team_id', width: 90 },
            {
              title: '状态',
              dataIndex: 'status',
              render: (v: string, r: teamApi.Registration) => (
                <RegistrationStatusTag status={v} waitlistAhead={r.waitlist_ahead} />
              ),
            },
            { title: '用时', dataIndex: 'duration' },
            {
              title: '操作',
              render: (_, r) => (
                <Space>
                  {r.status === RegistrationStatus.PENDING && (
                    <Button size="small" type="primary" onClick={async () => { await teamApi.approveRegistration(r.id); message.success('已通过'); onRegistrations(); }}>通过</Button>
                  )}
                  {(r.status === RegistrationStatus.PENDING || r.status === RegistrationStatus.WAITLIST) && (
                    <Button size="small" danger onClick={async () => { await teamApi.rejectRegistration(r.id); message.success(r.status === RegistrationStatus.PENDING ? '已拒绝，最早候补已自动补为待审核' : '已移除候补'); onRegistrations(); }}>拒绝</Button>
                  )}
                  {r.status !== RegistrationStatus.PENDING && r.status !== RegistrationStatus.WAITLIST && (
                    <span style={{ color: '#999' }}>—</span>
                  )}
                </Space>
              ),
            },
          ]}
        />
      ),
    })
  }

  return (
    <Modal title={`打卡点管理：${activity?.title || ''}`} open={!!activity} onCancel={onClose} footer={null} width={720}>
      {activity && (
        <>
          <Form form={form} layout="inline" style={{ marginBottom: 16 }}>
            <Form.Item name="name" rules={[{ required: true }]}><Input placeholder="名称" /></Form.Item>
            <Form.Item name="sequence" rules={[{ required: true }]}><InputNumber placeholder="序号" min={1} /></Form.Item>
            <Form.Item name="lat" rules={[{ required: true }]}><InputNumber placeholder="纬度" step={0.0001} /></Form.Item>
            <Form.Item name="lng" rules={[{ required: true }]}><InputNumber placeholder="经度" step={0.0001} /></Form.Item>
            <Form.Item name="task_type" rules={[{ required: true }]}>
              <Select style={{ width: 100 }} options={Object.entries(TaskType).map(([k, v]) => ({ label: k, value: v }))} />
            </Form.Item>
            <Form.Item name="radius_meters"><InputNumber placeholder="半径m" min={50} max={2000} /></Form.Item>
            <Form.Item name="clue"><Input placeholder="线索" /></Form.Item>
            <Form.Item name="task_content"><Input placeholder="任务内容" /></Form.Item>
            <Form.Item name="expected_answer"><Input placeholder="答案(答题)" /></Form.Item>
            <Form.Item name="qr_code" rules={[{ required: true }]}><Input placeholder="二维码串" /></Form.Item>
            <Button type="primary" onClick={onCreate}>添加</Button>
          </Form>
          <Table
            rowKey="id"
            size="small"
            loading={loading}
            dataSource={cps}
            pagination={false}
            columns={[
              { title: '序号', dataIndex: 'sequence', width: 60 },
              { title: '名称', dataIndex: 'name' },
              { title: '任务', dataIndex: 'task_type' },
              { title: '半径', dataIndex: 'radius_meters', width: 70 },
              {
                title: '操作',
                width: 120,
                render: (_, cp) => (
                  <Popconfirm title="删除该打卡点？" onConfirm={() => onDelete(cp.id)}>
                    <Button size="small" danger>删除</Button>
                  </Popconfirm>
                ),
              },
            ]}
          />
          <Button style={{ marginTop: 12 }} onClick={onRegistrations}>查看/审核报名</Button>
        </>
      )}
    </Modal>
  )
}
