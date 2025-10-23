<template>
  <div>
    <h2 class="text-xl font-semibold mb-4">Admin: Orders</h2>
    <div v-if="!orders.length" class="text-gray-600">No orders yet.</div>
    <div class="space-y-3">
      <div v-for="o in orders" :key="o.id" class="bg-white border rounded p-3">
        <div class="flex items-center justify-between">
          <div>
            <div class="font-medium">Order #{{ o.id }}</div>
            <div class="text-sm text-gray-600">User: {{ o.user_id }} | Status: {{ o.status }}</div>
          </div>
          <div class="font-semibold">Total: {{ (o.total_cents/100).toFixed(2) }}</div>
        </div>
        <ul class="mt-2 list-disc list-inside text-sm">
          <li v-for="it in o.items" :key="it.id">Product {{ it.product_id }} x {{ it.qty }} @ {{ (it.price_cents/100).toFixed(2) }}</li>
        </ul>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { api } from '../lib/api'
import { toast } from '../lib/toast'

const orders = ref<any[]>([])
let timer: any = null
const lastOk = ref(true)

async function load() {
  try {
    orders.value = await api.adminOrders()
    if (!lastOk.value) { toast('Admin orders reconnected','success'); lastOk.value = true }
  } catch (e) {
    if (lastOk.value) { toast('Failed to load admin orders','error'); lastOk.value = false }
  }
}

onMounted(() => {
  load()
  timer = setInterval(load, 7000)
})

onUnmounted(() => { if (timer) clearInterval(timer) })
</script>
