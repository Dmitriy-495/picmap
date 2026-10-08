<template>
  <div class="min-h-[80vh] flex items-center justify-center px-4">
    <div class="w-full max-w-md">
      <div
        class="bg-white/5 backdrop-blur-xl border border-white/10 rounded-2xl p-8 shadow-2xl"
      >
        <div class="text-center mb-8">
          <h1 class="text-3xl font-bold mb-2">Вход</h1>
          <p class="text-gray-400 text-sm">Войдите, чтобы создавать маршруты</p>
        </div>

        <form @submit.prevent="handleLogin" class="space-y-4">
          <div>
            <label class="block text-sm text-gray-400 mb-2">Логин</label>
            <input
              v-model="username"
              type="text"
              required
              class="w-full px-4 py-3 bg-white/5 border border-white/10 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-[#00dc82] transition-colors"
              placeholder="tda495"
            />
          </div>

          <div>
            <label class="block text-sm text-gray-400 mb-2">Пароль</label>
            <input
              v-model="password"
              type="password"
              required
              class="w-full px-4 py-3 bg-white/5 border border-white/10 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-[#00dc82] transition-colors"
              placeholder="••••••••"
            />
          </div>

          <div
            v-if="error"
            class="p-3 bg-red-500/10 border border-red-500/30 rounded-lg text-red-400 text-sm"
          >
            {{ error }}
          </div>

          <button
            type="submit"
            :disabled="loading"
            class="w-full py-3 bg-[#00dc82] text-[#020420] font-semibold rounded-lg hover:bg-[#00c070] disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            {{ loading ? "Вход..." : "Войти" }}
          </button>
        </form>

        <div class="mt-6 text-center text-sm text-gray-500">
          <NuxtLink to="/" class="hover:text-[#00dc82] transition-colors">
            ← На главную
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const { login, isAuthenticated } = useAuth();

const username = ref("");
const password = ref("");
const loading = ref(false);
const error = ref("");

onMounted(() => {
  if (isAuthenticated.value) {
    navigateTo("/");
  }
});

const handleLogin = async () => {
  loading.value = true;
  error.value = "";

  try {
    await login(username.value, password.value);
    await navigateTo("/");
  } catch (e: any) {
    error.value = "Неверный логин или пароль";
    console.error(e);
  } finally {
    loading.value = false;
  }
};
</script>
