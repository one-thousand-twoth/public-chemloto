<script setup lang="ts">
import { WebsocketConnector } from '@/api/websocket/websocket';
import obtain from "@/assets/sounds/notification.mp3";
import { ButtonPanelAdmin, ButtonPanelPlayer, CheckPlayer, FieldsTable, LeaderBoard, TradeExchange, UserElements } from '@/components/game';
import RaiseHandComp from '@/components/game/RaiseHandComp.vue';
import RoomSlots from '@/components/game/RoomSlots.vue';
import { NumKey } from '@/components/keyboard';
import { ElementImage, IconButton, IconButtonBackground, Modal, Timer, UserInfo } from '@/components/UI/';
import { Hand } from '@/models/Game';
import { Role } from '@/models/User';
import { useInterfaceStore } from '@/stores/RoomInterface';
import { useGameStore } from '@/stores/useGameStore';
import { useKeyboardStore } from '@/stores/useRaiseHand';
import { useUserStore } from '@/stores/useUserStore';
import {
    ArrowLeftStartOnRectangleIcon,
    ArrowsPointingOutIcon,
    CheckIcon,
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

const currPlayer = computed(() => {
    return GameStore.gameState.Players.find(player => player.Name === curInfoPlayer.value)
})

const curInfoPlayer = ref('')

const InterfaceStore = useInterfaceStore()
const { currentPlayerSelection } = storeToRefs(InterfaceStore)

watch(currentPlayerSelection, () => {
    if (currentPlayerSelection.value === undefined) {
        return
    }
    selectedTool.value = 'trade'
})


const showKeyboard = computed(() => {
    return InputName.value !== ''
})

let audio = new Audio(obtain);

watch(() => GameStore.gameState.Bag.LastElements, () => { audio.play() })


const selectedTool = ref<'puzzle' | 'trade'>('puzzle')
const selectedBtn = ref<"strip" | "list">('strip')

const click_selected_raiseHand = ref('')

const playerElems = GameStore.SelfPlayer



</script>
<template>

    <RoomSlots>
        <template #left>
            <div v-show="!showKeyboard" class="flex flex-col flex-1 min-h-[0]">
                <LeaderBoard class=" overflow-y-auto flex-1 min-h-[0]"></LeaderBoard>
                <!-- <FieldsTable class="w-fit self-end flex-shrink-0 mx-auto" /> -->
            </div>
            <div class="h-full flex flex-col justify-center " v-show="showKeyboard">
                <NumKey class="" />
            </div>
        </template>
        <template #center>
            <div v-if="GameStore.gameState.Status == 'STATUS_COMPLETED'" class="text-lg">
                Игра завершена
            </div>
            <Timer v-else class="w-full shadow-large " />

            <div class="relative flex-1  w-full px-2 py-4 flex flex-col gap-2   items-center justify-center 
                 bars border-0 border-b-2 border-t-2  
                ">
                <template v-if="selectedBtn === 'strip'">

                    <ElementImage class="flex-1 max-w-[80%] aspect-square" :elname="GameStore.currElement" />
                    <div class="relative  flex flex-grow-0   w-full flex-1 flex-row flex-nowrap gap-1 lg:gap-2 justify-center items-center"
                        id="lastElementsContainer">
                        <ElementImage class="w-full h-auto" v-for="el in GameStore.LastElements.slice(1, 5)"
                            :elname="el" />
                    </div>
                </template>
                <template v-if="selectedBtn === 'list'">
                    <UserElements :player="GameStore.SelfPlayer" />
                </template>
            </div>

            <template v-if="GameStore.gameState.Status !== 'STATUS_COMPLETED'">
                <ButtonPanelPlayer v-model:btn="selectedBtn" v-model:radio="selectedTool" />
            </template>
        </template>
        <template #right>
            <RaiseHandComp v-model:selectedElem="click_selected_raiseHand" v-show="selectedTool == 'puzzle'"
                v-if="GameStore.SelfPlayer" :player="GameStore.SelfPlayer" />

            <TradeExchange v-show="selectedTool == 'trade'" />
        </template>

    </RoomSlots>


</template>
