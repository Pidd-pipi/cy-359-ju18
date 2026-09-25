// 报名状态徽标：待审核/候补/已通过等；候补额外显示前面还有几队
import StatusBadge from './StatusBadge'
import { RegistrationStatus } from '../constants'

interface Props {
  status: string
  waitlistAhead?: number
}

export default function RegistrationStatusTag({ status, waitlistAhead }: Props) {
  if (status === RegistrationStatus.WAITLIST) {
    return (
      <span>
        <StatusBadge status={status} domain="registration" />
        <span style={{ marginLeft: 8, color: '#d48806' }}>
          前面 {waitlistAhead ?? 0} 队
        </span>
      </span>
    )
  }
  return <StatusBadge status={status} domain="registration" />
}
