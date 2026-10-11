<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
    <div v-if="loading" class="text-center py-20">
      <div class="inline-block w-12 h-12 border-4 border-[#00dc82]/30 border-t-[#00dc82] rounded-full animate-spin"></div>
      <p class="text-gray-400 mt-4">Загрузка маршрута...</p>
    </div>

    <div v-else-if="!route" class="text-center py-20">
      <div class="text-6xl mb-4">🚫</div>
      <h2 class="text-2xl font-bold mb-2">Маршрут не найден</h2>
      <NuxtLink to="/" class="text-[#00dc82] hover:underline">
        ← На главную
      </NuxtLink>
    </div>

    <template v-else>
      <div class="mb-6 text-sm">
        <NuxtLink to="/" class="text-gray-400 hover:text-white transition-colors">
          Маршруты
        </NuxtLink>
        <span class="text-gray-600 mx-2">/</span>
        <span class="text-white">{{ route.title }}</span>
      </div>

      <div class="mb-8">
        <div class="flex flex-wrap items-start justify-between gap-4 mb-4">
          <div>
            <h1 class="text-4xl font-bold mb-2">{{ route.title }}</h1>
            <p v-if="route.description" class="text-gray-400">{{ route.description }}</p>
          </div>
          <div :class="['px-4 py-2 rounded-full text-sm font-medium border', statusBadge.color]">
            {{ statusBadge.text }}
          </div>
        </div>

        <div class="flex flex-wrap items-center gap-4 text-sm text-gray-400">
          <div class="flex items-center gap-2">
            <span>📅</span>
            <span>{{ formatDate(route.start_date) }} — {{ formatDate(route.end_date) }}</span>
          </div>
          <div v-if="route.places?.length" class="flex items-center gap-2">
            <span>📍</span>
            <span>{{ route.places.length }} точек маршрута</span>
          </div>
        </div>
      </div>

      <div v-if="route.places?.length" class="mb-8">
        <MapView :places="route.places" />
      </div>

      <div v-if="route.places?.length" class="mb-12">
        <h2 class="text-2xl font-bold mb-4">Точки маршрута</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div 
            v-for="(place, index) in route.places" 
            :key="place.id"
            class="bg-white/5 backdrop-blur-xl border border-white/10 rounded-xl p-5 hover:border-[#00dc82]/50 transition-all"
          >
            <div class="flex items-start gap-3">
              <div class="text-3xl">{{ place.emoji || '📍' }}</div>
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2 mb-1">
                  <span class="text-xs text-[#00dc82] font-bold">{{ index + 1 }}</span>
                  <h3 class="font-semibold truncate">{{ place.name }}</h3>
                </div>
                <p v-if="place.description" class="text-sm text-gray-400 line-clamp-2">
                  {{ place.description }}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="text-center py-12 bg-white/5 rounded-2xl border border-white/10">
        <div class="text-4xl mb-2">📍</div>
        <p class="text-gray-400">В этом маршруте пока нет точек</p>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
const routeParams = useRoute()
const { fetchRoute, formatDate, getStatusBadge } = useRoutes()

const route = ref<any>(null)
const loading = ref(true)

const statusBadge = computed(() => getStatusBadge(route.value?.status))

onMounted(async () => {
  const id = Number(routeParams.params.id)
  route.value = await fetchRoute(id)
  loading.value = false
})

useHead({
  title: route.value?.title || 'Маршрут'
})
</script>
