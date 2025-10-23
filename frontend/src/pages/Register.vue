<template>
  <div class="max-w-md mx-auto">
    <h2 class="text-xl font-semibold mb-4">Register</h2>
    <form @submit.prevent="onSubmit" class="space-y-3">
      <div>
        <label class="block text-sm">Email</label>
        <input v-model="email" type="email" class="input" required placeholder="you@example.com" />
      </div>
      <div>
        <label class="block text-sm">Password</label>
        <input v-model="password" type="password" class="input" required placeholder="••••••••" minlength="8" />
        <p class="text-xs text-gray-500 mt-1">Use at least 8 characters. You'll be logged in automatically after registering.</p>
      </div>
      <button class="btn" :disabled="loading">{{ loading ? 'Creating…' : 'Create account' }}</button>
      <p class="text-sm mt-2">Have an account? <RouterLink to="/login" class="text-blue-600">Login</RouterLink></p>
      <p v-if="error" class="text-sm text-red-600">{{ error }}</p>
    </form>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from '../lib/toast'
import { useAuth } from '../stores/auth'

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
    await auth.register(email.value, password.value)
    toast('Account created!','success')
    router.push('/')
  } catch (e: any) {
    error.value = 'Registration failed'
    toast('Registration failed','error')
  }
  loading.value = false
}
</script>
<style scoped>
.input { @apply w-full border rounded px-3 py-2; }
.btn { @apply bg-blue-600 text-white px-4 py-2 rounded; }
</style>
