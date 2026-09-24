import { get, post } from '../utils/request'
import { API_PATHS } from '../constants/apiPaths'
import type { LoginRequest, LoginResponse, User } from '../types'

export function login(payload: LoginRequest) {
  return post<LoginResponse>(API_PATHS.login, payload)
}

export function me() {
  return get<{ userId: number; username: string; role: string }>(API_PATHS.me)
}

export function fetchUsers(params?: Record<string, unknown>) {
  return get<{ list: User[]; total: number; page: number; pageSize: number }>(API_PATHS.users, params)
}
