<script setup lang="ts">
import { storeToRefs } from 'pinia';
import { provide, watch } from 'vue';
import { APISettings } from './api/config';
import { WebsocketConnector } from './api/websocket/websocket';
import Toaster from './components/Toaster.vue';
import { useUserStore } from './stores/useUserStore';
defineOptions({
    inheritAttrs: false,
});
const userStore = useUserStore()
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
</script>
<template>
    <router-view v-bind="$attrs" />
    <Toaster />
</template>