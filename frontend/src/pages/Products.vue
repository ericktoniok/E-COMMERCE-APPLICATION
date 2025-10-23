<template>
  <div>
    <h2 class="text-xl font-semibold mb-4">Products</h2>
    <div class="bg-white border rounded p-3 mb-4 grid md:grid-cols-4 gap-2">
      <input v-model="q" placeholder="Search by name" class="input md:col-span-2" />
      <input v-model.number="min" type="number" placeholder="Min price (cents)" class="input" />
      <select v-model="sort" class="input">
        <option value="">Sort</option>
        <option value="price_asc">Price: Low to High</option>
        <option value="price_desc">Price: High to Low</option>
        <option value="rating_desc">Rating: High to Low</option>
      </select>
    </div>
    <div class="mb-4 flex flex-wrap gap-2" v-if="categories.length">
      <button class="px-3 py-1 rounded border" :class="{ 'bg-blue-600 text-white border-blue-600': selectedCat===''}" @click="selectedCat=''">All</button>
      <button v-for="c in categories" :key="c" class="px-3 py-1 rounded border" :class="{ 'bg-blue-600 text-white border-blue-600': selectedCat===c }" @click="selectedCat=c">{{ c }}</button>
    </div>
    <div v-if="!filtered.length" class="text-gray-600">
      No products found.
      <RouterLink v-if="isAdmin" to="/admin/products" class="text-blue-600 underline ml-1">Add one</RouterLink>.
    </div>
    <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div v-for="p in filtered" :key="p.id" class="border rounded p-3 bg-white">
        <img :src="imgSrc(p)" class="w-full h-40 object-cover mb-2 rounded" />
        <div class="flex items-center justify-between">
          <h3 class="font-medium">{{ p.name }}</h3>
          <span v-if="p.category" class="text-xs px-2 py-0.5 rounded bg-gray-100 border">{{ p.category }}</span>
        </div>
        <div class="text-yellow-500 text-sm">{{ stars(p.rating) }}</div>
        <p class="text-sm text-gray-600 mt-1">{{ p.description }}</p>
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
const sort = ref('')
const selectedCat = ref('')
const auth = useAuth()
const isAdmin = computed(() => auth.role === 'admin')

onMounted(async () => {
  products.value = await api.products()
})

const categories = computed(() => Array.from(new Set(products.value.map((p: any) => p.category).filter(Boolean))).sort())

const filtered = computed(() => {
  let arr = products.value.filter(p => {
    const nameOk = !q.value || (p.name || '').toLowerCase().includes(q.value.toLowerCase())
    const minOk = min.value == null || p.price_cents >= min.value
    const maxOk = max.value == null || p.price_cents <= max.value
    const catOk = !selectedCat.value || p.category === selectedCat.value
    return nameOk && minOk && maxOk && catOk
  })
  if (sort.value === 'price_asc') arr = arr.slice().sort((a:any,b:any)=>a.price_cents-b.price_cents)
  if (sort.value === 'price_desc') arr = arr.slice().sort((a:any,b:any)=>b.price_cents-a.price_cents)
  if (sort.value === 'rating_desc') arr = arr.slice().sort((a:any,b:any)=> (b.rating||0)-(a.rating||0))
  return arr
})

function addToCart(p: any) {
  const cart: any[] = JSON.parse(localStorage.getItem('cart') || '[]')
  const idx = cart.findIndex(i => i.product_id === p.id)
  if (idx >= 0) cart[idx].qty += 1
  else cart.push({ product_id: p.id, qty: 1, name: p.name })
  localStorage.setItem('cart', JSON.stringify(cart))
  alert('Added to cart')
}

function imgSrc(p: any) {
  if (!p.image_url) return 'https://images.unsplash.com/photo-1517059224940-d4af9eec41e5?w=800&q=80&auto=format&fit=crop'
  if (String(p.image_url).startsWith('http')) return p.image_url
  return apiBase + p.image_url
}
function stars(rating: number) {
  const r = Math.max(0, Math.min(5, Math.round(rating)))
  return '★★★★★'.slice(0, r) + '☆☆☆☆☆'.slice(0, 5 - r)
}
</script>
<style scoped>
.btn { @apply bg-blue-600 text-white px-3 py-2 rounded; }
.input { @apply border rounded px-3 py-2 w-full; }
</style>
