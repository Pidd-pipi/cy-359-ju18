import { create } from 'zustand'
import * as productApi from '../api/product'
import type { Product, Redemption } from '../api/product'

interface ProductState {
  products: Product[]
  total: number
  loading: boolean
  myRedemptions: Redemption[]
  fetchProducts: (params?: { page?: number; page_size?: number; status?: string; keyword?: string }) => Promise<void>
  fetchMyRedemptions: () => Promise<void>
}

export const useProductStore = create<ProductState>((set) => ({
  products: [],
  total: 0,
  loading: false,
  myRedemptions: [],
  fetchProducts: async (params = {}) => {
    set({ loading: true })
    try {
      const res = await productApi.listProducts({ page: 1, page_size: 12, ...params })
      set({ products: res.data.list, total: res.data.total })
    } finally {
      set({ loading: false })
    }
  },
  fetchMyRedemptions: async () => {
    const res = await productApi.listMyRedemptions({ page: 1, page_size: 50 })
    set({ myRedemptions: res.data.list })
  },
}))
