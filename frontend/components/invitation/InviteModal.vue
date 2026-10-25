<template>
  <UModal v-model="isOpen">
    <div class="p-6">
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-2xl font-bold">Пригласить пользователя</h2>
        <button @click="close" class="text-gray-400 hover:text-white transition-colors">
          ✕
        </button>
      </div>

      <!-- Форма -->
      <div v-if="!createdUrl">
        <div class="mb-4">
          <label class="block text-sm text-gray-400 mb-2">Email приглашаемого</label>
          <input
            v-model="email"
            type="email"
            placeholder="friend@example.com"
            class="w-full px-4 py-3 bg-white/5 border border-white/10 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:border-[#00dc82] transition-colors"
            @keydown.enter="create"
          />
        </div>

        <div v-if="error" class="mb-4 p-3 bg-red-500/10 border border-red-500/30 rounded-lg text-red-400 text-sm">
          {{ error }}
        </div>

        <div class="flex gap-3">
          <AppButton @click="create" :disabled="creating || !email" class="flex-1">
            {{ creating ? 'Создаём...' : 'Создать приглашение' }}
          </AppButton>
          <AppButton variant="secondary" @click="close">Отмена</AppButton>
        </div>
      </div>

      <!-- Ссылка -->
      <div v-else>
        <div class="text-center mb-6">
          <div class="text-5xl mb-3">🎉</div>
          <h3 class="text-xl font-bold text-green-400">Приглашение создано!</h3>
          <p class="text-sm text-gray-400 mt-2">Скопируйте ссылку и отправьте другу</p>
        </div>

        <div class="flex gap-2 mb-4">
          <input
            :value="createdUrl"
            readonly
            class="flex-1 px-4 py-3 bg-black/30 border border-green-500/30 rounded-xl text-green-400 font-mono text-sm"
          />
          <AppButton @click="copyUrl">
            {{ copied ? '✅' : '📋' }}
          </AppButton>
        </div>

        <div class="flex gap-3">
          <AppButton @click="reset" variant="secondary" class="flex-1">
            Создать ещё
          </AppButton>
          <AppButton @click="close" class="flex-1">
            Готово
          </AppButton>
        </div>
      </div>
    </div>
  </UModal>
</template>

<script setup lang="ts">
const props = defineProps<{
  modelValue: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

const api = useApi()

const isOpen = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
})

const email = ref('')
const creating = ref(false)
const error = ref('')
const createdUrl = ref('')
const copied = ref(false)

const create = async () => {
  if (!email.value) return
  creating.value = true
  error.value = ''

  try {
    const res = await api.post<{ invitation: any; url: string }>('/api/admin/invitations', {
      email: email.value,
    })
    createdUrl.value = res.url
  } catch (e: any) {
    error.value = e?.data?.error || 'Ошибка создания'
  } finally {
    creating.value = false
  }
}

const copyUrl = async () => {
  try {
    await navigator.clipboard.writeText(createdUrl.value)
    copied.value = true
    setTimeout(() => (copied.value = false), 2000)
  } catch {
    alert('Не удалось скопировать')
  }
}

const reset = () => {
  email.value = ''
  createdUrl.value = ''
  error.value = ''
}

const close = () => {
  isOpen.value = false
  setTimeout(reset, 300)
}
</script>
