import { ConfigProvider } from 'antd'
import zhCN from 'antd/locale/zh_CN'
import { RouterProvider } from 'react-router-dom'
import { ErrorBoundary } from './components/common/ErrorBoundary'
import { router } from './router'

export default function App() {
  return (
    <ConfigProvider locale={zhCN}>
      <ErrorBoundary>
        <RouterProvider router={router} />
      </ErrorBoundary>
    </ConfigProvider>
  )
}
