import type { ReactNode } from 'react'
import { Card } from 'antd'

interface StatCardProps {
  title: string
  value: ReactNode
  icon?: ReactNode
  color?: string
  footer?: ReactNode
}

export function StatCard({ title, value, icon, color, footer }: StatCardProps) {
  return (
    <Card size="small" styles={{ body: { padding: 16 } }}>
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
        <div>
          <div style={{ color: '#666', fontSize: 13 }}>{title}</div>
          <div style={{ fontSize: 26, fontWeight: 600, color: color || '#111' }}>{value}</div>
        </div>
        {icon ? <div style={{ fontSize: 30, color: color || '#1677ff' }}>{icon}</div> : null}
      </div>
      {footer ? <div style={{ marginTop: 12, color: '#999', fontSize: 12 }}>{footer}</div> : null}
    </Card>
  )
}

export default StatCard
