export default defineNuxtRouteMiddleware((to) => {
  // Только на клиенте (у нас ssr: false)
  if (import.meta.server) return

  const token = localStorage.getItem('picmap_token')

  // Публичные страницы (без авторизации)
  const publicPages = ['/login', '/invite']
  const isPublicPage = publicPages.some(page => to.path.startsWith(page))

  // Если не авторизован и идёт на защищённую страницу → /login
  if (!token && !isPublicPage) {
    return navigateTo('/login')
  }

  // Если авторизован и идёт на /login → /
  if (token && to.path === '/login') {
    return navigateTo('/')
  }
})
