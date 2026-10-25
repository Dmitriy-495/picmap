<template>
  <div class="min-h-[80vh] flex items-center justify-center px-4 py-16">
    <div class="w-full max-w-md">
      <div class="bg-white/5 backdrop-blur-xl border border-white/10 rounded-3xl p-10 shadow-2xl">
        <div class="text-center mb-10">
          <h1 class="text-4xl font-bold mb-3">Вход</h1>
          <p class="text-gray-400">Войдите, чтобы создавать маршруты</p>
        </div>

        <form @submit.prevent="handleLogin" class="space-y-6">
          <div>
            <label class="block text-sm text-gray-400 mb-3">Логин</label>
            <input
              v-model="username"
              type="text"
              required
              autocomplete="username"
              class="w-full px-5 py-4 bg-white/5 border border-white/10 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:border-[#00dc82] focus:ring-2 focus:ring-[#00dc82]/20 transition-all"
              placeholder="tda"
            />
          </div>

          <div>
            <label class="block text-sm text-gray-400 mb-3">Пароль</label>
            <input
              v-model="password"
              type="password"
              required
              autocomplete="current-password"
              class="w-full px-5 py-4 bg-white/5 border border-white/10 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:border-[#00dc82] focus:ring-2 focus:ring-[#00dc82]/20 transition-all"
              placeholder="••••••••"
            />
          </div>

          <div v-if="error" class="p-4 bg-red-500/10 border border-red-500/30 rounded-xl text-red-400 text-sm">
            {{ error }}
          </div>

          <AppButton
            type="submit"
            :disabled="loading"
            size="lg"
            class="w-full"
          >
            {{ loading ? 'Вход...' : 'Войти' }}
          </AppButton>
        </form>

        <div class="mt-8 text-center text-sm text-gray-500">
          <NuxtLink to="/" class="hover:text-[#00dc82] transition-colors">
            ← На главную
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const { login, isAuthenticated } = useAuth()

const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

onMounted(() => {
  if (isAuthenticated.value) {
    navigateTo('/')
  }
})

const handleLogin = async () => {
  loading.value = true
  error.value = ''

  try {
    await login(username.value, password.value)
    await navigateTo('/')
  } catch (e: any) {
    // Проверяем на блокировку
    if (e?.data?.error?.includes('blocked')) {
      error.value = 'Аккаунт заблокирован. Обратитесь к администратору.'
    } else {
      error.value = 'Неверный логин или пароль'
    }
    console.error(e)
  } finally {
    loading.value = false
  }
}

useHead({
  title: 'Вход'
})
</script>
