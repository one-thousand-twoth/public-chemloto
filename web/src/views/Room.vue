<script setup lang="ts">
import { useGameStore } from '@/stores/useGameStore';
import { useUserStore } from '@/stores/useUserStore';
import AdminRoom from './AdminRoom.vue';
import PlayerRoom from './PlayerRoom.vue';


const userStore = useUserStore()
const gameStore = useGameStore()

</script>
<template>
    <div v-if="gameStore.fetching"> Загрузка...</div>
    <AdminRoom v-if="userStore.hasPermision"/>
    <PlayerRoom v-else />

</template>
    