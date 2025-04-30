<script setup lang="ts">
import { useUserStore } from '@/stores/useUserStore';
import AdminRoom from './AdminRoom.vue';
import PlayerRoom from './PlayerRoom.vue';


const userStore = useUserStore()

</script>
<template>
    <AdminRoom v-if="userStore.hasPermision"/>
    <PlayerRoom v-else />

</template>
    