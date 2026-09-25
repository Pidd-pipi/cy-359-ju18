// 共享组件：报名状态标签。待审核/候补/已通过等状态用 StatusBadge 展示，
// 候补额外显示“前面还有 N 队”。团队详情页与活动管理审核弹窗复用。
import { Tag, Tooltip } from 'antd'
import StatusBadge from './StatusBadge'
import type { Registration } from '../api/team'
import { RegistrationStatus } from '../constants'

export default function RegistrationStatusTag({ status, waitlistAhead }: {
  status: string
  waitlistAhead?: number
}) {
  if (status === RegistrationStatus.WAITLISTED) {
    return (
      <Tooltip title="名额不足，按报名提交顺序排队；前面队伍被拒绝后会自动递补为待审核">
        <span>
          <StatusBadge status={status} domain="registration" />
          <Tag color="orange" style={{ marginLeft: 4 }}>
            前面 {waitlistAhead ?? 0} 队
          </Tag>
        </span>
      </Tooltip>
    )
  }
  return <StatusBadge status={status} domain="registration" />
}

export function registrationAhead(r: Pick<Registration, 'waitlist_ahead'>): number {
  return r.waitlist_ahead ?? 0
}
