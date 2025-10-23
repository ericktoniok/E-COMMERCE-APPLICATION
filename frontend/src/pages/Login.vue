<template>
  <div class="max-w-md mx-auto">
    <h2 class="text-xl font-semibold mb-4">Login</h2>
    <form @submit.prevent="onSubmit" class="space-y-3">
      <div>
        <label class="block text-sm">Email</label>
        <input v-model="email" type="email" class="input" required placeholder="you@example.com" />
      </div>
      <div>
        <label class="block text-sm">Password</label>
        <input v-model="password" type="password" class="input" required placeholder="••••••••" minlength="8" />
        <p class="text-xs text-gray-500 mt-1">Use the seeded admin for dashboard access: <code>admin@local.test</code> / <code>Admin123!</code></p>
      </div>
      <button class="btn" :disabled="loading">{{ loading ? 'Signing in…' : 'Login' }}</button>
      <p class="text-sm mt-2">No account? <RouterLink to="/register" class="text-blue-600">Register</RouterLink></p>
      <p v-if="error" class="text-sm text-red-600">{{ error }}</p>
    </form>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '../stores/auth'
import { toast } from '../lib/toast'

const email = ref('')
const password = ref('')
const error = ref('')
const auth = useAuth()
const router = useRouter()
const loading = ref(false)

async function onSubmit() {
  error.value = ''
  loading.value = true
  try {
    await auth.login(email.value, password.value)
    toast('Welcome back!','success')
    router.push('/')
  } catch (e: any) {
    error.value = 'Login failed'
    toast('Login failed','error')
  }
  loading.value = false
}
</script>
<style scoped>
.input { @apply w-full border rounded px-3 py-2; }
.btn { @apply bg-blue-600 text-white px-4 py-2 rounded; }
</style>
