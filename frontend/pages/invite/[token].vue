<template>
  <div class="min-h-[80vh] flex items-center justify-center px-4 py-16">
    <div class="w-full max-w-md">
      <div class="bg-white/5 backdrop-blur-xl border border-white/10 rounded-3xl p-10 shadow-2xl">
        
        <!-- Загрузка -->
        <div v-if="loading" class="text-center py-8">
          <div class="inline-block w-12 h-12 border-4 border-[#00dc82]/30 border-t-[#00dc82] rounded-full animate-spin"></div>
          <p class="text-gray-400 mt-4">Проверяем приглашение...</p>
        </div>

        <!-- Ошибка приглашения -->
        <div v-else-if="inviteError" class="text-center">
          <div class="text-6xl mb-6">🚫</div>
          <h1 class="text-2xl font-bold mb-3">Приглашение недействительно</h1>
          <p class="text-gray-400 mb-8">{{ inviteError }}</p>
          <AppButton to="/login" size="lg">Войти</AppButton>
        </div>

        <!-- Форма регистрации -->
        <template v-else>
          <div class="text-center mb-10">
            <div class="text-5xl mb-4">🎉</div>
            <h1 class="text-3xl font-bold mb-2">Регистрация</h1>
            <p class="text-gray-400 text-sm">Вы приглашены в PicMap</p>
          </div>

          <form @submit.prevent="handleRegister" class="space-y-6">
            <div>
              <label class="block text-sm text-gray-400 mb-3">Логин</label>
              <input
                v-model="form.username"
                type="text"
                required
                minlength="3"
                maxlength="50"
                class="w-full px-5 py-4 bg-white/5 border border-white/10 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:border-[#00dc82] focus:ring-2 focus:ring-[#00dc82]/20 transition-all"
                placeholder="vra"
              />
              <p class="text-xs text-gray-500 mt-2">Минимум 3 символа</p>
            </div>

            <div>
              <label class="block text-sm text-gray-400 mb-3">Пароль</label>
              <input
                v-model="form.password"
                type="password"
                required
                minlength="6"
                class="w-full px-5 py-4 bg-white/5 border border-white/10 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:border-[#00dc82] focus:ring-2 focus:ring-[#00dc82]/20 transition-all"
                placeholder="••••••••"
              />
              <p class="text-xs text-gray-500 mt-2">Минимум 6 символов</p>
            </div>

            <div>
              <label class="block text-sm text-gray-400 mb-3">Повторите пароль</label>
              <input
                v-model="form.passwordConfirm"
                type="password"
                required
                class="w-full px-5 py-4 bg-white/5 border border-white/10 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:border-[#00dc82] focus:ring-2 focus:ring-[#00dc82]/20 transition-all"
                placeholder="••••••••"
              />
            </div>

            <div v-if="error" class="p-4 bg-red-500/10 border border-red-500/30 rounded-xl text-red-400 text-sm">
              {{ error }}
            </div>

            <AppButton
              type="submit"
              :disabled="submitting"
              size="lg"
              class="w-full"
            >
              {{ submitting ? 'Создаём аккаунт...' : 'Зарегистрироваться' }}
            </AppButton>
          </form>

          <div class="mt-8 text-center text-sm text-gray-500">
            Уже есть аккаунт?
            <NuxtLink to="/login" class="text-[#00dc82] hover:underline">Войти</NuxtLink>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const api = useApi()
const { isAuthenticated } = useAuth()

const token = computed(() => route.params.token as string)

const loading = ref(true)
const submitting = ref(false)
const inviteError = ref('')
const error = ref('')

const form = reactive({
  username: '',
  password: '',
  passwordConfirm: '',
})

onMounted(async () => {
  if (isAuthenticated.value) {
    navigateTo('/')
    return
  }

  try {
    await api.get(`/api/invitations/${token.value}`)
  } catch (e: any) {
    if (e?.statusCode === 404) {
      inviteError.value = 'Приглашение не найдено'
    } else if (e?.statusCode === 410) {
      inviteError.value = 'Приглашение уже использовано или истекло'
    } else {
      inviteError.value = 'Ошибка проверки приглашения'
    }
  } finally {
    loading.value = false
  }
})

const handleRegister = async () => {
  error.value = ''

  if (form.password !== form.passwordConfirm) {
    error.value = 'Пароли не совпадают'
    return
  }

  if (form.password.length < 6) {
    error.value = 'Пароль должен быть минимум 6 символов'
    return
  }

  submitting.value = true

  try {
    const res = await api.post<{ token: string; user: any }>('/api/auth/register', {
      token: token.value,
      username: form.username,
      password: form.password,
    })

    if (import.meta.client) {
      localStorage.setItem('picmap_token', res.token)
      localStorage.setItem('picmap_user', JSON.stringify(res.user))
    }

    await navigateTo('/')
  } catch (e: any) {
    error.value = e?.data?.error || 'Ошибка регистрации'
    console.error(e)
  } finally {
    submitting.value = false
  }
}

useHead({
  title: 'Регистрация'
})
</script>
