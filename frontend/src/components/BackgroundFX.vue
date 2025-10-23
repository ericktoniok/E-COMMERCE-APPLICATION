<template>
  <div class="pointer-events-none fixed inset-0 -z-10 overflow-hidden">
    <!-- Soft gradient backdrop -->
    <div class="absolute inset-0 bg-gradient-to-br from-blue-50 via-white to-purple-50"></div>

    <!-- Animated glow blobs -->
    <div class="absolute -top-24 -left-24 w-[520px] h-[520px] rounded-full bg-blue-300/20 blur-3xl animate-float-slow"></div>
    <div class="absolute -bottom-24 -right-24 w-[520px] h-[520px] rounded-full bg-purple-300/20 blur-3xl animate-float-slower"></div>

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
