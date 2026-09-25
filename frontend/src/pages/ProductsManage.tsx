import { useEffect, useState } from 'react'
import { Button, Col, Form, Input, InputNumber, Modal, Popconfirm, Row, Select, Space, Table, message } from 'antd'
import * as productApi from '../api/product'
import PageHeader from '../components/PageHeader'
import StatusBadge from '../components/StatusBadge'
import { useAuth } from '../hooks/useAuth'
import { ProductStatus, ProductType } from '../constants'

export default function ProductsManage() {
  useAuth(true)
  const [list, setList] = useState<productApi.Product[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [open, setOpen] = useState(false)
  const [editing, setEditing] = useState<productApi.Product | null>(null)
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)

  const load = async (p = page) => {
    setLoading(true)
    try {
      const res = await productApi.listProducts({ page: p, page_size: 10 })
      setList(res.data.list)
      setTotal(res.data.total)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [page])

  const openCreate = () => {
    setEditing(null)
    form.resetFields()
    form.setFieldsValue({ type: ProductType.EQUIPMENT, status: ProductStatus.ON, stock: 10 })
    setOpen(true)
  }

  const openEdit = (p: productApi.Product) => {
    setEditing(p)
    form.setFieldsValue(p)
    setOpen(true)
  }

  const onSave = async () => {
    const values = await form.validateFields()
    setLoading(true)
    try {
      if (editing) await productApi.updateProduct(editing.id, values)
      else await productApi.createProduct(values)
      message.success('保存成功')
      setOpen(false)
      load()
    } finally {
      setLoading(false)
    }
  }

  const onStatus = async (p: productApi.Product) => {
    await productApi.setProductStatus(p.id, p.status === ProductStatus.ON ? ProductStatus.OFF : ProductStatus.ON)
    message.success('已更新')
    load()
  }

  return (
    <div>
      <PageHeader
        title="商品管理（管理员）"
        extra={<Button type="primary" onClick={openCreate}>新建商品</Button>}
      />
      <Table<productApi.Product>
        rowKey="id"
        size="small"
        loading={loading}
        dataSource={list}
        pagination={{ current: page, pageSize: 10, total, showTotal: (t) => `共 ${t} 条`, onChange: setPage }}
        columns={[
          { title: 'ID', dataIndex: 'id', width: 60 },
          { title: '名称', dataIndex: 'name' },
          { title: '类型', dataIndex: 'type' },
          { title: '积分', dataIndex: 'points_cost', width: 80 },
          { title: '库存', dataIndex: 'stock', width: 70 },
          { title: '状态', dataIndex: 'status', width: 90, render: (v: string) => <StatusBadge status={v} domain="product" /> },
          {
            title: '操作',
            width: 220,
            render: (_, p) => (
              <Space>
                <Button size="small" onClick={() => openEdit(p)}>编辑</Button>
                <Button size="small" onClick={() => onStatus(p)}>
                  {p.status === ProductStatus.ON ? '下架' : '上架'}
                </Button>
              </Space>
            ),
          },
        ]}
      />
      <Modal title={editing ? '编辑商品' : '新建商品'} open={open} onCancel={() => setOpen(false)} onOk={onSave} confirmLoading={loading}>
        <Form form={form} layout="vertical">
          <Row gutter={12}>
            <Col span={12}>
              <Form.Item name="name" label="名称" rules={[{ required: true }]}>
                <Input />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="type" label="类型" rules={[{ required: true }]}>
                <Select options={Object.entries(ProductType).map(([k, v]) => ({ label: k, value: v }))} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="points_cost" label="所需积分" rules={[{ required: true }]}>
                <InputNumber min={1} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="stock" label="库存" rules={[{ required: true }]}>
                <InputNumber min={0} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={24}>
              <Form.Item name="status" label="状态" rules={[{ required: true }]}>
                <Select options={Object.entries(ProductStatus).map(([k, v]) => ({ label: k, value: v }))} />
              </Form.Item>
            </Col>
            <Col span={24}>
              <Form.Item name="image_url" label="图片 URL">
                <Input />
              </Form.Item>
            </Col>
            <Col span={24}>
              <Form.Item name="description" label="描述">
                <Input.TextArea rows={2} />
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </Modal>
    </div>
  )
}
