<script setup lang="ts">
import { ref } from 'vue';
import RoomList from '@/views/RoomList.vue';
import ClientList from '@/components/ClientList.vue';

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
    .active{
        background-color: rgb(37 99 235 / var(--tw-bg-opacity)) 
    }
</style>