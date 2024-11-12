import './assets/output.css'

import { createApp, provide, watch } from 'vue'
import App from './App.vue'

import Hub from '@/views/Hub.vue'
import Login from '@/views/Login.vue'
import { createPinia, storeToRefs } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import { APISettings } from './api/config'
import { WebsocketConnector } from './api/websocket/websocket'
import { piniaWebsocketPlugin, websocketPlugin } from './api/websocket/websocketPlugin'
import { useUserStore } from './stores/useUserStore'

const pinia = createPinia()
const app = createApp(App)

const connector = new WebsocketConnector(APISettings.baseURL, '')
// pinia.use(piniaWebsocketPlugin(connector))

app.use(websocketPlugin, connector)
app.use(pinia);

const userStore = useUserStore();
const { UserCreds } = storeToRefs(userStore)

// provide('connector', connector)
watch(UserCreds, () => {
  if (UserCreds.value) {
    if (UserCreds.value.token != "") {
      connector.token = UserCreds.value.token
      connector.Run()
    }
  }

}, { deep: true })


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
