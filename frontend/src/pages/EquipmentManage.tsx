import { useCallback, useEffect, useMemo, useState } from 'react'
import {
  Button,
  Card,
  Col,
  Descriptions,
  Drawer,
  Form,
  Input,
  InputNumber,
  Modal,
  Row,
  Select,
  Space,
  Table,
  TreeSelect,
  message
} from 'antd'
import type { ColumnsType } from 'antd/es/table'
import type { DataNode } from 'antd/es/tree'
import { PlusOutlined, SearchOutlined } from '@ant-design/icons'
import dayjs from 'dayjs'
import { fetchCategories } from '../api/category'
import {
  createEquipment,
  fetchEquipment,
  fetchEquipmentDetail,
  retireEquipment,
  transferOwner,
  updateEquipment
} from '../api/equipment'
import { fetchBorrows } from '../api/borrow'
import { fetchMaintenance } from '../api/maintenance'
import { fetchUsers } from '../api/auth'
import type { BorrowRecord, Equipment, EquipmentCategory, EquipmentPayload, MaintenanceRecord, User } from '../types'
import { StatusBadge } from '../components/common/StatusBadge'
import { CalendarCell } from '../components/common/CalendarCell'
import { formatCurrency } from '../utils/formatCurrency'
import { usePagination } from '../hooks/usePagination'
import { useAuthStore } from '../stores/authStore'

const canManage = (role?: string) => role === 'Admin' || role === 'LabManager'

export function EquipmentManage() {
  const [items, setItems] = useState<Equipment[]>([])
  const [categories, setCategories] = useState<EquipmentCategory[]>([])
  const [users, setUsers] = useState<User[]>([])
  const [keyword, setKeyword] = useState('')
  const [categoryId, setCategoryId] = useState<number | undefined>()
  const [loading, setLoading] = useState(false)
  const [detail, setDetail] = useState<Equipment | null>(null)
  const [borrowHistory, setBorrowHistory] = useState<BorrowRecord[]>([])
  const [maintenanceList, setMaintenanceList] = useState<MaintenanceRecord[]>([])
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Equipment | null>(null)
  const [form] = Form.useForm<EquipmentPayload>()
  const pagination = usePagination()
  const user = useAuthStore((state) => state.user)

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const result = await fetchEquipment({
        page: pagination.page,
        page_size: pagination.pageSize,
        keyword: keyword || undefined,
        category_id: categoryId
      })
      setItems(result.list)
      pagination.setTotal(result.total)
    } finally {
      setLoading(false)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [keyword, categoryId, pagination.page, pagination.pageSize])

  useEffect(() => {
    load()
  }, [load])

  useEffect(() => {
    fetchCategories().then(setCategories).catch(() => setCategories([]))
  }, [])

  useEffect(() => {
    if (canManage(user?.roleCode)) {
      fetchUsers({ page: 1, page_size: 100 }).then((res) => setUsers(res.list)).catch(() => setUsers([]))
    }
  }, [user])

  const treeData = useMemo(() => {
    const map = (list: EquipmentCategory[]): DataNode[] =>
      list.map((item) => ({
        key: item.id,
        title: item.name,
        value: item.id,
        children: item.children && item.children.length ? map(item.children) : undefined
      }))
    return map(categories)
  }, [categories])

  const openDetail = async (id: number) => {
    const data = await fetchEquipmentDetail(id)
    setDetail(data)
    const [borrows, maintenances] = await Promise.all([
      fetchBorrows({ equipment_id: id, page: 1, page_size: 20 }),
      fetchMaintenance({ equipment_id: id, page: 1, page_size: 20 })
    ])
    setBorrowHistory(borrows.list)
    setMaintenanceList(maintenances.list)
  }

  const openCreate = () => {
    setEditing(null)
    form.resetFields()
    setModalOpen(true)
  }

  const openEdit = (record: Equipment) => {
    setEditing(record)
    form.setFieldsValue({
      name: record.name,
      code: record.code,
      categoryId: record.categoryId,
      brandModel: record.brandModel,
      serialNumber: record.serialNumber,
      purchaseDate: record.purchaseDate?.slice(0, 10),
      purchasePrice: record.purchasePrice,
      location: record.location,
      ownerId: record.ownerId,
      supplier: record.supplier,
      warrantyExpiry: record.warrantyExpiry?.slice(0, 10),
      imageUrl: record.imageUrl
    })
    setModalOpen(true)
  }

  const submit = async () => {
    const values = await form.validateFields()
    const payload: EquipmentPayload = {
      ...values,
      purchaseDate: values.purchaseDate ? dayjs(values.purchaseDate).format('YYYY-MM-DD') : undefined,
      warrantyExpiry: values.warrantyExpiry ? dayjs(values.warrantyExpiry).format('YYYY-MM-DD') : undefined
    }
    if (editing) {
      await updateEquipment(editing.id, payload)
      message.success('设备已更新')
    } else {
      await createEquipment(payload)
      message.success('设备已登记')
    }
    setModalOpen(false)
    load()
  }

  const onRetire = async (id: number) => {
    await retireEquipment(id)
    message.success('设备已报废')
    load()
  }

  const onTransfer = async (id: number) => {
    let ownerId = 0
    Modal.confirm({
      title: '转移责任人',
      content: (
        <Select
          style={{ width: '100%', marginTop: 12 }}
          placeholder="选择新责任人"
          onChange={(value) => (ownerId = value as number)}
          options={users.map((u) => ({ label: `${u.name}（${u.username}）`, value: u.id }))}
        />
      ),
      onOk: async () => {
        if (!ownerId) return
        await transferOwner(id, ownerId)
        message.success('责任人已转移')
        load()
      }
    })
  }

  const columns: ColumnsType<Equipment> = [
    { title: '设备编号', dataIndex: 'code', width: 120 },
    { title: '设备名称', dataIndex: 'name' },
    { title: '分类', dataIndex: 'categoryName', width: 120 },
    { title: '品牌型号', dataIndex: 'brandModel', width: 160, ellipsis: true },
    { title: '状态', dataIndex: 'status', width: 100, render: (value: string) => <StatusBadge status={value} /> },
    { title: '存放位置', dataIndex: 'location', width: 140 },
    { title: '责任人', dataIndex: 'ownerName', width: 100 },
    { title: '价格', dataIndex: 'purchasePrice', width: 120, render: (value: number) => formatCurrency(value) },
    {
      title: '操作',
      width: 220,
      render: (_, record) => (
        <Space>
          <Button size="small" type="link" onClick={() => openDetail(record.id)}>
            详情
          </Button>
          {canManage(user?.roleCode) ? (
            <>
              <Button size="small" type="link" onClick={() => openEdit(record)}>
                编辑
              </Button>
              <Button size="small" type="link" danger onClick={() => onRetire(record.id)}>
                报废
              </Button>
              <Button size="small" type="link" onClick={() => onTransfer(record.id)}>
                转移
              </Button>
            </>
          ) : null}
        </Space>
      )
    }
  ]

  const nextDays = useMemo(() => {
    return Array.from({ length: 7 }, (_, i) => dayjs().add(i, 'day').format('MM-DD'))
  }, [])

  return (
    <Card
      title="设备管理"
      extra={
        canManage(user?.roleCode) ? (
          <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
            新增设备
          </Button>
        ) : null
      }
    >
      <Space style={{ marginBottom: 16 }} wrap>
        <Input
          prefix={<SearchOutlined />}
          placeholder="搜索名称/编号/型号"
          allowClear
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          onPressEnter={() => pagination.setPage(1)}
        />
        <TreeSelect
          style={{ width: 220 }}
          allowClear
          placeholder="选择分类"
          treeData={treeData}
          value={categoryId}
          onChange={(value) => {
            setCategoryId(value as number | undefined)
            pagination.setPage(1)
          }}
        />
        <Button
          onClick={() => {
            pagination.setPage(1)
            load()
          }}
        >
          查询
        </Button>
      </Space>
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

      <Modal
        title={editing ? '编辑设备' : '新增设备'}
        open={modalOpen}
        onCancel={() => setModalOpen(false)}
        onOk={submit}
        width={640}
      >
        <Form form={form} layout="vertical" initialValues={{ purchasePrice: 0 }}>
          <Row gutter={12}>
            <Col span={12}>
              <Form.Item name="name" label="设备名称" rules={[{ required: true }]}>
                <Input />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="code" label="设备编号" rules={[{ required: true }]}>
                <Input />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="categoryId" label="分类" rules={[{ required: true }]}>
                <TreeSelect treeData={treeData} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="ownerId" label="责任人" rules={[{ required: true }]}>
                <Select options={users.map((u) => ({ label: u.name, value: u.id }))} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="brandModel" label="品牌型号">
                <Input />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="serialNumber" label="序列号">
                <Input />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="location" label="存放位置">
                <Input placeholder="楼栋-房间-柜号" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="supplier" label="供应商">
                <Input />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="purchasePrice" label="购买价格">
                <InputNumber style={{ width: '100%' }} min={0} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="purchaseDate" label="购买日期">
                <Input placeholder="YYYY-MM-DD" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="warrantyExpiry" label="保修到期日">
                <Input placeholder="YYYY-MM-DD" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="imageUrl" label="设备图片 URL">
                <Input />
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </Modal>

      <Drawer title="设备详情" open={!!detail} onClose={() => setDetail(null)} width={760}>
        {detail ? (
          <>
            <Descriptions column={2} bordered size="small">
              <Descriptions.Item label="设备名称">{detail.name}</Descriptions.Item>
              <Descriptions.Item label="设备编号">{detail.code}</Descriptions.Item>
              <Descriptions.Item label="分类">{detail.categoryName}</Descriptions.Item>
              <Descriptions.Item label="状态">
                <StatusBadge status={detail.status} />
              </Descriptions.Item>
              <Descriptions.Item label="品牌型号">{detail.brandModel}</Descriptions.Item>
              <Descriptions.Item label="序列号">{detail.serialNumber}</Descriptions.Item>
              <Descriptions.Item label="存放位置">{detail.location}</Descriptions.Item>
              <Descriptions.Item label="责任人">{detail.ownerName}</Descriptions.Item>
              <Descriptions.Item label="供应商">{detail.supplier}</Descriptions.Item>
              <Descriptions.Item label="购买价格">{formatCurrency(detail.purchasePrice)}</Descriptions.Item>
              <Descriptions.Item label="购买日期">{detail.purchaseDate?.slice(0, 10)}</Descriptions.Item>
              <Descriptions.Item label="保修到期日">{detail.warrantyExpiry?.slice(0, 10)}</Descriptions.Item>
            </Descriptions>

            <h4 style={{ marginTop: 16 }}>未来 7 天预约</h4>
            <Row gutter={6}>
              {nextDays.map((day) => (
                <Col span={3} key={day}>
                  <CalendarCell date={day} muted={day === dayjs().format('MM-DD') ? false : true} />
                </Col>
              ))}
            </Row>

            <h4 style={{ marginTop: 16 }}>借用历史</h4>
            <Table
              rowKey="id"
              size="small"
              pagination={false}
              dataSource={borrowHistory}
              columns={[
                { title: '借用人', dataIndex: 'borrowerName', width: 90 },
                { title: '借用日期', dataIndex: 'borrowDate', render: (v: string) => v?.slice(0, 10) },
                { title: '状态', dataIndex: 'status', render: (v: string) => <StatusBadge status={v} /> }
              ]}
            />

            <h4 style={{ marginTop: 16 }}>维护记录</h4>
            <Table
              rowKey="id"
              size="small"
              pagination={false}
              dataSource={maintenanceList}
              columns={[
                { title: '类型', dataIndex: 'type', width: 110 },
                { title: '内容', dataIndex: 'content' },
                { title: '结果', dataIndex: 'result', render: (v: string) => <StatusBadge status={v} /> }
              ]}
            />
          </>
        ) : null}
      </Drawer>
    </Card>
  )
}

export default EquipmentManage
