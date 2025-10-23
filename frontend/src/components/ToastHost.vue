<template>
  <div class="fixed bottom-4 right-4 z-50 space-y-2">
    <div v-for="t in toasts" :key="t.id" :class="klass(t.type)" class="px-3 py-2 rounded shadow text-white min-w-[240px]">
      {{ t.message }}
    </div>
  </div>
</template>
<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'

type Toast = { id: number; message: string; type: 'success'|'error'|'info'; timeout: number }
const toasts = ref<Toast[]>([])

function klass(type: Toast['type']) {
  if (type === 'success') return 'bg-green-600'
  if (type === 'error') return 'bg-red-600'
  return 'bg-gray-800'
}

function addToast(t: Toast) {
  toasts.value.push(t)
  setTimeout(() => {
    toasts.value = toasts.value.filter(x => x.id !== t.id)
  }, t.timeout || 3000)
}

function handler(e: Event) {
  const ce = e as CustomEvent<Toast>
  if (ce.detail) addToast(ce.detail)
}

onMounted(() => window.addEventListener('app:toast', handler as any))
onUnmounted(() => window.removeEventListener('app:toast', handler as any))
</script>
