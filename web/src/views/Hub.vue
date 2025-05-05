<script setup lang="ts">
import ClientList from '@/components/ClientList.vue';
import RoomList from '@/components/RoomList.vue';
import IconButton from '@/components/UI/IconButton.vue';
// import { emojiRole } from '@/models/User';s
import UserInfo from '@/components/UI/UserInfo.vue';
import { useGameStore } from '@/stores/useGameStore';
import { useUserStore } from '@/stores/useUserStore';
import {
    ArrowLeftStartOnRectangleIcon
} from "@heroicons/vue/24/outline";
import { storeToRefs } from 'pinia';
import { ref } from 'vue';
import Room from './Room.vue';
const tabs = {
    "Комнаты": RoomList,
    "Игроки": ClientList,
}

const picked = ref('Комнаты')
const userStore = useUserStore()
const gameStore = useGameStore()
const { connected, fetching } = storeToRefs(userStore)
async function Remove() {
    await userStore.Remove(userStore.UserCreds!.username)
    window.location.reload();
}
</script>
<template>
    <div v-if="fetching"> Загрузка...</div>
    <template v-else>
        <div class="relative md:p-8 flex flex-col h-lvh " v-if="!connected">
            <div class="flex px-10">
                <div class="relative mt-2 px-10 flex justify-center flex-wrap bg-opacity-0 overflow-y-hidden bars border-0 border-t-2">
                    <div class="relative  z-0 has-[:checked]:z-[2] bg-white   border-main
                    has-[:checked]:border-solid has-[:checked]:border-2 has-[:checked]:border-t-0 hover:bg-blue-50
                     border-blue-400 rounded-br-lg rounded-bl-lg border-r " v-for="(_, tab) in tabs">
                        <input class="peer absolute opacity-0" name="tabs" type="radio" :id="tab" :value="tab"
                            v-model="picked" />
                        <label
                            class="block h-full px-6 py-2 peer-checked:text-blue-800 w-full cursor-pointer font-bold  text-black transition-colors "
                            :for="tab">
                            {{ tab }}
                        </label>
                    </div>
                </div>
                <div
                    class="flex ml-auto items-center relative pl-2 py-[0.1rem] mt-2 mr-2 mb:p-8 text-lg font-semibold bars w-fit self-end">
                    <UserInfo :role="userStore.UserInfo.role" :name="userStore.UserCreds?.username ?? ''" />
                    <IconButton :icon="ArrowLeftStartOnRectangleIcon" @click="Remove()" />
                </div>
                
            </div>
            <div class="">
                <component :is="(<any>tabs)[picked]" class="relative"></component>
            </div>
        </div>
        <Room v-else />
    </template>
</template>
