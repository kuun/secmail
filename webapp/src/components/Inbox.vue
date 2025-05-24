<template>
  <div class="bg-white rounded-xl shadow-lg p-8 border border-blue-100">
    <div class="flex items-center justify-between mb-6">
      <div>
        <h2 class="text-lg font-semibold text-gray-900">Inbox</h2>
        <p class="text-sm text-gray-600">{{ emailStore.address }}</p>
      </div>
      <div class="flex gap-3">
        <button @click="emailStore.refreshMessages"
          class="text-blue-600 hover:text-blue-800 px-4 py-2 rounded-md border border-blue-200 hover:bg-blue-50">
          Refresh
        </button>
        <button @click="router.push({ name: 'create' })"
          class="text-gray-600 hover:text-gray-800 px-4 py-2 rounded-md border border-gray-200 hover:bg-gray-50">
          Back
        </button>
      </div>
    </div>

    <!-- Messages List -->
    <div class="space-y-2">
      <div v-if="emailStore.messages.length === 0" class="text-center py-8 text-gray-500">
        No messages yet
      </div>
      <div v-else v-for="message in emailStore.messages" :key="message.id"
        class="p-4 border border-gray-200 rounded-lg hover:bg-gray-50 cursor-pointer">
        <div class="flex justify-between items-start">
          <div>
            <p class="font-medium text-gray-900">{{ message.from }}</p>
            <p class="text-sm text-gray-600">{{ message.subject }}</p>
          </div>
          <time class="text-xs text-gray-500">{{ formatDate(message.receivedAt) }}</time>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useEmailStore } from '../stores/email'

const emailStore = useEmailStore()
const router = useRouter()

const formatDate = (date: string) => {
  return new Date(date).toLocaleString()
}
</script>
