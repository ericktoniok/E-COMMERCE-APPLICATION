<template>
  <div>
    <div class="flex items-center justify-between mb-4">
      <h2 class="text-xl font-semibold">My Orders</h2>
      <span class="text-xs px-2 py-1 rounded-full"
            :class="sseConnected ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-600'">
        {{ sseConnected ? 'Live' : 'Offline' }}
      </span>
    </div>
    <div v-if="loading" class="space-y-3">
      <div v-for="n in 3" :key="n" class="border rounded p-3 bg-white animate-pulse">
        <div class="h-4 bg-gray-200 rounded w-1/3 mb-2"></div>
        <div class="h-3 bg-gray-100 rounded w-1/4"></div>
      </div>
    </div>
    <div v-else-if="!orders.length" class="text-gray-600">
      No orders yet. <RouterLink to="/products" class="text-blue-600 underline">Browse products</RouterLink>.
    </div>
    <div v-for="o in orders" :key="o.id" class="border rounded p-3 mb-3 bg-white">
      <div class="flex items-center justify-between">
        <div>
          <div class="font-medium">Order #{{ o.id }}</div>
          <div class="text-sm text-gray-600">Status: {{ o.status }}</div>
        </div>
        <div class="font-semibold">Total: {{ money(o.total_cents) }}</div>
      </div>
      <ul class="mt-2 list-disc list-inside text-sm">
        <li v-for="it in o.items" :key="it.id">Product {{ it.product_id }} x {{ it.qty }} @ {{ money(it.price_cents) }}</li>
      </ul>
    </div>
  </div>
</template>
<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { api } from '../lib/api'
import { toast } from '../lib/toast'
import { money } from '../lib/format'

const orders = ref<any[]>([])
let timer: any = null
const lastOk = ref(true)
let es: EventSource | null = null
const apiBase = import.meta.env.VITE_API_BASE || 'http://localhost:8080'
const sseConnected = ref(false)
const loading = ref(false)

async function load() {
  try {
    if (!orders.value.length) loading.value = true
    orders.value = await api.myOrders()
    if (!lastOk.value) { toast('Orders reconnected','success'); lastOk.value = true }
  } catch (e) {
    if (lastOk.value) { toast('Failed to load orders','error'); lastOk.value = false }
  } finally { loading.value = false }
}

onMounted(() => {
  load()
  // SSE: try to subscribe for real-time updates
  const token = localStorage.getItem('token')
  if (token) {
    try {
      es = new EventSource(`${apiBase}/api/orders/stream?token=${encodeURIComponent(token)}`)
      es.onopen = () => { sseConnected.value = true }
      es.onmessage = () => { load() }
      es.onerror = () => { sseConnected.value = false /* fallback keeps running */ }
    } catch {}
  }
  // Fallback polling
  timer = setInterval(load, 7000)
})

onUnmounted(() => { if (timer) clearInterval(timer); if (es) es.close(); sseConnected.value = false })
</script>

