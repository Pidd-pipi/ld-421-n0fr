import { Alert, Button, Empty, Space, Steps, Table, Tag, Typography } from 'antd'
import type { ColumnsType } from 'antd/es/table'
import { StatusBadge } from './StatusBadge'
import type { DisposalApproval, DisposalBlocker } from '../../types'

const disposalStatusMap: Record<string, { color: string; text: string }> = {
  Pending: { color: 'gold', text: '待审批' },
  Approved: { color: 'green', text: '审批通过' },
  Returned: { color: 'volcano', text: '已退回' }
}

const blockerTypeMap: Record<string, { color: string; text: string }> = {
  Borrow: { color: 'blue', text: '未归还借用' },
  Reservation: { color: 'purple', text: '未结束预约' }
}

interface DisposalProgressProps {
  approval: DisposalApproval
  reviewing?: boolean
  canReview?: boolean
  onReview?: (approval: DisposalApproval) => void
}

export function DisposalProgress({ approval, reviewing, canReview, onReview }: DisposalProgressProps) {
  const current = approval.status === 'Approved' ? 2 : 1

  const blockerColumns: ColumnsType<DisposalBlocker> = [
    {
      title: '类型',
      dataIndex: 'type',
      width: 110,
      render: (value: string) => {
        const meta = blockerTypeMap[value] || { color: 'default', text: value }
        return <Tag color={meta.color}>{meta.text}</Tag>
      }
    },
    { title: '使用人', dataIndex: 'userName', width: 100 },
    { title: '说明', dataIndex: 'detail' },
    { title: '关联状态', dataIndex: 'status', width: 100, render: (value: string) => <StatusBadge status={value} /> }
  ]

  return (
    <div style={{ marginTop: 16 }}>
      <Space style={{ marginBottom: 12 }}>
        <Typography.Title level={5} style={{ margin: 0 }}>
          报废处置审批
        </Typography.Title>
        <StatusBadge status={approval.status} labelMap={disposalStatusMap} />
      </Space>

      <Steps
        size="small"
        current={current}
        status={approval.status === 'Returned' ? 'error' : approval.status === 'Approved' ? 'finish' : 'process'}
        items={[
          { title: '提交申请', description: approval.applicantName ? `申请人：${approval.applicantName}` : undefined },
          { title: '处置审批', description: approval.reviewerName ? `审批人：${approval.reviewerName}` : '等待管理员审批' },
          { title: '设备退役', description: approval.status === 'Approved' ? '状态已变更为 Retired' : '审批通过后退役' }
        ]}
      />

      <Typography.Paragraph type="secondary" style={{ marginTop: 12, marginBottom: 4 }}>
        报废原因：{approval.reason}
      </Typography.Paragraph>
      {approval.reviewedAt ? (
        <Typography.Paragraph type="secondary" style={{ marginBottom: 8 }}>
          审批时间：{approval.reviewedAt.replace('T', ' ').slice(0, 16)}
        </Typography.Paragraph>
      ) : null}

      {approval.status === 'Pending' && approval.blockers.length > 0 ? (
        <Alert
          type="warning"
          showIcon
          style={{ marginBottom: 8 }}
          message={`当前仍有 ${approval.blockers.length} 项未归还借用或未结束预约，审批将被退回`}
        />
      ) : null}
      {approval.status === 'Returned' ? (
        <Alert type="error" showIcon style={{ marginBottom: 8 }} message="审批因存在阻塞项被退回，设备未停用；请先归还借用或结束预约后重新提交。" />
      ) : null}
      {approval.status === 'Approved' ? (
        <Alert type="success" showIcon style={{ marginBottom: 8 }} message="审批通过，设备已报废（Retired）。原借用和预约记录均已保留。" />
      ) : null}

      {approval.blockers.length > 0 ? (
        <Table rowKey={(row) => `${row.type}-${row.recordId}`} size="small" pagination={false} columns={blockerColumns} dataSource={approval.blockers} />
      ) : approval.status === 'Pending' ? (
        <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="当前无阻塞项" />
      ) : null}

      {canReview && approval.status === 'Pending' ? (
        <Button type="primary" danger loading={reviewing} style={{ marginTop: 12 }} onClick={() => onReview?.(approval)}>
          审批处置申请
        </Button>
      ) : null}
    </div>
  )
}

export default DisposalProgress
