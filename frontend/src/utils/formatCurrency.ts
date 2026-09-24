export function formatCurrency(value: number | null | undefined): string {
  const num = Number(value ?? 0)
  return new Intl.NumberFormat('zh-CN', {
    style: 'currency',
    currency: 'CNY',
    minimumFractionDigits: 2
  }).format(num)
}
