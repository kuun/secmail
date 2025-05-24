import { defineStore } from 'pinia'
import axios from 'axios'

interface Message {
  id: number
  from: string
  subject: string
  receivedAt: string
  content: string
  htmlContent: string
  attachments: Array<{id: number, fileName: string}>
}

export const useEmailStore = defineStore('email', {
  state: () => ({
    address: '',
    expiresAt: null as Date | null,
    messages: [] as Message[],
    selectedMessage: null as Message | null,
    view: 'create' as 'create' | 'inbox'
  }),

  actions: {
    async generateEmail() {
      const response = await axios.post('/api/email')
      this.address = response.data.address
      this.expiresAt = new Date(response.data.expiresAt)
      this.messages = []
    },

    async refreshMessages() {
      if (!this.address) return
      try {
        const response = await fetch(`/api/email/${this.address}/messages`)
        const data = await response.json()
        this.messages = data.messages
      } catch (error) {
        console.error('Failed to fetch messages:', error)
      }
    },

    async selectMessage(messageId: number) {
      const response = await axios.get(`/api/message/${messageId}`)
      this.selectedMessage = response.data
    },

    async deleteEmail() {
      if (!this.address) return
      await axios.delete(`/api/email/${this.address}`)
      this.address = ''
      this.expiresAt = null
      this.messages = []
      this.selectedMessage = null
    },

    showInbox() {
      this.view = 'inbox'
      this.refreshMessages()
    },

    showCreate() {
      this.view = 'create'
    },
  }
})
