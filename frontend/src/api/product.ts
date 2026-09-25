import request from '../utils/request'

export interface Product {
  id: number
  name: string
  description?: string
  type: string
  points_cost: number
  stock: number
  status: string
  image_url?: string
  created_at: string
}

export interface ProductPayload {
  name: string
  description?: string
  type: string
  points_cost: number
  stock: number
  status: string
  image_url?: string
}

export interface Redemption {
  id: number
  user_id: number
  username: string
  product_id: number
  product_name: string
  quantity: number
  points_cost: number
  total_points: number
  status: string
  redeemed_at: string
  completed_at?: string
}

export function listProducts(params: { page?: number; page_size?: number; status?: string; keyword?: string }) {
  return request.get<{ list: Product[]; total: number }>('/products', { params })
}

export function getProduct(id: number) {
  return request.get<Product>(`/products/${id}`)
}

export function createProduct(data: ProductPayload) {
  return request.post('/products', data)
}

export function updateProduct(id: number, data: ProductPayload) {
  return request.put(`/products/${id}`, data)
}

export function setProductStatus(id: number, status: string) {
  return request.patch(`/products/${id}/status`, null, { params: { status } })
}

export function redeem(data: { product_id: number; quantity: number }) {
  return request.post('/redemptions', data)
}

export function listMyRedemptions(params: { page?: number; page_size?: number }) {
  return request.get<{ list: Redemption[]; total: number }>('/redemptions/mine', { params })
}

export function listAllRedemptions(params: { page?: number; page_size?: number; status?: string }) {
  return request.get<{ list: Redemption[]; total: number }>('/redemptions', { params })
}

export function updateRedemptionStatus(id: number, status: string) {
  return request.put(`/redemptions/${id}/status`, { status })
}
