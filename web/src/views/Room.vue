<script setup lang="ts">
import { WebsocketConnector } from '@/api/websocket/websocket';
import obtain from "@/assets/sounds/notification.mp3";
import { ButtonPanelAdmin, ButtonPanelPlayer, CheckPlayer, DesignButton, LeaderBoard, UserElements } from '@/components/game';
import { ElementImage, IconButton, IconButtonBackground, Modal, Timer, UserInfo } from '@/components/UI/';
import { Hand } from '@/models/Game';
import { Role } from '@/models/User';
import { useGameStore } from '@/stores/useGameStore';
import { useUserStore } from '@/stores/useUserStore';
import {
    ArrowLeftStartOnRectangleIcon,
    ArrowsPointingOutIcon,
    CheckIcon,
    EllipsisVerticalIcon
} from "@heroicons/vue/24/outline";
import { computed, inject, ref, watch, } from 'vue';
import FieldsTable from './FieldsTable.vue';

import { useFullscreen } from '@vueuse/core';
import { useTemplateRef } from 'vue';


const GameStore = useGameStore()
const userStore = useUserStore()
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
function EXITGame() {
    ws.Send(
        {
            "Type": "HUB_EXITGAME",
            "Name": userStore.UserInfo.room
        }
    )
}
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

const currPlayer = computed(() => {
    return GameStore.gameState.Players.find(player => player.Name === curInfoPlayer.value)
})

const curInfoPlayer = ref('')
const score = ref(0)
const curCheckPlayer = ref<Hand>()
const AdditionallyButton = ref(false)
const RemainsButton = ref(false)

let audio = new Audio(obtain);

watch(() => GameStore.gameState.Bag.LastElements, () => { audio.play() })


const el = useTemplateRef('el')
const { toggle } = useFullscreen(el)
</script>
<template>

    <div ref='el'>
        <div
            class="relative p-2 gap-12 grid grid-cols-[1.5fr_1.5fr_2fr] h-[100svh]  overflow-y-scroll bg-bg  w-dvw  items-center">



            <!-- #region LEFT -->
            <div class="relative   flex h-full flex-col gap-2">
                <div class="bars shadow-large p-3 min-w-[8.5rem]  grow-[1] bg-gray-50">
                    <LeaderBoard @selectPlayer="(name: string) => { curInfoPlayer = name }"></LeaderBoard>
                </div>
                <IconButtonBackground v-if="GameStore.gameState.Status !== 'STATUS_STARTED'"
                    class="w-full bg-red-700 text-white  rounded-lg" :icon="ArrowLeftStartOnRectangleIcon"
                    @click="DisconnectGame()">Выйти</IconButtonBackground>
            </div>
            <!-- #endregion LEFT -->


            <!-- #region CENTER -->
            <div class=" relative flex flex-col gap-4 lg:gap-20  items-center justify-center 
             h-[100%]
             lg:h-[80%]

            pb-4
            lg:pb-0

             ">
                <!-- <div class="w-full p-6  bg-red-500"> Привет</div>
                <div class="relative flex-col flex flex-initial items-center justify-center flex-shrink p-2 h-[30vh] w-full bg-red-500">
                   
                    <ElementImage class="flex-1 max-w-[60%] aspect-square" :elname="GameStore.currElement" />
                    <div class="relative  flex flex-grow-0   w-full flex-1 flex-row flex-nowrap gap-3 items-center" id="lastElementsContainer">
                        <ElementImage class="w-full h-auto" v-for="el in GameStore.LastElements.slice(1, 5)" :elname="el" />
                    </div>
                </div>
                <div class="w-full p-6  bg-red-500"> Привет</div> -->
                <div v-if="GameStore.gameState.Status == 'STATUS_COMPLETED'" class="text-lg">
                    Игра завершена
                </div>
                <Timer class="w-full shadow-large " v-else />

                <div class="relative flex-1  w-full px-2 py-4 flex flex-col gap-2   items-center justify-center 
                 bars border-0 border-b-2 border-t-2  
                ">
                    <ElementImage class="flex-1 max-w-[80%] aspect-square" :elname="GameStore.currElement" />
                    <div class="relative  flex flex-grow-0   w-full flex-1 flex-row flex-nowrap gap-3 items-center" id="lastElementsContainer">
                        <ElementImage class="w-full h-auto" v-for="el in GameStore.LastElements.slice(1, 5)" :elname="el" />
                    </div>
                </div>
                
                <template v-if="GameStore.gameState.Status !== 'STATUS_COMPLETED'">
                    <ButtonPanelAdmin v-if="userStore.UserInfo.role != Role.Player" />
                    <ButtonPanelPlayer v-else />
                </template>

            </div>
            <!-- #endregion CENTER -->


            <!-- #region RIGHT -->
            <div class='relative  flex flex-col h-full gap-2'>
                <IconButton class="absolute left-[-45px]" :icon="ArrowsPointingOutIcon" @click="toggle" />
                <div class="bars  shadow-large p-3 min-w-[8.5rem]  grow-[1] bg-gray-50">
                    <ul class="list-none p-0 font-bold m-0">
                        <li @click="curCheckPlayer = pl" class="break-words flex justify-between items-center p-2 hover:underline rounded-md my-2 mx-0
                            border-solid border-2 border-gray-600 m-3" v-for="pl in GameStore.gameState.RaisedHands">
                            <div class=" inline-flex">
                                <CheckIcon v-if="pl.Checked" class="text-lg size-6" />
                                <UserInfo :role="pl.Player.Role" :name="pl.Player.Name" />
                            </div>
                            {{ pl.Field }}
                        </li>
                    </ul>

                </div>
                <div v-if="AdditionallyButton"
                    class="relative z-[2]  top-[14px] border-solid border-2 text-sm border-blue-400 rounded-lg rounded-b-none p-3 ">
                    <div @click="EXITGame()" class="underline hover:text-blue-500">Закрыть игру</div>
                </div>
                <IconButtonBackground v-if="userStore.UserInfo.role != Role.Player"
                    class="w-full z-[3] bg-blue-500 text-white  rounded-lg" :icon="EllipsisVerticalIcon"
                    @click="AdditionallyButton = !AdditionallyButton">Дополнительно</IconButtonBackground>
            </div>
            <!-- #endregion RIGHT -->







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
                    <CheckPlayer :player="curCheckPlayer" />
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

        </div>
    </div>
</template>

<style scoped>
body {
    background-color: #3b4e5e;
}
</style>
