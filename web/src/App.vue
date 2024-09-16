<script setup lang="ts">
import { storeToRefs } from 'pinia';
import { provide, watch } from 'vue';
import { useRouter } from 'vue-router';
import { APISettings } from './api/config';
import { WebsocketConnector } from './api/websocket/websocket';
import Toaster from './components/Toaster.vue';
import { useGameStore } from './stores/useGameStore';
import { useUserStore } from './stores/useUserStore';
defineOptions({
    inheritAttrs: false,
});

const userStore = useUserStore()
// const gameStore = useGameStore()
// const router = useRouter()
const connector = new WebsocketConnector(APISettings.baseURL, '')
provide('connector', connector)

const { UserCreds } = storeToRefs(userStore)
watch(UserCreds, () => {
    if (userStore.UserCreds) {
        if (userStore.UserCreds.token != "") {
            connector.token = userStore.UserCreds.token
            connector.Run()
        }
    }

})
// const { connected } = storeToRefs(gameStore)
// watch(connected, () => {
//     console.log("Hello")
//     if (connected.value) {
//         router.replace({ name: "Room" })
//     } else {
//         router.replace({ name: "Hub" })
//     }
// })
</script>
<template>
    <router-view v-bind="$attrs" />
    <Toaster />
</template>