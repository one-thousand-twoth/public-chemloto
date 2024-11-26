<script setup lang="ts">
import ClientList from '@/components/ClientList.vue';
import RoomList from '@/components/RoomList.vue';
import IconButton from '@/components/UI/IconButton.vue';
// import { emojiRole } from '@/models/User';s
import UserInfo from '@/components/UI/UserInfo.vue';
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
const { connected } = storeToRefs(userStore)
async function Remove() {
    await userStore.Remove(userStore.UserCreds!.username)
    window.location.reload();
}
</script>
<template>
    <div class="relative md:p-8 flex flex-col " v-if="!connected">
        <div
            class="flex  items-center relative pl-2 py-[0.1rem]  mb:p-8 text-lg font-semibold border-[5px] border-blue-400 rounded-lg w-fit self-end">
            <UserInfo :role="userStore.UserCreds!.role" :name="userStore.UserCreds!.username" />
            <IconButton :icon="ArrowLeftStartOnRectangleIcon" @click="Remove()" />
        </div>
        <div class="">
            <div class="relative pl-[8px] top-[0.5rem] flex flex-wrap bg-opacity-0 overflow-y-hidden">
                <div class="relative  z-0 has-[:checked]:z-[2] bg-white has-[:checked]:border-solid has-[:checked]:border-[5px] has-[:checked]:border-b-0  border-blue-400 rounded-br-lg rounded-bl-lg border-r "
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
                class="relative border-solid border-[5px]  border-blue-400 rounded-lg"></component>
        </div>
    </div>
    <Room v-else />
</template>
