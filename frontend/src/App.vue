<template>
  <div class="min-h-screen bg-gray-50 text-gray-900">
    <NavBar />
    <main class="p-4">
      <RouterView />
    </main>
    <ToastHost />
  </div>
</template>
<script setup lang="ts">
import NavBar from './components/NavBar.vue'
import ToastHost from './components/ToastHost.vue'
</script>
