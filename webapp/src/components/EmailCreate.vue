<template>
  <div class="bg-white rounded-xl shadow-lg p-8 border border-blue-100 relative">
    <!-- Success Notification -->
    <Transition enter-active-class="transform ease-out duration-300 transition"
      enter-from-class="translate-y-2 opacity-0" enter-to-class="translate-y-0 opacity-100"
      leave-active-class="transition ease-in duration-100" leave-from-class="opacity-100" leave-to-class="opacity-0">
      <div v-if="showNotification" class="fixed inset-0 flex items-center justify-center z-50 pointer-events-none">
        <div class="bg-green-100 border border-green-400 text-green-700 px-6 py-3 rounded-lg shadow-lg text-center">
          {{ notificationMessage }}
        </div>
      </div>
    </Transition>

    <div v-if="!emailStore.address" class="space-y-8">
      <!-- Generate New Email -->
      <div class="text-center py-4">
        <p class="text-gray-600 mb-6">Create your security temporary email address instantly</p>
        <button @click="emailStore.generateEmail"
          class="w-60 bg-blue-600 text-white px-8 py-4 rounded-lg hover:bg-blue-700 transition-colors duration-200 shadow-md">
          Generate Email Address
        </button>
      </div>

      <!-- Or Divider -->
      <div class="flex items-center">
        <div class="flex-1 border-t border-gray-200"></div>
        <span class="px-4 text-gray-500 text-sm">OR</span>
        <div class="flex-1 border-t border-gray-200"></div>
      </div>

      <!-- Access Existing Email -->
      <div class="space-y-4">
        <p class="text-center text-gray-600">Access your existing temporary email</p>
        <form @submit.prevent="accessExistingEmail" class="max-w-lg mx-auto space-y-4">
          <div>
            <input type="email" v-model="existingEmail"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Enter your temporary email address"
              :class="{ 'border-red-500': showError }"
            >
            <p v-if="showError" class="mt-1 text-sm text-red-600">{{ errorMessage }}</p>
          </div>
          <div class="text-center">
            <button type="submit"
              class="w-60 bg-gray-100 text-gray-800 px-8 py-4 rounded-lg hover:bg-gray-200 transition-colors">
              Access Inbox
            </button>
          </div>
        </form>
      </div>
    </div>

    <template v-else>
      <div class="space-y-6">
        <div class="bg-blue-50 p-6 rounded-lg border border-blue-200">
          <label class="text-sm text-blue-700 font-medium mb-2 block">Your Security Temporary Email</label>
          <div class="flex flex-wrap items-center gap-2">
            <span class="text-base sm:text-xl font-mono text-blue-900 break-all">{{ emailStore.address }}</span>
            <button @click="handleCopy" title="Copy to clipboard"
              class="shrink-0 text-blue-600 hover:text-blue-800 p-1.5 sm:p-2 rounded-md border border-blue-200 hover:bg-blue-200 transition-colors">
              <ClipboardIcon class="w-4 h-4 sm:w-5 sm:h-5" />
            </button>
          </div>
          <ExpirationTimer :expires-at="emailStore.expiresAt" class="mt-3" />
          <label class="text-sm text-red-700 font-light mb-2 block">Please save this email address - it will only be
            shown once!</label>
        </div>

        <div class="flex justify-between items-center">
          <button @click="router.push({ name: 'inbox' })"
            class="text-blue-600 hover:text-blue-800 px-4 py-2 rounded-md border border-blue-200 hover:bg-blue-50">
            <span>Goto Inbox</span>
          </button>
          <button @click="emailStore.deleteEmail"
            class="bg-red-500 text-white px-4 py-2 rounded-md hover:bg-red-600 transition-colors">
            Delete Email
          </button>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useEmailStore } from '../stores/email'
import { ClipboardIcon } from '@heroicons/vue/24/solid'

const emailStore = useEmailStore()
const router = useRouter()
const showNotification = ref(false)
const notificationMessage = ref('')

const showSuccess = (message: string) => {
  notificationMessage.value = message
  showNotification.value = true
  setTimeout(() => {
    showNotification.value = false
  }, 3000)
}

const handleCopy = async () => {
  await navigator.clipboard.writeText(emailStore.address)
  showSuccess('Email address copied to clipboard!')
}

// Watch for new email creation
watch(() => emailStore.address, (newValue) => {
  if (newValue) {
    showSuccess('Email address created successfully!')
  }
})

const existingEmail = ref('')
const showError = ref(false)
const errorMessage = ref('')

const accessExistingEmail = async () => {
  if (!isValidEmailFormat(existingEmail.value)) {
    showError.value = true
    errorMessage.value = 'Invalid email format'
    return
  }

  try {
    // Verify email exists and is valid
    const response = await fetch(`/api/email/${existingEmail.value}`)
    if (response.ok) {
      emailStore.address = existingEmail.value
      router.push({ name: 'inbox' })
    } else if (response.status === 410) {
      showError.value = true
      errorMessage.value = 'Email has expired'
    } else {
      showError.value = true
      errorMessage.value = 'Email not found'
    }
  } catch (error) {
    showError.value = true
    errorMessage.value = 'Failed to verify email'
  }
}

// Reset error when input changes
watch(existingEmail, () => {
  showError.value = false
  errorMessage.value = ''
})

const isValidEmailFormat = (email: string) => {
  // Basic email format validation regex
  const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return re.test(email)
}
</script>

<style></style>
