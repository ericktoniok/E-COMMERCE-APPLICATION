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
import { onMounted, ref } from 'vue'
import { api } from '../lib/api'

const orders = ref<any[]>([])

async function load() {
  try {
    orders.value = await api.adminOrders()
  } catch (e) {
    // likely unauthorized
  }
}

onMounted(load)
</script>
