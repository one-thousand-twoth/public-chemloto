import './assets/output.css'

import { createApp } from 'vue'
import App from './App.vue'

import { createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import Login from '@/views/Login.vue'
import RoomList from '@/views/RoomList.vue'
import { useUserStore } from './stores/useUserStore'

const pinia = createPinia()

const app = createApp(App)
app.use(pinia)
const userStore = useUserStore()

const router = createRouter({
    routes: [{
      name: 'RoomList',
      path: '/',
      component: RoomList,
    },{
      name: 'Login',
      path: '/login',
      component: Login,
    }
  ],
    history: createWebHistory()
  })

router.beforeEach(async (to, from) => {
    if (
      // make sure the user is authenticated
      !userStore.UserCreds &&
      // ❗️ Avoid an infinite redirect
      to.name !== 'Login'
    ) {
      // redirect the user to the login page
      return { name: 'Login' }
    }
    if(userStore.UserCreds && to.name == 'Login'){
      return from
    }
})


app.use(router)
app.mount('#app')
