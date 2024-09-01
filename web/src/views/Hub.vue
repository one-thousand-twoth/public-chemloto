<script setup lang="ts">
import RoomList from '@/components/RoomList.vue';
import ClientList from '@/components/game/ClientList.vue';
import { ref } from 'vue';

const currentTab = ref('RoomList')
const tabs = {
    RoomList,
    ClientList,
}
</script>
<template>
    <div class="absolute">
        <button v-for="(_, tab) in tabs" :key="tab" :class="['tab-button', { active: currentTab === tab }]"
            @click="currentTab = tab"> {{ tab }} </button>
    </div>
    <component :is="(<any>tabs)[currentTab]" class="tab"></component>
</template>


<style scoped>
.active {
    background-color: rgb(37 99 235 / var(--tw-bg-opacity))
}
</style>