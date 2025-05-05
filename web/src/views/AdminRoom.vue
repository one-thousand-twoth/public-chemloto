<script setup lang="ts">
import { WebsocketConnector } from '@/api/websocket/websocket';
import obtain from "@/assets/sounds/notification.mp3";
import { ButtonPanelAdmin, ButtonPanelPlayer, CheckPlayer, FieldsTable, LeaderBoard, UserElements } from '@/components/game';
import RaiseHandComp from '@/components/game/RaiseHandComp.vue';
import RoomSlots from '@/components/game/RoomSlots.vue';
import { NumKey } from '@/components/keyboard';
import { ElementImage, IconButton, IconButtonBackground, Modal, Timer, UserInfo } from '@/components/UI/';
import { Hand } from '@/models/Game';
import { Role } from '@/models/User';
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
import { openModal } from 'jenesius-vue-modal';
import { storeToRefs } from 'pinia';
import { computed, inject, ref, useTemplateRef, watch } from 'vue';


const GameStore = useGameStore()
const userStore = useUserStore()
const keyboardStore = useKeyboardStore()
const { InputName } = storeToRefs(keyboardStore)
const ws = inject('connector') as WebsocketConnector

function AddScore(score: number, player: string) {
    ws.Send(
        {
            "Type": "ENGINE_ACTION",
            "Action": "AddScore",
            "Score": score,
            "Player": player,
        }
    )
}

function EXITGame() {
    const ok = confirm(`Вы уверены что хотите закрыть игру?
     Игроки смогут выйти из комнаты`)
    if (!ok) {
        return
    }
    ws.Send(
        {
            "Type": "HUB_EXITGAME",
            "Name": userStore.UserInfo.room
        }
    )
}

const currPlayer = computed(() => {
    return GameStore.gameState.Players.find(player => player.Name === curInfoPlayer.value)
})

const curInfoPlayer = ref('')
const score = ref(0)
const curCheckPlayer = ref<Hand>()

const RemainsButton = ref(false)


const showKeyboard = computed(() => {
    return InputName.value !== ''
})





const selectedTool = ref<'puzzle' | 'trade'>('puzzle')
const selectedBtn = ref<"strip" | "list">('strip')



</script>
<template>

    <RoomSlots>

        <template #left>
            <div v-show="!showKeyboard" class="flex flex-col flex-1 min-h-[0]">
                <LeaderBoard class=" overflow-y-auto flex-1 min-h-[0]" @selectPlayer="(name: string) => { curInfoPlayer = name }"></LeaderBoard>
                <FieldsTable class="w-fit self-end flex-shrink-0 mx-auto" />
            </div>
            <div class="h-full flex flex-col justify-center " v-show="showKeyboard">
                <NumKey class="" />
            </div> 
            <!-- <div class="relative flex flex-col flex-1" v-show="!showKeyboard"> -->
                <!-- <LeaderBoard class=" overflow-y-scroll flex-1 min-h-[0]" @selectPlayer="(name: string) => { curInfoPlayer = name }"></LeaderBoard> -->
                <!-- <div class="flex grow flex-col overflow-scroll flex-1 min-h-[min-content] bg-blue-300">
                    <div class="bars">Привет</div>
                    <div class="bars">Привет</div>
                    <div class="bars">Привет</div>
                    <div class="bars">Привет</div>
                    <div class="bars">Привет</div>
                    <div class="bars">Привет</div>
                    <div class="bars">Привет</div>
                    <div class="bars">Привет</div>
                    <div class="bars">Привет</div>
                    <div class="bars">Привет</div>
                    <div class="bars">Привет</div>
                </div> -->
                <!-- <FieldsTable class="w-fit self-end flex-shrink-0 mx-auto" /> -->
            <!-- </div> -->
        </template>

        <template #center>
            <div v-if="GameStore.gameState.Status == 'STATUS_COMPLETED'" class="text-lg">
                Игра завершена
            </div>
            <Timer v-else class="w-full shadow-large " />

            <div class="relative flex-1  w-full px-2 py-4 flex flex-col gap-2   items-center justify-center 
                 bars border-0 border-b-2 border-t-2  
                ">
                <ElementImage class="flex-1 max-w-[80%] aspect-square" :elname="GameStore.currElement" />
                <div class="relative  flex flex-grow-0   w-full flex-1 flex-row flex-nowrap gap-1 lg:gap-2 justify-center items-center"
                    id="lastElementsContainer">
                    <ElementImage class="w-full h-auto" v-for="el in GameStore.LastElements.slice(1, 5)" :elname="el" />
                </div>
            </div>

            <template v-if="GameStore.gameState.Status !== 'STATUS_COMPLETED'">
                <ButtonPanelAdmin class="mb-2" v-model:btn="selectedBtn" v-model:radio="selectedTool" />
            </template>
        </template>

        <template #right>
            <div class="h-full flex flex-col">
                <p class="text-sm md:text-base md:font-semibold">Поднятые руки:</p>
                <div class="flex-1">
                    <ul class="list-none p-0 font-bold m-0">
                        <li @click="openModal(CheckPlayer, { hand: hand })" class=" cursor-pointer break-words flex justify-between items-center p-2 hover:underline rounded-md my-2 mx-0
                                border-solid border-2 border-gray-600 m-3"
                            v-for="hand in GameStore.gameState.RaisedHands">
                            <div class=" inline-flex">
                                <CheckIcon v-if="hand.Checked" class="text-lg size-6" />
                                <UserInfo :role="hand.Player.Role" :name="hand.Player.Name" />
                            </div>
                            {{ hand.Field }}
                        </li>
                    </ul>
                    <div class=""></div>
                </div>
                <details>
                    <summary class="mb-2">
                        <!-- <div
                            class="inline-flex w-5 h-5 text-[1rem] items-center justify-center rounded-full bg-blue-400 text-white  font-bold">
                            {{ Object.entries(selfStock.Requests).length }}
                        </div> -->
                        Дополнительно
                    </summary>
                    <!--                <IconButton class="" :icon="XMarkIcon" />-->
                    <button @click="EXITGame()" class=" bg-denial border-gray-600">Закрыть игру</button>
                </details>
            </div class="flex">
        </template>

    </RoomSlots>

    <Modal :show="currPlayer !== undefined" @close="curInfoPlayer = ''; score = 0">
        <template #header>
            <h3 class="font-bold text-center">Информация о игроке {{ curInfoPlayer }}</h3>
        </template>
        <template #body>
            <UserElements v-if="currPlayer" :player="currPlayer" />

            <form v-if="userStore.UserInfo.role != Role.Player" class="flex flex-col gap-1"
                @submit.prevent="AddScore(score, curInfoPlayer)">
                <div class="text-lg">Добавить очки: </div>
                <input type="number" min="0" v-model="score" />
                <button type="submit">Отправить</button>
            </form>
        </template>

    </Modal>
    <Modal v-if="userStore.UserInfo.role != Role.Player" :show="curCheckPlayer !== undefined"
        @close="curCheckPlayer = undefined">
        <template #header>
            <h3 class="font-bold text-center">Проверка структуры {{ curCheckPlayer?.Player.Name }}</h3>
        </template>
        <template v-if="curCheckPlayer !== undefined" #body>
            <CheckPlayer :hand="curCheckPlayer" />
        </template>
    </Modal>

    <Modal :show="RemainsButton !== false" @close="RemainsButton = false">
        <template #header>
            <h3 class="font-bold text-center"> Выпавшие элементы</h3>
        </template>
        <template #body>
            <div class="flex flex-wrap  justify-start my-3 gap-3">
                <div v-for="key in Object.keys(GameStore.gameState.Bag.DraftedElements)"
                    class="flex items-center mb-3 gap-1">
                    <ElementImage class="w-8" :elname="key" />
                    <div>{{ GameStore.gameState.Bag.DraftedElements[key] }}</div>
                </div>
            </div>
        </template>
    </Modal>


</template>
