import { Empty, Button } from 'antd'

interface EmptyStateProps {
  description?: string
  actionText?: string
  onAction?: () => void
}

export function EmptyState({ description = '暂无数据', actionText, onAction }: EmptyStateProps) {
  return (
    <Empty
      description={description}
      style={{ padding: 40 }}
    >
      {actionText && onAction ? (
        <Button type="primary" onClick={onAction}>
          {actionText}
        </Button>
      ) : null}
    </Empty>
  )
}

export default EmptyState
