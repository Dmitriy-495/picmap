<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
    <!-- Заголовок -->
    <div class="mb-10">
      <h1 class="text-4xl font-bold mb-2">Маршруты</h1>
      <p class="text-gray-400">Твои путешествия и воспоминания</p>
    </div>

    <!-- Загрузка -->
    <div v-if="loading" class="text-center py-20">
      <div class="inline-block w-12 h-12 border-4 border-[#00dc82]/30 border-t-[#00dc82] rounded-full animate-spin"></div>
      <p class="text-gray-400 mt-4">Загрузка маршрутов...</p>
    </div>

    <!-- Ошибка -->
    <div v-else-if="error" class="bg-red-500/10 border border-red-500/30 rounded-xl p-6 text-center">
      <p class="text-red-400">{{ error }}</p>
      <button 
        @click="load"
        class="mt-4 px-4 py-2 text-sm bg-red-500/20 text-red-400 rounded-lg hover:bg-red-500/30 transition-colors"
      >
        Попробовать снова
      </button>
    </div>

    <!-- Пусто -->
    <div v-else-if="routes.length === 0" class="text-center py-20">
      <div class="text-6xl mb-4">🗺️</div>
      <h2 class="text-2xl font-bold mb-2">Пока нет маршрутов</h2>
      <p class="text-gray-400 mb-6">Создай свой первый маршрут!</p>
      <NuxtLink 
        v-if="isAuthenticated"
        to="/routes/new"
        class="inline-flex items-center gap-2 px-6 py-3 bg-[#00dc82] text-[#020420] font-semibold rounded-lg hover:bg-[#00c070] transition-colors"
      >
        + Создать маршрут
      </NuxtLink>
      <NuxtLink 
        v-else
        to="/login"
        class="inline-flex items-center gap-2 px-6 py-3 bg-[#00dc82] text-[#020420] font-semibold rounded-lg hover:bg-[#00c070] transition-colors"
      >
        Войти
      </NuxtLink>
    </div>

    <!-- Список маршрутов -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <RouteCard 
        v-for="route in routes" 
        :key="route.id" 
        :route="route" 
      />
    </div>
  </div>
</template>

<script setup lang="ts">
const { routes, loading, error, fetchRoutes } = useRoutes()
const { isAuthenticated } = useAuth()

const load = async () => {
  await fetchRoutes()
}

onMounted(() => {
  load()
})

useHead({
  title: 'Маршруты'
})
</script>
