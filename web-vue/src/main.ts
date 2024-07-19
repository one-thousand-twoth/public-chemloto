import './assets/output.css'

import { createApp } from 'vue'
import App from './App.vue'

import { createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import RoomList from '@/views/RoomList.vue'

const router = createRouter({
    routes: [{
      path: '/',
      component: RoomList
    }],
    history: createWebHistory()
  })

const pinia = createPinia()
const app = createApp(App)
app.use(pinia)
app.use(router)
app.mount('#app')
