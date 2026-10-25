<template>
  <div class="p-10">
    <div class="mb-10">
      <h1 class="text-4xl font-bold mb-3">Пользователи</h1>
      <p class="text-gray-400">Управление пользователями PicMap</p>
    </div>

    <!-- Загрузка -->
    <div v-if="loading" class="text-center py-20">
      <div class="inline-block w-12 h-12 border-4 border-[#00dc82]/30 border-t-[#00dc82] rounded-full animate-spin"></div>
    </div>

    <!-- Таблица -->
    <div v-else class="bg-white/5 backdrop-blur-xl border border-white/10 rounded-2xl overflow-hidden">
      <table class="w-full">
        <thead class="bg-white/5 border-b border-white/10">
          <tr>
            <th class="text-left px-6 py-4 text-sm font-semibold text-gray-400">ID</th>
            <th class="text-left px-6 py-4 text-sm font-semibold text-gray-400">Username</th>
            <th class="text-left px-6 py-4 text-sm font-semibold text-gray-400">Email</th>
            <th class="text-left px-6 py-4 text-sm font-semibold text-gray-400">Роли</th>
            <th class="text-left px-6 py-4 text-sm font-semibold text-gray-400">Статус</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="u in users" :key="u.id" class="border-b border-white/5 last:border-0">
            <td class="px-6 py-4 text-sm text-gray-400">{{ u.id }}</td>
            <td class="px-6 py-4">
              <div class="flex items-center gap-2">
                <div class="w-8 h-8 rounded-full bg-gradient-to-br from-[#00dc82] to-[#00b86b] flex items-center justify-center text-[#020420] font-bold text-sm">
                  {{ u.username[0].toUpperCase() }}
                </div>
                <span class="font-medium">{{ u.username }}</span>
                <span v-if="u.is_owner" class="text-lg">👑</span>
              </div>
            </td>
            <td class="px-6 py-4 text-sm text-gray-400">{{ u.email }}</td>
            <td class="px-6 py-4">
              <div class="flex flex-wrap gap-1">
                <span v-if="u.is_owner" class="px-2 py-0.5 text-xs bg-yellow-500/20 text-yellow-400 rounded">owner</span>
                <span v-if="u.is_superuser" class="px-2 py-0.5 text-xs bg-[#00dc82]/20 text-[#00dc82] rounded">superuser</span>
                <span v-if="u.is_admin" class="px-2 py-0.5 text-xs bg-blue-500/20 text-blue-400 rounded">admin</span>
                <span v-if="u.is_guide" class="px-2 py-0.5 text-xs bg-purple-500/20 text-purple-400 rounded">guide</span>
              </div>
            </td>
            <td class="px-6 py-4">
              <span v-if="u.is_blocked" class="px-2 py-0.5 text-xs bg-red-500/20 text-red-400 rounded">заблокирован</span>
              <span v-else class="px-2 py-0.5 text-xs bg-green-500/20 text-green-400 rounded">активен</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'administrator',
})

const api = useApi()

const loading = ref(true)
const users = ref<any[]>([])

onMounted(async () => {
  try {
    const data = await api.get<any[]>('/api/admin/users')
    users.value = data || []
  } catch (e) {
    console.error('Failed to load users', e)
  } finally {
    loading.value = false
  }
})

useHead({
  title: 'Админка | Пользователи'
})
</script>
