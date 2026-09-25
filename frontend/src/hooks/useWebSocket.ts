// 排行榜 WebSocket hook：订阅活动实时榜单更新
import { useEffect, useRef, useState } from 'react'
import { getToken } from '../utils/auth'
import type { LeaderboardRow } from '../api/checkin'

export function useWebSocket(activityId?: number) {
  const [rows, setRows] = useState<LeaderboardRow[] | null>(null)
  const [connected, setConnected] = useState(false)
  const wsRef = useRef<WebSocket | null>(null)

  useEffect(() => {
    if (!activityId) return
    const token = getToken()
    const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
    const url = `${protocol}://${window.location.host}/ws/leaderboard?activity_id=${activityId}&token=${token || ''}`
    const ws = new WebSocket(url)
    wsRef.current = ws
    ws.onopen = () => setConnected(true)
    ws.onclose = () => setConnected(false)
    ws.onerror = () => setConnected(false)
    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data)
        if (msg.type === 'leaderboard' || msg.type === 'leaderboard_update') {
          if (msg.data) {
            setRows(msg.data)
          }
        }
      } catch {
        // ignore
      }
    }
    return () => {
      ws.close()
      wsRef.current = null
    }
  }, [activityId])

  return { rows, connected }
}
