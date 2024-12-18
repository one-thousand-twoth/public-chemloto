<script setup lang="ts">
import { WebsocketConnector } from '@/api/websocket/websocket';
import obtain from "@/assets/sounds/notification.mp3";
import { ButtonPanelAdmin, ButtonPanelPlayer, CheckPlayer, LeaderBoard, UserElements } from '@/components/game';
import { ElementImage, IconButtonBackground, Modal, Timer, UserInfo } from '@/components/UI/';
import { Role } from '@/models/User';
import { Hand, useGameStore } from '@/stores/useGameStore';
import { useUserStore } from '@/stores/useUserStore';
import {
    ArrowLeftStartOnRectangleIcon,
    CheckIcon,
    EllipsisVerticalIcon
} from "@heroicons/vue/24/outline";
import { computed, inject, ref, watch, } from 'vue';
import FieldsTable from './FieldsTable.vue';

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


const currPlayer = computed(() => {
    return GameStore.gameState.Players.find(player => player.Name === curInfoPlayer.value)
})

const curInfoPlayer = ref('')
const curCheckPlayer = ref<Hand>()
const AdditionallyButton = ref(false)
let audio = new Audio(obtain);

watch(() => GameStore.gameState.Bag.LastElements, () => { audio.play() })
</script>
<template>
    <div class="relative flex bg-gray-100  flex-col items-center overflow-x-hidden">
        <div class="flex justify-between w-lvw grow gap-10 ">
            <!-- #region LEFT -->
            <div class="flex flex-col m-3 w-[20%] gap-2">
                <div class="bars p-3 min-w-[8.5rem]  grow-[1] bg-gray-50">
                    <LeaderBoard @selectPlayer="(name: string) => { curInfoPlayer = name }"></LeaderBoard>
                </div>
                <IconButtonBackground v-if="GameStore.gameState.Status !== 'STATUS_STARTED'"
                    class="w-full bg-red-700 text-white  rounded-lg" :icon="ArrowLeftStartOnRectangleIcon"
                    @click="DisconnectGame()">Выйти</IconButtonBackground>
            </div>
            <!-- #endregion LEFT -->
            <!-- #region CENTER -->
            <div class="flex flex-col items-center  h-lvh max-w-[900px] gap-2 p-5 pb-60 grow-[3]">
                <div  v-if="GameStore.gameState.Status == 'STATUS_COMPLETED'" class="text-lg">
                    Игра завершена
                </div>
                <Timer v-else />

                <div class=" h-auto w-full max-w-[50lvh] gap-2 flex flex-wrap items-center justify-center">
                    <ElementImage class="grow-[2] center" :elname="GameStore.currElement" />
                    <div class="flex flex-col flex-wrap gap-1 items-center" id="lastElementsContainer">
                        <ElementImage v-for="el in GameStore.LastElements" :elname="el" />
                    </div>
                </div>

                <FieldsTable />
                <template v-if="GameStore.gameState.Status !=='STATUS_COMPLETED'">
                <ButtonPanelAdmin v-if="userStore.UserInfo.role != Role.Player" />
                <ButtonPanelPlayer v-else />
            </template>
            </div>
            <!-- #endregion CENTER -->
            <!-- #region RIGHT -->
        <div class='relative flex flex-col m-3 w-[20%] gap-2'>
                <div class="bars p-3 min-w-[8.5rem]  grow-[1] bg-gray-50">
                    <h2 class="text-clip overflow-hidden ">Поднятые руки</h2>
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
                <div v-if="AdditionallyButton" class="relative z-[2]  top-[14px] border-solid border-2 text-sm border-blue-400 rounded-lg rounded-b-none p-3 ">
                    <div @click="EXITGame()" class="underline hover:text-blue-500">Закрыть игру</div>
                </div>
                <IconButtonBackground v-if="userStore.UserInfo.role != Role.Player"
                        class="w-full z-[3] bg-blue-500 text-white  rounded-lg" :icon="EllipsisVerticalIcon"
                        @click="AdditionallyButton = !AdditionallyButton">Дополнительно</IconButtonBackground>
        </div>
            <!-- #endregion RIGHT -->
       
        </div>


        <div class=" h-96"></div>


        <Modal :show="currPlayer !== undefined" @close="curInfoPlayer = ''">
            <template #header>
                <h3 class="font-bold text-center">Информация о игроке {{ curInfoPlayer }}</h3>
            </template>
            <template #body>
                <UserElements v-if="currPlayer" :player="currPlayer" />
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

    </div>
</template>

<style scoped>
body {
    background-color: #3b4e5e;
}
</style>
