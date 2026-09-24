import type { ReactNode } from 'react'
import { Navigate, useLocation } from 'react-router-dom'
import { useAuthStore } from '../stores/authStore'

export function AuthGuard({ children }: { children: ReactNode }) {
  const token = useAuthStore((state) => state.token)
  const location = useLocation()
  if (!token) {
    return <Navigate to="/login" replace state={{ from: location.pathname }} />
  }
  return <>{children}</>
}

export function RoleGuard({ roles, children }: { roles: string[]; children: ReactNode }) {
  const user = useAuthStore((state) => state.user)
  if (!user || !roles.includes(user.roleCode)) {
    return <Navigate to="/dashboard" replace />
  }
  return <>{children}</>
}
