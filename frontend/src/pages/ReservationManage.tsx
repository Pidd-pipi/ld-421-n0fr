import { useCallback, useEffect, useMemo, useState } from 'react'
import { Button, Card, Col, Form, Input, Modal, Row, Select, Space, Table, message } from 'antd'
import type { ColumnsType } from 'antd/es/table'
import { PlusOutlined } from '@ant-design/icons'
import dayjs from 'dayjs'
import { approveReservation, cancelReservation, createReservation, fetchReservations, rejectReservation } from '../api/reservation'
import { fetchEquipment } from '../api/equipment'
import { StatusBadge } from '../components/common/StatusBadge'
import { CalendarCell } from '../components/common/CalendarCell'
import { usePagination } from '../hooks/usePagination'
import { useAuthStore } from '../stores/authStore'
import type { CreateReservationPayload, Equipment, Reservation } from '../types'

const canApprove = (role?: string) => role === 'Admin' || role === 'LabManager'

export function ReservationManage() {
  const [items, setItems] = useState<Reservation[]>([])
  const [equipment, setEquipment] = useState<Equipment[]>([])
  const [selectedEquipment, setSelectedEquipment] = useState<number | undefined>()
  const [modalOpen, setModalOpen] = useState(false)
  const [form] = Form.useForm<CreateReservationPayload>()
  const [loading, setLoading] = useState(false)
  const pagination = usePagination()
  const user = useAuthStore((state) => state.user)

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const result = await fetchReservations({
        page: pagination.page,
        page_size: pagination.pageSize,
        equipment_id: selectedEquipment || undefined
      })
      setItems(result.list)
      pagination.setTotal(result.total)
    } finally {
      setLoading(false)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedEquipment, pagination.page, pagination.pageSize])

  useEffect(() => {
    load()
  }, [load])

  useEffect(() => {
    fetchEquipment({ page: 1, page_size: 100 }).then((res) => setEquipment(res.list))
  }, [])

  const nextDays = useMemo(() => {
    return Array.from({ length: 7 }, (_, i) => dayjs().add(i, 'day').format('MM-DD'))
  }, [])

  const submit = async () => {
    const values = await form.validateFields()
    await createReservation({
      equipmentId: values.equipmentId,
      startTime: dayjs(values.startTime).format('YYYY-MM-DD HH:mm'),
      endTime: dayjs(values.endTime).format('YYYY-MM-DD HH:mm'),
      purpose: values.purpose
    })
    message.success('预约已创建')
    setModalOpen(false)
    load()
  }

  const columns: ColumnsType<Reservation> = [
    { title: '设备', dataIndex: 'equipmentName', width: 150 },
    { title: '预约人', dataIndex: 'userName', width: 100 },
    { title: '开始时间', dataIndex: 'startTime', width: 160, render: (v: string) => dayjs(v).format('YYYY-MM-DD HH:mm') },
    { title: '结束时间', dataIndex: 'endTime', width: 160, render: (v: string) => dayjs(v).format('YYYY-MM-DD HH:mm') },
    { title: '使用目的', dataIndex: 'purpose', ellipsis: true },
    { title: '状态', dataIndex: 'status', width: 100, render: (v: string) => <StatusBadge status={v} /> },
    {
      title: '操作',
      width: 220,
      render: (_, record) => (
        <Space>
          {record.status === 'Pending' && canApprove(user?.roleCode) ? (
            <>
              <Button
                size="small"
                type="primary"
                onClick={async () => {
                  await approveReservation(record.id)
                  message.success('已通过')
                  load()
                }}
              >
                通过
              </Button>
              <Button
                size="small"
                danger
                onClick={async () => {
                  await rejectReservation(record.id)
                  message.success('已驳回')
                  load()
                }}
              >
                驳回
              </Button>
            </>
          ) : null}
          {(record.status === 'Pending' || record.status === 'Approved') &&
          (record.userId === user?.id || canApprove(user?.roleCode)) ? (
            <Button
              size="small"
              type="link"
              onClick={async () => {
                await cancelReservation(record.id)
                message.success('已取消')
                load()
              }}
            >
              取消
            </Button>
          ) : null}
        </Space>
      )
    }
  ]

  return (
    <Card
      title="预约管理"
      extra={
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)}>
          创建预约
        </Button>
      }
    >
      <Space style={{ marginBottom: 16 }} wrap>
        <Select
          style={{ width: 240 }}
          allowClear
          placeholder="按设备筛选"
          value={selectedEquipment}
          onChange={(value) => {
            setSelectedEquipment(value as number | undefined)
            pagination.setPage(1)
          }}
          options={equipment.map((item) => ({ label: `${item.name}（${item.code}）`, value: item.id }))}
        />
        <span style={{ color: '#666' }}>设备预约日历视图（未来 7 天）</span>
      </Space>

      <Row gutter={6} style={{ marginBottom: 16 }}>
        {nextDays.map((day) => (
          <Col span={3} key={day}>
            <CalendarCell date={day}>
              {items
                .filter((item) => dayjs(item.startTime).format('MM-DD') === day)
                .slice(0, 2)
                .map((item) => (
                  <div key={item.id} style={{ fontSize: 10, color: '#1677ff' }}>
                    {item.equipmentName}
                  </div>
                ))}
            </CalendarCell>
          </Col>
        ))}
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

      <Modal title="创建预约" open={modalOpen} onCancel={() => setModalOpen(false)} onOk={submit}>
        <Form form={form} layout="vertical">
          <Form.Item name="equipmentId" label="设备" rules={[{ required: true }]}>
            <Select
              placeholder="选择设备"
              options={equipment.map((item) => ({ label: `${item.name}（${item.code}）`, value: item.id }))}
            />
          </Form.Item>
          <Form.Item name="startTime" label="开始时间" rules={[{ required: true }]}>
            <Input placeholder="YYYY-MM-DD HH:mm" />
          </Form.Item>
          <Form.Item name="endTime" label="结束时间" rules={[{ required: true }]}>
            <Input placeholder="YYYY-MM-DD HH:mm" />
          </Form.Item>
          <Form.Item name="purpose" label="使用目的">
            <Input.TextArea rows={3} />
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  )
}

export default ReservationManage
