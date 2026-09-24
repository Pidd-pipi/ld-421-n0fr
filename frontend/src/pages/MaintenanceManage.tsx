import { useCallback, useEffect, useMemo, useState } from 'react'
import { Button, Card, Col, Form, Input, InputNumber, Modal, Row, Select, Space, Table, message } from 'antd'
import type { ColumnsType } from 'antd/es/table'
import { PlusOutlined, ToolOutlined, DollarOutlined } from '@ant-design/icons'
import dayjs from 'dayjs'
import { createMaintenance, executeMaintenance, fetchMaintenance, fetchMaintenanceStats } from '../api/maintenance'
import { fetchEquipment } from '../api/equipment'
import { fetchUsers } from '../api/auth'
import { StatusBadge } from '../components/common/StatusBadge'
import { StatCard } from '../components/common/StatCard'
import { CalendarCell } from '../components/common/CalendarCell'
import { formatCurrency } from '../utils/formatCurrency'
import { usePagination } from '../hooks/usePagination'
import { useAuthStore } from '../stores/authStore'
import type { CreateMaintenancePayload, Equipment, MaintenanceRecord, MaintenanceStats, User } from '../types'

const canWrite = (role?: string) => role === 'Admin' || role === 'LabManager' || role === 'Researcher'

export function MaintenanceManage() {
  const [items, setItems] = useState<MaintenanceRecord[]>([])
  const [equipment, setEquipment] = useState<Equipment[]>([])
  const [users, setUsers] = useState<User[]>([])
  const [stats, setStats] = useState<MaintenanceStats | null>(null)
  const [modalOpen, setModalOpen] = useState(false)
  const [form] = Form.useForm<CreateMaintenancePayload>()
  const [loading, setLoading] = useState(false)
  const pagination = usePagination()
  const user = useAuthStore((state) => state.user)

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const result = await fetchMaintenance({ page: pagination.page, page_size: pagination.pageSize })
      setItems(result.list)
      pagination.setTotal(result.total)
      const stat = await fetchMaintenanceStats()
      setStats(stat)
    } finally {
      setLoading(false)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pagination.page, pagination.pageSize])

  useEffect(() => {
    load()
  }, [load])

  useEffect(() => {
    fetchEquipment({ page: 1, page_size: 100 }).then((res) => setEquipment(res.list))
    if (canWrite(user?.roleCode)) {
      fetchUsers({ page: 1, page_size: 100 }).then((res) => setUsers(res.list)).catch(() => setUsers([]))
    }
  }, [user])

  const nextDays = useMemo(() => {
    return Array.from({ length: 7 }, (_, i) => dayjs().add(i, 'day').format('MM-DD'))
  }, [])

  const submit = async () => {
    const values = await form.validateFields()
    await createMaintenance({
      equipmentId: values.equipmentId,
      type: values.type,
      content: values.content,
      maintenanceDate: dayjs(values.maintenanceDate).format('YYYY-MM-DD'),
      nextMaintenanceDate: values.nextMaintenanceDate
        ? dayjs(values.nextMaintenanceDate).format('YYYY-MM-DD')
        : undefined,
      cost: values.cost ?? 0,
      maintainerId: values.maintainerId
    })
    message.success('维护记录已创建')
    setModalOpen(false)
    load()
  }

  const execute = (id: number) => {
    let result = 'Pass'
    Modal.confirm({
      title: '执行维护并记录结果',
      content: (
        <Select
          style={{ width: '100%', marginTop: 12 }}
          defaultValue="Pass"
          onChange={(value) => (result = value)}
          options={[
            { label: '通过', value: 'Pass' },
            { label: '失败', value: 'Fail' },
            { label: '需跟进', value: 'NeedsFollowUp' }
          ]}
        />
      ),
      onOk: async () => {
        await executeMaintenance(id, result)
        message.success('维护结果已记录')
        load()
      }
    })
  }

  const columns: ColumnsType<MaintenanceRecord> = [
    { title: '设备', dataIndex: 'equipmentName', width: 150 },
    { title: '类型', dataIndex: 'type', width: 120, render: (v: string) => <StatusBadge status={v} /> },
    { title: '维护内容', dataIndex: 'content', ellipsis: true },
    { title: '维护日期', dataIndex: 'maintenanceDate', width: 110, render: (v: string) => v?.slice(0, 10) },
    { title: '下次维护', dataIndex: 'nextMaintenanceDate', width: 110, render: (v?: string) => v?.slice(0, 10) },
    { title: '费用', dataIndex: 'cost', width: 110, render: (v: number) => formatCurrency(v) },
    { title: '维护人', dataIndex: 'maintainerName', width: 100 },
    { title: '结果', dataIndex: 'result', width: 100, render: (v: string) => <StatusBadge status={v} /> },
    {
      title: '操作',
      width: 130,
      render: (_, record) =>
        canWrite(user?.roleCode) ? (
          <Button size="small" type="link" onClick={() => execute(record.id)}>
            执行维护
          </Button>
        ) : null
    }
  ]

  return (
    <div>
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={8}>
          <StatCard title="维护记录数" value={stats?.totalRecords ?? 0} icon={<ToolOutlined />} color="#1677ff" />
        </Col>
        <Col xs={24} sm={8}>
          <StatCard title="累计维护费用" value={formatCurrency(stats?.totalCost)} icon={<DollarOutlined />} color="#fa8c16" />
        </Col>
        <Col xs={24} sm={8}>
          <StatCard title="需跟进维护" value={stats?.resultCount?.NeedsFollowUp ?? 0} color="#f5222d" />
        </Col>
      </Row>

      <Card
        title="维护管理"
        style={{ marginTop: 16 }}
        extra={
          canWrite(user?.roleCode) ? (
            <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)}>
              创建维护记录
            </Button>
          ) : null
        }
      >
        <Row gutter={6} style={{ marginBottom: 16 }}>
          {nextDays.map((day) => {
            const upcoming = items.filter((item) => item.nextMaintenanceDate?.slice(5, 10) === day)
            return (
              <Col span={3} key={day}>
                <CalendarCell date={day} highlight={upcoming.length > 0}>
                  {upcoming.map((item) => (
                    <div key={item.id} style={{ fontSize: 10, color: '#fa8c16' }}>
                      {item.equipmentName}
                    </div>
                  ))}
                </CalendarCell>
              </Col>
            )
          })}
        </Row>

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
      </Card>

      <Modal title="创建维护记录" open={modalOpen} onCancel={() => setModalOpen(false)} onOk={submit}>
        <Form form={form} layout="vertical">
          <Form.Item name="equipmentId" label="设备" rules={[{ required: true }]}>
            <Select
              placeholder="选择设备"
              options={equipment.map((item) => ({ label: `${item.name}（${item.code}）`, value: item.id }))}
            />
          </Form.Item>
          <Form.Item name="type" label="维护类型" rules={[{ required: true }]}>
            <Select
              options={[
                { label: '预防性维护', value: 'Preventive' },
                { label: '纠正性维护', value: 'Corrective' },
                { label: '校准', value: 'Calibration' },
                { label: '清洁', value: 'Cleaning' }
              ]}
            />
          </Form.Item>
          <Form.Item name="maintenanceDate" label="维护日期" rules={[{ required: true }]}>
            <Input placeholder="YYYY-MM-DD" />
          </Form.Item>
          <Form.Item name="nextMaintenanceDate" label="下次维护日期">
            <Input placeholder="YYYY-MM-DD" />
          </Form.Item>
          <Form.Item name="maintainerId" label="维护人" rules={[{ required: true }]}>
            <Select options={users.map((u) => ({ label: u.name, value: u.id }))} />
          </Form.Item>
          <Form.Item name="cost" label="维护费用">
            <InputNumber style={{ width: '100%' }} min={0} />
          </Form.Item>
          <Form.Item name="content" label="维护内容">
            <Input.TextArea rows={3} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default MaintenanceManage
