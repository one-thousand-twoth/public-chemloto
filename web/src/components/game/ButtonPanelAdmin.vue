<script setup lang="ts">
import { WebsocketConnector } from '@/api/websocket/websocket';
import ObtainStateController from '@/state_controllers/obtain';
import TradeStateController from '@/state_controllers/trade';
import { useGameStore } from '@/stores/useGameStore';
import { useUserStore } from '@/stores/useUserStore';
import { storeToRefs } from "pinia";
import { inject, ref } from 'vue';

const ws = inject('connector') as WebsocketConnector

const gameStore = useGameStore()
const { gameState } = storeToRefs(gameStore)
const userStore = useUserStore()


function StartGame() {
    ws.Send({
        Type: 'HUB_STARTGAME',
        Name: userStore.UserInfo!.room
    })
}


const ObtainContollerr = new ObtainStateController(ws)
const TradeController = new TradeStateController(ws)


function GetElement() {
    if (!ObtainContollerr.isValid()) return;
    ObtainContollerr.getElement()
}
function sendContinue() {
    if (ObtainContollerr.isValid()) {
        ObtainContollerr.sendContinue()
    }
    if (TradeController.isValid()) {
        TradeController.sendContinue()
    }
    throw new Error("State doesn`t provide sendContinue action")
}



</script>
<template>
    <div class="flex flex-wrap justify-center gap-2">
        <template v-if="gameStore.gameState.Status === 'STATUS_WAITING'">
            <button @click="StartGame()">
                Начать игру
            </button>
        </template>
        <template v-if="gameStore.gameState.Status === 'STATUS_STARTED'">
            <!-- <div class="flex border-slate-300 hover:bg-slate-100 border-b-main border-b-2 
         shadow-large items-center cursor-pointer px-2 py-2 rounded border bg-white
         border-main-dark text-main" @click="swap()">
                <component :is="selectedBtn == 'strip' ? ShoppingBagIcon : ArrowDownCircleIcon"
                    class="size-7 lg:size-10 text-slate-500" />
            </div>


            <div class="flex rounded shadow-large border border-b-main border-b-2">
                <DesignButton class="rounded-none rounded-l" v-model="selectedRadio" value="puzzle" label="Puzzle">
                    <PuzzlePieceIcon class="size-7 lg:size-10" />
                </DesignButton>

                <DesignButton class="rounded-none rounded-r" v-model="selectedRadio" value="hand" label="Hand">
                    <HandRaisedIcon class="size-7 lg:size-10 -rotate-90" />
                </DesignButton>
            </div> -->
            <button :disabled="!ObtainContollerr.isValid()" class="text-sm" @click="GetElement()">Достать
                элемент</button>
        </template>
    </div>

    <!-- <template v-if="gameStore.gameState.Status === 'STATUS_WAITING'">
        <button @click="StartGame()">
            Начать игру
        </button>
    </template>
<template v-else>
        <template v-if="gameStore.gameState.State == 'OBTAIN'">
            <button @click="GetElement()">Достать
                элемент</button>
        </template>
<template v-if="gameStore.gameState.State == 'HAND'">
            <button disabled @click="GetElement()">Ждем
                проверки</button>
        </template>
<template v-if="gameStore.gameState.State == 'TRADE'">
            <button @click="TradeButton = !TradeButton">Обменять</button>
            <button @click="sendContinue()">Продолжить</button>
        </template>
</template>
<Modal :show="TradeButton" @close="TradeButton = false">
    <template #header>
            <h3 class="font-bold text-center">Обменять</h3>
        </template>
    <template v-if="gameStore.SelfPlayer" #body>
            <Trade :players="gameStore.gameState.Players" />
        </template>
</Modal> -->
</template>