<template>
  <div class="p-10">
    <div class="mb-10">
      <h1 class="text-4xl font-bold mb-3">Обзор</h1>
      <p class="text-gray-400">Статистика PicMap</p>
    </div>

    <!-- Загрузка -->
    <div v-if="loading" class="text-center py-20">
      <div class="inline-block w-12 h-12 border-4 border-[#00dc82]/30 border-t-[#00dc82] rounded-full animate-spin"></div>
    </div>

    <!-- Статистика -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div class="bg-white/5 backdrop-blur-xl border border-white/10 rounded-2xl p-6">
        <div class="text-3xl mb-3">👥</div>
        <div class="text-3xl font-bold text-[#00dc82]">{{ stats.users }}</div>
        <div class="text-sm text-gray-400">Пользователей</div>
      </div>

      <div class="bg-white/5 backdrop-blur-xl border border-white/10 rounded-2xl p-6">
        <div class="text-3xl mb-3">🗺️</div>
        <div class="text-3xl font-bold text-[#00dc82]">{{ stats.routes }}</div>
        <div class="text-sm text-gray-400">Маршрутов</div>
      </div>

      <div class="bg-white/5 backdrop-blur-xl border border-white/10 rounded-2xl p-6">
        <div class="text-3xl mb-3">📍</div>
        <div class="text-3xl font-bold text-[#00dc82]">{{ stats.places }}</div>
        <div class="text-sm text-gray-400">Точек</div>
      </div>

      <div class="bg-white/5 backdrop-blur-xl border border-white/10 rounded-2xl p-6">
        <div class="text-3xl mb-3">💬</div>
        <div class="text-3xl font-bold text-[#00dc82]">{{ stats.comments }}</div>
        <div class="text-sm text-gray-400">Комментариев</div>
      </div>
    </div>

    <!-- Активные приглашения -->
    <div class="mt-10">
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-2xl font-bold">Последние приглашения</h2>
        <AppButton to="/administrator/invitations" variant="secondary" size="sm">
          Все приглашения →
        </AppButton>
      </div>

      <div v-if="recentInvitations.length === 0" class="text-center py-12 bg-white/5 rounded-2xl border border-white/10">
        <div class="text-4xl mb-2">🎟️</div>
        <p class="text-gray-400">Нет приглашений</p>
      </div>

      <div v-else class="bg-white/5 backdrop-blur-xl border border-white/10 rounded-2xl overflow-hidden">
        <table class="w-full">
          <thead class="bg-white/5 border-b border-white/10">
            <tr>
              <th class="text-left px-6 py-4 text-sm font-semibold text-gray-400">Токен</th>
              <th class="text-left px-6 py-4 text-sm font-semibold text-gray-400">Создано</th>
              <th class="text-left px-6 py-4 text-sm font-semibold text-gray-400">Истекает</th>
              <th class="text-left px-6 py-4 text-sm font-semibold text-gray-400">Статус</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="inv in recentInvitations" :key="inv.id" class="border-b border-white/5 last:border-0">
              <td class="px-6 py-4 text-sm font-mono text-gray-300">{{ inv.token.slice(0, 16) }}...</td>
              <td class="px-6 py-4 text-sm text-gray-400">{{ formatDate(inv.created_at) }}</td>
              <td class="px-6 py-4 text-sm text-gray-400">{{ formatDate(inv.expires_at) }}</td>
              <td class="px-6 py-4">
                <AppBadge :variant="getInvitationStatus(inv).variant">
                  {{ getInvitationStatus(inv).text }}
                </AppBadge>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'administrator',
})

const api = useApi()

const loading = ref(true)
const stats = ref({
  users: 0,
  routes: 0,
  places: 0,
  comments: 0,
})

const recentInvitations = ref<any[]>([])

onMounted(async () => {
  try {
    // Загружаем приглашения
    const invitations = await api.get<any[]>('/api/admin/invitations')
    recentInvitations.value = invitations.slice(0, 5)
  } catch (e) {
    console.error('Failed to load invitations', e)
  }

  // Статистика (заглушка — потом добавим API)
  stats.value = {
    users: 2,
    routes: 1,
    places: 3,
    comments: 1,
  }

  loading.value = false
})

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  })
}

const getInvitationStatus = (inv: any) => {
  if (inv.used_at) return { text: 'Использовано', variant: 'past' as const }
  if (new Date(inv.expires_at) < new Date()) return { text: 'Истекло', variant: 'past' as const }
  return { text: 'Активно', variant: 'now' as const }
}

useHead({
  title: 'Админка | Обзор'
})
</script>
