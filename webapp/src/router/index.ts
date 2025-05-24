import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import EmailCreate from '../components/EmailCreate.vue'
import Inbox from '../components/Inbox.vue'
import { useEmailStore } from '../stores/email'



const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'create',
    component: EmailCreate
  },
  {
    path: '/inbox',
    name: 'inbox',
    component: Inbox,
    beforeEnter: (to, from) => {
        const emailStore = useEmailStore()
        if (!emailStore.address) {
          return { name: 'create' }
        }
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
