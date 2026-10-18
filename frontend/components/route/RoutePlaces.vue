<template>
  <div>
    <h2 class="text-3xl font-bold mb-8">Точки маршрута</h2>
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div
        v-for="(place, index) in places"
        :key="place.id"
        class="bg-white/5 backdrop-blur-xl border border-white/10 rounded-2xl p-6 hover:border-[#00dc82]/50 hover:bg-white/[0.07] transition-all"
      >
        <div class="flex items-start gap-4">
          <div class="text-4xl flex-shrink-0">{{ place.emoji || '📍' }}</div>
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 mb-2">
              <span class="text-xs text-[#00dc82] font-bold bg-[#00dc82]/10 px-2 py-0.5 rounded">
                {{ index + 1 }}
              </span>
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
</template>

<script setup lang="ts">
defineProps<{
  places: any[]
}>()
</script>
