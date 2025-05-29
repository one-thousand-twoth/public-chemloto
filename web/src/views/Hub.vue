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

type TabsNames = keyof typeof tabs;

const picked = ref<TabsNames>('Комнаты')
const userStore = useUserStore()

const { connected, fetching } = storeToRefs(userStore)
async function Remove() {
    await userStore.Remove(userStore.UserCreds!.username)
    window.location.reload();
}
window.onload = function () {

}

</script>
<template>
    <div v-if="fetching">Загрузка...</div>
    <template v-else>
        <div v-if="!connected" class="flex flex-col gap-4 h-lvh p-2 overflow-y-auto w-full
         sm:max-w-[80%] lg:max-w-[60%] mx-auto
         " style="scrollbar-gutter: both-edges">
            <!-- Header with tabs -->
            <div class="flex">
                <!-- Desktop tabs -->
                <div class=" hidden sm:flex relative mt-2 px-10 justify-center flex-wrap bg-opacity-0
                         overflow-y-hidden bars border-0 border-t-2">
                    <div class="relative  bg-white border-main
                        has-[:checked]:border-solid has-[:checked]:border-2 has-[:checked]:border-t-0 hover:bg-blue-50
                         border-blue-400 rounded-br-lg rounded-bl-lg border-r " v-for="(_, tabname) in tabs">
                        <input class="peer absolute opacity-0" name="tabs" type="radio" :id="tabname" :value="tabname"
                            v-model="picked" />
                        <label class="block h-full px-6 py-2 peer-checked:text-blue-800 w-full
                                 cursor-pointer font-bold  text-black transition-colors " :for="tabname">
                            {{ tabname }}
                        </label>
                    </div>
                </div>

                <!-- Mobile dropdown -->
                <div class="sm:hidden ">
                    <select class="bars" v-model="picked">
                        <option v-for="(_, tab) in tabs" :value="tab">{{ tab }}</option>
                    </select>
                </div>

                <!-- User info and logout -->
                <div class="flex bars px-2 ml-auto">
                    <UserInfo :role="userStore.UserInfo.role" :name="userStore.UserCreds?.username ?? ''" />
                    <IconButton :icon="ArrowLeftStartOnRectangleIcon" @click="Remove()" />
                </div>
            </div>

            <!-- Content area -->
            <div class="content-area">
                <component :is="tabs[picked]"></component>
            </div>

            <div> </div>
        </div>

        <Room v-else />
    </template>
</template>
