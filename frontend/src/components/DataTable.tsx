// 共享组件：通用分页表格封装
import { Table } from 'antd'
import type { TableProps } from 'antd'

interface Props<T> {
  columns: TableProps<T>['columns']
  dataSource: T[]
  rowKey: string | ((record: T) => string | number)
  loading?: boolean
  total?: number
  page?: number
  pageSize?: number
  onPageChange?: (page: number, pageSize: number) => void
  scroll?: { x?: number; y?: number }
}

export default function DataTable<T extends object>({
  columns,
  dataSource,
  rowKey,
  loading,
  total,
  page,
  pageSize,
  onPageChange,
  scroll,
}: Props<T>) {
  return (
    <Table<T>
      rowKey={rowKey}
      columns={columns}
      dataSource={dataSource}
      loading={loading}
      scroll={scroll}
      pagination={
        total === undefined
          ? false
          : {
              current: page || 1,
              pageSize: pageSize || 10,
              total,
              showSizeChanger: true,
              showTotal: (t) => `共 ${t} 条`,
              onChange: onPageChange,
            }
      }
    />
  )
}
