<template>
  <div>
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </div>
</template>

<script setup lang="ts">
const { loadFromStorage } = useAuth();

onMounted(() => {
  loadFromStorage();
});

useHead({
  titleTemplate: (title) => (title ? `${title} | PicMap` : "PicMap"),
});
</script>
