<script setup lang="ts">
import ClientList from '@/components/ClientList.vue';
import RoomList from '@/components/RoomList.vue';
import { useUserStore } from '@/stores/useUserStore';
import { storeToRefs } from 'pinia';
import { ref } from 'vue';
import Room from './Room.vue';

const tabs = {
    "Комнаты": RoomList,
    "Игроки": ClientList,
}

const picked = ref('Комнаты')
const userStore = useUserStore()
const { connected } = storeToRefs(userStore)
</script>
<template>
    <div v-if="!connected" class="p-2 md:p-8 ">
        <div class="relative left-[8px] top-[0.5rem] flex flex-wrap bg-opacity-0 overflow-y-hidden">
            <div class="relative  z-0 has-[:checked]:z-[3] bg-white has-[:checked]:border-solid has-[:checked]:border-[5px] has-[:checked]:border-b-0  border-blue-400 rounded-br-lg rounded-bl-lg border-r "
                v-for="(_, tab) in tabs">
                <input class="peer absolute opacity-0" name="tabs" type="radio" :id="tab" :value="tab"
                    v-model="picked" />
                <label
                    class="block  px-6 py-2 peer-checked:text-blue-800 w-full cursor-pointer bg-white font-bold  text-black transition-colors hover:bg-blue-50"
                    :for="tab">
                    {{ tab }}
                </label>
            </div>
        </div>
        <component :is="(<any>tabs)[picked]"
            class="relative border-solid border-[5px] z-[2] border-blue-400 rounded-lg"></component>
    </div>
    <Room v-else/>
</template>
