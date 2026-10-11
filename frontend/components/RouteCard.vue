<template>
  <NuxtLink 
    :to="`/routes/${route.id}`"
    class="group block bg-white/5 backdrop-blur-xl border border-white/10 rounded-2xl overflow-hidden hover:border-[#00dc82]/50 hover:bg-white/[0.07] transition-all duration-300"
  >
    <!-- Обложка -->
    <div class="aspect-video bg-gradient-to-br from-[#00dc82]/20 to-[#00b86b]/10 relative overflow-hidden">
      <div class="absolute inset-0 flex items-center justify-center">
        <span class="text-6xl opacity-50">🗺️</span>
      </div>
      
      <!-- Статус -->
      <div 
        v-if="route.status"
        :class="['absolute top-3 right-3 px-3 py-1 rounded-full text-xs font-medium border backdrop-blur-sm', statusBadge.color]"
      >
        {{ statusBadge.text }}
      </div>
    </div>

    <!-- Контент -->
    <div class="p-5">
      <h3 class="text-lg font-semibold mb-2 group-hover:text-[#00dc82] transition-colors line-clamp-1">
        {{ route.title }}
      </h3>
      
      <p v-if="route.description" class="text-sm text-gray-400 mb-4 line-clamp-2">
        {{ route.description }}
      </p>

      <div class="flex items-center justify-between text-xs text-gray-500">
        <div class="flex items-center gap-1">
          <span>📅</span>
          <span>{{ formatDate(route.start_date) }}</span>
        </div>
        <div v-if="route.places?.length" class="flex items-center gap-1">
          <span>📍</span>
          <span>{{ route.places.length }}</span>
        </div>
      </div>
    </div>
  </NuxtLink>
</template>

<script setup lang="ts">
const props = defineProps<{
  route: any
}>()

const { formatDate, getStatusBadge } = useRoutes()

const statusBadge = computed(() => getStatusBadge(props.route.status))
</script>
