// 统一请求封装：axios 实例 + 拦截器（对应后端统一响应 {code,message,data}）
import axios, { AxiosError } from 'axios'
import { message } from 'antd'

export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

const request = axios.create({
  baseURL: '/api/v1',
  timeout: 15000,
})

request.interceptors.request.use((config) => {
  const token = localStorage.getItem('orienteering_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

request.interceptors.response.use(
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  (response): any => {
    const body = response.data as ApiResponse
    if (body.code !== 0) {
      message.error(body.message || '请求失败')
      return Promise.reject(new Error(body.message))
    }
    return body
  },
  (error: AxiosError<ApiResponse>) => {
    const body = error.response?.data
    const msg = body?.message || error.message || '网络错误'
    if (error.response?.status === 401) {
      localStorage.removeItem('orienteering_token')
      localStorage.removeItem('orienteering_user')
      if (!window.location.pathname.startsWith('/login')) {
        window.location.href = '/login'
      }
    }
    message.error(msg)
    return Promise.reject(new Error(msg))
  },
)

export default request
