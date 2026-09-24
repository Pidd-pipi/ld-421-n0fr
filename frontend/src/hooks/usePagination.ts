import { useMemo, useState } from 'react'

export function usePagination(initialPage = 1, initialPageSize = 10, initialTotal = 0) {
  const [page, setPage] = useState(initialPage)
  const [pageSize, setPageSize] = useState(initialPageSize)
  const [total, setTotal] = useState(initialTotal)

  const pagination = useMemo(
    () => ({
      page,
      pageSize,
      total,
      setPage,
      setPageSize,
      setTotal
    }),
    [page, pageSize, total]
  )

  return pagination
}

export default usePagination
