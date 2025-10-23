<template>
  <div>
    <h2 class="text-xl font-semibold mb-4">Products</h2>
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div v-for="p in products" :key="p.id" class="border rounded p-3 bg-white">
        <img v-if="p.image_url" :src="apiBase + p.image_url" class="w-full h-40 object-cover mb-2" />
        <h3 class="font-medium">{{ p.name }}</h3>
        <p class="text-sm text-gray-600">{{ p.description }}</p>
        <p class="font-semibold mt-1">{{ (p.price_cents/100).toFixed(2) }}</p>
        <button class="btn mt-2" @click="addToCart(p)">Add to cart</button>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { api } from '../lib/api'

const products = ref<any[]>([])
const apiBase = import.meta.env.VITE_API_BASE || 'http://localhost:8080'

onMounted(async () => {
  products.value = await api.products()
})

function addToCart(p: any) {
  const cart: any[] = JSON.parse(localStorage.getItem('cart') || '[]')
  const idx = cart.findIndex(i => i.product_id === p.id)
  if (idx >= 0) cart[idx].qty += 1
  else cart.push({ product_id: p.id, qty: 1, name: p.name })
  localStorage.setItem('cart', JSON.stringify(cart))
  alert('Added to cart')
}
</script>
<style scoped>
.btn { @apply bg-blue-600 text-white px-3 py-2 rounded; }
</style>
