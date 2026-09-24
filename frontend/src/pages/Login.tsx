import { useState } from 'react'
import { Button, Card, Form, Input, Typography, message } from 'antd'
import { LockOutlined, UserOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { useAuthStore } from '../stores/authStore'
import type { LoginRequest } from '../types'

export function LoginPage() {
  const navigate = useNavigate()
  const login = useAuthStore((state) => state.login)
  const [loading, setLoading] = useState(false)

  const onFinish = async (values: LoginRequest) => {
    setLoading(true)
    try {
      await login(values)
      message.success('登录成功')
      navigate('/dashboard')
    } catch {
      // 错误提示已由请求拦截器处理
    } finally {
      setLoading(false)
    }
  }

  return (
    <div
      style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: 'linear-gradient(135deg,#1677ff 0%,#06b6d4 100%)'
      }}
    >
      <Card style={{ width: 380 }} title="实验室设备管理系统">
        <Form<LoginRequest> initialValues={{ username: 'admin', password: 'admin123' }} onFinish={onFinish}>
          <Form.Item name="username" rules={[{ required: true, message: '请输入用户名' }]}>
            <Input prefix={<UserOutlined />} placeholder="用户名" size="large" />
          </Form.Item>
          <Form.Item name="password" rules={[{ required: true, message: '请输入密码' }]}>
            <Input.Password prefix={<LockOutlined />} placeholder="密码" size="large" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block size="large" loading={loading}>
              登录
            </Button>
          </Form.Item>
        </Form>
        <Typography.Paragraph type="secondary" style={{ fontSize: 12, marginBottom: 0 }}>
          预置账号：admin/admin123、labmanager/lab123、researcher/res123、student/stu123
        </Typography.Paragraph>
      </Card>
    </div>
  )
}

export default LoginPage
