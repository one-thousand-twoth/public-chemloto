<script setup lang="ts">
import Toaster from './components/Toaster.vue'
import { WebsocketConnector } from './api/websocket/websocket';
import { APISettings } from './api/config';
import { provide } from 'vue';
import { useUserStore } from './stores/useUserStore';
import { storeToRefs } from 'pinia';
import { watch } from 'vue';
defineOptions({
  inheritAttrs: false,
});
const userStore = useUserStore()
const connector = new WebsocketConnector(APISettings.baseURL, '')
provide('connector', connector)
if (userStore.UserCreds) {
    connector.token = userStore.UserCreds?.token
    connector.Run()
}
const { UserCreds } = storeToRefs(userStore)
watch(UserCreds, () => {
    if (userStore.UserCreds) {
        connector.token = userStore.UserCreds?.token
        connector.Run()
    }

})
</script>
<template>
    <router-view v-bind="$attrs" />
    <Toaster />
</template>