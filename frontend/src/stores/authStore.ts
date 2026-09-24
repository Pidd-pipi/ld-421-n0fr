import { create } from 'zustand'
import { login as loginApi } from '../api/auth'
import { TOKEN_KEY } from '../utils/request'
import type { LoginRequest, User } from '../types'

export const USER_KEY = 'lab_equipment_user'

function readUser(): User | null {
  const raw = localStorage.getItem(USER_KEY)
  if (!raw) return null
  try {
    return JSON.parse(raw) as User
  } catch {
    return null
  }
}

interface AuthState {
  token: string | null
  user: User | null
  login: (payload: LoginRequest) => Promise<void>
  logout: () => void
}

export const useAuthStore = create<AuthState>((set) => ({
  token: localStorage.getItem(TOKEN_KEY),
  user: readUser(),
  login: async (payload) => {
    const result = await loginApi(payload)
    localStorage.setItem(TOKEN_KEY, result.token)
    localStorage.setItem(USER_KEY, JSON.stringify(result.user))
    set({ token: result.token, user: result.user })
  },
  logout: () => {
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USER_KEY)
    set({ token: null, user: null })
  }
}))
