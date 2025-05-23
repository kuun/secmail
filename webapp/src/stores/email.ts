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
    emailAddress: '',
    expiresAt: null as Date | null,
    messages: [] as Message[],
    selectedMessage: null as Message | null,
  }),

  actions: {
    async generateEmail() {
      const response = await axios.post('/api/email')
      this.emailAddress = response.data.emailAddress
      this.expiresAt = new Date(response.data.expiresAt)
      this.messages = []
    },

    async refreshMessages() {
      if (!this.emailAddress) return
      const response = await axios.get(`/api/email/${this.emailAddress}/messages`)
      this.messages = response.data
    },

    async selectMessage(messageId: number) {
      const response = await axios.get(`/api/message/${messageId}`)
      this.selectedMessage = response.data
    },

    async deleteEmail() {
      if (!this.emailAddress) return
      await axios.delete(`/api/email/${this.emailAddress}`)
      this.emailAddress = ''
      this.expiresAt = null
      this.messages = []
      this.selectedMessage = null
    }
  }
})
