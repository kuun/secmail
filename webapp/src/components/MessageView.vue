<template>
  <div class="mt-6 border-t border-gray-200 pt-6">
    <div class="space-y-4">
      <div class="space-y-1">
        <div class="text-sm text-gray-500">From:</div>
        <div class="font-medium">{{ message.from }}</div>
      </div>
      
      <div class="space-y-1">
        <div class="text-sm text-gray-500">Subject:</div>
        <div class="font-medium">{{ message.subject || '(No subject)' }}</div>
      </div>

      <div class="space-y-1">
        <div class="text-sm text-gray-500">Received:</div>
        <div>{{ new Date(message.receivedAt).toLocaleString() }}</div>
      </div>

      <div v-if="message.attachments.length > 0" class="space-y-2">
        <div class="text-sm text-gray-500">Attachments:</div>
        <div class="flex flex-wrap gap-2">
          <a
            v-for="attachment in message.attachments"
            :key="attachment.id"
            :href="`/api/message/${message.id}/attachment/${attachment.id}`"
            download
            class="inline-flex items-center px-3 py-1 rounded-full bg-gray-100 hover:bg-gray-200 text-sm"
          >
            <span class="truncate max-w-xs">{{ attachment.fileName }}</span>
          </a>
        </div>
      </div>

      <div class="border-t border-gray-200 pt-4">
        <div v-if="message.htmlContent" v-html="sanitizedHtml" class="prose max-w-none"></div>
        <pre v-else class="whitespace-pre-wrap">{{ message.content }}</pre>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import DOMPurify from 'dompurify'

interface Attachment {
  id: number
  fileName: string
}

interface Message {
  id: number
  from: string
  subject: string
  receivedAt: string
  content: string
  htmlContent: string
  attachments: Attachment[]
}

const props = defineProps<{
  message: Message
}>()

const sanitizedHtml = computed(() => {
  return DOMPurify.sanitize(props.message.htmlContent)
})
</script>
