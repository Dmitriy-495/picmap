export default defineNuxtRouteMiddleware((to) => {
  if (import.meta.server) return

  // Проверяем только /administrator
  if (!to.path.startsWith('/administrator')) return

  const token = localStorage.getItem('picmap_token')
  const userStr = localStorage.getItem('picmap_user')

  // Не авторизован
  if (!token || !userStr) {
    return navigateTo('/login')
  }

  try {
    const user = JSON.parse(userStr)
    // Только owner
    if (!user.is_owner) {
      return navigateTo('/')
    }
  } catch {
    return navigateTo('/login')
  }
})
