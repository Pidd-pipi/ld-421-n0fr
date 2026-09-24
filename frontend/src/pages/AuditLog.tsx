import { useCallback, useEffect, useState } from 'react'
import { Card, Input, Space, Table } from 'antd'
import type { ColumnsType } from 'antd/es/table'
import { SearchOutlined } from '@ant-design/icons'
import { fetchAuditLogs, type AuditLog } from '../api/audit'
import { usePagination } from '../hooks/usePagination'

export function AuditLogPage() {
  const [items, setItems] = useState<AuditLog[]>([])
  const [keyword, setKeyword] = useState('')
  const [loading, setLoading] = useState(false)
  const pagination = usePagination()

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const result = await fetchAuditLogs({
        page: pagination.page,
        page_size: pagination.pageSize,
        keyword: keyword || undefined
      })
      setItems(result.list)
      pagination.setTotal(result.total)
    } finally {
      setLoading(false)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [keyword, pagination.page, pagination.pageSize])

  useEffect(() => {
    load()
  }, [load])

  const columns: ColumnsType<AuditLog> = [
    { title: 'ID', dataIndex: 'id', width: 70 },
    { title: '操作人', dataIndex: 'userName', width: 120 },
    { title: '动作', dataIndex: 'action', width: 220 },
    { title: '资源类型', dataIndex: 'resourceType', width: 120 },
    { title: '资源 ID', dataIndex: 'resourceId', width: 90 },
    { title: '详情', dataIndex: 'detail', ellipsis: true },
    { title: 'IP', dataIndex: 'ip', width: 130 },
    { title: '时间', dataIndex: 'createdAt', width: 180, render: (v: string) => new Date(v).toLocaleString() }
  ]

  return (
    <Card title="操作日志">
      <Space style={{ marginBottom: 16 }}>
        <Input
          prefix={<SearchOutlined />}
          placeholder="搜索动作/操作人/资源"
          allowClear
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          onPressEnter={() => pagination.setPage(1)}
        />
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
    </Card>
  )
}

export default AuditLogPage
