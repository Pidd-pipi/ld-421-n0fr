import { createBrowserRouter, Navigate } from 'react-router-dom'
import { AppLayout } from '../components/layout/AppLayout'
import { AuthGuard, RoleGuard } from './guards'
import { LoginPage } from '../pages/Login'
import { Dashboard } from '../pages/Dashboard'
import { EquipmentManage } from '../pages/EquipmentManage'
import { BorrowManage } from '../pages/BorrowManage'
import { ReservationManage } from '../pages/ReservationManage'
import { MaintenanceManage } from '../pages/MaintenanceManage'
import { AuditLogPage } from '../pages/AuditLog'

export const router = createBrowserRouter([
  {
    path: '/login',
    element: <LoginPage />
  },
  {
    path: '/',
    element: (
      <AuthGuard>
        <AppLayout />
      </AuthGuard>
    ),
    children: [
      { index: true, element: <Navigate to="/dashboard" replace /> },
      { path: 'dashboard', element: <Dashboard /> },
      { path: 'equipment', element: <EquipmentManage /> },
      { path: 'borrow', element: <BorrowManage /> },
      { path: 'reservations', element: <ReservationManage /> },
      { path: 'maintenance', element: <MaintenanceManage /> },
      {
        path: 'audit-logs',
        element: (
          <RoleGuard roles={['Admin']}>
            <AuditLogPage />
          </RoleGuard>
        )
      }
    ]
  }
])
