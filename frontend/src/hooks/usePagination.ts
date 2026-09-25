// 分页 hook：封装页码/页大小状态与变更
import { useState } from 'react'

export function usePagination(initialPage = 1, initialPageSize = 10) {
  const [page, setPage] = useState(initialPage)
  const [pageSize, setPageSize] = useState(initialPageSize)

  const reset = () => {
    setPage(1)
  }

  const onChange = (p: number, ps: number) => {
    setPage(p)
    setPageSize(ps)
  }

  return { page, pageSize, setPage, setPageSize, reset, onChange }
}
