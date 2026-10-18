<template>
  <AppContainer>
    <div class="py-16">
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
        <AppButton variant="danger" @click="load">
          Попробовать снова
        </AppButton>
      </div>

      <!-- Пусто -->
      <AppEmpty
        v-else-if="routes.length === 0"
        icon="🗺️"
        title="Пока нет маршрутов"
        description="Создай свой первый маршрут!"
      >
        <AppButton v-if="isAuthenticated" to="/routes/new" size="lg">
          + Создать маршрут
        </AppButton>
        <AppButton v-else to="/login" size="lg">
          Войти
        </AppButton>
      </AppEmpty>

      <!-- Список маршрутов -->
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        <RouteCard
          v-for="route in routes"
          :key="route.id"
          :route="route"
        />
      </div>
    </div>
  </AppContainer>
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
