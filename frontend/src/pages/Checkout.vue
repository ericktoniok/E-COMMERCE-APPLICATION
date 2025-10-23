<template>
  <div class="max-w-lg mx-auto">
    <h2 class="text-xl font-semibold mb-4">Checkout</h2>
    <div v-if="!items.length" class="text-gray-600">Your cart is empty.</div>
    <div v-else class="space-y-3">
      <div class="bg-white border rounded p-3">
        <h3 class="font-medium mb-2">Items</h3>
        <ul class="list-disc list-inside text-sm">
          <li v-for="it in items" :key="it.product_id">{{ it.name || ('#'+it.product_id) }} x {{ it.qty }}</li>
        </ul>
        <div class="mt-3 text-sm">Subtotal: <span class="font-semibold">{{ money(subtotal) }}</span></div>
      </div>
      <button class="btn-primary" @click="placeOrder" :disabled="loading">
        {{ loading ? 'Placing order...' : 'Place Order' }}
      </button>
      <p v-if="error" class="text-sm text-red-600">{{ error }}</p>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../lib/api'
import { toast } from '../lib/toast'
import { money } from '../lib/format'

const items = ref<Array<{product_id:number; qty:number; name?:string}>>([])
const loading = ref(false)
const error = ref('')
const router = useRouter()
const prices = ref<Record<number, number>>({})
const subtotal = ref(0)

onMounted(() => {
  items.value = JSON.parse(localStorage.getItem('cart') || '[]')
  // Fetch product prices to compute subtotal
  api.products().then((list: any[]) => {
    prices.value = Object.fromEntries(list.map(p => [p.id, p.price_cents]))
    subtotal.value = items.value.reduce((sum, it) => sum + (prices.value[it.product_id] || 0) * it.qty, 0)
  }).catch(()=>{})
})

async function placeOrder() {
  error.value = ''
  loading.value = true
  try {
    await api.checkout(items.value.map((i: {product_id:number; qty:number}) => ({ product_id: i.product_id, qty: i.qty })))
    localStorage.removeItem('cart')
    toast('Order placed!','success')
    router.push('/orders')
  } catch (e:any) {
    error.value = 'Checkout failed (are you logged in?)'
    toast('Checkout failed','error')
  } finally {
    loading.value = false
  }
}
</script>
<style scoped>
.btn-primary { @apply bg-blue-600 text-white px-4 py-2 rounded; }
</style>
