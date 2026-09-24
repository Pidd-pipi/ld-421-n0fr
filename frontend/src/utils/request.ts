import axios, { AxiosError, type AxiosRequestConfig } from 'axios'
import { message } from 'antd'

interface ApiEnvelope<T> {
  code: number
  message: string
  data: T
}

export const TOKEN_KEY = 'lab_equipment_token'

const instance = axios.create({
  baseURL: '/',
  timeout: 15000
})

instance.interceptors.request.use((config) => {
  const token = localStorage.getItem(TOKEN_KEY)
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

instance.interceptors.response.use(
  (response) => response,
  (error: AxiosError<ApiEnvelope<unknown>>) => {
    const status = error.response?.status
    const msg = error.response?.data?.message || error.message || '网络错误'
    if (status === 401) {
      localStorage.removeItem(TOKEN_KEY)
      message.error('登录已过期，请重新登录')
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    } else if (status === 403) {
      message.error('无权限执行此操作')
    } else {
      message.error(msg)
    }
    return Promise.reject(error)
  }
)

async function request<T>(config: AxiosRequestConfig): Promise<T> {
  const response = await instance.request<ApiEnvelope<T>>(config)
  const payload = response.data
  if (payload && typeof payload === 'object' && 'code' in payload) {
    if (payload.code === 0) {
      return payload.data
    }
    message.error(payload.message || '请求失败')
    throw new Error(payload.message || '请求失败')
  }
  return payload as unknown as T
}

export const get = <T>(url: string, params?: Record<string, unknown>) =>
  request<T>({ method: 'GET', url, params })

export const post = <T>(url: string, data?: unknown) =>
  request<T>({ method: 'POST', url, data })

export const put = <T>(url: string, data?: unknown) =>
  request<T>({ method: 'PUT', url, data })

export const patch = <T>(url: string, data?: unknown) =>
  request<T>({ method: 'PATCH', url, data })

export const del = <T>(url: string) => request<T>({ method: 'DELETE', url })

export default instance
