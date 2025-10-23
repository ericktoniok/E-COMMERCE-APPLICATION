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
            <th class="p-2">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(it, idx) in items" :key="it.product_id" class="border-b">
            <td class="p-2">{{ it.name || ('#'+it.product_id) }}</td>
            <td class="p-2">{{ it.qty }}</td>
            <td class="p-2 space-x-2">
              <button class="btn" @click="inc(idx)">+</button>
              <button class="btn" @click="dec(idx)">-</button>
              <button class="btn-danger" @click="remove(idx)">Remove</button>
            </td>
          </tr>
        </tbody>
      </table>
      <div class="mt-4 flex gap-2">
        <RouterLink class="btn" to="/products">Continue Shopping</RouterLink>
        <RouterLink class="btn-primary" to="/checkout">Checkout</RouterLink>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'

const items = ref<Array<{product_id:number; qty:number; name?:string}>>([])

function load() {
  items.value = JSON.parse(localStorage.getItem('cart') || '[]')
}
function save() {
  localStorage.setItem('cart', JSON.stringify(items.value))
}
function inc(i:number) { items.value[i].qty++; save() }
function dec(i:number) { items.value[i].qty = Math.max(1, items.value[i].qty-1); save() }
function remove(i:number) { items.value.splice(i,1); save() }

onMounted(load)
</script>
<style scoped>
.btn { @apply bg-gray-200 px-3 py-1 rounded; }
.btn-primary { @apply bg-blue-600 text-white px-3 py-1 rounded; }
.btn-danger { @apply bg-red-600 text-white px-3 py-1 rounded; }
</style>
