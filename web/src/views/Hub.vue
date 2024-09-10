<script setup lang="ts">
import RoomList from '@/components/RoomList.vue';
import ClientList from '@/components/ClientList.vue';
import { ref } from 'vue';
import Tab from '@/components/tabs/Tab.vue';

const currentTab = ref('RoomList')
const tabs = {
    RoomList,
    ClientList,
}
</script>
<template>
    <!-- <div class="absolute">
        <button v-for="(_, tab) in tabs" :key="tab" :class="['tab-button', { active: currentTab === tab }]"
            @click="currentTab = tab"> {{ tab }} </button>
    </div> -->
    <!-- <component :is="(<any>tabs)[currentTab]" class="tab"></component> -->
    <Tab />
</template>


<style scoped>
.active {
    background-color: rgb(37 99 235 / var(--tw-bg-opacity))
}
</style>