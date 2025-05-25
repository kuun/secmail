import { defineStore } from 'pinia'
import axios from 'axios'
import router  from '../router'

interface StoredEmail {
  address: string
  expiresAt: string
}

export interface Message {
  id: string
  from: string
  subject: string
  receivedAt: string
  content: string
  htmlContent: string
  attachments: Array<{id: string, fileName: string}>
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
    async loadStoredEmail() {
      const stored = localStorage.getItem('tempEmail')
      if (stored) {
        const data: StoredEmail = JSON.parse(stored)
        const expiresAt = new Date(data.expiresAt)
        
        // First check local expiration
        if (expiresAt > new Date()) {
          try {
            // Validate with server
            const response = await axios.get(`/api/email/${data.address}`)
            if (response.status === 200) {
              this.address = data.address
              this.expiresAt = expiresAt
              return true
            }
          } catch (error: any) {
            if (error.response?.status === 410 || error.response?.status === 404) {
              // Email has expired on server
              localStorage.removeItem('tempEmail')
              router.push({ name: 'create' })
            }
          }
        } else {
          // Locally expired
          localStorage.removeItem('tempEmail')
        }
      }
      return false
    },

    saveEmail() {
      if (this.address && this.expiresAt) {
        const data: StoredEmail = {
          address: this.address,
          expiresAt: this.expiresAt.toISOString()
        }
        localStorage.setItem('tempEmail', JSON.stringify(data))
      }
    },

    setEmail(address: string, expiresAt: Date) {
      this.address = address
      this.expiresAt = expiresAt
      this.saveEmail()
    },

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
        if (!response.ok) {
          if (response.status === 410) {
            // Email expired, redirect to create page
            router.push({ name: 'create' })
            return
          }
          throw new Error('Failed to fetch messages')
        }
        const data = await response.json()
        this.messages = data.messages
      } catch (error) {
        console.error('Failed to fetch messages:', error)
        this.messages = []
      }
    },

    async selectMessage(messageId: string) {
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
