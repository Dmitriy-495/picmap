<template>
  <AppContainer>
    <div class="py-16 max-w-2xl mx-auto">
      <!-- Хлебные крошки -->
      <div class="mb-8 text-sm">
        <NuxtLink to="/" class="text-gray-400 hover:text-white transition-colors">
          Маршруты
        </NuxtLink>
        <span class="text-gray-600 mx-3">/</span>
        <span class="text-white">Новый маршрут</span>
      </div>

      <div class="bg-white/5 backdrop-blur-xl border border-white/10 rounded-3xl p-8 md:p-10 shadow-2xl">
        <h1 class="text-3xl font-bold mb-8">Создать маршрут</h1>

        <form @submit.prevent="handleSubmit" class="space-y-6">
          <!-- Название -->
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-2">
              Название <span class="text-red-400">*</span>
            </label>
            <input
              v-model="form.title"
              type="text"
              required
              placeholder="Например: Карелия 2026"
              class="w-full px-5 py-3.5 bg-white/5 border border-white/10 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:border-[#00dc82] focus:ring-2 focus:ring-[#00dc82]/20 transition-all"
            />
          </div>

          <!-- Описание -->
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-2">
              Описание
            </label>
            <textarea
              v-model="form.description"
              rows="3"
              placeholder="Кратко о путешествии..."
              class="w-full px-5 py-3.5 bg-white/5 border border-white/10 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:border-[#00dc82] focus:ring-2 focus:ring-[#00dc82]/20 transition-all resize-none"
            ></textarea>
          </div>

          <!-- Даты -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label class="block text-sm font-medium text-gray-300 mb-2">
                Дата начала <span class="text-red-400">*</span>
              </label>
              <input
                v-model="form.start_date"
                type="date"
                required
                class="w-full px-5 py-3.5 bg-white/5 border border-white/10 rounded-xl text-white focus:outline-none focus:border-[#00dc82] focus:ring-2 focus:ring-[#00dc82]/20 transition-all [color-scheme:dark]"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-300 mb-2">
                Дата окончания <span class="text-red-400">*</span>
              </label>
              <input
                v-model="form.end_date"
                type="date"
                required
                class="w-full px-5 py-3.5 bg-white/5 border border-white/10 rounded-xl text-white focus:outline-none focus:border-[#00dc82] focus:ring-2 focus:ring-[#00dc82]/20 transition-all [color-scheme:dark]"
              />
            </div>
          </div>

          <!-- Ошибка -->
          <div v-if="error" class="p-4 bg-red-500/10 border border-red-500/30 rounded-xl text-red-400 text-sm">
            {{ error }}
          </div>

          <!-- Кнопки -->
          <div class="flex gap-4 pt-4">
            <AppButton
              type="submit"
              :disabled="loading"
              size="lg"
              class="flex-1"
            >
              {{ loading ? 'Создание...' : 'Создать маршрут' }}
            </AppButton>
            <AppButton
              to="/"
              variant="secondary"
              size="lg"
            >
              Отмена
            </AppButton>
          </div>
        </form>
      </div>
    </div>
  </AppContainer>
</template>

<script setup lang="ts">
const { createRoute } = useRoutes()
const { isAuthenticated } = useAuth()

const form = reactive({
  title: '',
  description: '',
  start_date: '',
  end_date: '',
})

const loading = ref(false)
const error = ref('')

// Защита: только авторизованные могут создавать
onMounted(() => {
  if (!isAuthenticated.value) {
    navigateTo('/login')
  }
})

const handleSubmit = async () => {
  loading.value = true
  error.value = ''

  try {
    const newRoute = await createRoute({
      title: form.title,
      description: form.description,
      start_date: form.start_date,
      end_date: form.end_date,
    })
    await navigateTo(`/routes/${newRoute.id}`)
  } catch (e: any) {
    error.value = e?.data?.error || 'Ошибка создания маршрута'
    console.error(e)
  } finally {
    loading.value = false
  }
}

useHead({
  title: 'Новый маршрут'
})
</script>
