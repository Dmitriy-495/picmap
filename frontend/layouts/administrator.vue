<template>
  <div class="min-h-screen bg-[#020420] flex">
    <!-- Sidebar -->
    <aside class="w-64 bg-[#0a0e1f] border-r border-white/10 flex flex-col">
      <!-- Логотип -->
      <div class="h-20 flex items-center px-6 border-b border-white/10">
        <NuxtLink to="/" class="flex items-center gap-3 group">
          <div class="w-10 h-10 rounded-xl bg-gradient-to-br from-[#00dc82] to-[#00b86b] flex items-center justify-center shadow-lg shadow-[#00dc82]/20">
            <span class="text-[#020420] text-xl font-bold">🌍</span>
          </div>
          <div>
            <div class="text-lg font-bold">PicMap</div>
            <div class="text-xs text-[#00dc82]">Админка</div>
          </div>
        </NuxtLink>
      </div>

      <!-- Навигация -->
      <nav class="flex-1 p-4 space-y-2">
        <NuxtLink
          to="/administrator"
          class="flex items-center gap-3 px-4 py-3 rounded-xl text-gray-400 hover:text-white hover:bg-white/5 transition-colors"
          active-class="bg-[#00dc82]/10 text-[#00dc82]"
          exact-active-class="bg-[#00dc82]/10 text-[#00dc82]"
        >
          <span>📊</span>
          <span>Обзор</span>
        </NuxtLink>

        <NuxtLink
          to="/administrator/invitations"
          class="flex items-center gap-3 px-4 py-3 rounded-xl text-gray-400 hover:text-white hover:bg-white/5 transition-colors"
          active-class="bg-[#00dc82]/10 text-[#00dc82]"
        >
          <span>🎟️</span>
          <span>Приглашения</span>
        </NuxtLink>

        <NuxtLink
          to="/administrator/users"
          class="flex items-center gap-3 px-4 py-3 rounded-xl text-gray-400 hover:text-white hover:bg-white/5 transition-colors"
          active-class="bg-[#00dc82]/10 text-[#00dc82]"
        >
          <span>👥</span>
          <span>Пользователи</span>
        </NuxtLink>
      </nav>

      <!-- Футер -->
      <div class="p-4 border-t border-white/10">
        <div class="text-xs text-gray-500 mb-3">
          {{ user?.username }} <span v-if="user?.is_owner">👑</span>
        </div>
        <NuxtLink
          to="/"
          class="block px-4 py-2 text-sm text-gray-400 hover:text-white transition-colors"
        >
          ← На сайт
        </NuxtLink>
      </div>
    </aside>

    <!-- Контент -->
    <main class="flex-1 overflow-y-auto">
      <slot />
    </main>
  </div>
</template>

<script setup lang="ts">
const { user } = useAuth()
</script>
