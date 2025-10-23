<template>
  <div>
    <h2 class="text-xl font-semibold mb-4">My Orders</h2>
    <div v-if="!orders.length" class="text-gray-600">No orders yet.</div>
    <div v-for="o in orders" :key="o.id" class="border rounded p-3 mb-3 bg-white">
      <div class="flex items-center justify-between">
        <div>
          <div class="font-medium">Order #{{ o.id }}</div>
          <div class="text-sm text-gray-600">Status: {{ o.status }}</div>
        </div>
        <div class="font-semibold">Total: {{ (o.total_cents/100).toFixed(2) }}</div>
      </div>
      <ul class="mt-2 list-disc list-inside text-sm">
        <li v-for="it in o.items" :key="it.id">Product {{ it.product_id }} x {{ it.qty }} @ {{ (it.price_cents/100).toFixed(2) }}</li>
      </ul>
    </div>
  </div>
</template>
<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { api } from '../lib/api'

const orders = ref<any[]>([])
let timer: any = null

async function load() {
  try {
    orders.value = await api.myOrders()
  } catch (e) {
    // likely unauthenticated
  }
}

onMounted(() => {
  load()
  timer = setInterval(load, 7000)
})

onUnmounted(() => { if (timer) clearInterval(timer) })
</script>
