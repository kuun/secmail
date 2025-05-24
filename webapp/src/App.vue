<template>
  <div class="min-h-screen p-4 relative">
    <AnimatedBackground />
    <div class="max-w-4xl mx-auto relative">
      <header class="text-center mb-12">
        <h1 class="text-4xl font-bold text-blue-800">SecMail</h1>
        <p class="text-blue-600 mt-2">Secure Temporary Email Service</p>
      </header>
<!-- 
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view> -->
      <RouterView/>
    </div>
  </div>
</template>

<script setup lang="ts">
import AnimatedBackground from './components/AnimatedBackground.vue'
import { RouterView } from 'vue-router';
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
