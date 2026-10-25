<template>
  <header class="sticky top-0 z-50 backdrop-blur-xl bg-[#020420]/80 border-b border-white/10">
    <AppContainer>
      <div class="flex items-center justify-between h-20">
        <NuxtLink to="/" class="flex items-center gap-3 group">
          <div class="w-10 h-10 rounded-xl bg-gradient-to-br from-[#00dc82] to-[#00b86b] flex items-center justify-center shadow-lg shadow-[#00dc82]/20 group-hover:shadow-[#00dc82]/40 transition-shadow">
            <span class="text-[#020420] text-xl font-bold">🌍</span>
          </div>
          <span class="text-2xl font-bold">PicMap</span>
        </NuxtLink>

        <nav class="flex items-center gap-6">
          <NuxtLink
            to="/"
            class="text-sm text-gray-400 hover:text-white transition-colors"
            active-class="text-white"
          >
            Маршруты
          </NuxtLink>

          <template v-if="isAuthenticated">
            <NuxtLink
              to="/routes/new"
              class="text-sm text-gray-400 hover:text-white transition-colors"
            >
              + Создать
            </NuxtLink>

            <!-- Админка (для superuser) -->
            <NuxtLink
              v-if="isSuperuser"
              to="/administrator"
              class="text-sm text-gray-400 hover:text-white transition-colors"
              active-class="text-white"
            >
              ⚙️ Администратор
            </NuxtLink>

            <!-- Пригласить (для owner) -->
            <button
              v-if="isOwner"
              @click="showInviteModal = true"
              class="text-sm text-[#00dc82] hover:text-[#00c070] transition-colors"
            >
              🎟️ Пригласить
            </button>

            <div class="flex items-center gap-3 pl-6 border-l border-white/10">
              <span class="text-sm text-gray-400">{{ user?.username }}</span>
              <button
                @click="handleLogout"
                class="text-sm text-gray-400 hover:text-white transition-colors"
              >
                Выйти
              </button>
            </div>
          </template>

          <template v-else>
            <NuxtLink
              to="/login"
              class="px-5 py-2.5 text-sm rounded-xl bg-[#00dc82] text-[#020420] font-semibold hover:bg-[#00c070] transition-colors"
            >
              Войти
            </NuxtLink>
          </template>
        </nav>
      </div>
    </AppContainer>

    <!-- Модалка приглашения -->
    <InviteModal v-model="showInviteModal" />
  </header>
</template>

<script setup lang="ts">
const { isAuthenticated, isSuperuser, isOwner, user, logout } = useAuth()

const showInviteModal = ref(false)

const handleLogout = () => {
  logout()
  navigateTo('/login')
}
</script>
