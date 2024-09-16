import './assets/output.css'

import { createApp } from 'vue'
import App from './App.vue'

import Hub from '@/views/Hub.vue'
import Login from '@/views/Login.vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from './stores/useUserStore'

const pinia = createPinia()
const app = createApp(App)
app.use(pinia);

const userStore = useUserStore();


const router = createRouter({
  routes: [
    {
      name: 'Hub',
      path: '/',
      component: Hub,
    }, {
      name: 'Login',
      path: '/login',
      component: Login,
    },
  ],
  history: createWebHistory()
});

(async () => {
  const ok = await userStore.check()
  if (!ok) {
    router.replace({ name: "Login" })
  }
})();


router.beforeEach(async (to, from) => {
  // userStore.check()
  if (
    !userStore.UserCreds && to.name !== 'Login') {
    console.log("to Logon")
    return { name: 'Login' }
  }
  if (userStore.UserCreds && to.name == 'Login') {
    return from
  }
})


app.use(router)
app.mount('#app')
