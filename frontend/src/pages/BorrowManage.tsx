import { useCallback, useEffect, useState } from 'react'
import { Button, Card, Form, Input, Modal, Select, Space, Table, Tabs, message } from 'antd'
import type { ColumnsType } from 'antd/es/table'
import { PlusOutlined } from '@ant-design/icons'
import dayjs from 'dayjs'
import { createBorrow, fetchBorrows } from '../api/borrow'
import { fetchEquipment } from '../api/equipment'
import { useBorrowFlow } from '../hooks/useBorrowFlow'
import { StatusBadge } from '../components/common/StatusBadge'
import { StepIndicator } from '../components/common/StepIndicator'
import { usePagination } from '../hooks/usePagination'
import type { BorrowRecord, CreateBorrowPayload, Equipment } from '../types'
import { useAuthStore } from '../stores/authStore'

const canApprove = (role?: string) => role === 'Admin' || role === 'LabManager'

export function BorrowManage() {
  const [items, setItems] = useState<BorrowRecord[]>([])
  const [equipment, setEquipment] = useState<Equipment[]>([])
  const [status, setStatus] = useState<string>('')
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [form] = Form.useForm<CreateBorrowPayload>()
  const pagination = usePagination()
  const user = useAuthStore((state) => state.user)
  const flow = useBorrowFlow()

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const result = await fetchBorrows({
        page: pagination.page,
        page_size: pagination.pageSize,
        status: status || undefined
      })
      setItems(result.list)
      pagination.setTotal(result.total)
    } finally {
      setLoading(false)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [status, pagination.page, pagination.pageSize])

  useEffect(() => {
    load()
  }, [load])

  useEffect(() => {
    fetchEquipment({ page: 1, page_size: 100 }).then((res) => setEquipment(res.list))
  }, [])

  const submit = async () => {
    const values = await form.validateFields()
    await flow.submit({
      equipmentId: values.equipmentId,
      borrowDate: dayjs(values.borrowDate).format('YYYY-MM-DD'),
      expectedReturnDate: dayjs(values.expectedReturnDate).format('YYYY-MM-DD'),
      reason: values.reason
    })
    message.success('借用申请已提交')
    setModalOpen(false)
    load()
  }

  const approve = async (id: number) => {
    await flow.approve(id)
    message.success('已审批通过')
    load()
  }

  const reject = async (id: number) => {
    await flow.reject(id)
    message.success('已驳回')
    load()
  }

  const confirmReturn = (id: number) => {
    let condition = 'Good'
    Modal.confirm({
      title: '确认归还',
      content: (
        <Select
          style={{ width: '100%', marginTop: 12 }}
          defaultValue="Good"
          onChange={(value) => (condition = value)}
          options={[
            { label: '良好', value: 'Good' },
            { label: '损坏', value: 'Damaged' },
            { label: '丢失', value: 'Lost' }
          ]}
        />
      ),
      onOk: async () => {
        await flow.confirmReturn(id, {
          actualReturnDate: dayjs().format('YYYY-MM-DD'),
          returnCondition: condition
        })
        message.success('归还已确认')
        load()
      }
    })
  }

  const columns: ColumnsType<BorrowRecord> = [
    { title: '设备', dataIndex: 'equipmentName', width: 160 },
    { title: '借用人', dataIndex: 'borrowerName', width: 100 },
    { title: '借用日期', dataIndex: 'borrowDate', width: 110, render: (v: string) => v?.slice(0, 10) },
    { title: '预计归还', dataIndex: 'expectedReturnDate', width: 110, render: (v: string) => v?.slice(0, 10) },
    { title: '事由', dataIndex: 'reason', ellipsis: true },
    { title: '状态', dataIndex: 'status', width: 110, render: (v: string) => <StatusBadge status={v} /> },
    { title: '审批人', dataIndex: 'approverName', width: 100 },
    {
      title: '流程',
      width: 210,
      render: (_, record) => <StepIndicator current={flow.stepOf(record.status)} />
    },
    {
      title: '操作',
      width: 220,
      render: (_, record) => (
        <Space>
          {record.status === 'Pending' && canApprove(user?.roleCode) ? (
            <>
              <Button size="small" type="primary" onClick={() => approve(record.id)}>
                通过
              </Button>
              <Button size="small" danger onClick={() => reject(record.id)}>
                驳回
              </Button>
            </>
          ) : null}
          {(record.status === 'Approved' || record.status === 'Overdue') && canApprove(user?.roleCode) ? (
            <Button size="small" type="link" onClick={() => confirmReturn(record.id)}>
              确认归还
            </Button>
          ) : null}
        </Space>
      )
    }
  ]

  return (
    <Card
      title="借用管理"
      extra={
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)}>
          提交借用
        </Button>
      }
    >
      <Tabs
        activeKey={status}
        onChange={(key) => {
          setStatus(key)
          pagination.setPage(1)
        }}
        items={[
          { key: '', label: '全部' },
          { key: 'Pending', label: '待审批' },
          { key: 'Approved', label: '已通过' },
          { key: 'Rejected', label: '已驳回' },
          { key: 'Returned', label: '已归还' },
          { key: 'Overdue', label: '已逾期' }
        ]}
      />
      <Table
        rowKey="id"
        loading={loading}
        columns={columns}
        dataSource={items}
        pagination={{
          current: pagination.page,
          pageSize: pagination.pageSize,
          total: pagination.total,
          showSizeChanger: true,
          onChange: (page, pageSize) => {
            pagination.setPage(page)
            pagination.setPageSize(pageSize)
          }
        }}
      />

      <Modal title="提交借用申请" open={modalOpen} onCancel={() => setModalOpen(false)} onOk={submit}>
        <Form form={form} layout="vertical">
          <Form.Item name="equipmentId" label="选择设备" rules={[{ required: true }]}>
            <Select
              placeholder="选择可用设备"
              options={equipment
                .filter((item) => item.status === 'Available')
                .map((item) => ({ label: `${item.name}（${item.code}）`, value: item.id }))}
            />
          </Form.Item>
          <Form.Item name="borrowDate" label="借用日期" rules={[{ required: true }]}>
            <Input placeholder="YYYY-MM-DD" />
          </Form.Item>
          <Form.Item name="expectedReturnDate" label="预计归还日期" rules={[{ required: true }]}>
            <Input placeholder="YYYY-MM-DD" />
          </Form.Item>
          <Form.Item name="reason" label="借用事由">
            <Input.TextArea rows={3} />
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  )
}

export default BorrowManage
