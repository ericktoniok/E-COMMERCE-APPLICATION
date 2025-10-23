<script setup lang="ts">
import { ref, onMounted } from 'vue'

const health = ref('checking...')
const base = import.meta.env.VITE_API_BASE || 'http://localhost:8080'

onMounted(async () => {
  try {
    const res = await fetch(`${base}/api/health`)
    const data = await res.json()
    health.value = data.status || 'ok'
  } catch (e) {
    health.value = 'unavailable'
  }
})
</script>

<template>
  <div>
    <p class="text-sm text-gray-600">API health: {{ health }}</p>
    <div class="mt-4 p-4 bg-white rounded shadow">Welcome! Storefront coming next.</div>
  </div>
</template>
