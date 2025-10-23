<template>
  <div class="pointer-events-none fixed inset-0 -z-10 overflow-hidden">
    <!-- Soft gradient backdrop -->
    <div class="absolute inset-0 bg-gradient-to-br" :class="gradient"></div>

    <!-- Animated glow blobs -->
    <div class="absolute -top-24 -left-24 w-[520px] h-[520px] rounded-full blur-3xl animate-float-slow" :class="blob1"></div>
    <div class="absolute -bottom-24 -right-24 w-[520px] h-[520px] rounded-full blur-3xl animate-float-slower" :class="blob2"></div>

    <!-- Subtle grid overlay -->
    <div class="absolute inset-0 opacity-[0.08]" aria-hidden="true">
      <svg class="w-full h-full" xmlns="http://www.w3.org/2000/svg">
        <defs>
          <pattern id="bg-grid" width="32" height="32" patternUnits="userSpaceOnUse">
            <path d="M 32 0 L 0 0 0 32" fill="none" stroke="currentColor" stroke-width="0.5" />
          </pattern>
          <radialGradient id="fade" cx="50%" cy="30%" r="70%">
            <stop offset="0%" stop-color="#000" stop-opacity="0.25"/>
            <stop offset="100%" stop-color="#000" stop-opacity="0"/>
          </radialGradient>
        </defs>
        <rect width="100%" height="100%" fill="url(#bg-grid)" />
        <rect width="100%" height="100%" fill="url(#fade)" />
      </svg>
    </div>
  </div>
</template>
<script setup lang="ts">
import { computed } from 'vue'

type Scheme = 'blue' | 'emerald' | 'amber' | 'rose'
const props = defineProps<{ scheme?: Scheme }>()
const scheme = computed<Scheme>(() => (props.scheme ?? 'blue'))

const gradients: Record<Scheme, string> = {
  blue:   'from-blue-50 via-white to-purple-50',
  emerald:'from-emerald-50 via-white to-teal-50',
  amber:  'from-amber-50 via-white to-rose-50',
  rose:   'from-rose-50 via-white to-pink-50',
}
const blob1Map: Record<Scheme, string> = {
  blue:   'bg-blue-300/20',
  emerald:'bg-emerald-300/20',
  amber:  'bg-amber-300/20',
  rose:   'bg-rose-300/20',
}
const blob2Map: Record<Scheme, string> = {
  blue:   'bg-purple-300/20',
  emerald:'bg-teal-300/20',
  amber:  'bg-rose-300/20',
  rose:   'bg-pink-300/20',
}

const gradient = computed(() => gradients[scheme.value])
const blob1 = computed(() => blob1Map[scheme.value])
const blob2 = computed(() => blob2Map[scheme.value])
</script>
<style scoped>
@keyframes float-slow {
  0% { transform: translate3d(0,0,0) scale(1); }
  50% { transform: translate3d(20px, -10px, 0) scale(1.05); }
  100% { transform: translate3d(0,0,0) scale(1); }
}
@keyframes float-slower {
  0% { transform: translate3d(0,0,0) scale(1); }
  50% { transform: translate3d(-24px, 12px, 0) scale(1.05); }
  100% { transform: translate3d(0,0,0) scale(1); }
}
.animate-float-slow { animation: float-slow 12s ease-in-out infinite; }
.animate-float-slower { animation: float-slower 16s ease-in-out infinite; }
</style>
