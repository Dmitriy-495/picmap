<template>
  <div class="min-h-screen bg-[#020420] flex flex-col">
    <!-- Хедер -->
    <header class="sticky top-0 z-50 backdrop-blur-xl bg-[#020420]/80 border-b border-white/10">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between h-20">
          <!-- Логотип -->
          <NuxtLink to="/" class="flex items-center gap-3 group">
            <div class="w-10 h-10 rounded-xl bg-gradient-to-br from-[#00dc82] to-[#00b86b] flex items-center justify-center shadow-lg shadow-[#00dc82]/20 group-hover:shadow-[#00dc82]/40 transition-shadow">
              <span class="text-[#020420] text-xl font-bold">🌍</span>
            </div>
            <span class="text-2xl font-bold">PicMap</span>
          </NuxtLink>

          <!-- Навигация -->
          <nav class="flex items-center gap-8">
            <NuxtLink 
              to="/" 
              class="text-sm text-gray-400 hover:text-white transition-colors"
              active-class="text-white"
            >
              Маршруты
            </NuxtLink>

            <template v-if="isAuthenticated">
              <NuxtLink 
                to="/routes/new" 
                class="text-sm text-gray-400 hover:text-white transition-colors"
              >
                + Создать
              </NuxtLink>
              <div class="flex items-center gap-3 pl-6 border-l border-white/10">
                <span class="text-sm text-gray-400">{{ user?.username }}</span>
                <button 
                  @click="handleLogout"
                  class="text-sm text-gray-400 hover:text-white transition-colors"
                >
                  Выйти
                </button>
              </div>
            </template>

            <template v-else>
              <NuxtLink 
                to="/login" 
                class="px-5 py-2.5 text-sm rounded-xl bg-[#00dc82] text-[#020420] font-semibold hover:bg-[#00c070] transition-colors"
              >
                Войти
              </NuxtLink>
            </template>
          </nav>
        </div>
      </div>
    </header>

    <!-- Контент -->
    <main class="flex-1">
      <slot />
    </main>

    <!-- Футер -->
    <footer class="border-t border-white/10 mt-24">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
        <div class="flex flex-col sm:flex-row justify-between items-center gap-6 text-sm text-gray-500">
          <div>© 2026 PicMap. Твои путешествия. Твои воспоминания.</div>
          <div class="flex items-center gap-6">
            <a href="https://github.com/Dmitriy-495/picmap" target="_blank" class="hover:text-white transition-colors">
              GitHub
            </a>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
const { isAuthenticated, user, logout } = useAuth()

const handleLogout = () => {
  logout()
  navigateTo('/login')
}
</script>
