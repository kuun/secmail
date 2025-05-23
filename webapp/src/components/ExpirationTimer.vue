<template>
  <div class="text-sm text-gray-500">
    <span>Expires in: </span>
    <span :class="{ 'text-red-500': timeLeft < 60000 }">{{ formattedTimeLeft }}</span>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'

const props = defineProps<{
  expiresAt: Date | null
}>()

const emit = defineEmits<{
  (e: 'expired'): void
}>()

const timeLeft = ref(0)
let timer: number

const updateTimeLeft = () => {
  if (!props.expiresAt) return
  
  const now = new Date().getTime()
  const expiry = new Date(props.expiresAt).getTime()
  timeLeft.value = Math.max(0, expiry - now)

  if (timeLeft.value === 0) {
    emit('expired')
  }
}

const formattedTimeLeft = computed(() => {
  const minutes = Math.floor(timeLeft.value / 60000)
  const seconds = Math.floor((timeLeft.value % 60000) / 1000)
  
  if (minutes > 0) {
    return `${minutes}m ${seconds}s`
  }
  return `${seconds}s`
})

onMounted(() => {
  updateTimeLeft()
  timer = window.setInterval(updateTimeLeft, 1000)
})

onUnmounted(() => {
  clearInterval(timer)
})
</script>
