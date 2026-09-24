import { Tag } from 'antd'

interface StatusBadgeProps {
  status: string
  labelMap?: Record<string, { color: string; text?: string }>
}

const defaultMap: Record<string, { color: string; text?: string }> = {
  Available: { color: 'green', text: '可用' },
  InUse: { color: 'blue', text: '使用中' },
  Maintenance: { color: 'orange', text: '维护中' },
  Retired: { color: 'default', text: '已报废' },
  Lost: { color: 'red', text: '已丢失' },
  Pending: { color: 'gold', text: '待审批' },
  Approved: { color: 'green', text: '已通过' },
  Rejected: { color: 'red', text: '已驳回' },
  Returned: { color: 'cyan', text: '已归还' },
  Overdue: { color: 'volcano', text: '已逾期' },
  Cancelled: { color: 'default', text: '已取消' },
  Preventive: { color: 'blue', text: '预防性维护' },
  Corrective: { color: 'orange', text: '纠正性维护' },
  Calibration: { color: 'purple', text: '校准' },
  Cleaning: { color: 'cyan', text: '清洁' },
  Pass: { color: 'green', text: '通过' },
  Fail: { color: 'red', text: '失败' },
  NeedsFollowUp: { color: 'gold', text: '需跟进' },
  Good: { color: 'green', text: '良好' },
  Damaged: { color: 'orange', text: '损坏' }
}

export function StatusBadge({ status, labelMap }: StatusBadgeProps) {
  const meta = (labelMap && labelMap[status]) || defaultMap[status] || { color: 'default' }
  return <Tag color={meta.color}>{meta.text || status}</Tag>
}

export default StatusBadge
