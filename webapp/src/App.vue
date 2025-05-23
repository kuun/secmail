<template>
  <div class="min-h-screen bg-gray-100 p-4">
    <div class="max-w-4xl mx-auto">
      <header class="text-center mb-8">
        <h1 class="text-3xl font-bold text-gray-800">SecMail</h1>
        <p class="text-gray-600">Temporary Email Service</p>
      </header>

      <div class="bg-white rounded-lg shadow p-6">
        <div v-if="!emailStore.emailAddress" class="text-center">
          <button
            @click="emailStore.generateEmail"
            class="bg-blue-500 text-white px-6 py-3 rounded-lg hover:bg-blue-600"
          >
            Generate Temporary Email
          </button>
        </div>

        <template v-else>
          <div class="flex items-center justify-between mb-6">
            <div>
              <div class="flex items-center gap-2">
                <span class="text-lg font-semibold">{{ emailStore.emailAddress }}</span>
                <button
                  @click="copyToClipboard"
                  class="text-sm bg-gray-100 px-2 py-1 rounded"
                >
                  Copy
                </button>
              </div>
              <ExpirationTimer
                :expires-at="emailStore.expiresAt"
                @expired="emailStore.deleteEmail"
              />
            </div>
            <div class="space-x-2">
              <button
                @click="emailStore.refreshMessages"
                class="bg-gray-100 px-4 py-2 rounded hover:bg-gray-200"
              >
                Refresh
              </button>
              <button
                @click="emailStore.deleteEmail"
                class="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600"
              >
                Delete
              </button>
            </div>
          </div>

          <MessageList
            :messages="emailStore.messages"
            @select="emailStore.selectMessage"
          />

          <MessageView
            v-if="emailStore.selectedMessage"
            :message="emailStore.selectedMessage"
          />
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useEmailStore } from './stores/email'
import ExpirationTimer from './components/ExpirationTimer.vue'
import MessageList from './components/MessageList.vue'
import MessageView from './components/MessageView.vue'

const emailStore = useEmailStore()

const copyToClipboard = async () => {
  await navigator.clipboard.writeText(emailStore.emailAddress)
}
</script>
