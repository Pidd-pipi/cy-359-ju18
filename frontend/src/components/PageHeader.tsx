// 共享组件：页面标题
import { Typography } from 'antd'

interface Props {
  title: string
  extra?: React.ReactNode
  children?: React.ReactNode
}

export default function PageHeader({ title, extra, children }: Props) {
  return (
    <div style={{ marginBottom: 16 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography.Title level={4} style={{ margin: 0 }}>
          {title}
        </Typography.Title>
        {extra}
      </div>
      {children && <div style={{ marginTop: 8 }}>{children}</div>}
    </div>
  )
}
