import { create } from 'zustand'
import * as favoriteApi from '../api/favorite'
import type { Favorite } from '../api/favorite'

interface FavoriteState {
  favorites: Favorite[]
  total: number
  fetchFavorites: () => Promise<void>
}

export const useFavoriteStore = create<FavoriteState>((set) => ({
  favorites: [],
  total: 0,
  fetchFavorites: async () => {
    const res = await favoriteApi.listFavorites({ page: 1, page_size: 50 })
    set({ favorites: res.data.list, total: res.data.total })
  },
}))
