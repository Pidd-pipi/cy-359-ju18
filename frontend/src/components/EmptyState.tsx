// 共享组件：空状态占位
import { Empty, Button } from 'antd'

interface Props {
  description?: string
  actionText?: string
  onAction?: () => void
}

export default function EmptyState({ description = '暂无数据', actionText, onAction }: Props) {
  return (
    <Empty description={description} style={{ padding: '32px 0' }}>
      {actionText && onAction && (
        <Button type="primary" onClick={onAction}>
          {actionText}
        </Button>
      )}
    </Empty>
  )
}
