// 与后端 internal/constants/enums.go 一一对应的业务枚举

export const RoleType = {
  USER: 'user',
  ADMIN: 'admin',
} as const

export const ActivityStatus = {
  DRAFT: 'draft',
  PUBLISHED: 'published',
  ONGOING: 'ongoing',
  FINISHED: 'finished',
  CANCELLED: 'cancelled',
} as const

export const Difficulty = {
  FAMILY: 'family',
  ADULT: 'adult',
  PRO: 'pro',
} as const

export const TaskType = {
  NONE: 'none',
  QUIZ: 'quiz',
  PHOTO: 'photo',
} as const

export const CheckinType = {
  GPS: 'gps',
  QRCODE: 'qrcode',
} as const

export const RegistrationStatus = {
  PENDING: 'pending',
  WAITLIST: 'waitlist',
  APPROVED: 'approved',
  REJECTED: 'rejected',
  FINISHED: 'finished',
} as const

export const ProductType = {
  EQUIPMENT: 'equipment',
  COUPON: 'coupon',
  BADGE: 'badge',
} as const

export const ProductStatus = {
  ON: 'on',
  OFF: 'off',
} as const

export const RedemptionStatus = {
  PENDING: 'pending',
  COMPLETED: 'completed',
  CANCELLED: 'cancelled',
} as const

export const TeamRole = {
  CAPTAIN: 'captain',
  MEMBER: 'member',
} as const

export const CheckinResult = {
  CORRECT: 'correct',
  WRONG: 'wrong',
  NONE: 'none',
} as const

// 状态 -> 徽标颜色/文案（与后端 util/formatters.go StatusColor 保持一致）
// 不同域之间可能存在相同状态字符串（如 finished/pending/cancelled），故按域拆分。
export type StatusDomain = 'activity' | 'registration' | 'redemption' | 'product' | 'checkin'

export const statusConfig: Record<StatusDomain, Record<string, { color: string; text: string }>> = {
  activity: {
    [ActivityStatus.DRAFT]: { color: 'default', text: '草稿' },
    [ActivityStatus.PUBLISHED]: { color: 'processing', text: '已发布' },
    [ActivityStatus.ONGOING]: { color: 'success', text: '进行中' },
    [ActivityStatus.FINISHED]: { color: 'cyan', text: '已结束' },
    [ActivityStatus.CANCELLED]: { color: 'error', text: '已取消' },
  },
  registration: {
    [RegistrationStatus.PENDING]: { color: 'default', text: '待审核' },
    [RegistrationStatus.WAITLIST]: { color: 'warning', text: '候补' },
    [RegistrationStatus.APPROVED]: { color: 'processing', text: '已通过' },
    [RegistrationStatus.REJECTED]: { color: 'error', text: '已拒绝' },
    [RegistrationStatus.FINISHED]: { color: 'success', text: '已完成' },
  },
  redemption: {
    [RedemptionStatus.PENDING]: { color: 'default', text: '待处理' },
    [RedemptionStatus.COMPLETED]: { color: 'cyan', text: '已完成' },
    [RedemptionStatus.CANCELLED]: { color: 'error', text: '已取消' },
  },
  product: {
    [ProductStatus.ON]: { color: 'processing', text: '上架中' },
    [ProductStatus.OFF]: { color: 'default', text: '已下架' },
  },
  checkin: {
    [CheckinResult.CORRECT]: { color: 'success', text: '正确' },
    [CheckinResult.WRONG]: { color: 'error', text: '错误' },
    [CheckinResult.NONE]: { color: 'default', text: '无任务' },
  },
}

export function getStatusInfo(domain: StatusDomain, status: string) {
  return statusConfig[domain]?.[status]
}

export const difficultyConfig: Record<string, string> = {
  [Difficulty.FAMILY]: '亲子',
  [Difficulty.ADULT]: '成人',
  [Difficulty.PRO]: '专业',
}

export const taskTypeConfig: Record<string, string> = {
  [TaskType.NONE]: '无任务',
  [TaskType.QUIZ]: '答题',
  [TaskType.PHOTO]: '拍照',
}

export const productTypeConfig: Record<string, string> = {
  [ProductType.EQUIPMENT]: '户外装备',
  [ProductType.COUPON]: '活动优惠券',
  [ProductType.BADGE]: '虚拟勋章',
}
