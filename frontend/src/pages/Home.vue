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
  <div class="space-y-8">
    <section class="bg-white rounded p-6 border">
      <div class="flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
        <div>
          <h1 class="text-2xl md:text-3xl font-semibold">Mini E‑Commerce</h1>
          <p class="text-gray-600 mt-2">Go + Fiber backend, Vue 3 frontend, PostgreSQL, and a simulated M‑Pesa checkout.</p>
          <p class="text-sm text-gray-500 mt-1">API health: <span :class="health==='ok' ? 'text-green-600' : 'text-red-600'">{{ health }}</span></p>
          <div class="mt-4 flex gap-3">
            <RouterLink to="/products" class="px-4 py-2 bg-blue-600 text-white rounded">Browse Products</RouterLink>
            <RouterLink to="/login" class="px-4 py-2 border rounded">Sign In</RouterLink>
          </div>
        </div>
      </div>
    </section>

    <section class="grid md:grid-cols-4 gap-4">
      <div class="bg-white rounded p-4 border">
        <div class="font-medium">1. Browse</div>
        <div class="text-sm text-gray-600">See products, add items to your cart.</div>
      </div>
      <div class="bg-white rounded p-4 border">
        <div class="font-medium">2. Checkout</div>
        <div class="text-sm text-gray-600">Place an order with a simulated M‑Pesa request.</div>
      </div>
      <div class="bg-white rounded p-4 border">
        <div class="font-medium">3. Webhook</div>
        <div class="text-sm text-gray-600">Mock M‑Pesa calls the backend via a webhook.</div>
      </div>
      <div class="bg-white rounded p-4 border">
        <div class="font-medium">4. Track</div>
        <div class="text-sm text-gray-600">Orders update automatically as payment status changes.</div>
      </div>
    </section>

    <section class="bg-white rounded p-6 border">
      <div class="text-sm text-gray-600 mb-3">Tech</div>
      <div class="flex flex-wrap gap-2">
        <span class="px-2 py-1 rounded bg-gray-100 border">Go Fiber</span>
        <span class="px-2 py-1 rounded bg-gray-100 border">Vue 3 + TS</span>
        <span class="px-2 py-1 rounded bg-gray-100 border">PostgreSQL</span>
        <span class="px-2 py-1 rounded bg-gray-100 border">Docker Compose</span>
        <span class="px-2 py-1 rounded bg-gray-100 border">TailwindCSS</span>
      </div>
    </section>
  </div>
</template>
