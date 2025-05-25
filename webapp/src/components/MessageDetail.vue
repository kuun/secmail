<template>
  <div class="bg-white rounded-xl shadow-lg p-8 border border-blue-100">
    <div class="flex items-center justify-between mb-6">
      <div>
        <h2 class="text-lg font-semibold text-gray-900">{{ message?.subject }}</h2>
        <p class="text-sm text-gray-600">{{ message?.from }}</p>
        <time class="text-xs text-gray-500">{{ formatDate(message?.receivedAt) }}</time>
      </div>
      <button @click="router.back()"
        class="text-gray-600 hover:text-gray-800 px-4 py-2 rounded-md border border-gray-200 hover:bg-gray-50">
        Back
      </button>
    </div>

    <div class="space-y-4">
      <!-- HTML Content -->
      <div v-if="message?.htmlContent" class="prose max-w-none p-4 bg-white rounded-lg border border-gray-200"
        v-html="message.htmlContent"></div>
      <!-- Plain Text Content -->
      <div v-else class="whitespace-pre-wrap p-4 bg-white rounded-lg border border-gray-200">
        {{ message?.content }}
      </div>

      <!-- Attachments -->
      <div v-if="message?.attachments?.length" class="mt-6">
        <h3 class="text-sm font-medium text-gray-700 mb-2">Attachments</h3>
        <div class="space-y-2">
          <div v-for="attachment in message.attachments" :key="attachment.id"
            class="flex items-center gap-2 p-2 border border-gray-200 rounded-md">
            <span class="text-sm text-gray-600">{{ attachment.fileName }}</span>
            <a :href="`/api/message/${message.id}/attachment/${attachment.id}`" download
              class="text-blue-600 hover:text-blue-800 text-sm">
              Download
            </a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useEmailStore } from '../stores/email'
import { Message } from '../stores/email'

const route = useRoute()
const router = useRouter()
const emailStore = useEmailStore()
const message = ref<Message| null>(null)

onMounted(async () => {
  const messageId = route.params.id as string
  await emailStore.selectMessage(messageId)
  message.value = emailStore.selectedMessage
})

const formatDate = (date: string|undefined) => {
  return date ? new Date(date).toLocaleString() : ''
}
</script>
