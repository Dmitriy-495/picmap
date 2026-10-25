<template>
  <div class="bg-white/5 backdrop-blur-xl border border-white/10 rounded-2xl p-6">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-2xl font-bold flex items-center gap-2">
        <span>💬</span>
        <span>Воспоминания</span>
        <span v-if="comments.length" class="text-sm font-normal text-gray-500">
          ({{ comments.length }})
        </span>
      </h2>
    </div>

    <div v-if="isAuthenticated" class="mb-6">
      <div class="flex gap-3">
        <div class="flex-1">
          <textarea
            v-model="newComment"
            placeholder="Напиши воспоминание..."
            rows="3"
            class="w-full px-4 py-3 bg-white/5 border border-white/10 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:border-[#00dc82] transition-colors resize-none"
            @keydown.ctrl.enter="submitComment"
            @keydown.meta.enter="submitComment"
          ></textarea>
          <div class="text-xs text-gray-500 mt-1">Ctrl+Enter — отправить</div>
        </div>
        <button
          @click="submitComment"
          :disabled="!newComment.trim() || submitting"
          class="px-6 py-3 bg-[#00dc82] text-[#020420] font-semibold rounded-xl hover:bg-[#00c070] disabled:opacity-50 disabled:cursor-not-allowed transition-colors self-start"
        >
          {{ submitting ? '...' : 'Отправить' }}
        </button>
      </div>
    </div>

    <div v-else class="mb-6 p-4 bg-white/5 border border-white/10 rounded-xl text-center">
      <p class="text-gray-400 text-sm">
        <NuxtLink to="/login" class="text-[#00dc82] hover:underline">Войдите</NuxtLink>, чтобы оставить воспоминание
      </p>
    </div>

    <div v-if="loading" class="text-center py-8">
      <div class="inline-block w-8 h-8 border-3 border-[#00dc82]/30 border-t-[#00dc82] rounded-full animate-spin"></div>
    </div>

    <div v-else-if="comments.length === 0" class="text-center py-8 text-gray-500">
      <div class="text-4xl mb-2">💭</div>
      <p>Пока нет воспоминаний</p>
      <p class="text-sm mt-1">Стань первым, кто поделится!</p>
    </div>

    <div v-else class="space-y-4">
      <div v-for="comment in comments" :key="comment.id" class="group bg-white/[0.03] border border-white/5 rounded-xl p-4 hover:bg-white/[0.05] transition-colors">
        <div class="flex items-start gap-3">
          <div class="w-10 h-10 rounded-full bg-gradient-to-br from-[#00dc82] to-[#00b86b] flex items-center justify-center text-[#020420] font-bold flex-shrink-0">
            {{ comment.username[0].toUpperCase() }}
          </div>
          <div class="flex-1 min-w-0">
            <div class="flex items-center justify-between gap-2 mb-1">
              <div class="flex items-center gap-2 flex-wrap">
                <span class="font-semibold text-sm">{{ comment.username }}</span>
                <span class="text-xs text-gray-500">{{ formatTime(comment.created_at) }}</span>
                <span v-if="comment.updated_at" class="text-xs text-gray-600 italic">(изменено)</span>
              </div>
              <div v-if="canEdit(comment)" class="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                <button @click="startEdit(comment)" class="p-1.5 text-gray-400 hover:text-white transition-colors" title="Редактировать">✏️</button>
                <button @click="deleteComment(comment.id)" class="p-1.5 text-gray-400 hover:text-red-400 transition-colors" title="Удалить">🗑️</button>
              </div>
            </div>
            <div v-if="editingId === comment.id">
              <textarea v-model="editText" rows="2" class="w-full px-3 py-2 bg-white/5 border border-[#00dc82]/50 rounded-lg text-white text-sm focus:outline-none resize-none" @keydown.ctrl.enter="saveEdit(comment.id)" @keydown.escape="cancelEdit"></textarea>
              <div class="flex gap-2 mt-2">
                <button @click="saveEdit(comment.id)" class="px-3 py-1 text-xs bg-[#00dc82] text-[#020420] font-medium rounded-lg hover:bg-[#00c070] transition-colors">Сохранить</button>
                <button @click="cancelEdit" class="px-3 py-1 text-xs bg-white/5 text-gray-400 rounded-lg hover:bg-white/10 transition-colors">Отмена</button>
              </div>
            </div>
            <p v-else class="text-sm text-gray-300 whitespace-pre-wrap break-words">{{ comment.text }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Comment {
  id: number
  route_id: number
  user_id: number
  username: string
  text: string
  created_at: string
  updated_at?: string | null
}

const props = defineProps<{ routeId: number }>()

const api = useApi()
const { user, isAuthenticated, isSuperuser, isAdmin } = useAuth()

const comments = ref<Comment[]>([])
const loading = ref(true)
const submitting = ref(false)
const newComment = ref('')
const editingId = ref<number | null>(null)
const editText = ref('')

const loadComments = async () => {
  loading.value = true
  try {
    const data = await api.get<Comment[]>(`/api/routes/${props.routeId}/comments`)
    comments.value = data || []
  } catch (e) {
    console.error('Failed to load comments', e)
    comments.value = []
  } finally {
    loading.value = false
  }
}

const submitComment = async () => {
  if (!newComment.value.trim() || submitting.value) return
  submitting.value = true
  try {
    const comment = await api.post<Comment>(`/api/routes/${props.routeId}/comments`, { text: newComment.value.trim() })
    comments.value.push(comment)
    newComment.value = ''
  } catch (e) {
    console.error('Failed to post comment', e)
    alert('Ошибка отправки комментария')
  } finally {
    submitting.value = false
  }
}

const canEdit = (comment: Comment) => {
  if (!user.value) return false
  return comment.user_id === user.value.id || isSuperuser.value || isAdmin.value
}

const startEdit = (comment: Comment) => {
  editingId.value = comment.id
  editText.value = comment.text
}

const cancelEdit = () => {
  editingId.value = null
  editText.value = ''
}

const saveEdit = async (id: number) => {
  if (!editText.value.trim()) return
  try {
    const updated = await api.put<Comment>(`/api/comments/${id}`, { text: editText.value.trim() })
    const index = comments.value.findIndex(c => c.id === id)
    if (index !== -1) comments.value[index] = updated
    cancelEdit()
  } catch (e) {
    console.error('Failed to update comment', e)
    alert('Ошибка редактирования')
  }
}

const deleteComment = async (id: number) => {
  if (!confirm('Удалить это воспоминание?')) return
  try {
    await api.delete(`/api/comments/${id}`)
    comments.value = comments.value.filter(c => c.id !== id)
  } catch (e) {
    console.error('Failed to delete comment', e)
    alert('Ошибка удаления')
  }
}

const formatTime = (dateStr: string) => {
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)
  if (minutes < 1) return 'только что'
  if (minutes < 60) return `${minutes} мин. назад`
  if (hours < 24) return `${hours} ч. назад`
  if (days < 7) return `${days} дн. назад`
  return date.toLocaleDateString('ru-RU', { day: 'numeric', month: 'long', year: date.getFullYear() !== now.getFullYear() ? 'numeric' : undefined })
}

onMounted(() => { loadComments() })
defineExpose({ reload: loadComments })
</script>
