<template>
  <div class="divide-y divide-gray-200">
    <div v-if="messages.length === 0" class="py-4 text-center text-gray-500">
      No messages yet
    </div>
    
    <div
      v-for="message in sortedMessages"
      :key="message.id"
      @click="$emit('select', message.id)"
      class="py-4 px-2 hover:bg-gray-50 cursor-pointer transition-colors"
    >
      <div class="flex justify-between items-start">
        <div class="space-y-1">
          <div class="font-medium text-gray-900">{{ message.from }}</div>
          <div class="text-gray-600">{{ message.subject || '(No subject)' }}</div>
        </div>
        <time
          :datetime="message.receivedAt"
          class="text-sm text-gray-500"
        >
          {{ formatDate(message.receivedAt) }}
        </time>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Message {
  id: number
  from: string
  subject: string
  receivedAt: string
}

const props = defineProps<{
  messages: Message[]
}>()

defineEmits<{
  (e: 'select', messageId: number): void
}>()

const sortedMessages = computed(() => {
  return [...props.messages].sort((a, b) => 
    new Date(b.receivedAt).getTime() - new Date(a.receivedAt).getTime()
  )
})

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  const now = new Date()
  
  // If today, show time only
  if (date.toDateString() === now.toDateString()) {
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  }
  
  // Otherwise show date
  return date.toLocaleDateString([], {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>
