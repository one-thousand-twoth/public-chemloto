<script setup lang="ts">
import { WebsocketConnector } from '@/api/websocket/websocket';
import { ButtonPanelAdmin, ButtonPanelPlayer, CheckPlayer, LeaderBoard, RaiseHandComp, Trade, UserElements } from '@/components/game';
import { ElementImage, IconButtonBackground, Modal, Timer } from '@/components/UI/';
import { Role } from '@/models/User';
import { Hand, useGameStore } from '@/stores/useGameStore';
import { useUserStore } from '@/stores/useUserStore';
import {
    ArrowLeftStartOnRectangleIcon
} from "@heroicons/vue/24/outline";
import { computed, inject, ref, } from 'vue';
import FieldsTable from './FieldsTable.vue';

const GameStore = useGameStore()
const userStore = useUserStore()
const ws = inject('connector') as WebsocketConnector

function DisconnectGame() {
    ws.Send(
        {
            "Type": "HUB_UNSUBSCRIBE",
            "Target": "room",
            "Name": userStore.UserCreds!.room
        }
    )
}

const currPlayer = computed(() => {
    return GameStore.gameState.Players.find(player => player.Name === curInfoPlayer.value)
})

const curInfoPlayer = ref('')
const curCheckPlayer = ref<Hand>()
const RaiseHandButton = ref(false)
const TradeButton = ref(false)


</script>
<template>
    <div class="relative flex max-h-lvh flex-col items-center overflow-x-hidden">
        <main class="flex justify-between w-lvw grow gap-10 bg-gray-100">
            <!-- #region LEFT -->
            <div class="flex flex-col m-3 w-[20%] gap-2">
                <div class="bars p-3 min-w-[8.5rem]  grow-[1] bg-gray-50">
                    <LeaderBoard @selectPlayer="(name: string) => { curInfoPlayer = name }"></LeaderBoard>
                </div>
                <IconButtonBackground v-if="!GameStore.gameState.Started || userStore.UserCreds?.role != Role.Player"
                    class="w-full bg-red-700 text-white  rounded-lg" :icon="ArrowLeftStartOnRectangleIcon"
                    @click="DisconnectGame()">Выйти</IconButtonBackground>
            </div>
            <!-- #endregion LEFT -->
            <!-- #region CENTER -->
            <div class="flex flex-col items-center  h-lvh max-w-[900px] gap-2 p-5 pb-60 grow-[3]">
                <Timer />
                <div class=" h-auto w-full max-w-[50lvh] gap-2 flex flex-wrap items-center justify-center">
                    <ElementImage class="grow-[2] center" :elname="GameStore.currElement" />
                    <div class="flex flex-col flex-wrap gap-1 items-center" id="lastElementsContainer">
                        <!-- <h3 class="mb-2">Последние элементы</h3> -->
                        <ElementImage v-for="el in GameStore.LastElements" :elname="el" />
                    </div>
                </div>
                <FieldsTable />
                <ButtonPanelAdmin v-if="userStore.UserCreds?.role != Role.Player" />
                <ButtonPanelPlayer v-else />

            </div>
            <!-- #endregion CENTER -->
            <!-- #region RIGHT -->
            <div class="bars min-w-[135px] p-3 bg-gray-50 w-[20%] m-3">
                <h2>Поднятые руки</h2>
                <ul class="list-none p-0 font-bold m-0">
                    <li @click="curCheckPlayer = pl" class="break-words flex justify-between items-center p-2 hover:underline rounded-md my-2 mx-0
                    border-solid border-2 border-gray-600 m-3" v-for="pl in GameStore.gameState.RaisedHands">
                        {{ pl.Player.Name }} - {{ pl.Field }}
                    </li>
                </ul>
            </div>
            <!-- #endregion RIGHT -->
        </main>


      


        <Modal :show="currPlayer !== undefined" @close="curInfoPlayer = ''">
            <template #header>
                <h3 class="font-bold text-center">Информация о игроке {{ curInfoPlayer }}</h3>
            </template>
            <template #body>
                <UserElements v-if="currPlayer" :player="currPlayer" />
            </template>
        </Modal>
        <Modal v-if="userStore.UserCreds?.role != Role.Player" :show="curCheckPlayer !== undefined"
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
