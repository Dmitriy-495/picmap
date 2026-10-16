<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
    <!-- Заголовок -->
    <div class="mb-14">
      <h1 class="text-5xl font-bold mb-3">Маршруты</h1>
      <p class="text-lg text-gray-400">Твои путешествия и воспоминания</p>
    </div>

    <!-- Загрузка -->
    <div v-if="loading" class="text-center py-24">
      <div class="inline-block w-14 h-14 border-4 border-[#00dc82]/30 border-t-[#00dc82] rounded-full animate-spin"></div>
      <p class="text-gray-400 mt-6">Загрузка маршрутов...</p>
    </div>

    <!-- Ошибка -->
    <div v-else-if="error" class="bg-red-500/10 border border-red-500/30 rounded-2xl p-8 text-center">
      <p class="text-red-400 text-lg mb-6">{{ error }}</p>
      <button 
        @click="load"
        class="px-6 py-3 text-sm bg-red-500/20 text-red-400 rounded-xl hover:bg-red-500/30 transition-colors"
      >
        Попробовать снова
      </button>
    </div>

    <!-- Пусто -->
    <div v-else-if="routes.length === 0" class="text-center py-24">
      <div class="text-7xl mb-6">🗺️</div>
      <h2 class="text-3xl font-bold mb-3">Пока нет маршрутов</h2>
      <p class="text-gray-400 mb-10 text-lg">Создай свой первый маршрут!</p>
      <NuxtLink 
        v-if="isAuthenticated"
        to="/routes/new"
        class="inline-flex items-center gap-2 px-8 py-4 bg-[#00dc82] text-[#020420] font-semibold rounded-xl hover:bg-[#00c070] transition-colors"
      >
        + Создать маршрут
      </NuxtLink>
      <NuxtLink 
        v-else
        to="/login"
        class="inline-flex items-center gap-2 px-8 py-4 bg-[#00dc82] text-[#020420] font-semibold rounded-xl hover:bg-[#00c070] transition-colors"
      >
        Войти
      </NuxtLink>
    </div>

    <!-- Список маршрутов -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
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
