import type { ReactNode } from 'react'

interface CalendarCellProps {
  date: string | number
  children?: ReactNode
  highlight?: boolean
  muted?: boolean
  onClick?: () => void
}

export function CalendarCell({ date, children, highlight, muted, onClick }: CalendarCellProps) {
  return (
    <div
      onClick={onClick}
      style={{
        minHeight: 56,
        padding: 6,
        border: '1px solid #f0f0f0',
        borderRadius: 6,
        cursor: onClick ? 'pointer' : 'default',
        background: highlight ? '#e6f4ff' : muted ? '#fafafa' : '#fff',
        opacity: muted ? 0.5 : 1,
        overflow: 'hidden'
      }}
    >
      <div style={{ fontSize: 12, color: '#999', marginBottom: 4 }}>{date}</div>
      {children}
    </div>
  )
}

export default CalendarCell
