// 共享组件：状态徽标，按域+状态映射颜色与文案
import { Badge } from 'antd'
import { getStatusInfo, type StatusDomain } from '../constants'

interface Props {
  status: string
  domain?: StatusDomain
  fallback?: string
}

export default function StatusBadge({ status, domain = 'activity', fallback = '未知' }: Props) {
  const cfg = getStatusInfo(domain, status)
  return <Badge status={(cfg?.color as any) || 'default'} text={cfg?.text || fallback} />
}
