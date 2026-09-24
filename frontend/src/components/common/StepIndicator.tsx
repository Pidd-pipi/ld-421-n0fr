import { Steps } from 'antd'

interface StepIndicatorProps {
  current: number
  steps?: string[]
}

const defaultSteps = ['提交申请', '审批', '使用中', '确认归还']

export function StepIndicator({ current, steps = defaultSteps }: StepIndicatorProps) {
  return <Steps size="small" current={current} items={steps.map((title) => ({ title }))} />
}

export default StepIndicator
