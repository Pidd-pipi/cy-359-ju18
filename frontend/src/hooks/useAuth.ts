// 认证守卫 hook：校验登录态并支持管理员鉴权
import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuthStore } from '../stores/authStore'
import { isAdmin } from '../utils/auth'

export function useAuth(requireAdmin = false) {
  const navigate = useNavigate()
  const { token, user, init } = useAuthStore()

  useEffect(() => {
    init()
  }, [init])

  useEffect(() => {
    if (!token) {
      navigate('/login', { replace: true })
      return
    }
    if (requireAdmin && !isAdmin()) {
      navigate('/', { replace: true })
    }
  }, [token, requireAdmin, navigate])

  return { token, user, isAdmin: user?.role === 'admin' || isAdmin() }
}
