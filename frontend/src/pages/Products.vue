<template>
  <div>
    <h2 class="text-xl font-semibold mb-4">Products</h2>
    <div class="bg-white border rounded p-3 mb-4 grid md:grid-cols-4 gap-2">
      <input v-model="q" placeholder="Search by name" class="input md:col-span-2" />
      <input v-model.number="min" type="number" placeholder="Min price (cents)" class="input" />
      <input v-model.number="max" type="number" placeholder="Max price (cents)" class="input" />
    </div>
    <div v-if="!filtered.length" class="text-gray-600">
      No products found.
      <RouterLink v-if="isAdmin" to="/admin/products" class="text-blue-600 underline ml-1">Add one</RouterLink>.
    </div>
    <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div v-for="p in filtered" :key="p.id" class="border rounded p-3 bg-white">
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
import { computed, onMounted, ref } from 'vue'
import { api } from '../lib/api'
import { useAuth } from '../stores/auth'

const products = ref<any[]>([])
const apiBase = import.meta.env.VITE_API_BASE || 'http://localhost:8080'
const q = ref('')
const min = ref<number | null>(null)
const max = ref<number | null>(null)
const auth = useAuth()
const isAdmin = computed(() => auth.role === 'admin')

onMounted(async () => {
  products.value = await api.products()
})

const filtered = computed(() => {
  return products.value.filter(p => {
    const nameOk = !q.value || (p.name || '').toLowerCase().includes(q.value.toLowerCase())
    const minOk = min.value == null || p.price_cents >= min.value
    const maxOk = max.value == null || p.price_cents <= max.value
    return nameOk && minOk && maxOk
  })
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
.input { @apply border rounded px-3 py-2 w-full; }
</style>
