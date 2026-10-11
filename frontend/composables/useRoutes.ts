export interface Route {
  id: number
  owner_id: number
  title: string
  description: string
  start_date: string
  end_date: string
  cover_photo_id?: number | null
  created_at: string
  status?: string
  places?: Place[]
}

export interface Place {
  id: number
  route_id: number
  name: string
  lat: number
  lng: number
  description: string
  emoji: string
  order_index: number
  created_at: string
}

export interface CreateRouteInput {
  title: string
  description?: string
  start_date: string
  end_date: string
}

export const useRoutes = () => {
  const api = useApi()

  const routes = ref<Route[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const fetchRoutes = async () => {
    loading.value = true
    error.value = null
    try {
      const data = await api.get<Route[]>('/api/routes')
      routes.value = data || []
      return routes.value
    } catch (e: any) {
      error.value = e.message || 'Ошибка загрузки маршрутов'
      console.error('fetchRoutes error:', e)
      return []
    } finally {
      loading.value = false
    }
  }

  const fetchRoute = async (id: number) => {
    loading.value = true
    error.value = null
    try {
      const data = await api.get<Route>(`/api/routes/${id}`)
      return data
    } catch (e: any) {
      error.value = e.message || 'Маршрут не найден'
      console.error('fetchRoute error:', e)
      return null
    } finally {
      loading.value = false
    }
  }

  const createRoute = async (input: CreateRouteInput) => {
    const data = await api.post<Route>('/api/routes', input)
    routes.value = [data, ...routes.value]
    return data
  }

  const deleteRoute = async (id: number) => {
    await api.delete(`/api/routes/${id}`)
    routes.value = routes.value.filter(r => r.id !== id)
  }

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr)
    return date.toLocaleDateString('ru-RU', {
      day: 'numeric',
      month: 'long',
      year: 'numeric'
    })
  }

  const getStatusBadge = (status?: string) => {
    switch (status) {
      case 'будет':
        return { text: 'Будет', color: 'bg-blue-500/20 text-blue-400 border-blue-500/30' }
      case 'сейчас':
        return { text: 'Сейчас', color: 'bg-green-500/20 text-green-400 border-green-500/30' }
      case 'было':
        return { text: 'Было', color: 'bg-gray-500/20 text-gray-400 border-gray-500/30' }
      default:
        return { text: status || '—', color: 'bg-gray-500/20 text-gray-400 border-gray-500/30' }
    }
  }

  return {
    routes,
    loading,
    error,
    fetchRoutes,
    fetchRoute,
    createRoute,
    deleteRoute,
    formatDate,
    getStatusBadge,
  }
}
