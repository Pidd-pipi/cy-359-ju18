// 格式化工具（对应后端 internal/util/formatters.go）
import dayjs from 'dayjs'

export function formatDateTime(value?: string | null): string {
  if (!value) return '-'
  return dayjs(value).format('YYYY-MM-DD HH:mm:ss')
}

export function formatDuration(totalSeconds?: number): string {
  if (!totalSeconds || totalSeconds <= 0) return '--:--:--'
  const h = Math.floor(totalSeconds / 3600)
  const m = Math.floor((totalSeconds % 3600) / 60)
  const s = totalSeconds % 60
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${pad(h)}:${pad(m)}:${pad(s)}`
}

export function maskPhone(phone?: string): string {
  if (!phone || phone.length !== 11) return phone || '-'
  return `${phone.slice(0, 3)}****${phone.slice(7)}`
}
