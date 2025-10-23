<template>
  <header class="p-4 border-b bg-white flex items-center justify-between">
    <h1 class="text-xl font-semibold">
      <RouterLink to="/">Mini E-Commerce</RouterLink>
    </h1>
    <nav class="flex items-center gap-4 text-sm">
      <RouterLink to="/products" class="link">Products</RouterLink>
      <RouterLink to="/cart" class="link">Cart <span v-if="cartCount" class="ml-1 bg-blue-600 text-white rounded px-2">{{ cartCount }}</span></RouterLink>
      <RouterLink to="/orders" class="link" v-if="isAuthed">My Orders</RouterLink>
      <RouterLink to="/admin/products" class="link" v-if="isAdmin">Admin</RouterLink>
      <template v-if="!isAuthed">
        <RouterLink to="/login" class="link">Login</RouterLink>
        <RouterLink to="/register" class="link">Register</RouterLink>
      </template>
      <button v-else class="text-red-600" @click="logout">Logout</button>
    </nav>
  </header>
</template>
<script setup lang="ts">
import { computed, onMounted, onUpdated, ref } from 'vue'
import { useAuth } from '../stores/auth'

const auth = useAuth()
const isAuthed = computed(() => !!auth.token)
const isAdmin = computed(() => auth.role === 'admin')
const cartCount = ref(0)

function refreshCart() {
  const items = JSON.parse(localStorage.getItem('cart') || '[]')
  cartCount.value = items.reduce((a:any, b:any) => a + (b.qty||1), 0)
}

function logout() { auth.logout() }

onMounted(refreshCart)
// keep it up-to-date on re-renders
onUpdated(refreshCart)
</script>
<style scoped>
.link { @apply text-gray-700 hover:text-blue-600; }
</style>
