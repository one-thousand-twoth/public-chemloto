<script setup lang="ts">
import { WebsocketConnector } from '@/api/websocket/websocket';
import obtain from "@/assets/sounds/notification.mp3";
import { IconButton, IconButtonBackground } from '@/components/UI';
import { Role } from '@/models/User';
import { useGameStore } from '@/stores/useGameStore';
import { useKeyboardStore } from '@/stores/useRaiseHand';
import { useUserStore } from '@/stores/useUserStore';
import {
    ArrowLeftStartOnRectangleIcon,
    ArrowsPointingInIcon,
    ArrowsPointingOutIcon,
    EllipsisVerticalIcon
} from "@heroicons/vue/24/outline";
import { useFullscreen } from '@vueuse/core';
import { storeToRefs } from 'pinia';
import { computed, inject, ref, useTemplateRef, watch } from 'vue';


const GameStore = useGameStore()
const userStore = useUserStore()
const keyboardStore = useKeyboardStore()
const { InputName } = storeToRefs(keyboardStore)
const ws = inject('connector') as WebsocketConnector

function DisconnectGame() {
    ws.Send(
        {
            "Type": "HUB_UNSUBSCRIBE",
            "Target": "room",
            "Name": userStore.UserInfo.room
        }
    )
}



const AdditionallyButton = ref(false)


let audio = new Audio(obtain);

watch(() => GameStore.gameState.Bag.LastElements, () => { audio.play() })

// @ts-ignore
const el = useTemplateRef('el')
const { toggle, isFullscreen } = useFullscreen(el)

const FullscreenIcon = computed(() => { return isFullscreen.value ? ArrowsPointingInIcon : ArrowsPointingOutIcon })

const click_selected_raiseHand = ref('')


function click(e: Event) {
    click_selected_raiseHand.value = ''
    InputName.value = ""

}

</script>
<template>

    <div @click="click">
        <div
            class="relative p-2 gap-4 md:gap-12 grid grid-cols-[1.2fr_1.0fr_2fr] h-lvh grid-rows-[calc(100lvh-1rem)] bg-bg  w-dvw  items-stretch">

            <!-- #region LEFT -->
            <div class="relative flex flex-col gap-2">
                <div
                    class="bars shadow-large grow flex flex-col justify-stretch min-h-[0] p-3 min-w-[8.5rem] bg-gray-50">
                    <slot name="left" />
                </div>
                <IconButtonBackground v-if="GameStore.gameState.Status !== 'STATUS_STARTED'"
                    class="w-full bg-red-700 text-white  rounded-lg" :icon="ArrowLeftStartOnRectangleIcon"
                    @click="DisconnectGame()">
                    Выйти</IconButtonBackground>
            </div>
            <!-- #endregion LEFT -->

            <!-- #region CENTER -->
            <div class=" relative flex flex-col gap-2 lg:gap-20  items-center justify-center 
             h-[100%]  lg:h-[80%]  lg:pb-0
             ">
                <slot name="center" />
            </div>
            <!-- #endregion CENTER -->

            <!-- #region RIGHT -->
            <div class='relative ml-8 md:ml-0 flex flex-col h-full gap-2'>
                <IconButton class="absolute left-[-45px]" :icon="FullscreenIcon" @click="toggle" />
                <div class="bars  shadow-large p-3 min-w-[8.5rem]  grow-[1] bg-gray-50">

                    <slot name="right"></slot>

                </div>
                <!-- <IconButtonBackground v-if="userStore.UserInfo.role != Role.Player"
                    class="w-full  bg-blue-500 text-white  rounded-lg" :icon="EllipsisVerticalIcon"
                    @click="AdditionallyButton = !AdditionallyButton">Дополнительно</IconButtonBackground> -->
                <!-- <div v-if="AdditionallyButton"
                    class="absolute z-[2]  top-[-14px] border-solid border-2 text-sm border-blue-400 rounded-lg rounded-b-none p-3 ">
                    <div @click="EXITGame()" class="underline hover:text-blue-500">Закрыть игру</div>
                </div> -->
            </div>
            <!-- #endregion RIGHT -->
        </div>
    </div>
</template>
