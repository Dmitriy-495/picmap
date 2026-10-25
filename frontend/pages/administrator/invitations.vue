<template>
  <div class="p-10">
    <div class="mb-10 flex items-center justify-between">
      <div>
        <h1 class="text-4xl font-bold mb-3">Приглашения</h1>
        <p class="text-gray-400">Управление доступом к PicMap</p>
      </div>
      <AppButton @click="showCreateForm = true" size="lg">
        + Создать приглашение
      </AppButton>
    </div>

    <!-- Форма создания -->
    <div v-if="showCreateForm" class="mb-8 bg-white/5 backdrop-blur-xl border border-white/10 rounded-2xl p-6">
      <h2 class="text-xl font-bold mb-4">Новое приглашение</h2>
      <form @submit.prevent="createInvitation" class="flex gap-4 items-end">
        <div class="flex-1">
          <label class="block text-sm text-gray-400 mb-2">Email приглашаемого</label>
          <input
            v-model="newEmail"
            type="email"
            required
            placeholder="friend@example.com"
            class="w-full px-4 py-3 bg-white/5 border border-white/10 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:border-[#00dc82] transition-colors"
          />
        </div>
        <AppButton type="submit" :disabled="creating">
          {{ creating ? 'Создаём...' : 'Создать' }}
        </AppButton>
        <AppButton variant="secondary" @click="showCreateForm = false">
          Отмена
        </AppButton>
      </form>
    </div>

    <!-- Ссылка после создания -->
    <div v-if="createdUrl" class="mb-8 bg-green-500/10 border border-green-500/30 rounded-2xl p-6">
      <div class="flex items-center justify-between">
        <div>
          <h3 class="text-lg font-bold text-green-400 mb-2">✅ Приглашение создано!</h3>
          <p class="text-sm text-gray-400 mb-3">Отправьте эту ссылку другу:</p>
          <div class="flex items-center gap-2">
            <input
              :value="createdUrl"
              readonly
              class="flex-1 px-4 py-2 bg-black/30 border border-green-500/30 rounded-xl text-green-400 font-mono text-sm"
            />
            <AppButton size="sm" @click="copyUrl">
              Копировать
            </AppButton>
          </div>
        </div>
      </div>
    </div>

    <!-- Загрузка -->
    <div v-if="loading" class="text-center py-20">
      <div class="inline-block w-12 h-12 border-4 border-[#00dc82]/30 border-t-[#00dc82] rounded-full animate-spin"></div>
    </div>

    <!-- Пусто -->
    <div v-else-if="invitations.length === 0" class="text-center py-20 bg-white/5 rounded-2xl border border-white/10">
      <div class="text-5xl mb-4">🎟️</div>
      <h2 class="text-2xl font-bold mb-2">Нет приглашений</h2>
      <p class="text-gray-400 mb-6">Создайте первое приглашение</p>
    </div>

    <!-- Таблица -->
    <div v-else class="bg-white/5 backdrop-blur-xl border border-white/10 rounded-2xl overflow-hidden">
      <table class="w-full">
        <thead class="bg-white/5 border-b border-white/10">
          <tr>
            <th class="text-left px-6 py-4 text-sm font-semibold text-gray-400">Токен</th>
            <th class="text-left px-6 py-4 text-sm font-semibold text-gray-400">Создано</th>
            <th class="text-left px-6 py-4 text-sm font-semibold text-gray-400">Истекает</th>
            <th class="text-left px-6 py-4 text-sm font-semibold text-gray-400">Статус</th>
            <th class="text-right px-6 py-4 text-sm font-semibold text-gray-400">Действия</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="inv in invitations" :key="inv.id" class="border-b border-white/5 last:border-0">
            <td class="px-6 py-4 text-sm font-mono text-gray-300">{{ inv.token.slice(0, 16) }}...</td>
            <td class="px-6 py-4 text-sm text-gray-400">{{ formatDate(inv.created_at) }}</td>
            <td class="px-6 py-4 text-sm text-gray-400">{{ formatDate(inv.expires_at) }}</td>
            <td class="px-6 py-4">
              <AppBadge :variant="getInvitationStatus(inv).variant">
                {{ getInvitationStatus(inv).text }}
              </AppBadge>
            </td>
            <td class="px-6 py-4 text-right">
              <div class="flex items-center justify-end gap-2">
                <button
                  v-if="!inv.used_at && new Date(inv.expires_at) > new Date()"
                  @click="copyInviteUrl(inv.token)"
                  class="px-3 py-1.5 text-xs bg-[#00dc82]/20 text-[#00dc82] rounded-lg hover:bg-[#00dc82]/30 transition-colors"
                >
                  📋 Копировать
                </button>
                <button
                  @click="deleteInvitation(inv.id)"
                  class="px-3 py-1.5 text-xs bg-red-500/20 text-red-400 rounded-lg hover:bg-red-500/30 transition-colors"
                >
                  🗑️
                </button>
              </div>
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
const creating = ref(false)
const showCreateForm = ref(false)
const newEmail = ref('')
const createdUrl = ref('')
const invitations = ref<any[]>([])

const loadInvitations = async () => {
  loading.value = true
  try {
    const data = await api.get<any[]>('/api/admin/invitations')
    invitations.value = data || []
  } catch (e) {
    console.error('Failed to load invitations', e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadInvitations()
})

const createInvitation = async () => {
  creating.value = true
  createdUrl.value = ''

  try {
    const res = await api.post<{ invitation: any; url: string }>('/api/admin/invitations', {
      email: newEmail.value,
    })
    createdUrl.value = res.url
    newEmail.value = ''
    showCreateForm.value = false
    await loadInvitations()
  } catch (e: any) {
    alert(e?.data?.error || 'Ошибка создания')
  } finally {
    creating.value = false
  }
}

const copyUrl = async () => {
  try {
    await navigator.clipboard.writeText(createdUrl.value)
    alert('✅ Ссылка скопирована!')
  } catch {
    alert('Не удалось скопировать')
  }
}

const copyInviteUrl = async (token: string) => {
  const url = window.location.origin.replace('localhost:3000', 'localhost:3000') + '/invite/' + token
  try {
    await navigator.clipboard.writeText(url)
    alert('✅ Ссылка скопирована!')
  } catch {
    alert('Не удалось скопировать')
  }
}

const deleteInvitation = async (id: number) => {
  if (!confirm('Удалить приглашение?')) return

  try {
    await api.delete(`/api/admin/invitations/${id}`)
    invitations.value = invitations.value.filter(i => i.id !== id)
  } catch (e) {
    alert('Ошибка удаления')
  }
}

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
  title: 'Админка | Приглашения'
})
</script>
