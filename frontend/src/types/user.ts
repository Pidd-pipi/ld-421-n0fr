export interface Role {
  id: number
  code: string
  name: string
  description?: string
}

export interface User {
  id: number
  username: string
  name: string
  email: string
  phone: string
  roleId: number
  roleCode: string
  roleName: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface PageResult<T> {
  list: T[]
  total: number
  page: number
  pageSize: number
}
