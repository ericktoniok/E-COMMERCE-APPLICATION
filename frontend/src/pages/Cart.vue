<template>
  <div>
    <h2 class="text-xl font-semibold mb-4">Cart</h2>
    <div v-if="!items.length" class="text-gray-600">Your cart is empty.</div>
    <div v-else>
      <table class="w-full text-left bg-white border rounded">
        <thead>
          <tr class="border-b">
            <th class="p-2">Product</th>
            <th class="p-2">Qty</th>
            <th class="p-2">Price</th>
            <th class="p-2">Line Total</th>
            <th class="p-2">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(it, idx) in items" :key="it.product_id" class="border-b">
            <td class="p-2">{{ it.name || ('#'+it.product_id) }}</td>
            <td class="p-2">{{ it.qty }}</td>
            <td class="p-2">{{ price(it.product_id) }}</td>
            <td class="p-2">{{ lineTotal(it) }}</td>
            <td class="p-2 space-x-2">
              <button class="btn" @click="inc(idx)">+</button>
              <button class="btn" @click="dec(idx)">-</button>
              <button class="btn-danger" @click="remove(idx)">Remove</button>
            </td>
          </tr>
        </tbody>
      </table>
      <div class="mt-4 flex items-center justify-between">
        <div class="text-lg">Subtotal: <span class="font-semibold">{{ money(subtotal) }}</span></div>
        <div class="flex gap-2">
          <RouterLink class="btn" to="/products">Continue Shopping</RouterLink>
          <RouterLink class="btn-primary" :class="{ 'opacity-60 pointer-events-none': !items.length }" to="/checkout">Checkout</RouterLink>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { money as formatMoney } from '../lib/format'
import { api } from '../lib/api'

const items = ref<Array<{product_id:number; qty:number; name?:string}>>([])
const prices = ref<Record<number, number>>({}) // product_id -> price_cents
const subtotal = ref(0)

function load() {
  items.value = JSON.parse(localStorage.getItem('cart') || '[]')
  compute()
}
function save() {
  localStorage.setItem('cart', JSON.stringify(items.value))
  compute()
}
function inc(i:number) { items.value[i].qty++; save() }
function dec(i:number) { items.value[i].qty = Math.max(1, items.value[i].qty-1); save() }
function remove(i:number) { items.value.splice(i,1); save() }

function compute() {
  subtotal.value = items.value.reduce((sum: number, it: {product_id:number; qty:number}) => sum + (prices.value[it.product_id] || 0) * it.qty, 0)
}
function money(c:number) { return formatMoney(c) }
function price(id:number) { return money(prices.value[id] || 0) }
function lineTotal(it:{product_id:number; qty:number}) { return money((prices.value[it.product_id] || 0) * it.qty) }

onMounted(async () => {
  load()
  try {
    const list:any[] = await api.products()
    prices.value = Object.fromEntries(list.map(p => [p.id, p.price_cents]))
    compute()
  } catch {}
})
</script>
<style scoped>
.btn { @apply bg-gray-200 px-3 py-1 rounded; }
.btn-primary { @apply bg-blue-600 text-white px-3 py-1 rounded; }
.btn-danger { @apply bg-red-600 text-white px-3 py-1 rounded; }
</style>
